package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"tradeshield-backend/internal/domain"
	"tradeshield-backend/internal/handler"
	"tradeshield-backend/internal/repository"
)

func setupTestRouter() http.Handler {
	repository.DB.SeedData()
	r := chi.NewRouter()
	h := handler.NewHandler(repository.DB)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		})
		r.Post("/auth/login", h.Login)
		r.Post("/auth/register", h.Register)
		r.Get("/directory", h.ListDirectory)
		r.Get("/escrow/summary", h.GetSummary)

		r.Get("/proposals", h.ListProposals)
		r.Post("/proposals", h.CreateProposal)
		r.Get("/proposals/{id}", h.GetProposal)
		r.Post("/proposals/{id}/approve", h.ApproveProposal)
		r.Post("/proposals/{id}/request-modification", h.RequestModification)
		r.Post("/proposals/{id}/ship", h.ShipProposal)
		r.Post("/proposals/{id}/confirm-delivery", h.ConfirmDeliveryProposal)

		r.Get("/contracts", h.ListContracts)
		r.Post("/contracts", h.CreateContract)
		r.Get("/contracts/{id}", h.GetContract)
		r.Post("/milestones/{id}/approve-release", h.ApproveAndReleaseMilestone)

		r.Get("/connections", h.ListConnections)
		r.Get("/connections/requests", h.ListConnectionRequests)
		r.Post("/connections/requests", h.SendConnectionRequest)
		r.Post("/connections/requests/{id}/accept", h.AcceptConnectionRequest)

		r.Get("/disputes", h.ListDisputes)
		r.Post("/disputes", h.RaiseDispute)
		r.Post("/disputes/{id}/arbitrate", h.ArbitrateDispute)

		r.Get("/ledger/accounts", h.ListLedgerAccounts)
		r.Get("/ledger/journals", h.ListLedgerJournals)

		r.Get("/support/tickets", h.ListSupportTickets)
		r.Post("/support/tickets", h.SubmitSupportTicket)

		r.Post("/mock/penny-drop", h.MockPennyDrop)
		r.Get("/mock/gst-verify/{gstin}", h.MockGSTVerify)

		// Admin routes
		r.Get("/admin/stats", h.GetAdminStats)
		r.Get("/admin/buyers", h.ListAdminBuyers)
		r.Get("/admin/suppliers", h.ListAdminSuppliers)
		r.Get("/admin/finance", h.GetAdminFinance)
		r.Post("/admin/kyc/{userId}/verify", h.VerifyAdminKYC)
		r.Post("/admin/settlements/{id}/payout", h.ProcessAdminSettlement)
	})
	return r
}

