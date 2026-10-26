package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"tradeshield-backend/internal/domain"
	"tradeshield-backend/internal/handler"
	"tradeshield-backend/internal/repository"
)

// TestAuditEdgeCases tests the edge cases identified during the deep security and logic audit.
func TestAuditEdgeCases(t *testing.T) {
	repository.DB = repository.NewStore()
	repository.DB.SeedData()
	r := setupTestRouter()
	h := handler.NewHandler(repository.DB)

	// -------------------------------------------------------------
	// 1. Auth Edge Cases
	// -------------------------------------------------------------
	t.Run("1.1 Empty Password Bypasses Login for Existing User", func(t *testing.T) {
		// Existing seeded user: admin@payshieldx.in with password "password123"
		body, _ := json.Marshal(map[string]string{
			"email":    "admin@payshieldx.in",
			"password": "", // empty password
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Logf("CRITICAL VULNERABILITY CONFIRMED: Empty password allowed login to admin account! Status: %d", w.Code)
		} else {
			t.Logf("Expected rejection, got: %d", w.Code)
		}
	})

	t.Run("1.2 Unregistered Email Dynamically Logs In", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    "completely_random_unregistered_hacker@evil.com",
			"password": "wrongpassword",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Logf("VULNERABILITY CONFIRMED: Unregistered email automatically authenticated! Status: %d", w.Code)
		}
	})

	t.Run("1.3 OTP Bypass Code 849201 Verifies Any Email Without Generation", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":   "victim@corporate.in",
			"otp":     "849201",
			"purpose": "register",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/verify-otp", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.VerifyOTP(w, req)

		if w.Code == http.StatusOK {
			t.Logf("CRITICAL VULNERABILITY CONFIRMED: Bypass code 849201 verified victim email without OTP request! Status: %d", w.Code)
		}
	})

	t.Run("1.4 Account Hijack via ResetPassword with Bypass Code", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":        "admin@payshieldx.in",
			"otp":          "849201",
			"new_password": "HackedPassword123!",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/reset-password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.ResetPassword(w, req)

		if w.Code == http.StatusOK {
			t.Logf("CRITICAL VULNERABILITY CONFIRMED: Admin account password hijacked using bypass code! Status: %d", w.Code)
		}
	})

	// -------------------------------------------------------------
	// 2. Proposal & Dual Approval Edge Cases
	// -------------------------------------------------------------
	t.Run("2.1 Duplicate Approval Re-creates Contracts & Duplicates Ledger Postings", func(t *testing.T) {
		// Create proposal
		propPayload := handler.CreateProposalReq{
			BuyerID:    "user_buyer_1",
			SupplierID: "user_supplier_1",
			Amount:     100000,
		}
		body, _ := json.Marshal(propPayload)
		req := httptest.NewRequest("POST", "/api/v1/proposals", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var res map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &res)
		propData := res["data"].(map[string]interface{})
		propID := propData["id"].(string)

		// First approval
		appReq1 := httptest.NewRequest("POST", "/api/v1/proposals/"+propID+"/approve", bytes.NewReader([]byte(`{"role":"buyer"}`)))
		appReq1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		r.ServeHTTP(w1, appReq1)
		contractsCount1 := len(repository.DB.Contracts)
		journalsCount1 := len(repository.DB.Journals)

		// Second duplicate approval
		appReq2 := httptest.NewRequest("POST", "/api/v1/proposals/"+propID+"/approve", bytes.NewReader([]byte(`{"role":"buyer"}`)))
		appReq2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, appReq2)
		contractsCount2 := len(repository.DB.Contracts)
		journalsCount2 := len(repository.DB.Journals)

		if contractsCount2 > contractsCount1 {
			t.Logf("BUG CONFIRMED: Duplicate approval created a duplicate contract (%d -> %d) and duplicate ledger entries (%d -> %d)",
				contractsCount1, contractsCount2, journalsCount1, journalsCount2)
		}
	})

	t.Run("2.2 Modification Request on Approved Proposal Overwrites State", func(t *testing.T) {
		propPayload := handler.CreateProposalReq{
			BuyerID:    "user_buyer_1",
			SupplierID: "user_supplier_1",
			Amount:     50000,
		}
		body, _ := json.Marshal(propPayload)
		req := httptest.NewRequest("POST", "/api/v1/proposals", bytes.NewReader(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		var res map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &res)
		propID := res["data"].(map[string]interface{})["id"].(string)

		// Approve
		appReq := httptest.NewRequest("POST", "/api/v1/proposals/"+propID+"/approve", bytes.NewReader([]byte(`{"role":"buyer"}`)))
		wApp := httptest.NewRecorder()
		r.ServeHTTP(wApp, appReq)

		// Request modification on approved proposal
		modReq := httptest.NewRequest("POST", "/api/v1/proposals/"+propID+"/request-modification", nil)
		wMod := httptest.NewRecorder()
		r.ServeHTTP(wMod, modReq)

		p := repository.DB.Proposals[propID]
		if p.Status == "Modification Requested" && p.ContractID != "" {
			t.Logf("BUG CONFIRMED: Approved funded proposal was reverted to 'Modification Requested' leaving orphaned contract: %s", p.ContractID)
		}
	})

	t.Run("2.3 Approving Expired Proposal", func(t *testing.T) {
		pastTime := time.Now().Add(-48 * time.Hour)
		p := domain.Proposal{
			ID:             "prop_expired_01",
			ProposalNumber: "999999",
			Amount:         100000,
			Status:         "PENDING_BUYER_APPROVAL",
			ValidTillDate:  &pastTime,
			CreatedAt:      time.Now().Add(-72 * time.Hour),
		}
		repository.DB.Proposals[p.ID] = p

		appReq := httptest.NewRequest("POST", "/api/v1/proposals/"+p.ID+"/approve", bytes.NewReader([]byte(`{"role":"buyer"}`)))
		wApp := httptest.NewRecorder()
		r.ServeHTTP(wApp, appReq)

		if wApp.Code == http.StatusOK {
			t.Logf("BUG CONFIRMED: Expired proposal approved successfully without validity date validation! Status: %d", wApp.Code)
		}
	})

	t.Run("2.4 XSS in Proposal PDF Generation", func(t *testing.T) {
		xssPayload := "<script>alert('pwned')</script>"
		p := domain.Proposal{
			ID:              "prop_xss",
			ProposalNumber:  "888888",
			BuyerSignatory:  xssPayload,
			ItemDescription: "<img src=x onerror=alert(1)>",
			Amount:          10000,
			CreatedAt:       time.Now(),
		}
		repository.DB.Proposals[p.ID] = p

		pdfReq := httptest.NewRequest("GET", "/api/v1/proposals/"+p.ID+"/pdf", nil)
		wPdf := httptest.NewRecorder()
		r.ServeHTTP(wPdf, pdfReq)

		htmlBody := wPdf.Body.String()
		if strings.Contains(htmlBody, "<script>alert('pwned')</script>") {
			t.Logf("SECURITY VULNERABILITY CONFIRMED: Unescaped HTML/XSS payload reflected directly in PDF HTML output!")
		}
	})

	// -------------------------------------------------------------
	// 3. Milestone & Vault Edge Cases
	// -------------------------------------------------------------
	t.Run("3.1 Milestone Proof Submission on Released Milestone", func(t *testing.T) {
		m := domain.Milestone{
			ID:         "ms_test_released",
			ContractID: "cntr_test",
			Status:     domain.MilestoneStatusReleased,
			Amount:     50000,
		}
		repository.DB.Milestones[m.ID] = m

		subReq := httptest.NewRequest("POST", "/api/v1/milestones/"+m.ID+"/submit-proof", bytes.NewReader([]byte(`{"deliverable_proof_url":"https://evil.com/fake.pdf"}`)))
		subReq.Header.Set("Content-Type", "application/json")
		wSub := httptest.NewRecorder()
		r.ServeHTTP(wSub, subReq)

		updatedM := repository.DB.Milestones[m.ID]
		if updatedM.Status == domain.MilestoneStatusInInspection {
			t.Logf("BUG CONFIRMED: Released milestone reverted to IN_INSPECTION upon proof submission!")
		}
	})

	t.Run("3.2 Releasing Milestone Without Proof", func(t *testing.T) {
		m := domain.Milestone{
			ID:                  "ms_no_proof",
			ContractID:          "cntr_test",
			Status:              domain.MilestoneStatusFunded,
			DeliverableProofURL: "",
			Amount:              50000,
		}
		repository.DB.Milestones[m.ID] = m

		appReq := httptest.NewRequest("POST", "/api/v1/milestones/"+m.ID+"/approve-release", nil)
		wApp := httptest.NewRecorder()
		r.ServeHTTP(wApp, appReq)

		if wApp.Code == http.StatusOK {
			t.Logf("LOGIC FLAW CONFIRMED: Milestone released without any deliverable proof uploaded! Status: %d", wApp.Code)
		}
	})

	// -------------------------------------------------------------
	// 4. Admin Dashboard Edge Cases
	// -------------------------------------------------------------
	t.Run("4.1 Verify KYC for Non-Existent User ID Returns 200 OK", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/admin/kyc/non_existent_user_9999/verify", bytes.NewReader([]byte(`{"status":"Approved"}`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Logf("BUG CONFIRMED: Non-existent user ID returns 200 OK with 'updated: true': %s", w.Body.String())
		}
	})

	t.Run("4.2 Settlement Payout for Non-Existent Settlement ID Returns 200 OK", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/v1/admin/settlements/SETTL_DOES_NOT_EXIST/payout", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Logf("BUG CONFIRMED: Non-existent settlement returns 200 OK with fabricated UTR: %s", w.Body.String())
		}
	})

	t.Run("4.3 Concurrency Deadlock in ProcessAdminSettlement", func(t *testing.T) {
		isolatedStore := repository.NewStore()
		isolatedStore.SeedData()
		isolatedHandler := handler.NewHandler(isolatedStore)
		settlementID := "SETTL-2026-901"

		done := make(chan bool, 1)
		go func() {
			req := httptest.NewRequest("POST", "/api/v1/admin/settlements/"+settlementID+"/payout", nil)
			w := httptest.NewRecorder()
			// Manually simulate what handler does:
			isolatedHandler.ProcessAdminSettlement(w, req)
			done <- true
		}()

		select {
		case <-done:
			t.Logf("Settlement processed")
		case <-time.After(150 * time.Millisecond):
			t.Logf("CRITICAL CONCURRENCY DEADLOCK CONFIRMED: ProcessAdminSettlement deadlocks trying to acquire s.mu.Lock() inside h.store.Lock()!")
		}
	})

	t.Run("4.4 Dispute Arbitration with Mismatched / Negative Shares", func(t *testing.T) {
		disp := domain.Dispute{
			ID:          "disp_audit_test",
			ClaimAmount: 100000,
			Status:      domain.DisputeStatusUnderReview,
		}
		repository.DB.Disputes[disp.ID] = disp

		// Claims total 100k, but arbitration awards 500k to buyer and -100k to seller
		arbBody, _ := json.Marshal(handler.ArbitrateReq{
			Verdict:            "REFUND_BUYER",
			BuyerRefundShare:   500000,
			SellerReleaseShare: -100000,
		})
		req := httptest.NewRequest("POST", "/api/v1/disputes/"+disp.ID+"/arbitrate", bytes.NewReader(arbBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Logf("BUG CONFIRMED: Arbitration allowed invalid shares (Buyer: 500k, Seller: -100k on 100k claim)!")
		}
	})
}
