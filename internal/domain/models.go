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
	ID               string    `json:"id"`
	OrderID          string    `json:"order_id"`
	BuyerID          string    `json:"buyer_id"`
	BuyerName        string    `json:"buyer_name,omitempty"`
	SupplierID       string    `json:"supplier_id"`
	SupplierName     string    `json:"supplier_name,omitempty"`
	Amount           float64   `json:"amount"`
	Currency         string    `json:"currency"`
	Terms            string    `json:"terms"`
	DeliveryTimeline string    `json:"delivery_timeline"`
	Notes            string    `json:"notes"`
	Status           string    `json:"status"`
	PaymentStatus    string    `json:"payment_status"`
	BuyerApproved    bool      `json:"buyer_approved"`
	SupplierApproved bool      `json:"supplier_approved"`
	LRNumber         string    `json:"lr_number,omitempty"`
	TransporterName  string    `json:"transporter_name,omitempty"`
	ProofURL         string    `json:"proof_url,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
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
