package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"tradeshield-backend/internal/domain"
	"tradeshield-backend/internal/pkg/response"
	"tradeshield-backend/internal/pkg/token"
	"tradeshield-backend/internal/repository"
)

type Handler struct {
	store *repository.Store
}

func NewHandler(s *repository.Store) *Handler {
	return &Handler{store: s}
}

// -------------------------------------------------------------
// 1. Auth & Users
// -------------------------------------------------------------

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	h.store.RLock()
	defer h.store.RUnlock()

	reqEmail := strings.ToLower(strings.TrimSpace(req.Email))

	for _, u := range h.store.Users {
		if strings.ToLower(u.Email) == reqEmail {
			if req.Password != "" && u.Password != "" && u.Password != req.Password {
				response.Error(w, http.StatusUnauthorized, "Invalid password credentials")
				return
			}
			t, _ := token.GenerateToken(u.ID, u.Email, u.Role, u.OrganizationID)
			response.JSON(w, http.StatusOK, map[string]interface{}{
				"token": t,
				"user":  u,
			}, "Login successful")
			return
		}
	}

	// Demo fallback
	u := h.store.Users["user_buyer_1"]
	t, _ := token.GenerateToken(u.ID, u.Email, u.Role, u.OrganizationID)
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"token": t,
		"user":  u,
	}, "Login successful (Pre-seeded Demo Mode)")
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req domain.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user payload")
		return
	}

	if req.Email == "" || req.BusinessName == "" {
		response.Error(w, http.StatusBadRequest, "Email and Business Name are required")
		return
	}

	h.store.Lock()
	defer h.store.Unlock()

	userID := "user_" + uuid.New().String()[:8]
	req.ID = userID
	req.Verified = true
	req.IsVerified = true
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	h.store.Users[req.ID] = req
	if repository.PG != nil {
		repository.PG.SaveUser(req)
	}
	t, _ := token.GenerateToken(req.ID, req.Email, req.Role, req.OrganizationID)

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"token": t,
		"user":  req,
	}, "Account registered successfully")
}

func (h *Handler) ListDirectory(w http.ResponseWriter, r *http.Request) {
	roleQuery := r.URL.Query().Get("role")
	searchQuery := strings.ToLower(r.URL.Query().Get("q"))

	h.store.RLock()
	defer h.store.RUnlock()

	results := make([]domain.User, 0)
	for _, u := range h.store.Users {
		if roleQuery != "" && u.Role != roleQuery {
			continue
		}
		if searchQuery != "" {
			nameMatch := strings.Contains(strings.ToLower(u.BusinessName), searchQuery)
			gstMatch := strings.Contains(strings.ToLower(u.GST), searchQuery)
			cityMatch := strings.Contains(strings.ToLower(u.City), searchQuery)
			if !nameMatch && !gstMatch && !cityMatch {
				continue
			}
		}
		results = append(results, u)
	}

	response.JSON(w, http.StatusOK, results, "Business directory results")
}

// -------------------------------------------------------------
// 2. Escrow Summary & Metrics
// -------------------------------------------------------------

func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	summary := h.store.GetEscrowSummary()
	response.JSON(w, http.StatusOK, summary, "Escrow vault summary retrieved")
}

// -------------------------------------------------------------
// 3. Proposals & Dual-Approval Engine
// -------------------------------------------------------------

func (h *Handler) ListProposals(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	defer h.store.RUnlock()

	list := make([]domain.Proposal, 0, len(h.store.Proposals))
	for _, p := range h.store.Proposals {
		list = append(list, p)
	}
	response.JSON(w, http.StatusOK, list, "Proposals retrieved successfully")
}

func (h *Handler) GetProposal(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	h.store.RLock()
	defer h.store.RUnlock()

	p, ok := h.store.Proposals[id]
	if !ok {
		response.Error(w, http.StatusNotFound, "Proposal not found")
		return
	}
	response.JSON(w, http.StatusOK, p, "Proposal details retrieved")
}

