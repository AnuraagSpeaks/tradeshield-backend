package domain

import (
	"time"
)

const (
	RoleBuyer              = "buyer"
	RoleSupplier           = "supplier"
	RoleAdmin              = "admin"
	RoleBuyerAdmin         = "BUYER_ADMIN"
	RoleBuyerMember        = "BUYER_MEMBER"
	RoleSellerAdmin        = "SELLER_ADMIN"
	RoleSellerMember       = "SELLER_MEMBER"
	RoleArbitrator         = "ARBITRATOR"
	RolePlatformSuperAdmin = "SUPER_ADMIN"
)

const (
	OrgTypeBuyer    = "BUYER"
	OrgTypeSupplier = "SUPPLIER"
	OrgTypeBoth     = "BOTH"
)

const (
	ContractStatusDraft             = "DRAFT"
	ContractStatusPendingAcceptance = "PENDING_ACCEPTANCE"
	ContractStatusAwaitingFunding   = "AWAITING_FUNDING"
	ContractStatusFunded            = "FUNDED"
	ContractStatusInProgress        = "IN_PROGRESS"
	ContractStatusCompleted         = "COMPLETED"
	ContractStatusDisputed          = "DISPUTED"
	ContractStatusCancelled         = "CANCELLED"
)

const (
	MilestoneStatusPending       = "PENDING"
	MilestoneStatusFunded        = "FUNDED"
	MilestoneStatusWorkSubmitted = "WORK_SUBMITTED"
	MilestoneStatusInInspection  = "IN_INSPECTION"
	MilestoneStatusApproved      = "APPROVED"
	MilestoneStatusReleased      = "RELEASED"
	MilestoneStatusDisputed      = "DISPUTED"
	MilestoneStatusRefunded      = "REFUNDED"
)

const (
	DisputeStatusRaised           = "RAISED"
	DisputeStatusUnderReview      = "Under Review"
	DisputeStatusEvidenceReq      = "Evidence Requested"
	DisputeStatusResolvedRefund   = "Resolved (Refund)"
	DisputeStatusResolvedReleased = "Resolved (Released to Supplier)"
	DisputeStatusClosed           = "CLOSED"
)

