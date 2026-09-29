package repository

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"tradeshield-backend/internal/domain"
)

type Store struct {
	mu           sync.RWMutex
	Users        map[string]domain.User
	Orgs         map[string]domain.Organization
	BankAccounts map[string]domain.BankAccount
	Contracts    map[string]domain.Contract
	Milestones   map[string]domain.Milestone
	Proposals    map[string]domain.Proposal
	Approvals    map[string]domain.ProposalApproval
	Connections  map[string]domain.Connection
	ConnRequests map[string]domain.ConnectionRequest
	Tickets      map[string]domain.SupportTicket
	LedgerAccts  map[string]domain.LedgerAccount
	Journals     []domain.JournalEntry
	Disputes     map[string]domain.Dispute
}

var DB = NewStore()

func NewStore() *Store {
	return &Store{
		Users:        make(map[string]domain.User),
		Orgs:         make(map[string]domain.Organization),
		BankAccounts: make(map[string]domain.BankAccount),
		Contracts:    make(map[string]domain.Contract),
		Milestones:   make(map[string]domain.Milestone),
		Proposals:    make(map[string]domain.Proposal),
		Approvals:    make(map[string]domain.ProposalApproval),
		Connections:  make(map[string]domain.Connection),
		ConnRequests: make(map[string]domain.ConnectionRequest),
		Tickets:      make(map[string]domain.SupportTicket),
		LedgerAccts:  make(map[string]domain.LedgerAccount),
		Journals:     make([]domain.JournalEntry, 0),
		Disputes:     make(map[string]domain.Dispute),
	}
}

