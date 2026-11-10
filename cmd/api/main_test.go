package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"tradeshield-backend/internal/domain"
	"tradeshield-backend/internal/handler"
	"tradeshield-backend/internal/middleware"
	"tradeshield-backend/internal/pkg/token"
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

		// Public endpoints
		r.Post("/auth/login", h.Login)
		r.Post("/auth/register", h.Register)
		r.Post("/auth/send-otp", h.SendOTP)
		r.Post("/auth/verify-otp", h.VerifyOTP)
		r.Post("/auth/reset-password", h.ResetPassword)
		r.Get("/directory", h.ListDirectory)
		r.Get("/escrow/summary", h.GetSummary)
		r.Get("/proposals/{id}", h.GetProposal)
		r.Get("/proposals/{id}/pdf", h.GetProposalPDFHTML)
		r.Get("/suppliers", h.ListAdminSuppliers)
		r.Post("/mock/penny-drop", h.MockPennyDrop)
		r.Get("/mock/gst-verify/{gstin}", h.MockGSTVerify)
		r.Get("/documents/*", h.GetDocumentPreview)

		// Authenticated endpoints
		r.Group(func(auth chi.Router) {
			auth.Use(middleware.AuthMiddleware)

			auth.Get("/proposals", h.ListProposals)
			auth.Post("/proposals", h.CreateProposal)
			auth.Post("/proposals/{id}/approve", h.ApproveProposal)
			auth.Post("/proposals/{id}/request-modification", h.RequestModification)
			auth.Post("/proposals/{id}/ship", h.ShipProposal)
			auth.Post("/proposals/{id}/confirm-delivery", h.ConfirmDeliveryProposal)

			auth.Get("/contracts", h.ListContracts)
			auth.Post("/contracts", h.CreateContract)
			auth.Get("/contracts/{id}", h.GetContract)
			auth.Post("/milestones/{id}/submit-proof", h.SubmitMilestoneProof)
			auth.Post("/milestones/{id}/approve-release", h.ApproveAndReleaseMilestone)

			auth.Get("/connections", h.ListConnections)
			auth.Get("/connections/requests", h.ListConnectionRequests)
			auth.Post("/connections/requests", h.SendConnectionRequest)
			auth.Post("/connections/requests/{id}/accept", h.AcceptConnectionRequest)

			auth.Get("/disputes", h.ListDisputes)
			auth.Post("/disputes", h.RaiseDispute)

			auth.Get("/support/tickets", h.ListSupportTickets)
			auth.Post("/support/tickets", h.SubmitSupportTicket)

			// Admin endpoints
			auth.Group(func(admin chi.Router) {
				admin.Use(middleware.RequireRole("admin"))

				admin.Get("/admin/stats", h.GetAdminStats)
				admin.Get("/admin/buyers", h.ListAdminBuyers)
				admin.Get("/admin/suppliers", h.ListAdminSuppliers)
				admin.Get("/admin/finance", h.GetAdminFinance)
				admin.Post("/admin/kyc/{userId}/verify", h.VerifyAdminKYC)
				admin.Post("/admin/settlements/{id}/payout", h.ProcessAdminSettlement)

				admin.Get("/ledger/accounts", h.ListLedgerAccounts)
				admin.Get("/ledger/journals", h.ListLedgerJournals)

				admin.Post("/disputes/{id}/arbitrate", h.ArbitrateDispute)
			})
		})
	})
	return r
}

func getTestToken(id, email, role, org string) string {
	t, _ := token.GenerateToken(id, email, role, org)
	return "Bearer " + t
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

func TestSecurity_AuthenticationEnforced(t *testing.T) {
	r := setupTestRouter()

	// 1. Unauthenticated request to /proposals should fail with 401
	reqUnauth := httptest.NewRequest("GET", "/api/v1/proposals", nil)
	wUnauth := httptest.NewRecorder()
	r.ServeHTTP(wUnauth, reqUnauth)
	if wUnauth.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized for unauthenticated /proposals, got %d", wUnauth.Code)
	}

	// 2. Unauthenticated request to /admin/stats should fail with 401
	reqAdminUnauth := httptest.NewRequest("GET", "/api/v1/admin/stats", nil)
	wAdminUnauth := httptest.NewRecorder()
	r.ServeHTTP(wAdminUnauth, reqAdminUnauth)
	if wAdminUnauth.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized for unauthenticated /admin/stats, got %d", wAdminUnauth.Code)
	}

	// 3. Authenticated request as buyer to /admin/stats should fail with 403 Forbidden
	buyerToken := getTestToken("user_buyer_1", "buyer@example.com", "buyer", "org_1")
	reqBuyerAdmin := httptest.NewRequest("GET", "/api/v1/admin/stats", nil)
	reqBuyerAdmin.Header.Set("Authorization", buyerToken)
	wBuyerAdmin := httptest.NewRecorder()
	r.ServeHTTP(wBuyerAdmin, reqBuyerAdmin)
	if wBuyerAdmin.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 Forbidden for buyer accessing /admin/stats, got %d", wBuyerAdmin.Code)
	}

	// 4. Authenticated request as admin to /admin/stats should succeed with 200
	adminToken := getTestToken("user_admin_1", "admin@payshieldx.in", "admin", "org_admin")
	reqAdmin := httptest.NewRequest("GET", "/api/v1/admin/stats", nil)
	reqAdmin.Header.Set("Authorization", adminToken)
	wAdmin := httptest.NewRecorder()
	r.ServeHTTP(wAdmin, reqAdmin)
	if wAdmin.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for admin accessing /admin/stats, got %d", wAdmin.Code)
	}
}