type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Password       string    `json:"password,omitempty"`
	FullName       string    `json:"full_name"`
	ContactPerson  string    `json:"contact_person,omitempty"`
	BusinessName   string    `json:"business_name"`
	GST            string    `json:"gst"`
	PAN            string    `json:"pan,omitempty"`
	Mobile         string    `json:"mobile"`
	City           string    `json:"city"`
	Role           string    `json:"role"`
	OrganizationID string    `json:"organization_id,omitempty"`
	Verified       bool      `json:"verified"`
	IsVerified     bool      `json:"is_verified"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Organization struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	LegalName     string    `json:"legal_name"`
	GSTIN         string    `json:"gstin"`
	PAN           string    `json:"pan"`
	OrgType       string    `json:"org_type"`
	TrustScore    int       `json:"trust_score"`
	IsKYCVerified bool      `json:"is_kyc_verified"`
	Address       string    `json:"address"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Pincode       string    `json:"pincode"`
	BankAccountID string    `json:"bank_account_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type BankAccount struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	AccountNumber  string    `json:"account_number"`
	IFSCCode       string    `json:"ifsc_code"`
	BankName       string    `json:"bank_name"`
	AccountHolder  string    `json:"account_holder"`
	IsPennyDropped bool      `json:"is_penny_dropped"`
	IsVerified     bool      `json:"is_verified"`
	CreatedAt      time.Time `json:"created_at"`
}

type Proposal struct {
	ID                 string     `json:"id"`
	ProposalNumber     string     `json:"proposal_number"`
	OrderID            string     `json:"order_id"`
	BuyerID            string     `json:"buyer_id"`
	BuyerName          string     `json:"buyer_name"`
	BuyerSignatory     string     `json:"buyer_signatory,omitempty"`
	BuyerEmail         string     `json:"buyer_email,omitempty"`
	BuyerGSTIN         string     `json:"buyer_gstin,omitempty"`
	BuyerAddress       string     `json:"buyer_address,omitempty"`
	SupplierID         string     `json:"supplier_id"`
	SupplierName       string     `json:"supplier_name"`
	SupplierSignatory   string     `json:"supplier_signatory,omitempty"`
	SupplierEmail       string     `json:"supplier_email,omitempty"`
	SupplierGSTIN       string     `json:"supplier_gstin,omitempty"`
	SupplierAddress     string     `json:"supplier_address,omitempty"`
	ItemDescription    string     `json:"item_description,omitempty"`
	BaseAmount         float64    `json:"base_amount"`
	DiscountPercent    float64    `json:"discount_percent"`
	DiscountAmount     float64    `json:"discount_amount"`
	TaxPercent         float64    `json:"tax_percent"`
	TaxAmount          float64    `json:"tax_amount"`
	TotalPayableAmount float64    `json:"total_payable_amount"`
	Amount             float64    `json:"amount"` // Total payable / deal amount
	Currency           string     `json:"currency"`
	Terms              string     `json:"terms"`
	MilestonesSummary  string     `json:"milestones_summary,omitempty"`
	DeliveryTimeline   string     `json:"delivery_timeline"`
	Notes              string     `json:"notes"`
	Status             string     `json:"status"` // PENDING_BUYER_APPROVAL, APPROVED, REJECTED, DISPATCHED, DELIVERED
	PaymentStatus      string     `json:"payment_status"` // Pending Escrow, Payment Held in Escrow, Released
	BuyerApproved      bool       `json:"buyer_approved"`
	SupplierApproved   bool       `json:"supplier_approved"`
	ContractID         string     `json:"contract_id,omitempty"`
	LRNumber           string     `json:"lr_number,omitempty"`
	TransporterName    string     `json:"transporter_name,omitempty"`
	ProofURL           string     `json:"proof_url,omitempty"`
	ValidTillDate      *time.Time `json:"valid_till_date,omitempty"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type ProposalApproval struct {
	ID         string    `json:"id"`
	ProposalID string    `json:"proposal_id"`
	UserID     string    `json:"user_id"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	Timestamp  time.Time `json:"timestamp"`
}

type Connection struct {
	ID           string    `json:"id"`
	BuyerID      string    `json:"buyer_id"`
	BuyerName    string    `json:"buyer_name,omitempty"`
	SupplierID   string    `json:"supplier_id"`
	SupplierName string    `json:"supplier_name,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

type ConnectionRequest struct {
	ID        string    `json:"id"`
	FromID    string    `json:"from_id"`
	FromName  string    `json:"from_name,omitempty"`
	ToID      string    `json:"to_id"`
	ToName    string    `json:"to_name,omitempty"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type SupportTicket struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Company   string    `json:"company"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Contract struct {
	ID                   string      `json:"id"`
	Title                string      `json:"title"`
	ContractNumber       string      `json:"contract_number"`
	BuyerOrgID           string      `json:"buyer_org_id"`
	BuyerOrgName         string      `json:"buyer_org_name,omitempty"`
	SupplierOrgID        string      `json:"supplier_org_id"`
	SupplierOrgName      string      `json:"supplier_org_name,omitempty"`
	TotalAmount          float64     `json:"total_amount"`
	Currency             string      `json:"currency"`
	PlatformFeePercent   float64     `json:"platform_fee_percent"`
	PlatformFeeAmount    float64     `json:"platform_fee_amount"`
	EscrowVirtualAccount string      `json:"escrow_virtual_account"`
	Status               string      `json:"status"`
	Description          string      `json:"description"`
	DeliveryTerms        string      `json:"delivery_terms"`
	InspectionPeriodDays int         `json:"inspection_period_days"`
	Milestones           []Milestone `json:"milestones,omitempty"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

type Milestone struct {
	ID                  string     `json:"id"`
	ContractID          string     `json:"contract_id"`
	Sequence            int        `json:"sequence"`
	Title               string     `json:"title"`
	Description         string     `json:"description"`
	Percentage          float64    `json:"percentage"`
	Amount              float64    `json:"amount"`
	Status              string     `json:"status"`
	DeliverableProofURL string     `json:"deliverable_proof_url,omitempty"`
	InspectionNotes     string     `json:"inspection_notes,omitempty"`
	SubmittedAt         *time.Time `json:"submitted_at,omitempty"`
	ApprovedAt          *time.Time `json:"approved_at,omitempty"`
	ReleasedAt          *time.Time `json:"released_at,omitempty"`
	DueDate             time.Time  `json:"due_date"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type LedgerAccount struct {
	ID        string    `json:"id"`
	OrgID     string    `json:"org_id,omitempty"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type JournalEntry struct {
	ID          string              `json:"id"`
	ContractID  string              `json:"contract_id,omitempty"`
	MilestoneID string              `json:"milestone_id,omitempty"`
	ProposalID  string              `json:"proposal_id,omitempty"`
	Reference   string              `json:"reference"`
	Description string              `json:"description"`
	CreatedAt   time.Time           `json:"created_at"`
	Postings    []LedgerTransaction `json:"postings"`
}

type LedgerTransaction struct {
	ID             string    `json:"id"`
	JournalEntryID string    `json:"journal_entry_id"`
	AccountID      string    `json:"account_id"`
	AccountName    string    `json:"account_name"`
	DebitAmount    float64   `json:"debit_amount"`
	CreditAmount   float64   `json:"credit_amount"`
	CreatedAt      time.Time `json:"created_at"`
}

type Dispute struct {
	ID                 string     `json:"id"`
	ProposalID         string     `json:"proposal_id,omitempty"`
	ContractID         string     `json:"contract_id,omitempty"`
	ContractTitle      string     `json:"contract_title,omitempty"`
	MilestoneID        string     `json:"milestone_id,omitempty"`
	MilestoneTitle     string     `json:"milestone_title,omitempty"`
	InitiatorOrgID     string     `json:"initiator_org_id,omitempty"`
	InitiatorOrgName   string     `json:"initiator_org_name,omitempty"`
	RespondentOrgID    string     `json:"respondent_org_id,omitempty"`
	RespondentOrgName  string     `json:"respondent_org_name,omitempty"`
	RaisedBy           string     `json:"raised_by,omitempty"`
	Reason             string     `json:"reason"`
	ClaimAmount        float64    `json:"claim_amount"`
	Status             string     `json:"status"`
	EvidenceURLs       []string   `json:"evidence_urls,omitempty"`
	ArbitratorID       string     `json:"arbitrator_id,omitempty"`
	ResolutionVerdict  string     `json:"resolution_verdict,omitempty"`
	BuyerRefundShare   float64    `json:"buyer_refund_share"`
	SellerReleaseShare float64    `json:"seller_release_share"`
	CreatedAt          time.Time  `json:"created_at"`
	ResolvedAt         *time.Time `json:"resolved_at,omitempty"`
}

type EscrowVaultSummary struct {
	TotalLockedINR   float64 `json:"total_locked_inr"`
	TotalReleasedINR float64 `json:"total_released_inr"`
	TotalDisputedINR float64 `json:"total_disputed_inr"`
	PlatformFeeINR   float64 `json:"platform_fee_inr"`
	ActiveDealsCount int     `json:"active_deals_count"`
}

type EmailOTP struct {
	Email     string    `json:"email"`
	OTP       string    `json:"otp"`
	Purpose   string    `json:"purpose"`
	ExpiresAt time.Time `json:"expires_at"`
}

type KYCDocument struct {
	Type        string    `json:"type"`
	DocumentNo  string    `json:"document_no"`
	URL         string    `json:"url"`
	Status      string    `json:"status"` // Approved, Pending, Rejected
	UploadedAt  time.Time `json:"uploaded_at"`
	VerifiedAt  *time.Time `json:"verified_at,omitempty"`
}

type RefundRecord struct {
	RefundID      string    `json:"refund_id"`
	TransactionID string    `json:"transaction_id"`
	Amount        float64   `json:"amount"`
	Reason        string    `json:"reason"`
	Status        string    `json:"status"`
	UTRNumber     string    `json:"utr_number,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type SettlementDetails struct {
	AccountNumber  string `json:"account_number"`
	BankName       string `json:"bank_name"`
	IFSCCode       string `json:"ifsc_code"`
	AccountHolder  string `json:"account_holder"`
	IsPennyDropped bool   `json:"is_penny_dropped"`
	PennyDropStatus string `json:"penny_drop_status"`
	SettlementCycle string `json:"settlement_cycle"`
}

type BuyerRecord struct {
	BuyerID                string         `json:"buyer_id"`
	Name                   string         `json:"name"`
	CompanyName            string         `json:"company_name"`
	Mobile                 string         `json:"mobile"`
	Email                  string         `json:"email"`
	GSTIN                  string         `json:"gstin"`
	PAN                    string         `json:"pan"`
	KYCStatus              string         `json:"kyc_status"` // Approved, Pending, Rejected
	CompletedTransactions  int            `json:"completed_transactions"`
	Disputes               int            `json:"disputes"`
	TotalTransactionValue  float64        `json:"total_transaction_value"`
	RefundHistory          []RefundRecord `json:"refund_history"`
	AccountStatus          string         `json:"account_status"` // Active, Suspended, Flagged
	CreatedAt              time.Time      `json:"created_at"`
}

type SupplierRecord struct {
	SupplierID            string            `json:"supplier_id"`
	CompanyName           string            `json:"company_name"`
	ContactPerson         string            `json:"contact_person"`
	Mobile                string            `json:"mobile"`
	Email                 string            `json:"email"`
	GSTIN                 string            `json:"gstin"`
	PAN                   string            `json:"pan"`
	KYCDocuments          []KYCDocument     `json:"kyc_documents"`
	VerificationStatus    string            `json:"verification_status"` // Approved, Pending, Rejected
	SupplierPlan          string            `json:"supplier_plan"`       // Growth, Business, Enterprise
	PlanExpiry            time.Time         `json:"plan_expiry"`
	TotalProposals        int               `json:"total_proposals"`
	AcceptedProposals     int               `json:"accepted_proposals"`
	RejectedProposals     int               `json:"rejected_proposals"`
	PendingProposals      int               `json:"pending_proposals"`
	CompletedTransactions int               `json:"completed_transactions"`
	Disputes              int               `json:"disputes"`
	TotalTransactionValue float64           `json:"total_transaction_value"`
	Refunds               float64           `json:"refunds"`
	SettlementInfo        SettlementDetails `json:"settlement_info"`
	AccountStatus         string            `json:"account_status"` // Active, Suspended, Under Review
	CreatedAt             time.Time         `json:"created_at"`
}

type SettlementItem struct {
	SettlementID string     `json:"settlement_id"`
	SupplierID   string     `json:"supplier_id"`
	SupplierName string     `json:"supplier_name"`
	BankName     string     `json:"bank_name"`
	AccountNo    string     `json:"account_no"`
	IFSCCode     string     `json:"ifsc_code"`
	Amount       float64    `json:"amount"`
	DealRef      string     `json:"deal_ref"`
	Status       string     `json:"status"` // PENDING, PROCESSING, SETTLED
	UTRNumber    string     `json:"utr_number,omitempty"`
	DisbursedAt  *time.Time `json:"disbursed_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type AdminBusinessHealthStats struct {
	TotalBuyers            int     `json:"total_buyers"`
	TotalSuppliers         int     `json:"total_suppliers"`
	ActiveUsers            int     `json:"active_users"`
	NewRegistrationsToday  int     `json:"new_registrations_today"`
	NewRegistrationsWeek   int     `json:"new_registrations_week"`
	NewRegistrationsMonth  int     `json:"new_registrations_month"`
	KYCPending             int     `json:"kyc_pending"`
	KYCApproved            int     `json:"kyc_approved"`
	KYCRejected            int     `json:"kyc_rejected"`
	TotalPaymentProposals  int     `json:"total_payment_proposals"`
	PendingProposals       int     `json:"pending_proposals"`
	ApprovedProposals      int     `json:"approved_proposals"`
	DisputedTransactions   int     `json:"disputed_transactions"`
	FailedTransactions     int     `json:"failed_transactions"`
	PlatformRevenue        float64 `json:"platform_revenue"`
	MembershipRevenue      float64 `json:"membership_revenue"`
	RefundAmount           float64 `json:"refund_amount"`
	PendingSettlements     float64 `json:"pending_settlements"`
	TodayCollection        float64 `json:"today_collection"`
	MonthlyRevenue         float64 `json:"monthly_revenue"`
}

type AdminFinanceStats struct {
	MembershipRevenue       float64          `json:"membership_revenue"`
	MembershipGrowth        float64          `json:"membership_growth"`
	MembershipBusiness      float64          `json:"membership_business"`
	MembershipEnterprise    float64          `json:"membership_enterprise"`
	GrowthCount             int              `json:"growth_count"`
	BusinessCount           int              `json:"business_count"`
	EnterpriseCount         int              `json:"enterprise_count"`
	OtherRevenue            float64          `json:"other_revenue"`
	TotalRevenue            float64          `json:"total_revenue"`
	EscrowNodalBalance      float64          `json:"escrow_nodal_balance"`
	PendingSettlementAmount float64          `json:"pending_settlement_amount"`
	RefundsTotal            float64          `json:"refunds_total"`
	TodayCollection         float64          `json:"today_collection"`
	MonthlyRevenue          float64          `json:"monthly_revenue"`
	PendingSettlements      []SettlementItem `json:"pending_settlements"`
}
