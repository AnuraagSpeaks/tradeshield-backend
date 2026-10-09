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
	OTPs         map[string]domain.EmailOTP
	Buyers       map[string]domain.BuyerRecord
	Suppliers    map[string]domain.SupplierRecord
	Settlements  map[string]domain.SettlementItem
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
		OTPs:         make(map[string]domain.EmailOTP),
		Buyers:       make(map[string]domain.BuyerRecord),
		Suppliers:    make(map[string]domain.SupplierRecord),
		Settlements:  make(map[string]domain.SettlementItem),
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
		Email:        "court@payshieldx.in",
		Password:     "Arbiter@Shield2026",
		FullName:     "Justice (Retd.) K. N. Verma",
		BusinessName: "PayShield Neutral Arbitration Panel",
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
	validDate := now.Add(14 * 24 * time.Hour)
	p1 := domain.Proposal{
		ID:                 "prop_4776008",
		ProposalNumber:     "4776008",
		OrderID:            "ORD-4776008",
		BuyerID:            uApexBuyer.ID,
		BuyerName:          uApexBuyer.BusinessName,
		BuyerSignatory:     "Vikram Malhotra",
		BuyerEmail:         "procurement@apexauto.in",
		BuyerGSTIN:         "27AAACA1234A1Z5",
		BuyerAddress:       "Plot 42, MIDC Bhosari Industrial Area, Pune, Maharashtra 411026",
		SupplierID:         uBharatSupplier.ID,
		SupplierName:       uBharatSupplier.BusinessName,
		SupplierSignatory:   "Rajesh Singhania",
		SupplierEmail:       "sales@bharatcastings.com",
		SupplierGSTIN:       "24AABCB5678B1Z2",
		SupplierAddress:     "Survey No. 118, GIDC Makarpura Industrial Estate, Vadodara, Gujarat 390010",
		ItemDescription:    "Supply of 10,000 Precision Cast Flanges (Grade ASTM A105) & Hydrostatic Testing",
		BaseAmount:         250000.00,
		DiscountPercent:    20.00,
		DiscountAmount:     50000.00,
		TaxPercent:         18.00,
		TaxAmount:          36000.00,
		TotalPayableAmount: 236000.00,
		Amount:             236000.00,
		Currency:           "INR",
		Terms:              "3-Tranche Milestone: 20% Advance, 40% Dispatch, 40% Delivery Inspection",
		MilestonesSummary:  "20% Advance (QC Cert) • 40% Dispatch (LR Proof) • 40% Delivery (Warehouse Signoff)",
		DeliveryTimeline:   "21 Business Days",
		Notes:              "100% Escrow Backed via PayShieldX Nodal Trust Account. 48hr inspection window.",
		Status:             "PENDING_BUYER_APPROVAL",
		PaymentStatus:      "Pending Escrow",
		BuyerApproved:      false,
		SupplierApproved:   true,
		ValidTillDate:      &validDate,
		CreatedAt:          now.Add(-1 * 24 * time.Hour),
		UpdatedAt:          now.Add(-1 * 24 * time.Hour),
	}
	p2 := domain.Proposal{
		ID:                 "prop_4830735",
		ProposalNumber:     "4830735",
		OrderID:            "ORD-4830735",
		BuyerID:            uApexBuyer.ID,
		BuyerName:          uApexBuyer.BusinessName,
		BuyerSignatory:     "Vikram Malhotra",
		BuyerEmail:         "procurement@apexauto.in",
		BuyerGSTIN:         "27AAACA1234A1Z5",
		BuyerAddress:       "Plot 42, MIDC Bhosari Industrial Area, Pune, Maharashtra 411026",
		SupplierID:         uSupplier1.ID,
		SupplierName:       uSupplier1.BusinessName,
		SupplierSignatory:   "Kiran Patel",
		SupplierEmail:       "patel@rajesh.com",
		SupplierGSTIN:       "03CCCCC3333C3ZC",
		SupplierAddress:     "Industrial Area B, Focal Point, Ludhiana, Punjab 141010",
		ItemDescription:    "High-Tensile Fasteners & Cold Forged Bolt Assemblies (Grade 10.9)",
		BaseAmount:         1500000.00,
		DiscountPercent:    0.00,
		DiscountAmount:     0.00,
		TaxPercent:         18.00,
		TaxAmount:          270000.00,
		TotalPayableAmount: 1770000.00,
		Amount:             1770000.00,
		Currency:           "INR",
		Terms:              "Milestone based payment released strictly after metallurgy lab approval",
		MilestonesSummary:  "20% Advance • 40% Dispatch • 40% Delivery",
		DeliveryTimeline:   "15 Business Days",
		Notes:              "Full escrow protection enabled. Auto-arbitration clause active.",
		Status:             "PENDING_BUYER_APPROVAL",
		PaymentStatus:      "Pending Escrow",
		BuyerApproved:      false,
		SupplierApproved:   true,
		ValidTillDate:      &validDate,
		CreatedAt:          now.Add(-2 * 24 * time.Hour),
		UpdatedAt:          now.Add(-2 * 24 * time.Hour),
	}
	p3 := domain.Proposal{
		ID:                 "prop_1001",
		ProposalNumber:     "4690112",
		OrderID:            "ORD-55410",
		BuyerID:            uBuyer1.ID,
		BuyerName:          uBuyer1.BusinessName,
		BuyerSignatory:     "Sanjay Kumar",
		BuyerEmail:         "buyer1@acme.com",
		BuyerGSTIN:         "27AAAAA1111A1ZA",
		BuyerAddress:       "Andheri East, Mumbai, Maharashtra 400069",
		SupplierID:         uSupplier1.ID,
		SupplierName:       uSupplier1.BusinessName,
		SupplierSignatory:   "Kiran Patel",
		SupplierEmail:       "patel@rajesh.com",
		SupplierGSTIN:       "03CCCCC3333C3ZC",
		SupplierAddress:     "Industrial Area B, Focal Point, Ludhiana, Punjab 141010",
		ItemDescription:    "Industrial Ceramic Tiles & Refractory Castables (Grade AZS-33)",
		BaseAmount:         175000.00,
		DiscountPercent:    0.00,
		DiscountAmount:     0.00,
		TaxPercent:         18.00,
		TaxAmount:          31500.00,
		TotalPayableAmount: 206500.00,
		Amount:             206500.00,
		Currency:           "INR",
		Terms:              "50% advance in escrow, 50% upon receipt of QC report",
		MilestonesSummary:  "50% Advance • 50% Final Handover",
		DeliveryTimeline:   "2026-06-25",
		Notes:              "Funds released upon verified physical inspection.",
		Status:             "APPROVED",
		PaymentStatus:      "Payment Held in Escrow",
		BuyerApproved:      true,
		SupplierApproved:   true,
		ValidTillDate:      &validDate,
		CreatedAt:          now.Add(-5 * 24 * time.Hour),
		UpdatedAt:          now.Add(-4 * 24 * time.Hour),
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
		DeliverableProofURL: "https://docs.payshieldx.in/proofs/raw_material_cert_089.pdf",
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
		DeliverableProofURL: "https://docs.payshieldx.in/proofs/eway_bill_lr_089.pdf",
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

	// 7. Seed Buyers
	b1 := domain.BuyerRecord{
		BuyerID:               "BYR-2026-001",
		Name:                  "Vikram Malhotra",
		CompanyName:           "Apex Auto Components Pvt Ltd",
		Mobile:                "+91 9820123456",
		Email:                 "procurement@apexauto.in",
		GSTIN:                 "27AAACA1234A1Z5",
		PAN:                   "AAACA1234A",
		KYCStatus:             "Approved",
		CompletedTransactions: 42,
		Disputes:              1,
		TotalTransactionValue: 48500000.00,
		AccountStatus:         "Active",
		CreatedAt:             now.AddDate(0, -6, -10),
		RefundHistory: []domain.RefundRecord{
			{
				RefundID:      "RF-8921",
				TransactionID: "TX-7712",
				Amount:        120000.00,
				Reason:        "Minor specification mismatch on batch 4 - Mutual credit adjustment",
				Status:        "Credited",
				UTRNumber:     "ICICR520260911002341",
				CreatedAt:     now.AddDate(0, -1, -5),
			},
		},
	}
	b2 := domain.BuyerRecord{
		BuyerID:               "BYR-2026-002",
		Name:                  "Amit Rawat",
		CompanyName:           "Rawat Handlooms & Textiles Ltd",
		Mobile:                "+91 9876543210",
		Email:                 "procurement@rawathandlooms.in",
		GSTIN:                 "07AAAAA1111A1ZA",
		PAN:                   "AAAAA1111A",
		KYCStatus:             "Approved",
		CompletedTransactions: 28,
		Disputes:              0,
		TotalTransactionValue: 19200000.00,
		AccountStatus:         "Active",
		CreatedAt:             now.AddDate(0, -4, -15),
		RefundHistory:         []domain.RefundRecord{},
	}
	b3 := domain.BuyerRecord{
		BuyerID:               "BYR-2026-003",
		Name:                  "Kavita Deshmukh",
		CompanyName:           "Maharastra Agro Processing Corp",
		Mobile:                "+91 9845112233",
		Email:                 "kavita@mahaagro.gov.in",
		GSTIN:                 "27AAACM9988C1Z4",
		PAN:                   "AAACM9988C",
		KYCStatus:             "Pending",
		CompletedTransactions: 5,
		Disputes:              0,
		TotalTransactionValue: 3400000.00,
		AccountStatus:         "Under Review",
		CreatedAt:             now.AddDate(0, 0, -3),
		RefundHistory:         []domain.RefundRecord{},
	}
	b4 := domain.BuyerRecord{
		BuyerID:               "BYR-2026-004",
		Name:                  "Rohan Singhal",
		CompanyName:           "Singhal Steel & Fabrication",
		Mobile:                "+91 9811882233",
		Email:                 "rohan@singhalsteels.com",
		GSTIN:                 "08AAACS4455S1Z1",
		PAN:                   "AAACS4455S",
		KYCStatus:             "Approved",
		CompletedTransactions: 64,
		Disputes:              2,
		TotalTransactionValue: 72000000.00,
		AccountStatus:         "Active",
		CreatedAt:             now.AddDate(0, -8, 0),
		RefundHistory: []domain.RefundRecord{
			{
				RefundID:      "RF-9941",
				TransactionID: "TX-4401",
				Amount:        300000.00,
				Reason:        "Transit damage on structural angle bars - Arbitration Refund",
				Status:        "Credited",
				UTRNumber:     "HDFCR520260815009124",
				CreatedAt:     now.AddDate(0, -2, -12),
			},
		},
	}
	s.Buyers[b1.BuyerID] = b1
	s.Buyers[b2.BuyerID] = b2
	s.Buyers[b3.BuyerID] = b3
	s.Buyers[b4.BuyerID] = b4

	// 8. Seed Suppliers
	supVerifiedAt := now.AddDate(0, -3, 0)
	sp1 := domain.SupplierRecord{
		SupplierID:    "SUP-2026-001",
		CompanyName:   "Bharat Precision Castings Ltd",
		ContactPerson: "Rajesh Singhania",
		Mobile:        "+91 9898123456",
		Email:         "sales@bharatcastings.com",
		GSTIN:         "24AABCB5678B1Z2",
		PAN:           "AABCB5678B",
		VerificationStatus: "Approved",
		KYCDocuments: []domain.KYCDocument{
			{ Type: "GST Registration Certificate", DocumentNo: "GST-24AABCB5678B1Z2", URL: "https://docs.payshieldx.in/kyc/gst_bharat.pdf", Status: "Approved", UploadedAt: now.AddDate(0, -3, 0), VerifiedAt: &supVerifiedAt },
			{ Type: "Company PAN Card", DocumentNo: "PAN-AABCB5678B", URL: "https://docs.payshieldx.in/kyc/pan_bharat.pdf", Status: "Approved", UploadedAt: now.AddDate(0, -3, 0), VerifiedAt: &supVerifiedAt },
			{ Type: "MSME Udyam Certificate", DocumentNo: "UDYAM-GJ-01-008912", URL: "https://docs.payshieldx.in/kyc/msme_bharat.pdf", Status: "Approved", UploadedAt: now.AddDate(0, -3, 0), VerifiedAt: &supVerifiedAt },
			{ Type: "Bank Cancelled Cheque", DocumentNo: "HDFC-000405012938", URL: "https://docs.payshieldx.in/kyc/cheque_bharat.pdf", Status: "Approved", UploadedAt: now.AddDate(0, -3, 0), VerifiedAt: &supVerifiedAt },
		},
		SupplierPlan:          "Business Plan",
		PlanExpiry:            now.AddDate(0, 9, 15),
		TotalProposals:        58,
		AcceptedProposals:     52,
		PendingProposals:      4,
		RejectedProposals:     2,
		CompletedTransactions: 49,
		Disputes:              1,
		TotalTransactionValue: 62400000.00,
		Refunds:               0,
		SettlementInfo: domain.SettlementDetails{
			BankName:        "HDFC Bank Corporate",
			AccountNumber:   "000405012938",
			IFSCCode:        "HDFC0000004",
			AccountHolder:   "Bharat Precision Castings Ltd",
			IsPennyDropped:  true,
			PennyDropStatus: "Verified (Penny Drop Active)",
			SettlementCycle: "T+1 Next Day",
		},
		AccountStatus: "Active",
		CreatedAt:     now.AddDate(0, -9, 0),
	}
	sp2 := domain.SupplierRecord{
		SupplierID:    "SUP-2026-002",
		CompanyName:   "S.S. Enterprises",
		ContactPerson: "Sunil Sharma",
		Mobile:        "+91 9822334455",
		Email:         "sunil@ssenterprises.co.in",
		GSTIN:         "20KBIPS8898M1ZG",
		PAN:           "KBIPS8898M",
		VerificationStatus: "Pending",
		KYCDocuments: []domain.KYCDocument{
			{ Type: "GST Certificate", DocumentNo: "20KBIPS8898M1ZG", URL: "https://docs.payshieldx.in/kyc/gst_ss.pdf", Status: "Pending", UploadedAt: now.AddDate(0, 0, -2) },
			{ Type: "Company PAN Card", DocumentNo: "KBIPS8898M", URL: "https://docs.payshieldx.in/kyc/pan_ss.pdf", Status: "Pending", UploadedAt: now.AddDate(0, 0, -2) },
			{ Type: "Bank Passbook Copy", DocumentNo: "ICIC-5020008819", URL: "https://docs.payshieldx.in/kyc/bank_ss.pdf", Status: "Pending", UploadedAt: now.AddDate(0, 0, -2) },
		},
		SupplierPlan:          "Growth Plan",
		PlanExpiry:            now.AddDate(0, 1, 10),
		TotalProposals:        12,
		AcceptedProposals:     9,
		PendingProposals:      2,
		RejectedProposals:     1,
		CompletedTransactions: 8,
		Disputes:              0,
		TotalTransactionValue: 9800000.00,
		Refunds:               0,
		SettlementInfo: domain.SettlementDetails{
			BankName:        "ICICI Bank Commercial",
			AccountNumber:   "50200088192410",
			IFSCCode:        "ICIC0000104",
			AccountHolder:   "S.S. Enterprises",
			IsPennyDropped:  false,
			PennyDropStatus: "Pending Penny Drop",
			SettlementCycle: "T+1 Next Day",
		},
		AccountStatus: "Under Review",
		CreatedAt:     now.AddDate(0, 0, -5),
	}
	sp3 := domain.SupplierRecord{
		SupplierID:    "SUP-2026-003",
		CompanyName:   "Kalyani Industrial Forgings",
		ContactPerson: "Anand Kalyani",
		Mobile:        "+91 9844001122",
		Email:         "anand@kalyaniforgings.in",
		GSTIN:         "29AAACK5544K1Z8",
		PAN:           "AAACK5544K",
		VerificationStatus: "Approved",
		KYCDocuments: []domain.KYCDocument{
			{ Type: "GST Certificate", DocumentNo: "29AAACK5544K1Z8", URL: "https://docs.payshieldx.in/kyc/gst_kalyani.pdf", Status: "Approved", UploadedAt: now.AddDate(0, -3, 0), VerifiedAt: &supVerifiedAt },
			{ Type: "Company PAN Card", DocumentNo: "AAACK5544K", URL: "https://docs.payshieldx.in/kyc/pan_kalyani.pdf", Status: "Approved", UploadedAt: now.AddDate(0, -3, 0), VerifiedAt: &supVerifiedAt },
			{ Type: "Bank Cancelled Cheque", DocumentNo: "SBI-30910049281", URL: "https://docs.payshieldx.in/kyc/cheque_kalyani.pdf", Status: "Approved", UploadedAt: now.AddDate(0, -3, 0), VerifiedAt: &supVerifiedAt },
		},
		SupplierPlan:          "Enterprise Plan",
		PlanExpiry:            now.AddDate(1, 2, 0),
		TotalProposals:        114,
		AcceptedProposals:     108,
		PendingProposals:      5,
		RejectedProposals:     1,
		CompletedTransactions: 102,
		Disputes:              3,
		TotalTransactionValue: 142000000.00,
		Refunds:               1,
		SettlementInfo: domain.SettlementDetails{
			BankName:        "State Bank of India Corporate",
			AccountNumber:   "309100492810",
			IFSCCode:        "SBIN0000301",
			AccountHolder:   "Kalyani Industrial Forgings",
			IsPennyDropped:  true,
			PennyDropStatus: "Verified (Penny Drop Active)",
			SettlementCycle: "T+0 Instant",
		},
		AccountStatus: "Active",
		CreatedAt:     now.AddDate(0, -14, 0),
	}
	sp4 := domain.SupplierRecord{
		SupplierID:    "SUP-2026-004",
		CompanyName:   "Vanguard Electricals & Switchgear",
		ContactPerson: "Deepak Mehta",
		Mobile:        "+91 9820998877",
		Email:         "d.mehta@vanguardelec.com",
		GSTIN:         "27AAACV9911V1Z3",
		PAN:           "AAACV9911V",
		VerificationStatus: "Approved",
		KYCDocuments: []domain.KYCDocument{
			{ Type: "GST Registration Certificate", DocumentNo: "27AAACV9911V1Z3", URL: "https://docs.payshieldx.in/kyc/gst_vanguard.pdf", Status: "Approved", UploadedAt: now.AddDate(0, -5, 0), VerifiedAt: &supVerifiedAt },
			{ Type: "Company PAN Card", DocumentNo: "AAACV9911V", URL: "https://docs.payshieldx.in/kyc/pan_vanguard.pdf", Status: "Approved", UploadedAt: now.AddDate(0, -5, 0), VerifiedAt: &supVerifiedAt },
		},
		SupplierPlan:          "Business Plan",
		PlanExpiry:            now.AddDate(0, 6, 20),
		TotalProposals:        36,
		AcceptedProposals:     32,
		PendingProposals:      3,
		RejectedProposals:     1,
		CompletedTransactions: 30,
		Disputes:              0,
		TotalTransactionValue: 38500000.00,
		Refunds:               0,
		SettlementInfo: domain.SettlementDetails{
			BankName:        "Kotak Mahindra Bank",
			AccountNumber:   "0029104000129",
			IFSCCode:        "KKBK0000182",
			AccountHolder:   "Vanguard Electricals & Switchgear",
			IsPennyDropped:  true,
			PennyDropStatus: "Verified (Penny Drop Active)",
			SettlementCycle: "T+0 Instant",
		},
		AccountStatus: "Active",
		CreatedAt:     now.AddDate(0, -7, 0),
	}
	s.Suppliers[sp1.SupplierID] = sp1
	s.Suppliers[sp2.SupplierID] = sp2
	s.Suppliers[sp3.SupplierID] = sp3
	s.Suppliers[sp4.SupplierID] = sp4

	// 9. Seed Settlements Queue
	st1 := domain.SettlementItem{
		SettlementID: "SETTL-2026-901",
		SupplierID:   "SUP-2026-001",
		SupplierName: "Bharat Precision Castings Ltd",
		DealRef:      "TS-CTR-2026-089 (Milestone 2 QC)",
		BankName:     "HDFC Bank Corporate (A/C ...2938)",
		AccountNo:    "000405012938",
		IFSCCode:     "HDFC0000004",
		Amount:       600000.00,
		Status:       "PENDING",
		CreatedAt:    now.AddDate(0, 0, -1),
	}
	st2 := domain.SettlementItem{
		SettlementID: "SETTL-2026-902",
		SupplierID:   "SUP-2026-002",
		SupplierName: "S.S. Enterprises",
		DealRef:      "TS-CTR-2026-104 (Milestone 1 Dispatch)",
		BankName:     "ICICI Bank Commercial (A/C ...2410)",
		AccountNo:    "50200088192410",
		IFSCCode:     "ICIC0000104",
		Amount:       450000.00,
		Status:       "PENDING",
		CreatedAt:    now.AddDate(0, 0, -2),
	}
	st3 := domain.SettlementItem{
		SettlementID: "SETTL-2026-903",
		SupplierID:   "SUP-2026-004",
		SupplierName: "Vanguard Electricals & Switchgear",
		DealRef:      "TS-CTR-2026-042 (Final Acceptance)",
		BankName:     "Kotak Mahindra Bank (A/C ...0129)",
		AccountNo:    "0029104000129",
		IFSCCode:     "KKBK0000182",
		Amount:       800000.00,
		Status:       "PENDING",
		CreatedAt:    now.AddDate(0, 0, -1),
	}
	s.Settlements[st1.SettlementID] = st1
	s.Settlements[st2.SettlementID] = st2
	s.Settlements[st3.SettlementID] = st3
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