func (s *Store) SeedData() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Seed Organizations
	buyerOrg := domain.Organization{
		ID:            "org_buyer_01",
		Name:          "Acme Industrial Ltd",
		LegalName:     "Acme Industrial Limited",
		GSTIN:         "07AAAAA1111A1ZA",
		PAN:           "AAAAA1111A",
		OrgType:       domain.OrgTypeBuyer,
		TrustScore:    95,
		IsKYCVerified: true,
		Address:       "Okhla Phase 3",
		City:          "Delhi",
		State:         "Delhi",
		Pincode:       "110020",
		CreatedAt:     time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:     time.Now(),
	}

	sellerOrg := domain.Organization{
		ID:            "org_seller_01",
		Name:          "Rajesh Exports",
		LegalName:     "Rajesh Exports Limited",
		GSTIN:         "03CCCCC3333C3ZC",
		PAN:           "CCCCC3333C",
		OrgType:       domain.OrgTypeSupplier,
		TrustScore:    98,
		IsKYCVerified: true,
		Address:       "Focal Point",
		City:          "Ludhiana",
		State:         "Punjab",
		Pincode:       "141010",
		CreatedAt:     time.Now().Add(-60 * 24 * time.Hour),
		UpdatedAt:     time.Now(),
	}

	s.Orgs[buyerOrg.ID] = buyerOrg
	s.Orgs[sellerOrg.ID] = sellerOrg

	// 2. Seed Users (Matching PayShieldX)
	uBuyer1 := domain.User{
		ID:           "user_buyer_1",
		Email:        "buyer1@acme.com",
		Password:     "password123",
		FullName:     "Vikram Malhotra",
		BusinessName: "Acme Industrial Ltd",
		GST:          "07AAAAA1111A1ZA",
		Mobile:       "9876543210",
		City:         "Delhi",
		Role:         domain.RoleBuyer,
		Verified:     true,
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-30 * 24 * time.Hour),
		UpdatedAt:    time.Now(),
	}
	uBuyer2 := domain.User{
		ID:           "user_buyer_2",
		Email:        "buyer2@bhandari.com",
		Password:     "password123",
		FullName:     "Ramesh Bhandari",
		BusinessName: "Bhandari Textiles",
		GST:          "08BBBBB2222B2ZB",
		Mobile:       "9876543211",
		City:         "Jaipur",
		Role:         domain.RoleBuyer,
		Verified:     true,
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-20 * 24 * time.Hour),
		UpdatedAt:    time.Now(),
	}
	uBuyer3 := domain.User{
		ID:           "user_buyer_3",
		Email:        "buyer3@sujit.com",
		Password:     "password123",
		FullName:     "Sujit Kumar",
		BusinessName: "Sujit Enterprises",
		GST:          "07EEEEE5555E5ZE",
		Mobile:       "9876543214",
		City:         "Delhi",
		Role:         domain.RoleBuyer,
		Verified:     true,
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-10 * 24 * time.Hour),
		UpdatedAt:    time.Now(),
	}
	uSupplier1 := domain.User{
		ID:           "user_supplier_1",
		Email:        "supplier1@rajesh.com",
		Password:     "password123",
		FullName:     "Rajesh Patel",
		BusinessName: "Rajesh Exports",
		GST:          "03CCCCC3333C3ZC",
		Mobile:       "9876543212",
		City:         "Ludhiana",
		Role:         domain.RoleSupplier,
		Verified:     true,
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-60 * 24 * time.Hour),
		UpdatedAt:    time.Now(),
	}
	uSupplier2 := domain.User{
		ID:           "user_supplier_2",
		Email:        "supplier2@abc.com",
		Password:     "password123",
		FullName:     "Alok Sharma",
		BusinessName: "ABC Traders",
		GST:          "09DDDDD4444D4ZD",
		Mobile:       "9876543213",
		City:         "Delhi",
		Role:         domain.RoleSupplier,
		Verified:     true,
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-40 * 24 * time.Hour),
		UpdatedAt:    time.Now(),
	}
	uAdmin := domain.User{
		ID:           "user_admin",
		Email:        "admin@payshieldx.in",
		Password:     "password123",
		FullName:     "Compliance Head",
		BusinessName: "PayShieldX Arbiter Panel",
		GST:          "N/A",
		Mobile:       "N/A",
		City:         "Bangalore",
		Role:         domain.RoleAdmin,
		Verified:     true,
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-90 * 24 * time.Hour),
		UpdatedAt:    time.Now(),
	}

	uApexBuyer := domain.User{
		ID:           "user_buyer_apex",
		Email:        "procurement@apexauto.in",
		Password:     "Apex@Shield2026",
		FullName:     "Vikram Malhotra",
		BusinessName: "Apex Auto Components Pvt Ltd",
		GST:          "27AAACA1234A1Z5",
		PAN:          "AAACA1234A",
		Mobile:       "9820123456",
		City:         "Pune",
		Role:         domain.RoleBuyer,
		Verified:     true,
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-45 * 24 * time.Hour),
		UpdatedAt:    time.Now(),
	}
	uBharatSupplier := domain.User{
		ID:           "user_supplier_bharat",
		Email:        "sales@bharatcastings.com",
		Password:     "Bharat@Shield2026",
		FullName:     "Rajesh Singhania",
		BusinessName: "Bharat Precision Castings Ltd",
		GST:          "24AABCB5678B1Z2",
		PAN:          "AABCB5678B",
		Mobile:       "9898123456",
		City:         "Vadodara",
		Role:         domain.RoleSupplier,
		Verified:     true,
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-60 * 24 * time.Hour),
		UpdatedAt:    time.Now(),
	}
	uCourtAdmin := domain.User{
		ID:           "user_admin_court",
		Email:        "court@tradeshield.in",
		Password:     "Arbiter@Shield2026",
		FullName:     "Justice (Retd.) K. N. Verma",
		BusinessName: "TradeShield Neutral Arbitration Panel",
		GST:          "07AAACT0001A1Z9",
		PAN:          "AAACT0001A",
		Mobile:       "9811000001",
		City:         "New Delhi",
		Role:         domain.RoleAdmin,
		Verified:     true,
		IsVerified:   true,
		CreatedAt:    time.Now().Add(-90 * 24 * time.Hour),
		UpdatedAt:    time.Now(),
	}

	s.Users[uBuyer1.ID] = uBuyer1
	s.Users[uBuyer2.ID] = uBuyer2
	s.Users[uBuyer3.ID] = uBuyer3
	s.Users[uApexBuyer.ID] = uApexBuyer
	s.Users[uSupplier1.ID] = uSupplier1
	s.Users[uSupplier2.ID] = uSupplier2
	s.Users[uBharatSupplier.ID] = uBharatSupplier
	s.Users[uAdmin.ID] = uAdmin
	s.Users[uCourtAdmin.ID] = uCourtAdmin

	// 3. Seed Proposals
	now := time.Now()
	p1 := domain.Proposal{
		ID:               "prop_1001",
		OrderID:          "ORD-55410",
		BuyerID:          uBuyer1.ID,
		BuyerName:        uBuyer1.BusinessName,
		SupplierID:       uSupplier1.ID,
		SupplierName:     uSupplier1.BusinessName,
		Amount:           175000,
		Currency:         "INR",
		Terms:            "partial",
		DeliveryTimeline: "2026-06-25",
		Notes:            "50% advance in escrow, 50% upon receipt of QC report.",
		Status:           "Confirmed",
		PaymentStatus:    "Payment Held in Escrow",
		BuyerApproved:    true,
		SupplierApproved: true,
		CreatedAt:        now.Add(-5 * 24 * time.Hour),
		UpdatedAt:        now.Add(-4 * 24 * time.Hour),
	}
	p2 := domain.Proposal{
		ID:               "prop_1002",
		OrderID:          "ORD-90214",
		BuyerID:          uBuyer1.ID,
		BuyerName:        uBuyer1.BusinessName,
		SupplierID:       uSupplier1.ID,
		SupplierName:     uSupplier1.BusinessName,
		Amount:           450000,
		Currency:         "INR",
		Terms:            "full",
		DeliveryTimeline: "2026-07-10",
		Notes:            "Total funds locked in escrow. Releases upon delivery check at Delhi warehouse.",
		Status:           "Approved by Supplier",
		PaymentStatus:    "Payment Held in Escrow",
		BuyerApproved:    false,
		SupplierApproved: true,
		CreatedAt:        now.Add(-2 * 24 * time.Hour),
		UpdatedAt:        now.Add(-1 * 24 * time.Hour),
	}
	p3 := domain.Proposal{
		ID:               "prop_1003",
		OrderID:          "ORD-20412",
		BuyerID:          uBuyer2.ID,
		BuyerName:        uBuyer2.BusinessName,
		SupplierID:       uSupplier1.ID,
		SupplierName:     uSupplier1.BusinessName,
		Amount:           85000,
		Currency:         "INR",
		Terms:            "advance",
		DeliveryTimeline: "2026-06-30",
		Notes:            "Advance payment protection. Materials must match sample cotton count 40s.",
		Status:           "Modification Requested",
		PaymentStatus:    "Unpaid",
		BuyerApproved:    false,
		SupplierApproved: false,
		CreatedAt:        now.Add(-3 * 24 * time.Hour),
		UpdatedAt:        now.Add(-2 * 24 * time.Hour),
	}

	s.Proposals[p1.ID] = p1
	s.Proposals[p2.ID] = p2
	s.Proposals[p3.ID] = p3

	// 4. Seed Connections & Requests
	conn1 := domain.Connection{
		ID:           "conn_1",
		BuyerID:      uBuyer1.ID,
		BuyerName:    uBuyer1.BusinessName,
		SupplierID:   uSupplier1.ID,
		SupplierName: uSupplier1.BusinessName,
		Timestamp:    now.Add(-10 * 24 * time.Hour),
	}
	conn2 := domain.Connection{
		ID:           "conn_2",
		BuyerID:      uBuyer2.ID,
		BuyerName:    uBuyer2.BusinessName,
		SupplierID:   uSupplier1.ID,
		SupplierName: uSupplier1.BusinessName,
		Timestamp:    now.Add(-8 * 24 * time.Hour),
	}
	s.Connections[conn1.ID] = conn1
	s.Connections[conn2.ID] = conn2

	req1 := domain.ConnectionRequest{
		ID:        "req_1",
		FromID:    uBuyer3.ID,
		FromName:  uBuyer3.BusinessName,
		ToID:      uSupplier2.ID,
		ToName:    uSupplier2.BusinessName,
		Status:    "pending",
		Timestamp: now.Add(-1 * 24 * time.Hour),
	}
	s.ConnRequests[req1.ID] = req1

	// 5. Seed Contract & Milestones (For Milestone Ledger View)
	contractID := "cntr_2026_089"
	m1ID := "ms_089_1"
	m2ID := "ms_089_2"
	m3ID := "ms_089_3"

	subTime := now.Add(-3 * 24 * time.Hour)
	appTime := now.Add(-2 * 24 * time.Hour)
	relTime := now.Add(-2 * 24 * time.Hour)

	m1 := domain.Milestone{
		ID:                  m1ID,
		ContractID:          contractID,
		Sequence:            1,
		Title:               "20% Advance for Raw Material Procurement (Grade A SG Iron)",
		Description:         "Advance mobilization fund for raw material purchase with mill test certificates",
		Percentage:          20,
		Amount:              500000,
		Status:              domain.MilestoneStatusReleased,
		DeliverableProofURL: "https://docs.tradeshield.in/proofs/raw_material_cert_089.pdf",
		InspectionNotes:     "Mill test certificate verified against ASTM A536 standards.",
		SubmittedAt:         &subTime,
		ApprovedAt:          &appTime,
		ReleasedAt:          &relTime,
		DueDate:             now.Add(-10 * 24 * time.Hour),
		CreatedAt:           now.Add(-20 * 24 * time.Hour),
		UpdatedAt:           relTime,
	}

	sub2 := now.Add(-1 * 24 * time.Hour)
	m2 := domain.Milestone{
		ID:                  m2ID,
		ContractID:          contractID,
		Sequence:            2,
		Title:               "40% On Dispatch with Lorry Receipt (LR) & E-Way Bill",
		Description:         "Machining completed, dispatched via V-Trans logistics (LR #VT-982134)",
		Percentage:          40,
		Amount:              1000000,
		Status:              domain.MilestoneStatusInInspection,
		DeliverableProofURL: "https://docs.tradeshield.in/proofs/eway_bill_lr_089.pdf",
		InspectionNotes:     "Consignment in transit. Expected delivery in Pune warehouse by tomorrow.",
		SubmittedAt:         &sub2,
		DueDate:             now.Add(2 * 24 * time.Hour),
		CreatedAt:           now.Add(-20 * 24 * time.Hour),
		UpdatedAt:           now,
	}

	m3 := domain.Milestone{
		ID:          m3ID,
		ContractID:  contractID,
		Sequence:    3,
		Title:       "40% Final QC Inspection & Warehouse Receiving",
		Description: "Dimensional tolerance check (CMM report) and receiving in inventory",
		Percentage:  40,
		Amount:      1000000,
		Status:      domain.MilestoneStatusFunded,
		DueDate:     now.Add(10 * 24 * time.Hour),
		CreatedAt:   now.Add(-20 * 24 * time.Hour),
		UpdatedAt:   now,
	}

	contract := domain.Contract{
		ID:                   contractID,
		Title:                "Supply of 5,000 Precision Cast Flanges & Valve Housings",
		ContractNumber:       "TS-CTR-2026-089",
		BuyerOrgID:           buyerOrg.ID,
		BuyerOrgName:         buyerOrg.Name,
		SupplierOrgID:        sellerOrg.ID,
		SupplierOrgName:      sellerOrg.Name,
		TotalAmount:          2500000,
		Currency:             "INR",
		PlatformFeePercent:   0.75,
		PlatformFeeAmount:    18750,
		EscrowVirtualAccount: "ICIC0000104TS089",
		Status:               domain.ContractStatusInProgress,
		Description:          "High precision grade ASTM A536 ductile iron castings with CNC finish according to drawing #DWG-AF-2026.",
		DeliveryTerms:        "DAP Pune Warehouse (Incoterms 2020)",
		InspectionPeriodDays: 5,
		Milestones:           []domain.Milestone{m1, m2, m3},
		CreatedAt:            now.Add(-20 * 24 * time.Hour),
		UpdatedAt:            now,
	}

	s.Contracts[contract.ID] = contract
	s.Milestones[m1.ID] = m1
	s.Milestones[m2.ID] = m2
	s.Milestones[m3.ID] = m3

	// 6. Seed Ledger Chart of Accounts
	s.LedgerAccts["acc_vault"] = domain.LedgerAccount{ID: "acc_vault", Name: "ESCROW_VAULT", Type: "ASSET", Balance: 2000000}
	s.LedgerAccts["acc_buyer"] = domain.LedgerAccount{ID: "acc_buyer", OrgID: buyerOrg.ID, Name: "BUYER_DEPOSIT", Type: "LIABILITY", Balance: 2500000}
	s.LedgerAccts["acc_seller"] = domain.LedgerAccount{ID: "acc_seller", OrgID: sellerOrg.ID, Name: "SELLER_PAYOUT", Type: "LIABILITY", Balance: 500000}
	s.LedgerAccts["acc_fee"] = domain.LedgerAccount{ID: "acc_fee", Name: "PLATFORM_REVENUE", Type: "REVENUE", Balance: 18750}
}