type CreateProposalReq struct {
	BuyerID          string  `json:"buyer_id"`
	SupplierID       string  `json:"supplier_id"`
	Amount           float64 `json:"amount"`
	Currency         string  `json:"currency"`
	Terms            string  `json:"terms"`
	DeliveryTimeline string  `json:"delivery_timeline"`
	Notes            string  `json:"notes"`
}

func (h *Handler) CreateProposal(w http.ResponseWriter, r *http.Request) {
	var req CreateProposalReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.Amount <= 0 {
		response.Error(w, http.StatusBadRequest, "Amount must be greater than zero")
		return
	}

	propID := "prop_" + uuid.New().String()[:8]
	orderID := fmt.Sprintf("ORD-%d", 10000+time.Now().Unix()%90000)

	h.store.Lock()
	defer h.store.Unlock()

	buyerName := "Verified Buyer"
	if b, ok := h.store.Users[req.BuyerID]; ok {
		buyerName = b.BusinessName
	}
	supplierName := "Verified Supplier"
	if s, ok := h.store.Users[req.SupplierID]; ok {
		supplierName = s.BusinessName
	}

	curr := req.Currency
	if curr == "" {
		curr = "INR"
	}

	p := domain.Proposal{
		ID:               propID,
		OrderID:          orderID,
		BuyerID:          req.BuyerID,
		BuyerName:        buyerName,
		SupplierID:       req.SupplierID,
		SupplierName:     supplierName,
		Amount:           req.Amount,
		Currency:         curr,
		Terms:            req.Terms,
		DeliveryTimeline: req.DeliveryTimeline,
		Notes:            req.Notes,
		Status:           "Draft",
		PaymentStatus:    "Unpaid",
		BuyerApproved:    false,
		SupplierApproved: false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	h.store.Proposals[p.ID] = p
	if repository.PG != nil {
		repository.PG.SaveProposal(p)
	}
	response.JSON(w, http.StatusCreated, p, "Payment Confirmation Proposal created successfully")
}

type ApproveProposalReq struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"` // buyer or supplier
}

func (h *Handler) ApproveProposal(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ApproveProposalReq
	json.NewDecoder(r.Body).Decode(&req)

	h.store.Lock()
	defer h.store.Unlock()

	p, ok := h.store.Proposals[id]
	if !ok {
		response.Error(w, http.StatusNotFound, "Proposal not found")
		return
	}

	if req.Role == "buyer" {
		p.BuyerApproved = true
	} else if req.Role == "supplier" {
		p.SupplierApproved = true
	}

	if p.BuyerApproved && p.SupplierApproved {
		p.Status = "Confirmed"
		p.PaymentStatus = "Payment Held in Escrow"
	} else if p.SupplierApproved {
		p.Status = "Approved by Supplier"
	} else if p.BuyerApproved {
		p.Status = "Approved by Buyer"
	}

	p.UpdatedAt = time.Now()
	h.store.Proposals[id] = p
	if repository.PG != nil {
		repository.PG.SaveProposal(p)
	}

	response.JSON(w, http.StatusOK, p, "Proposal approved successfully")
}

func (h *Handler) RequestModification(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	h.store.Lock()
	defer h.store.Unlock()

	p, ok := h.store.Proposals[id]
	if !ok {
		response.Error(w, http.StatusNotFound, "Proposal not found")
		return
	}

	p.Status = "Modification Requested"
	p.BuyerApproved = false
	p.SupplierApproved = false
	p.UpdatedAt = time.Now()
	h.store.Proposals[id] = p
	if repository.PG != nil {
		repository.PG.SaveProposal(p)
	}

	response.JSON(w, http.StatusOK, p, "Modification request recorded")
}

type ShipProposalReq struct {
	LRNumber        string `json:"lr_number"`
	TransporterName string `json:"transporter_name"`
	ProofURL        string `json:"proof_url"`
}

