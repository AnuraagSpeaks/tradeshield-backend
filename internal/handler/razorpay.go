package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"tradeshield-backend/internal/pkg/razorpay"
	"tradeshield-backend/internal/pkg/response"
)

type CreateOrderReq struct {
	Amount   int    `json:"amount"`   // in paise (e.g. 100 = ₹1.00)
	Receipt  string `json:"receipt,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type CreateOrderRes struct {
	OrderID  string `json:"order_id"`
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
}

// CreateOrder handles POST /api/v1/create-order and /api/create-order
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Validate amount >= 100 paise (Razorpay requirement)
	if req.Amount < 100 {
		response.Error(w, http.StatusBadRequest, "Amount must be at least 100 paise (₹1.00)")
		return
	}
	if req.Amount > 1000000000 { // ₹1 Crore
		response.Error(w, http.StatusBadRequest, "Amount exceeds maximum single checkout limit (₹1,00,00,000)")
		return
	}

	if req.Currency != "" && req.Currency != "INR" {
		response.Error(w, http.StatusBadRequest, "Only INR currency is currently supported")
		return
	}
	if req.Currency == "" {
		req.Currency = "INR"
	}

	orderReq := razorpay.OrderRequest{
		Amount:   req.Amount,
		Currency: req.Currency,
		Receipt:  req.Receipt,
	}

	orderResp, err := razorpay.CreateOrder(orderReq)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create Razorpay order: %v", err))
		return
	}

	res := CreateOrderRes{
		OrderID:  orderResp.ID,
		Amount:   orderResp.Amount,
		Currency: orderResp.Currency,
	}
	response.JSON(w, http.StatusOK, res, "Razorpay order created successfully")
}

type VerifyPaymentReq struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Signature string `json:"signature"`
}

type VerifyPaymentRes struct {
	Success   bool   `json:"success"`
	PaymentID string `json:"payment_id,omitempty"`
	OrderID   string `json:"order_id,omitempty"`
	Message   string `json:"message"`
}

// VerifyPayment handles POST /api/v1/verify-payment and /api/verify-payment
func (h *Handler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	var req VerifyPaymentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.OrderID == "" || req.PaymentID == "" || req.Signature == "" {
		response.Error(w, http.StatusBadRequest, "order_id, payment_id and signature are all required")
		return
	}

	isValid, err := razorpay.VerifySignature(req.OrderID, req.PaymentID, req.Signature)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, fmt.Sprintf("Signature verification error: %v", err))
		return
	}

	if !isValid {
		response.JSON(w, http.StatusBadRequest, VerifyPaymentRes{
			Success: false,
			Message: "Razorpay signature verification failed. Tampered or invalid signature.",
		}, "Verification failed")
		return
	}

	response.JSON(w, http.StatusOK, VerifyPaymentRes{
		Success:   true,
		PaymentID: req.PaymentID,
		OrderID:   req.OrderID,
		Message:   "Payment verified successfully",
	}, "Payment confirmed")
}