func (s *Store) GetEscrowSummary() domain.EscrowVaultSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var locked, released, disputed, fee float64

	// Count from proposals
	for _, p := range s.Proposals {
		if p.PaymentStatus == "Payment Held in Escrow" || p.PaymentStatus == "Shipped" {
			locked += p.Amount
		} else if p.PaymentStatus == "Payment Released" || p.PaymentStatus == "Delivered" {
			released += p.Amount
		} else if p.PaymentStatus == "Dispute Raised" {
			disputed += p.Amount
		}
		fee += p.Amount * 0.0075
	}

	// Count from milestone contracts
	for _, m := range s.Milestones {
		switch m.Status {
		case domain.MilestoneStatusFunded, domain.MilestoneStatusWorkSubmitted, domain.MilestoneStatusInInspection:
			locked += m.Amount
		case domain.MilestoneStatusReleased:
			released += m.Amount
		case domain.MilestoneStatusDisputed:
			disputed += m.Amount
		}
	}

	for _, c := range s.Contracts {
		fee += c.PlatformFeeAmount
	}

	totalDeals := len(s.Contracts) + len(s.Proposals)

	return domain.EscrowVaultSummary{
		TotalLockedINR:   locked,
		TotalReleasedINR: released,
		TotalDisputedINR: disputed,
		PlatformFeeINR:   fee,
		ActiveDealsCount: totalDeals,
	}
}