func (h *Handler) ShipProposal(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ShipProposalReq
	json.NewDecoder(r.Body).Decode(&req)

	h.store.Lock()
	defer h.store.Unlock()

	p, ok := h.store.Proposals[id]
	if !ok {
		response.Error(w, http.StatusNotFound, "Proposal not found")
		return
	}

	p.PaymentStatus = "Shipped"
	p.LRNumber = req.LRNumber
	p.TransporterName = req.TransporterName
	p.ProofURL = req.ProofURL
	p.UpdatedAt = time.Now()
	h.store.Proposals[id] = p
	if repository.PG != nil {
		repository.PG.SaveProposal(p)
	}

	response.JSON(w, http.StatusOK, p, "Shipment details and LR submitted")
}

func (h *Handler) ConfirmDeliveryProposal(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	h.store.Lock()
	p, ok := h.store.Proposals[id]
	if !ok {
		h.store.Unlock()
		response.Error(w, http.StatusNotFound, "Proposal not found")
		return
	}

	p.PaymentStatus = "Payment Released"
	p.UpdatedAt = time.Now()
	h.store.Proposals[id] = p
	if repository.PG != nil {
		repository.PG.SaveProposal(p)
	}
	h.store.Unlock()

	h.store.RecordDoubleEntry("", "", "PROPOSAL_PAYOUT", fmt.Sprintf("Released payout for proposal %s (%s)", p.ID, p.OrderID), "ESCROW_VAULT", "SELLER_PAYOUT", p.Amount)

	response.JSON(w, http.StatusOK, p, "Delivery confirmed! Escrow funds released to supplier.")
}

// -------------------------------------------------------------
// 4. Milestone Escrow Contracts (High-Volume Multi-Stage)
// -------------------------------------------------------------

func (h *Handler) ListContracts(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	defer h.store.RUnlock()

	list := make([]domain.Contract, 0, len(h.store.Contracts))
	for _, c := range h.store.Contracts {
		msList := make([]domain.Milestone, 0)
		for _, m := range h.store.Milestones {
			if m.ContractID == c.ID {
				msList = append(msList, m)
			}
		}
		c.Milestones = msList
		list = append(list, c)
	}
	response.JSON(w, http.StatusOK, list, "Contracts retrieved successfully")
}

func (h *Handler) GetContract(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	h.store.RLock()
	defer h.store.RUnlock()

	c, ok := h.store.Contracts[id]
	if !ok {
		response.Error(w, http.StatusNotFound, "Contract not found")
		return
	}

	msList := make([]domain.Milestone, 0)
	for _, m := range h.store.Milestones {
		if m.ContractID == c.ID {
			msList = append(msList, m)
		}
	}
	c.Milestones = msList
	response.JSON(w, http.StatusOK, c, "Contract details retrieved")
}

type CreateContractReq struct {
	Title                string             `json:"title"`
	BuyerOrgID           string             `json:"buyer_org_id"`
	SupplierOrgID        string             `json:"supplier_org_id"`
	TotalAmount          float64            `json:"total_amount"`
	Description          string             `json:"description"`
	DeliveryTerms        string             `json:"delivery_terms"`
	InspectionPeriodDays int                `json:"inspection_period_days"`
	Milestones           []domain.Milestone `json:"milestones"`
}