func TestSecurity_AuthEdgeCasesAndHygine(t *testing.T) {
	r := setupTestRouter()

	// 1. Empty password login should be rejected with 400
	emptyPassBody := `{"email":"admin@payshieldx.in","password":"","role":"admin"}`
	reqEmptyPass := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(emptyPassBody))
	reqEmptyPass.Header.Set("Content-Type", "application/json")
	wEmptyPass := httptest.NewRecorder()
	r.ServeHTTP(wEmptyPass, reqEmptyPass)
	if wEmptyPass.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request for empty password login, got %d", wEmptyPass.Code)
	}

	// 2. Wrong password login should fail with 401
	wrongPassBody := `{"email":"admin@payshieldx.in","password":"wrongpassword","role":"admin"}`
	reqWrongPass := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(wrongPassBody))
	reqWrongPass.Header.Set("Content-Type", "application/json")
	wWrongPass := httptest.NewRecorder()
	r.ServeHTTP(wWrongPass, reqWrongPass)
	if wWrongPass.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized for wrong password, got %d", wWrongPass.Code)
	}

	// 3. Valid admin login succeeds and response JSON does NOT contain plaintext password
	validPassBody := `{"email":"admin@payshieldx.in","password":"Admin@Shield2026","role":"admin"}`
	reqValidPass := httptest.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(validPassBody))
	reqValidPass.Header.Set("Content-Type", "application/json")
	wValidPass := httptest.NewRecorder()
	r.ServeHTTP(wValidPass, reqValidPass)
	if wValidPass.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for valid admin login, got %d. Body: %s", wValidPass.Code, wValidPass.Body.String())
	}
	respStr := wValidPass.Body.String()
	if strings.Contains(respStr, `"password"`) || strings.Contains(respStr, `"password_hash"`) {
		t.Fatalf("Security loophole: Password field leaked in login JSON response! Body: %s", respStr)
	}

	// 4. Directory does not expose admin users or password fields
	reqDir := httptest.NewRequest("GET", "/api/v1/directory", nil)
	wDir := httptest.NewRecorder()
	r.ServeHTTP(wDir, reqDir)
	if wDir.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for directory, got %d", wDir.Code)
	}
	dirStr := wDir.Body.String()
	if strings.Contains(dirStr, `"password"`) || strings.Contains(dirStr, `"password_hash"`) {
		t.Fatalf("Security loophole: Password leaked in directory API! Body: %s", dirStr)
	}
	if strings.Contains(dirStr, "admin@payshieldx.in") {
		t.Fatalf("Security loophole: Admin user found in public directory! Body: %s", dirStr)
	}

	// 5. Backdoor OTP bypass codes (849201, 123456) must fail when no OTP was requested
	otpBypassBody := `{"email":"randomvictim@example.com","otp":"849201","purpose":"register"}`
	reqOTPBypass := httptest.NewRequest("POST", "/api/v1/auth/verify-otp", strings.NewReader(otpBypassBody))
	reqOTPBypass.Header.Set("Content-Type", "application/json")
	wOTPBypass := httptest.NewRecorder()
	r.ServeHTTP(wOTPBypass, reqOTPBypass)
	if wOTPBypass.Code != http.StatusBadRequest {
		t.Fatalf("Security loophole: Backdoor OTP bypass code was accepted! Expected 400, got %d", wOTPBypass.Code)
	}
}