func TestHealthCheck(t *testing.T) {
	r := setupTestRouter()
	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestProposalDualApprovalLifecycle(t *testing.T) {
	r := setupTestRouter()

	// 1. Create Proposal
	propPayload := handler.CreateProposalReq{
		BuyerID:          "user_buyer_1",
		SupplierID:       "user_supplier_1",
		Amount:           350000,
		Currency:         "INR",
		Terms:            "partial",
		DeliveryTimeline: "2026-07-20",
		Notes:            "50% advance in escrow, 50% on receipt of goods",
	}
	body, _ := json.Marshal(propPayload)
	req := httptest.NewRequest("POST", "/api/v1/proposals", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected proposal creation 201, got %d. Body: %s", w.Code, w.Body.String())
	}

	var created domain.Proposal
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	createdData, _ := json.Marshal(res["data"])
	json.Unmarshal(createdData, &created)

	// 2. Supplier Approves
	appSupReq := httptest.NewRequest("POST", "/api/v1/proposals/"+created.ID+"/approve", bytes.NewReader([]byte(`{"role":"supplier"}`)))
	appSupReq.Header.Set("Content-Type", "application/json")
	wSup := httptest.NewRecorder()
	r.ServeHTTP(wSup, appSupReq)
	if wSup.Code != http.StatusOK {
		t.Fatalf("Expected supplier approve 200, got %d", wSup.Code)
	}

	// 3. Buyer Approves -> Locks deal to Confirmed & Escrow Held
	appBuyReq := httptest.NewRequest("POST", "/api/v1/proposals/"+created.ID+"/approve", bytes.NewReader([]byte(`{"role":"buyer"}`)))
	appBuyReq.Header.Set("Content-Type", "application/json")
	wBuy := httptest.NewRecorder()
	r.ServeHTTP(wBuy, appBuyReq)
	if wBuy.Code != http.StatusOK {
		t.Fatalf("Expected buyer approve 200, got %d", wBuy.Code)
	}

	// 4. Supplier Ships with LR
	shipReq := httptest.NewRequest("POST", "/api/v1/proposals/"+created.ID+"/ship", bytes.NewReader([]byte(`{"lr_number":"LR-9988","transporter_name":"V-Trans","proof_url":"https://proof.pdf"}`)))
	shipReq.Header.Set("Content-Type", "application/json")
	wShip := httptest.NewRecorder()
	r.ServeHTTP(wShip, shipReq)
	if wShip.Code != http.StatusOK {
		t.Fatalf("Expected ship 200, got %d", wShip.Code)
	}

	// 5. Buyer Confirms Delivery -> Releases Payout
	delivReq := httptest.NewRequest("POST", "/api/v1/proposals/"+created.ID+"/confirm-delivery", nil)
	wDeliv := httptest.NewRecorder()
	r.ServeHTTP(wDeliv, delivReq)
	if wDeliv.Code != http.StatusOK {
		t.Fatalf("Expected delivery confirm 200, got %d", wDeliv.Code)
	}
}

func TestDirectoryAndConnectionRequests(t *testing.T) {
	r := setupTestRouter()

	// List directory
	reqDir := httptest.NewRequest("GET", "/api/v1/directory?q=Rajesh", nil)
	wDir := httptest.NewRecorder()
	r.ServeHTTP(wDir, reqDir)
	if wDir.Code != http.StatusOK {
		t.Fatalf("Expected directory 200, got %d", wDir.Code)
	}

	// Send connection request
	sendReq := httptest.NewRequest("POST", "/api/v1/connections/requests", bytes.NewReader([]byte(`{"from_id":"user_buyer_1","to_id":"user_supplier_2"}`)))
	sendReq.Header.Set("Content-Type", "application/json")
	wSend := httptest.NewRecorder()
	r.ServeHTTP(wSend, sendReq)
	if wSend.Code != http.StatusCreated {
		t.Fatalf("Expected connection request 201, got %d", wSend.Code)
	}
}

func TestDisputeAndArbitration(t *testing.T) {
	r := setupTestRouter()

	// 1. Raise Dispute
	dispPayload := handler.RaiseDisputeReq{
		ContractID:  "cntr_2026_089",
		MilestoneID: "ms_089_3",
		Reason:      "Dimensional variation in valve housings",
		ClaimAmount: 1000000,
	}
	body, _ := json.Marshal(dispPayload)
	req := httptest.NewRequest("POST", "/api/v1/disputes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Expected dispute creation 201, got %d", w.Code)
	}

	var disp domain.Dispute
	var res map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	dispData, _ := json.Marshal(res["data"])
	json.Unmarshal(dispData, &disp)

	// 2. Arbitrate Dispute
	arbReq := httptest.NewRequest("POST", "/api/v1/disputes/"+disp.ID+"/arbitrate", bytes.NewReader([]byte(`{"resolution_verdict":"REFUND_BUYER","buyer_refund_share":1000000,"seller_release_share":0}`)))
	arbReq.Header.Set("Content-Type", "application/json")
	wArb := httptest.NewRecorder()
	r.ServeHTTP(wArb, arbReq)
	if wArb.Code != http.StatusOK {
		t.Fatalf("Expected arbitration 200, got %d", wArb.Code)
	}
}

func TestDoubleEntryLedgerAudit(t *testing.T) {
	r := setupTestRouter()

	reqAccts := httptest.NewRequest("GET", "/api/v1/ledger/accounts", nil)
	wAccts := httptest.NewRecorder()
	r.ServeHTTP(wAccts, reqAccts)
	if wAccts.Code != http.StatusOK {
		t.Fatalf("Expected ledger accounts 200, got %d", wAccts.Code)
	}

	reqJournals := httptest.NewRequest("GET", "/api/v1/ledger/journals", nil)
	wJournals := httptest.NewRecorder()
	r.ServeHTTP(wJournals, reqJournals)
	if wJournals.Code != http.StatusOK {
		t.Fatalf("Expected ledger journals 200, got %d", wJournals.Code)
	}
}

func TestSupportTicketLogging(t *testing.T) {
	r := setupTestRouter()

	ticketPayload := domain.SupportTicket{
		Name:    "Sanjay Kumar",
		Company: "Kumar Ceramics Ltd",
		Email:   "sanjay@kumarceramics.com",
		Phone:   "9876543210",
		Message: "Inquiry about high-volume API integration for export consignments.",
	}
	body, _ := json.Marshal(ticketPayload)
	req := httptest.NewRequest("POST", "/api/v1/support/tickets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Expected ticket 201, got %d", w.Code)
	}
}

func TestPennyDropSimulation(t *testing.T) {
	r := setupTestRouter()
	payload := map[string]string{
		"account_number": "50200045892134",
		"ifsc_code":       "HDFC0000240",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/v1/mock/penny-drop", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
}

func TestAdminEndpoints(t *testing.T) {
	r := setupTestRouter()

	// 1. Stats
	req := httptest.NewRequest("GET", "/api/v1/admin/stats", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected stats status 200, got %d", w.Code)
	}

	// 2. Buyers
	req = httptest.NewRequest("GET", "/api/v1/admin/buyers", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected buyers status 200, got %d", w.Code)
	}

	// 3. Suppliers
	req = httptest.NewRequest("GET", "/api/v1/admin/suppliers", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected suppliers status 200, got %d", w.Code)
	}

	// 4. Finance
	req = httptest.NewRequest("GET", "/api/v1/admin/finance", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected finance status 200, got %d", w.Code)
	}
}

