package repository

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"tradeshield-backend/internal/domain"
)

type PostgresStore struct {
	db *sql.DB
}

var PG *PostgresStore

func InitPostgres() *PostgresStore {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Println("ℹ️  No DATABASE_URL provided. Running with In-Memory Repository.")
		return nil
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("⚠️  Failed to connect to PostgreSQL: %v. Falling back to memory.", err)
		return nil
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Printf("⚠️  PostgreSQL ping failed: %v. Falling back to memory.", err)
		return nil
	}

	log.Println("✅ Connected to PostgreSQL Database successfully!")

	store := &PostgresStore{db: db}
	if err := store.migrate(); err != nil {
		log.Printf("⚠️  Migration error: %v", err)
	}

	PG = store
	return store
}

func (p *PostgresStore) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(100) PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255),
		full_name VARCHAR(255),
		business_name VARCHAR(255),
		gst VARCHAR(50),
		pan VARCHAR(50),
		mobile VARCHAR(50),
		city VARCHAR(100),
		role VARCHAR(50),
		verified BOOLEAN DEFAULT true,
		data JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS proposals (
		id VARCHAR(100) PRIMARY KEY,
		order_id VARCHAR(100),
		buyer_id VARCHAR(100),
		buyer_name VARCHAR(255),
		supplier_id VARCHAR(100),
		supplier_name VARCHAR(255),
		amount NUMERIC(15,2),
		currency VARCHAR(10) DEFAULT 'INR',
		terms VARCHAR(50),
		delivery_timeline VARCHAR(100),
		notes TEXT,
		status VARCHAR(50),
		payment_status VARCHAR(50),
		buyer_approved BOOLEAN DEFAULT false,
		supplier_approved BOOLEAN DEFAULT false,
		lr_number VARCHAR(100),
		transporter_name VARCHAR(255),
		proof_url TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS contracts (
		id VARCHAR(100) PRIMARY KEY,
		title VARCHAR(255),
		contract_number VARCHAR(100),
		buyer_org_id VARCHAR(100),
		buyer_org_name VARCHAR(255),
		supplier_org_id VARCHAR(100),
		supplier_org_name VARCHAR(255),
		total_amount NUMERIC(15,2),
		currency VARCHAR(10) DEFAULT 'INR',
		platform_fee_percent NUMERIC(5,2),
		platform_fee_amount NUMERIC(15,2),
		escrow_virtual_account VARCHAR(100),
		status VARCHAR(50),
		description TEXT,
		delivery_terms VARCHAR(255),
		inspection_period_days INT,
		milestones JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS disputes (
		id VARCHAR(100) PRIMARY KEY,
		proposal_id VARCHAR(100),
		contract_id VARCHAR(100),
		milestone_id VARCHAR(100),
		initiator_org_id VARCHAR(100),
		respondent_org_id VARCHAR(100),
		reason TEXT,
		claim_amount NUMERIC(15,2),
		status VARCHAR(50),
		resolution_verdict VARCHAR(50),
		buyer_refund_share NUMERIC(15,2),
		seller_release_share NUMERIC(15,2),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		resolved_at TIMESTAMP WITH TIME ZONE
	);

	CREATE TABLE IF NOT EXISTS journal_entries (
		id VARCHAR(100) PRIMARY KEY,
		contract_id VARCHAR(100),
		milestone_id VARCHAR(100),
		reference VARCHAR(100),
		description TEXT,
		postings JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS support_tickets (
		id VARCHAR(100) PRIMARY KEY,
		name VARCHAR(255),
		company VARCHAR(255),
		email VARCHAR(255),
		phone VARCHAR(50),
		message TEXT,
		status VARCHAR(50) DEFAULT 'open',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS buyers (
		buyer_id VARCHAR(100) PRIMARY KEY,
		name VARCHAR(255),
		company_name VARCHAR(255),
		mobile VARCHAR(50),
		email VARCHAR(255),
		gstin VARCHAR(50),
		pan VARCHAR(50),
		kyc_status VARCHAR(50),
		completed_transactions INT DEFAULT 0,
		disputes INT DEFAULT 0,
		total_transaction_value NUMERIC(15,2) DEFAULT 0,
		account_status VARCHAR(50) DEFAULT 'Active',
		refund_history JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS suppliers (
		supplier_id VARCHAR(100) PRIMARY KEY,
		company_name VARCHAR(255),
		contact_person VARCHAR(255),
		mobile VARCHAR(50),
		email VARCHAR(255),
		gstin VARCHAR(50),
		pan VARCHAR(50),
		kyc_status VARCHAR(50),
		kyc_documents JSONB,
		plan_tier VARCHAR(50),
		plan_expiry TIMESTAMP WITH TIME ZONE,
		total_proposals INT DEFAULT 0,
		accepted_proposals INT DEFAULT 0,
		pending_proposals INT DEFAULT 0,
		rejected_proposals INT DEFAULT 0,
		completed_transactions INT DEFAULT 0,
		disputes INT DEFAULT 0,
		total_transaction_value NUMERIC(15,2) DEFAULT 0,
		refunds INT DEFAULT 0,
		settlement_details JSONB,
		account_status VARCHAR(50) DEFAULT 'Active',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS settlements (
		settlement_id VARCHAR(100) PRIMARY KEY,
		supplier_id VARCHAR(100),
		supplier_name VARCHAR(255),
		deal_ref VARCHAR(255),
		bank_name VARCHAR(255),
		account_number VARCHAR(100),
		ifsc VARCHAR(50),
		amount NUMERIC(15,2),
		status VARCHAR(50),
		utr_number VARCHAR(100),
		disbursed_at TIMESTAMP WITH TIME ZONE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`
	_, err := p.db.Exec(schema)
	if err == nil {
		log.Println("✅ PostgreSQL Schema & Tables Verified (Users, Proposals, Contracts, Disputes, Buyers, Suppliers, Settlements)!")
	}
	return err
}

func (p *PostgresStore) SaveUser(u domain.User) error {
	if p == nil || p.db == nil {
		return nil
	}
	dataBytes, _ := json.Marshal(u)
	query := `
		INSERT INTO users (id, email, password, full_name, business_name, gst, mobile, city, role, verified, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			business_name = EXCLUDED.business_name,
			gst = EXCLUDED.gst,
			mobile = EXCLUDED.mobile,
			city = EXCLUDED.city,
			role = EXCLUDED.role,
			data = EXCLUDED.data,
			updated_at = EXCLUDED.updated_at;
	`
	_, err := p.db.Exec(query, u.ID, u.Email, u.Password, u.FullName, u.BusinessName, u.GST, u.Mobile, u.City, u.Role, u.Verified, string(dataBytes), u.CreatedAt, u.UpdatedAt)
	return err
}

func (p *PostgresStore) SaveProposal(pr domain.Proposal) error {
	if p == nil || p.db == nil {
		return nil
	}
	query := `
		INSERT INTO proposals (id, order_id, buyer_id, buyer_name, supplier_id, supplier_name, amount, currency, terms, delivery_timeline, notes, status, payment_status, buyer_approved, supplier_approved, lr_number, transporter_name, proof_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			payment_status = EXCLUDED.payment_status,
			buyer_approved = EXCLUDED.buyer_approved,
			supplier_approved = EXCLUDED.supplier_approved,
			lr_number = EXCLUDED.lr_number,
			transporter_name = EXCLUDED.transporter_name,
			proof_url = EXCLUDED.proof_url,
			updated_at = EXCLUDED.updated_at;
	`
	_, err := p.db.Exec(query, pr.ID, pr.OrderID, pr.BuyerID, pr.BuyerName, pr.SupplierID, pr.SupplierName, pr.Amount, pr.Currency, pr.Terms, pr.DeliveryTimeline, pr.Notes, pr.Status, pr.PaymentStatus, pr.BuyerApproved, pr.SupplierApproved, pr.LRNumber, pr.TransporterName, pr.ProofURL, pr.CreatedAt, pr.UpdatedAt)
	return err
}

func (p *PostgresStore) SaveContract(c domain.Contract) error {
	if p == nil || p.db == nil {
		return nil
	}
	msBytes, _ := json.Marshal(c.Milestones)
	query := `
		INSERT INTO contracts (id, title, contract_number, buyer_org_id, buyer_org_name, supplier_org_id, supplier_org_name, total_amount, currency, platform_fee_percent, platform_fee_amount, escrow_virtual_account, status, description, delivery_terms, inspection_period_days, milestones, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			milestones = EXCLUDED.milestones,
			updated_at = EXCLUDED.updated_at;
	`
	_, err := p.db.Exec(query, c.ID, c.Title, c.ContractNumber, c.BuyerOrgID, c.BuyerOrgName, c.SupplierOrgID, c.SupplierOrgName, c.TotalAmount, c.Currency, c.PlatformFeePercent, c.PlatformFeeAmount, c.EscrowVirtualAccount, c.Status, c.Description, c.DeliveryTerms, c.InspectionPeriodDays, string(msBytes), c.CreatedAt, c.UpdatedAt)
	return err
}

func (p *PostgresStore) SaveDispute(d domain.Dispute) error {
	if p == nil || p.db == nil {
		return nil
	}
	query := `
		INSERT INTO disputes (id, proposal_id, contract_id, milestone_id, initiator_org_id, respondent_org_id, reason, claim_amount, status, resolution_verdict, buyer_refund_share, seller_release_share, created_at, resolved_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			resolution_verdict = EXCLUDED.resolution_verdict,
			buyer_refund_share = EXCLUDED.buyer_refund_share,
			seller_release_share = EXCLUDED.seller_release_share,
			resolved_at = EXCLUDED.resolved_at;
	`
	_, err := p.db.Exec(query, d.ID, d.ProposalID, d.ContractID, d.MilestoneID, d.InitiatorOrgID, d.RespondentOrgID, d.Reason, d.ClaimAmount, d.Status, d.ResolutionVerdict, d.BuyerRefundShare, d.SellerReleaseShare, d.CreatedAt, d.ResolvedAt)
	return err
}

func (p *PostgresStore) SaveJournal(j domain.JournalEntry) error {
	if p == nil || p.db == nil {
		return nil
	}
	postingsBytes, _ := json.Marshal(j.Postings)
	query := `
		INSERT INTO journal_entries (id, contract_id, milestone_id, reference, description, postings, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING;
	`
	_, err := p.db.Exec(query, j.ID, j.ContractID, j.MilestoneID, j.Reference, j.Description, string(postingsBytes), j.CreatedAt)
	return err
}

func (p *PostgresStore) SaveSupportTicket(t domain.SupportTicket) error {
	if p == nil || p.db == nil {
		return nil
	}
	query := `
		INSERT INTO support_tickets (id, name, company, email, phone, message, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status;
	`
	_, err := p.db.Exec(query, t.ID, t.Name, t.Company, t.Email, t.Phone, t.Message, t.Status, t.CreatedAt)
	return err
}

func (p *PostgresStore) SaveBuyer(b domain.BuyerRecord) error {
	if p == nil || p.db == nil {
		return nil
	}
	refundsBytes, _ := json.Marshal(b.RefundHistory)
	query := `
		INSERT INTO buyers (buyer_id, name, company_name, mobile, email, gstin, pan, kyc_status, completed_transactions, disputes, total_transaction_value, account_status, refund_history, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (buyer_id) DO UPDATE SET
			name = EXCLUDED.name,
			company_name = EXCLUDED.company_name,
			mobile = EXCLUDED.mobile,
			email = EXCLUDED.email,
			gstin = EXCLUDED.gstin,
			pan = EXCLUDED.pan,
			kyc_status = EXCLUDED.kyc_status,
			completed_transactions = EXCLUDED.completed_transactions,
			disputes = EXCLUDED.disputes,
			total_transaction_value = EXCLUDED.total_transaction_value,
			account_status = EXCLUDED.account_status,
			refund_history = EXCLUDED.refund_history;
	`
	_, err := p.db.Exec(query, b.BuyerID, b.Name, b.CompanyName, b.Mobile, b.Email, b.GSTIN, b.PAN, b.KYCStatus, b.CompletedTransactions, b.Disputes, b.TotalTransactionValue, b.AccountStatus, string(refundsBytes), b.CreatedAt)
	return err
}

func (p *PostgresStore) SaveSupplier(s domain.SupplierRecord) error {
	if p == nil || p.db == nil {
		return nil
	}
	docsBytes, _ := json.Marshal(s.KYCDocuments)
	settleBytes, _ := json.Marshal(s.SettlementInfo)
	query := `
		INSERT INTO suppliers (supplier_id, company_name, contact_person, mobile, email, gstin, pan, kyc_status, kyc_documents, plan_tier, plan_expiry, total_proposals, accepted_proposals, pending_proposals, rejected_proposals, completed_transactions, disputes, total_transaction_value, refunds, settlement_details, account_status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
		ON CONFLICT (supplier_id) DO UPDATE SET
			company_name = EXCLUDED.company_name,
			contact_person = EXCLUDED.contact_person,
			mobile = EXCLUDED.mobile,
			email = EXCLUDED.email,
			gstin = EXCLUDED.gstin,
			pan = EXCLUDED.pan,
			kyc_status = EXCLUDED.kyc_status,
			kyc_documents = EXCLUDED.kyc_documents,
			plan_tier = EXCLUDED.plan_tier,
			plan_expiry = EXCLUDED.plan_expiry,
			total_proposals = EXCLUDED.total_proposals,
			accepted_proposals = EXCLUDED.accepted_proposals,
			pending_proposals = EXCLUDED.pending_proposals,
			rejected_proposals = EXCLUDED.rejected_proposals,
			completed_transactions = EXCLUDED.completed_transactions,
			disputes = EXCLUDED.disputes,
			total_transaction_value = EXCLUDED.total_transaction_value,
			refunds = EXCLUDED.refunds,
			settlement_details = EXCLUDED.settlement_details,
			account_status = EXCLUDED.account_status;
	`
	_, err := p.db.Exec(query, s.SupplierID, s.CompanyName, s.ContactPerson, s.Mobile, s.Email, s.GSTIN, s.PAN, s.VerificationStatus, string(docsBytes), s.SupplierPlan, s.PlanExpiry, s.TotalProposals, s.AcceptedProposals, s.PendingProposals, s.RejectedProposals, s.CompletedTransactions, s.Disputes, s.TotalTransactionValue, s.Refunds, string(settleBytes), s.AccountStatus, s.CreatedAt)
	return err
}

func (p *PostgresStore) SaveSettlement(st domain.SettlementItem) error {
	if p == nil || p.db == nil {
		return nil
	}
	query := `
		INSERT INTO settlements (settlement_id, supplier_id, supplier_name, deal_ref, bank_name, account_number, ifsc, amount, status, utr_number, disbursed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (settlement_id) DO UPDATE SET
			status = EXCLUDED.status,
			utr_number = EXCLUDED.utr_number,
			disbursed_at = EXCLUDED.disbursed_at;
	`
	_, err := p.db.Exec(query, st.SettlementID, st.SupplierID, st.SupplierName, st.DealRef, st.BankName, st.AccountNo, st.IFSCCode, st.Amount, st.Status, st.UTRNumber, st.DisbursedAt, st.CreatedAt)
	return err
}

func (p *PostgresStore) LoadAllIntoStore(s *Store) error {
	if p == nil || p.db == nil {
		return nil
	}

	// 1. Load Users
	rows, err := p.db.Query("SELECT id, email, password, full_name, business_name, gst, mobile, city, role, verified, created_at, updated_at FROM users")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var u domain.User
			if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.FullName, &u.BusinessName, &u.GST, &u.Mobile, &u.City, &u.Role, &u.Verified, &u.CreatedAt, &u.UpdatedAt); err == nil {
				u.IsVerified = u.Verified
				s.Users[u.ID] = u
			}
		}
	}

	// 2. Load Proposals
	pRows, err := p.db.Query("SELECT id, order_id, buyer_id, buyer_name, supplier_id, supplier_name, amount, currency, terms, delivery_timeline, notes, status, payment_status, buyer_approved, supplier_approved, lr_number, transporter_name, proof_url, created_at, updated_at FROM proposals")
	if err == nil {
		defer pRows.Close()
		for pRows.Next() {
			var pr domain.Proposal
			if err := pRows.Scan(&pr.ID, &pr.OrderID, &pr.BuyerID, &pr.BuyerName, &pr.SupplierID, &pr.SupplierName, &pr.Amount, &pr.Currency, &pr.Terms, &pr.DeliveryTimeline, &pr.Notes, &pr.Status, &pr.PaymentStatus, &pr.BuyerApproved, &pr.SupplierApproved, &pr.LRNumber, &pr.TransporterName, &pr.ProofURL, &pr.CreatedAt, &pr.UpdatedAt); err == nil {
				s.Proposals[pr.ID] = pr
			}
		}
	}

	// 3. Load Contracts
	cRows, err := p.db.Query("SELECT id, title, contract_number, buyer_org_id, buyer_org_name, supplier_org_id, supplier_org_name, total_amount, currency, platform_fee_percent, platform_fee_amount, escrow_virtual_account, status, description, delivery_terms, inspection_period_days, milestones, created_at, updated_at FROM contracts")
	if err == nil {
		defer cRows.Close()
		for cRows.Next() {
			var c domain.Contract
			var msJSON string
			if err := cRows.Scan(&c.ID, &c.Title, &c.ContractNumber, &c.BuyerOrgID, &c.BuyerOrgName, &c.SupplierOrgID, &c.SupplierOrgName, &c.TotalAmount, &c.Currency, &c.PlatformFeePercent, &c.PlatformFeeAmount, &c.EscrowVirtualAccount, &c.Status, &c.Description, &c.DeliveryTerms, &c.InspectionPeriodDays, &msJSON, &c.CreatedAt, &c.UpdatedAt); err == nil {
				if msJSON != "" {
					json.Unmarshal([]byte(msJSON), &c.Milestones)
					for _, m := range c.Milestones {
						s.Milestones[m.ID] = m
					}
				}
				s.Contracts[c.ID] = c
			}
		}
	}

	// 4. Load Buyers
	bRows, err := p.db.Query("SELECT buyer_id, name, company_name, mobile, email, gstin, pan, kyc_status, completed_transactions, disputes, total_transaction_value, account_status, refund_history, created_at FROM buyers")
	if err == nil {
		defer bRows.Close()
		for bRows.Next() {
			var b domain.BuyerRecord
			var refundsJSON string
			if err := bRows.Scan(&b.BuyerID, &b.Name, &b.CompanyName, &b.Mobile, &b.Email, &b.GSTIN, &b.PAN, &b.KYCStatus, &b.CompletedTransactions, &b.Disputes, &b.TotalTransactionValue, &b.AccountStatus, &refundsJSON, &b.CreatedAt); err == nil {
				if refundsJSON != "" {
					json.Unmarshal([]byte(refundsJSON), &b.RefundHistory)
				}
				s.Buyers[b.BuyerID] = b
			}
		}
	}

	// 5. Load Suppliers
	sRows, err := p.db.Query("SELECT supplier_id, company_name, contact_person, mobile, email, gstin, pan, kyc_status, kyc_documents, plan_tier, plan_expiry, total_proposals, accepted_proposals, pending_proposals, rejected_proposals, completed_transactions, disputes, total_transaction_value, refunds, settlement_details, account_status, created_at FROM suppliers")
	if err == nil {
		defer sRows.Close()
		for sRows.Next() {
			var sp domain.SupplierRecord
			var docsJSON, settleJSON string
			if err := sRows.Scan(&sp.SupplierID, &sp.CompanyName, &sp.ContactPerson, &sp.Mobile, &sp.Email, &sp.GSTIN, &sp.PAN, &sp.VerificationStatus, &docsJSON, &sp.SupplierPlan, &sp.PlanExpiry, &sp.TotalProposals, &sp.AcceptedProposals, &sp.PendingProposals, &sp.RejectedProposals, &sp.CompletedTransactions, &sp.Disputes, &sp.TotalTransactionValue, &sp.Refunds, &settleJSON, &sp.AccountStatus, &sp.CreatedAt); err == nil {
				if docsJSON != "" {
					json.Unmarshal([]byte(docsJSON), &sp.KYCDocuments)
				}
				if settleJSON != "" {
					json.Unmarshal([]byte(settleJSON), &sp.SettlementInfo)
				}
				s.Suppliers[sp.SupplierID] = sp
			}
		}
	}

	// 6. Load Settlements
	stRows, err := p.db.Query("SELECT settlement_id, supplier_id, supplier_name, deal_ref, bank_name, account_number, ifsc, amount, status, utr_number, disbursed_at, created_at FROM settlements")
	if err == nil {
		defer stRows.Close()
		for stRows.Next() {
			var st domain.SettlementItem
			if err := stRows.Scan(&st.SettlementID, &st.SupplierID, &st.SupplierName, &st.DealRef, &st.BankName, &st.AccountNo, &st.IFSCCode, &st.Amount, &st.Status, &st.UTRNumber, &st.DisbursedAt, &st.CreatedAt); err == nil {
				s.Settlements[st.SettlementID] = st
			}
		}
	}

	log.Printf("📊 Restored from PostgreSQL: %d users, %d proposals, %d contracts, %d buyers, %d suppliers, %d settlements", len(s.Users), len(s.Proposals), len(s.Contracts), len(s.Buyers), len(s.Suppliers), len(s.Settlements))
	return nil
}