func TestProposalDualApprovalLifecycle(t *testing.T) {
	r := setupTestRouter()
	supplierToken := getTestToken("user_supplier_1", "sales@bharatcastings.com", "supplier", "org_sup_1")
	buyerToken := getTestToken("user_buyer_1", "purchase@apexauto.in", "buyer", "org_buy_1")

	// 1. Create Proposal (Authenticated)
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
	req.Header.Set("Authorization", supplierToken)
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
	appSupReq.Header.Set("Authorization", supplierToken)
	wSup := httptest.NewRecorder()
	r.ServeHTTP(wSup, appSupReq)
	if wSup.Code != http.StatusOK {
		t.Fatalf("Expected supplier approve 200, got %d", wSup.Code)
	}

	// 3. Buyer Approves -> Locks deal into Active Escrow Contract
	appBuyReq := httptest.NewRequest("POST", "/api/v1/proposals/"+created.ID+"/approve", bytes.NewReader([]byte(`{"role":"buyer"}`)))
	appBuyReq.Header.Set("Content-Type", "application/json")
	appBuyReq.Header.Set("Authorization", buyerToken)
	wBuy := httptest.NewRecorder()
	r.ServeHTTP(wBuy, appBuyReq)
	if wBuy.Code != http.StatusOK {
		t.Fatalf("Expected buyer approve 200, got %d", wBuy.Code)
	}

	// 4. Supplier Ships with LR
	shipReq := httptest.NewRequest("POST", "/api/v1/proposals/"+created.ID+"/ship", bytes.NewReader([]byte(`{"lr_number":"LR-9988","transporter_name":"V-Trans","proof_url":"https://proof.pdf"}`)))
	shipReq.Header.Set("Content-Type", "application/json")
	shipReq.Header.Set("Authorization", supplierToken)
	wShip := httptest.NewRecorder()
	r.ServeHTTP(wShip, shipReq)
	if wShip.Code != http.StatusOK {
		t.Fatalf("Expected ship 200, got %d", wShip.Code)
	}

	// 5. Buyer Confirms Delivery -> Releases Payout
	delivReq := httptest.NewRequest("POST", "/api/v1/proposals/"+created.ID+"/confirm-delivery", nil)
	delivReq.Header.Set("Authorization", buyerToken)
	wDeliv := httptest.NewRecorder()
	r.ServeHTTP(wDeliv, delivReq)
	if wDeliv.Code != http.StatusOK {
		t.Fatalf("Expected delivery confirm 200, got %d", wDeliv.Code)
	}
}

func TestDisputeAndArbitration(t *testing.T) {
	r := setupTestRouter()
	buyerToken := getTestToken("user_buyer_1", "buyer@example.com", "buyer", "org_1")
	adminToken := getTestToken("user_admin_1", "admin@payshieldx.in", "admin", "org_admin")

	// 1. Raise Dispute (Buyer)
	dispPayload := handler.RaiseDisputeReq{
		ContractID:  "cntr_2026_089",
		MilestoneID: "ms_089_3",
		Reason:      "Dimensional variation in valve housings",
		ClaimAmount: 1000000,
	}
	body, _ := json.Marshal(dispPayload)
	req := httptest.NewRequest("POST", "/api/v1/disputes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", buyerToken)
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

	// 2. Buyer tries to arbitrate dispute -> should be rejected with 403
	arbBuyerReq := httptest.NewRequest("POST", "/api/v1/disputes/"+disp.ID+"/arbitrate", bytes.NewReader([]byte(`{"resolution_verdict":"REFUND_BUYER","buyer_refund_share":1000000,"seller_release_share":0}`)))
	arbBuyerReq.Header.Set("Content-Type", "application/json")
	arbBuyerReq.Header.Set("Authorization", buyerToken)
	wArbBuyer := httptest.NewRecorder()
	r.ServeHTTP(wArbBuyer, arbBuyerReq)
	if wArbBuyer.Code != http.StatusForbidden {
		t.Fatalf("Expected 403 Forbidden for non-admin dispute arbitration, got %d", wArbBuyer.Code)
	}

	// 3. Admin arbitrates dispute with share > claim amount -> should be rejected with 400
	arbOverReq := httptest.NewRequest("POST", "/api/v1/disputes/"+disp.ID+"/arbitrate", bytes.NewReader([]byte(`{"resolution_verdict":"REFUND_BUYER","buyer_refund_share":800000,"seller_release_share":300000}`)))
	arbOverReq.Header.Set("Content-Type", "application/json")
	arbOverReq.Header.Set("Authorization", adminToken)
	wArbOver := httptest.NewRecorder()
	r.ServeHTTP(wArbOver, arbOverReq)
	if wArbOver.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request for arbitrated sum exceeding claim amount, got %d", wArbOver.Code)
	}

	// 4. Admin arbitrates dispute validly -> 200 OK
	arbReq := httptest.NewRequest("POST", "/api/v1/disputes/"+disp.ID+"/arbitrate", bytes.NewReader([]byte(`{"resolution_verdict":"REFUND_BUYER","buyer_refund_share":1000000,"seller_release_share":0}`)))
	arbReq.Header.Set("Content-Type", "application/json")
	arbReq.Header.Set("Authorization", adminToken)
	wArb := httptest.NewRecorder()
	r.ServeHTTP(wArb, arbReq)
	if wArb.Code != http.StatusOK {
		t.Fatalf("Expected arbitration 200, got %d", wArb.Code)
	}
}