func (s *Store) RecordDoubleEntry(contractID, milestoneID, ref, desc string, debitAcc, creditAcc string, amount float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	journalID := "jrn_" + uuid.New().String()[:8]
	postings := []domain.LedgerTransaction{
		{
			ID:             "tx_" + uuid.New().String()[:8],
			JournalEntryID: journalID,
			AccountID:      debitAcc,
			AccountName:    debitAcc,
			DebitAmount:    amount,
			CreditAmount:   0,
			CreatedAt:      time.Now(),
		},
		{
			ID:             "tx_" + uuid.New().String()[:8],
			JournalEntryID: journalID,
			AccountID:      creditAcc,
			AccountName:    creditAcc,
			DebitAmount:    0,
			CreditAmount:   amount,
			CreatedAt:      time.Now(),
		},
	}

	entry := domain.JournalEntry{
		ID:          journalID,
		ContractID:  contractID,
		MilestoneID: milestoneID,
		Reference:   ref,
		Description: desc,
		CreatedAt:   time.Now(),
		Postings:    postings,
	}

	s.Journals = append(s.Journals, entry)
	if PG != nil {
		PG.SaveJournal(entry)
	}

	return nil
}

func (s *Store) Lock()    { s.mu.Lock() }
func (s *Store) Unlock()  { s.mu.Unlock() }
func (s *Store) RLock()   { s.mu.RLock() }
func (s *Store) RUnlock() { s.mu.RUnlock() }