func (h *Handler) CreateContract(w http.ResponseWriter, r *http.Request) {
	var req CreateContractReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.TotalAmount <= 0 {
		response.Error(w, http.StatusBadRequest, "Total amount must be greater than zero")
		return
	}

	contractID := "cntr_" + uuid.New().String()[:8]
	cNum := fmt.Sprintf("TS-CTR-%d-%s", time.Now().Year(), uuid.New().String()[:4])
	van := fmt.Sprintf("ICIC0000104TS%s", uuid.New().String()[:4])

	c := domain.Contract{
		ID:                   contractID,
		Title:                req.Title,
		ContractNumber:       cNum,
		BuyerOrgID:           req.BuyerOrgID,
		BuyerOrgName:         "Acme Industrial Ltd",
		SupplierOrgID:        req.SupplierOrgID,
		SupplierOrgName:      "Rajesh Exports",
		TotalAmount:          req.TotalAmount,
		Currency:             "INR",
		PlatformFeePercent:   0.75,
		PlatformFeeAmount:    req.TotalAmount * 0.0075,
		EscrowVirtualAccount: van,
		Status:               domain.ContractStatusFunded,
		Description:          req.Description,
		DeliveryTerms:        req.DeliveryTerms,
		InspectionPeriodDays: req.InspectionPeriodDays,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	h.store.Lock()
	h.store.Contracts[c.ID] = c

	createdMilestones := make([]domain.Milestone, 0)
	for i, m := range req.Milestones {
		mID := fmt.Sprintf("ms_%s_%d", contractID[5:], i+1)
		m.ID = mID
		m.ContractID = contractID
		m.Sequence = i + 1
		m.Status = domain.MilestoneStatusFunded
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
		h.store.Milestones[m.ID] = m
		createdMilestones = append(createdMilestones, m)
	}
	h.store.Unlock()

	c.Milestones = createdMilestones
	if repository.PG != nil {
		repository.PG.SaveContract(c)
	}
	h.store.RecordDoubleEntry(c.ID, "", "ESCROW_DEPOSIT", "Buyer funded escrow contract via Virtual Account", "BUYER_WALLET", "ESCROW_VAULT", c.TotalAmount)

	response.JSON(w, http.StatusCreated, c, "Contract created and funded into Escrow Vault successfully")
}

type SubmitProofReq struct {
	ProofURL string `json:"deliverable_proof_url"`
	Notes    string `json:"inspection_notes"`
}

func (h *Handler) SubmitMilestoneProof(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req SubmitProofReq
	json.NewDecoder(r.Body).Decode(&req)

	h.store.Lock()
	m, ok := h.store.Milestones[id]
	if !ok {
		h.store.Unlock()
		response.Error(w, http.StatusNotFound, "Milestone not found")
		return
	}

	now := time.Now()
	m.Status = domain.MilestoneStatusInInspection
	m.DeliverableProofURL = req.ProofURL
	m.InspectionNotes = req.Notes
	m.SubmittedAt = &now
	m.UpdatedAt = now
	h.store.Milestones[id] = m
	h.store.Unlock()

	response.JSON(w, http.StatusOK, m, "Milestone proof submitted for inspection")
}

func (h *Handler) ApproveAndReleaseMilestone(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	h.store.Lock()
	m, ok := h.store.Milestones[id]
	if !ok {
		h.store.Unlock()
		response.Error(w, http.StatusNotFound, "Milestone not found")
		return
	}

	if m.Status == domain.MilestoneStatusReleased {
		h.store.Unlock()
		response.Error(w, http.StatusBadRequest, "Milestone funds already released")
		return
	}

	now := time.Now()
	m.Status = domain.MilestoneStatusReleased
	m.ApprovedAt = &now
	m.ReleasedAt = &now
	m.UpdatedAt = now
	h.store.Milestones[id] = m
	h.store.Unlock()

	h.store.RecordDoubleEntry(m.ContractID, m.ID, "MILESTONE_RELEASE", fmt.Sprintf("Released milestone %d funds to seller", m.Sequence), "ESCROW_VAULT", "SELLER_PAYOUT", m.Amount)

	response.JSON(w, http.StatusOK, m, "Milestone approved! Funds instantly released to seller account.")
}

// -------------------------------------------------------------
// 5. Connections & Networking
// -------------------------------------------------------------

func (h *Handler) ListConnections(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	defer h.store.RUnlock()

	list := make([]domain.Connection, 0, len(h.store.Connections))
	for _, c := range h.store.Connections {
		list = append(list, c)
	}
	response.JSON(w, http.StatusOK, list, "Connected trading partners retrieved")
}

func (h *Handler) ListConnectionRequests(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	defer h.store.RUnlock()

	list := make([]domain.ConnectionRequest, 0, len(h.store.ConnRequests))
	for _, req := range h.store.ConnRequests {
		list = append(list, req)
	}
	response.JSON(w, http.StatusOK, list, "Connection requests retrieved")
}

type SendConnReq struct {
	FromID string `json:"from_id"`
	ToID   string `json:"to_id"`
}

func (h *Handler) SendConnectionRequest(w http.ResponseWriter, r *http.Request) {
	var req SendConnReq
	json.NewDecoder(r.Body).Decode(&req)

	h.store.Lock()
	defer h.store.Unlock()

	fromName := "Business"
	if u, ok := h.store.Users[req.FromID]; ok {
		fromName = u.BusinessName
	}
	toName := "Partner"
	if u, ok := h.store.Users[req.ToID]; ok {
		toName = u.BusinessName
	}

	cr := domain.ConnectionRequest{
		ID:        "req_" + uuid.New().String()[:8],
		FromID:    req.FromID,
		FromName:  fromName,
		ToID:      req.ToID,
		ToName:    toName,
		Status:    "pending",
		Timestamp: time.Now(),
	}
	h.store.ConnRequests[cr.ID] = cr

	response.JSON(w, http.StatusCreated, cr, "Connection request sent successfully")
}

func (h *Handler) AcceptConnectionRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	h.store.Lock()
	defer h.store.Unlock()

	cr, ok := h.store.ConnRequests[id]
	if !ok {
		response.Error(w, http.StatusNotFound, "Request not found")
		return
	}

	cr.Status = "accepted"
	h.store.ConnRequests[id] = cr

	newConn := domain.Connection{
		ID:           "conn_" + uuid.New().String()[:8],
		BuyerID:      cr.FromID,
		BuyerName:    cr.FromName,
		SupplierID:   cr.ToID,
		SupplierName: cr.ToName,
		Timestamp:    time.Now(),
	}
	h.store.Connections[newConn.ID] = newConn

	response.JSON(w, http.StatusOK, newConn, "Connection request accepted and active")
}

