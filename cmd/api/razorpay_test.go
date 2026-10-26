package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"tradeshield-backend/internal/handler"
	"tradeshield-backend/internal/repository"
)

func TestRazorpayOrderAndVerification(t *testing.T) {
	os.Setenv("RAZORPAY_KEY_ID", "rzp_test_Tsa2ayJfSJsAw3")
	os.Setenv("RAZORPAY_KEY_SECRET", "vpK9sEBnK3I588zP0PRfIj3T")

	r := chi.NewRouter()
	h := handler.NewHandler(repository.DB)
	r.Post("/api/v1/create-order", h.CreateOrder)
	r.Post("/api/v1/verify-payment", h.VerifyPayment)
	r.Post("/api/v1/gst/lookup", h.GSTLookup)

	// 1. Test Create Order
	orderReqBody := []byte(`{"amount": 50000, "currency": "INR", "receipt": "test_rcpt_1"}`)
	req, _ := http.NewRequest("POST", "/api/v1/create-order", bytes.NewBuffer(orderReqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK from create-order, got %d: %s", w.Code, w.Body.String())
	}

	var res struct {
		Data struct {
			OrderID  string `json:"order_id"`
			Amount   int    `json:"amount"`
			Currency string `json:"currency"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to parse create-order response: %v", err)
	}
	if res.Data.OrderID == "" {
		t.Fatalf("Expected valid order_id, got empty")
	}
	if res.Data.Amount != 50000 {
		t.Fatalf("Expected amount 50000, got %d", res.Data.Amount)
	}

	// 2. Test Payment Verification with genuine HMAC signature
	testOrderID := res.Data.OrderID
	testPaymentID := "pay_test_987654321"
	mac := hmac.New(sha256.New, []byte("vpK9sEBnK3I588zP0PRfIj3T"))
	mac.Write([]byte(testOrderID + "|" + testPaymentID))
	validSig := hex.EncodeToString(mac.Sum(nil))

	verifyReqBody, _ := json.Marshal(map[string]string{
		"order_id":   testOrderID,
		"payment_id": testPaymentID,
		"signature":  validSig,
	})
	vReq, _ := http.NewRequest("POST", "/api/v1/verify-payment", bytes.NewBuffer(verifyReqBody))
	vReq.Header.Set("Content-Type", "application/json")
	vW := httptest.NewRecorder()
	r.ServeHTTP(vW, vReq)

	if vW.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for valid signature, got %d: %s", vW.Code, vW.Body.String())
	}

	// 3. Test Invalid Signature Rejection
	badVerifyReqBody, _ := json.Marshal(map[string]string{
		"order_id":   testOrderID,
		"payment_id": testPaymentID,
		"signature":  "invalid_tampered_signature",
	})
	badReq, _ := http.NewRequest("POST", "/api/v1/verify-payment", bytes.NewBuffer(badVerifyReqBody))
	badReq.Header.Set("Content-Type", "application/json")
	badW := httptest.NewRecorder()
	r.ServeHTTP(badW, badReq)

	if badW.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request for tampered signature, got %d", badW.Code)
	}

	// 4. Test GST Lookup
	gstReqBody := []byte(`{"gstin": "27AAACA1234A1Z5"}`)
	gstReq, _ := http.NewRequest("POST", "/api/v1/gst/lookup", bytes.NewBuffer(gstReqBody))
	gstReq.Header.Set("Content-Type", "application/json")
	gstW := httptest.NewRecorder()
	r.ServeHTTP(gstW, gstReq)

	if gstW.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for GST lookup, got %d: %s", gstW.Code, gstW.Body.String())
	}
}