func TestDoubleEntryLedgerAudit(t *testing.T) {
	r := setupTestRouter()
	adminToken := getTestToken("user_admin_1", "admin@payshieldx.in", "admin", "org_admin")

	initialVault := repository.DB.LedgerAccts["acc_vault"].Balance

	// Record a double entry deposit of 50,000 into vault
	err := repository.DB.RecordDoubleEntry("c_test", "m_test", "TEST_DEPOSIT", "Test Deposit", "ESCROW_VAULT", "BUYER_WALLET", 50000)
	if err != nil {
		t.Fatalf("Failed to record double entry: %v", err)
	}

	newVault := repository.DB.LedgerAccts["acc_vault"].Balance
	if newVault != initialVault+50000 {
		t.Fatalf("Expected vault balance to increase by 50000 from %.2f to %.2f, got %.2f", initialVault, initialVault+50000, newVault)
	}

	reqAccts := httptest.NewRequest("GET", "/api/v1/ledger/accounts", nil)
	reqAccts.Header.Set("Authorization", adminToken)
	wAccts := httptest.NewRecorder()
	r.ServeHTTP(wAccts, reqAccts)
	if wAccts.Code != http.StatusOK {
		t.Fatalf("Expected ledger accounts 200, got %d", wAccts.Code)
	}

	reqJournals := httptest.NewRequest("GET", "/api/v1/ledger/journals", nil)
	reqJournals.Header.Set("Authorization", adminToken)
	wJournals := httptest.NewRecorder()
	r.ServeHTTP(wJournals, reqJournals)
	if wJournals.Code != http.StatusOK {
		t.Fatalf("Expected ledger journals 200, got %d", wJournals.Code)
	}
}

func TestSupportTicketLogging(t *testing.T) {
	r := setupTestRouter()
	buyerToken := getTestToken("user_buyer_1", "buyer@example.com", "buyer", "org_1")

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
	req.Header.Set("Authorization", buyerToken)
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
	adminToken := getTestToken("user_admin_1", "admin@payshieldx.in", "admin", "org_admin")

	// 1. Stats
	req := httptest.NewRequest("GET", "/api/v1/admin/stats", nil)
	req.Header.Set("Authorization", adminToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected stats status 200, got %d", w.Code)
	}

	// 2. Buyers
	req = httptest.NewRequest("GET", "/api/v1/admin/buyers", nil)
	req.Header.Set("Authorization", adminToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected buyers status 200, got %d", w.Code)
	}

	// 3. Suppliers
	req = httptest.NewRequest("GET", "/api/v1/admin/suppliers", nil)
	req.Header.Set("Authorization", adminToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected suppliers status 200, got %d", w.Code)
	}

	// 4. Finance
	req = httptest.NewRequest("GET", "/api/v1/admin/finance", nil)
	req.Header.Set("Authorization", adminToken)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("Expected finance status 200, got %d", w.Code)
	}
}

func TestDocumentVaultPreviews(t *testing.T) {
	r := setupTestRouter()

	testPaths := []string{
		"/api/v1/documents/kyc/gst_bharat.pdf",
		"/api/v1/documents/kyc/pan_bharat.pdf",
		"/api/v1/documents/kyc/cheque_bharat.pdf",
		"/api/v1/documents/kyc/msme_bharat.pdf",
		"/api/v1/documents/proofs/eway_bill_lr_089.pdf",
		"/api/v1/documents/proofs/raw_material_cert_089.pdf",
	}

	for _, p := range testPaths {
		req := httptest.NewRequest("GET", p, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Path %s returned status %d, expected 200", p, w.Code)
		}
		contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("Path %s returned content-type %s, expected text/html", p, contentType)
		}
	}
}