// -------------------------------------------------------------
// 6. Dispute Resolution & Arbitration Court
// -------------------------------------------------------------

func (h *Handler) ListDisputes(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	defer h.store.RUnlock()

	list := make([]domain.Dispute, 0, len(h.store.Disputes))
	for _, d := range h.store.Disputes {
		list = append(list, d)
	}
	response.JSON(w, http.StatusOK, list, "Disputes list retrieved")
}

type RaiseDisputeReq struct {
	ProposalID  string  `json:"proposal_id,omitempty"`
	ContractID  string  `json:"contract_id,omitempty"`
	MilestoneID string  `json:"milestone_id,omitempty"`
	Reason      string  `json:"reason"`
	ClaimAmount float64 `json:"claim_amount"`
}

func (h *Handler) RaiseDispute(w http.ResponseWriter, r *http.Request) {
	var req RaiseDisputeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	dispID := "disp_" + uuid.New().String()[:8]
	disp := domain.Dispute{
		ID:                dispID,
		ProposalID:        req.ProposalID,
		ContractID:        req.ContractID,
		ContractTitle:     "Supply Contract",
		MilestoneID:       req.MilestoneID,
		InitiatorOrgID:    "org_buyer_01",
		InitiatorOrgName:  "Acme Industrial Ltd",
		RespondentOrgID:   "org_seller_01",
		RespondentOrgName: "Rajesh Exports",
		Reason:            req.Reason,
		ClaimAmount:       req.ClaimAmount,
		Status:            domain.DisputeStatusUnderReview,
		CreatedAt:         time.Now(),
	}

	h.store.Lock()
	h.store.Disputes[disp.ID] = disp

	if req.ProposalID != "" {
		if p, ok := h.store.Proposals[req.ProposalID]; ok {
			p.Status = "Disputed"
			p.PaymentStatus = "Dispute Raised"
			h.store.Proposals[req.ProposalID] = p
		}
	}
	if req.MilestoneID != "" {
		if m, ok := h.store.Milestones[req.MilestoneID]; ok {
			m.Status = domain.MilestoneStatusDisputed
			h.store.Milestones[req.MilestoneID] = m
		}
	}
	if repository.PG != nil {
		repository.PG.SaveDispute(disp)
	}
	h.store.Unlock()

	response.JSON(w, http.StatusCreated, disp, "Dispute registered and assigned to TradeShield Arbitration Panel")
}

type ArbitrateReq struct {
	Verdict            string  `json:"resolution_verdict"` // REFUND_BUYER, RELEASE_SELLER, SPLIT
	BuyerRefundShare   float64 `json:"buyer_refund_share"`
	SellerReleaseShare float64 `json:"seller_release_share"`
}

