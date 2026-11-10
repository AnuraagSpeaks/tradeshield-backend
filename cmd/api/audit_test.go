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

// TestAuditEdgeCases tests and verifies that security and logic loopholes are resolved.
func TestAuditEdgeCases(t *testing.T) {
	repository.DB = repository.NewStore()
	repository.DB.SeedData()
	r := setupTestRouter()
	h := handler.NewHandler(repository.DB)

	supplierToken := getTestToken("user_supplier_1", "sales@bharatcastings.com", "supplier", "org_sup_1")
	buyerToken := getTestToken("user_buyer_1", "purchase@apexauto.in", "buyer", "org_buy_1")
	adminToken := getTestToken("user_admin_1", "admin@payshieldx.in", "admin", "org_admin")

	// -------------------------------------------------------------
	// 1. Auth Edge Cases
	// -------------------------------------------------------------
	t.Run("1.1 Empty Password Bypasses Login for Existing User", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    "admin@payshieldx.in",
			"password": "",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected 400 rejection for empty password, got: %d", w.Code)
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
			t.Fatalf("Unregistered email was allowed to log in! Status: %d", w.Code)
		}
	})

	t.Run("1.3 OTP Bypass Code 849201 Verifies Any Email Without Generation", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":   "victim@example.com",
			"otp":     "849201",
			"purpose": "register",
		})
		req := httptest.NewRequest("POST", "/api/v1/auth/verify-otp", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Fatalf("Vulnerability: OTP bypass code 849201 succeeded! Status: %d", w.Code)
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
			t.Fatalf("Vulnerability: Admin account password hijacked using bypass code! Status: %d", w.Code)
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
		req.Header.Set("Authorization", supplierToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var res map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &res)
		propData, ok := res["data"].(map[string]interface{})
		if !ok || propData == nil {
			t.Fatalf("Failed to create proposal: %s", w.Body.String())
		}
		propID := propData["id"].(string)

		// First approval
		appReq1 := httptest.NewRequest("POST", "/api/v1/proposals/"+propID+"/approve", bytes.NewReader([]byte(`{"role":"buyer"}`)))
		appReq1.Header.Set("Content-Type", "application/json")
		appReq1.Header.Set("Authorization", buyerToken)
		w1 := httptest.NewRecorder()
		r.ServeHTTP(w1, appReq1)
		if w1.Code != http.StatusOK {
			t.Fatalf("Expected first approval 200, got: %d", w1.Code)
		}

		contractsCount1 := len(repository.DB.Contracts)

		// Second duplicate approval must be rejected with 400
		appReq2 := httptest.NewRequest("POST", "/api/v1/proposals/"+propID+"/approve", bytes.NewReader([]byte(`{"role":"buyer"}`)))
		appReq2.Header.Set("Content-Type", "application/json")
		appReq2.Header.Set("Authorization", buyerToken)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, appReq2)
		if w2.Code != http.StatusBadRequest {
			t.Fatalf("Expected duplicate approval to fail with 400, got: %d", w2.Code)
		}

		contractsCount2 := len(repository.DB.Contracts)
		if contractsCount2 > contractsCount1 {
			t.Fatalf("Duplicate approval created duplicate contract (%d -> %d)", contractsCount1, contractsCount2)
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
		appReq.Header.Set("Authorization", buyerToken)
		wApp := httptest.NewRecorder()
		r.ServeHTTP(wApp, appReq)

		if wApp.Code == http.StatusOK {
			t.Fatalf("Expired proposal approved without validation! Status: %d", wApp.Code)
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
			t.Fatalf("Unescaped HTML/XSS payload reflected directly in PDF HTML output!")
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
		subReq.Header.Set("Authorization", supplierToken)
		wSub := httptest.NewRecorder()
		r.ServeHTTP(wSub, subReq)

		if wSub.Code == http.StatusOK {
			t.Fatalf("Submitting proof on already released milestone should fail! Status: %d", wSub.Code)
		}
	})

	// -------------------------------------------------------------
	// 4. Admin Dashboard & Dispute Arbitration Edge Cases
	// -------------------------------------------------------------
	t.Run("4.4 Dispute Arbitration with Negative Shares", func(t *testing.T) {
		disp := domain.Dispute{
			ID:          "disp_audit_test",
			ClaimAmount: 100000,
			Status:      domain.DisputeStatusUnderReview,
		}
		repository.DB.Disputes[disp.ID] = disp

		arbBody, _ := json.Marshal(handler.ArbitrateReq{
			Verdict:            "REFUND_BUYER",
			BuyerRefundShare:   500000,
			SellerReleaseShare: -100000,
		})
		req := httptest.NewRequest("POST", "/api/v1/disputes/"+disp.ID+"/arbitrate", bytes.NewReader(arbBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", adminToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Fatalf("Arbitration allowed negative or over-claim shares! Status: %d", w.Code)
		}
	})
}