func (h *Handler) ArbitrateDispute(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ArbitrateReq
	json.NewDecoder(r.Body).Decode(&req)

	h.store.Lock()
	defer h.store.Unlock()

	disp, ok := h.store.Disputes[id]
	if !ok {
		response.Error(w, http.StatusNotFound, "Dispute record not found")
		return
	}

	now := time.Now()
	disp.ResolutionVerdict = req.Verdict
	disp.BuyerRefundShare = req.BuyerRefundShare
	disp.SellerReleaseShare = req.SellerReleaseShare
	disp.ResolvedAt = &now

	if req.Verdict == "REFUND_BUYER" {
		disp.Status = domain.DisputeStatusResolvedRefund
	} else {
		disp.Status = domain.DisputeStatusResolvedReleased
	}
	h.store.Disputes[id] = disp
	if repository.PG != nil {
		repository.PG.SaveDispute(disp)
	}

	response.JSON(w, http.StatusOK, disp, "Arbitration verdict executed and recorded")
}

// -------------------------------------------------------------
// 7. Double-Entry Ledger Endpoints
// -------------------------------------------------------------

func (h *Handler) ListLedgerAccounts(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	defer h.store.RUnlock()

	list := make([]domain.LedgerAccount, 0, len(h.store.LedgerAccts))
	for _, a := range h.store.LedgerAccts {
		list = append(list, a)
	}
	response.JSON(w, http.StatusOK, list, "Chart of ledger accounts")
}

func (h *Handler) ListLedgerJournals(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	defer h.store.RUnlock()

	response.JSON(w, http.StatusOK, h.store.Journals, "Ledger journal entries audit trail")
}

// -------------------------------------------------------------
// 8. Support Ticket Desk
// -------------------------------------------------------------

func (h *Handler) SubmitSupportTicket(w http.ResponseWriter, r *http.Request) {
	var req domain.SupportTicket
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid ticket payload")
		return
	}

	req.ID = "tkt_" + uuid.New().String()[:8]
	req.Status = "open"
	req.CreatedAt = time.Now()

	h.store.Lock()
	h.store.Tickets[req.ID] = req
	if repository.PG != nil {
		repository.PG.SaveSupportTicket(req)
	}
	h.store.Unlock()

	response.JSON(w, http.StatusCreated, req, "Support ticket logged. SLA response in 2 hours.")
}

func (h *Handler) ListSupportTickets(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	defer h.store.RUnlock()

	list := make([]domain.SupportTicket, 0, len(h.store.Tickets))
	for _, t := range h.store.Tickets {
		list = append(list, t)
	}
	response.JSON(w, http.StatusOK, list, "Support tickets list")
}

// -------------------------------------------------------------
// 9. KYC & Bank Verifier APIs
// -------------------------------------------------------------

type PennyDropReq struct {
	AccountNumber string `json:"account_number"`
	IFSC          string `json:"ifsc_code"`
}

func (h *Handler) MockPennyDrop(w http.ResponseWriter, r *http.Request) {
	var req PennyDropReq
	json.NewDecoder(r.Body).Decode(&req)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"status":            "SUCCESS",
		"account_number":    req.AccountNumber,
		"ifsc":              req.IFSC,
		"registered_name":   "RAJESH EXPORTS LIMITED",
		"penny_drop_ref_id": "IMPS" + uuid.New().String()[:10],
		"is_name_match":     true,
		"message":           "₹1.00 successfully credited and verified via IMPS",
	}, "Penny drop bank verification successful")
}

func (h *Handler) MockGSTVerify(w http.ResponseWriter, r *http.Request) {
	gstin := chi.URLParam(r, "gstin")
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"gstin":         gstin,
		"legal_name":    "ACME INDUSTRIAL LIMITED",
		"trade_name":    "ACME INDUSTRIAL",
		"status":        "ACTIVE",
		"taxpayer_type": "Regular",
		"state_code":    "07",
		"trust_score":   95,
		"is_valid":      true,
	}, "GSTIN verified successfully with GSTN Portal")
}
