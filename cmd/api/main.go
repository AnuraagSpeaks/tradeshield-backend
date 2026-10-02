package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"tradeshield-backend/internal/handler"
	"tradeshield-backend/internal/pkg/mailer"
	"tradeshield-backend/internal/repository"
)

func main() {
	mailer.InitMailer()
	
	// 1. Seed demo accounts (Apex Buyer, Bharat Supplier, PayShield Court Admin, etc.)
	repository.DB.SeedData()

	// 2. Load and overlay any existing users and proposals from PostgreSQL
	pg := repository.InitPostgres()
	if pg != nil {
		pg.LoadAllIntoStore(repository.DB)
		log.Println("📥 Syncing all accounts, proposals, and contracts into PostgreSQL...")
		for _, u := range repository.DB.Users {
			_ = pg.SaveUser(u)
		}
		for _, pr := range repository.DB.Proposals {
			_ = pg.SaveProposal(pr)
		}
		for _, c := range repository.DB.Contracts {
			_ = pg.SaveContract(c)
		}
	}

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Idempotency-Key"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","message":"PayShieldX Trade Protection Backend API is online","health_check":"/api/v1/health"}`))
	})
	r.Head("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	h := handler.NewHandler(repository.DB)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok","service":"PayShieldX Trade Protection API","version":"2.4.2-verified-production"}`))
		})

		// 1. Auth & Directory
		r.Post("/auth/login", h.Login)
		r.Post("/auth/register", h.Register)
		r.Post("/auth/send-otp", h.SendOTP)
		r.Post("/auth/verify-otp", h.VerifyOTP)
		r.Post("/auth/reset-password", h.ResetPassword)
		r.Get("/directory", h.ListDirectory)

		// 2. Escrow Metrics
		r.Get("/escrow/summary", h.GetSummary)

		// 3. Proposals & Dual-Approval System
		r.Get("/proposals", h.ListProposals)
		r.Post("/proposals", h.CreateProposal)
		r.Get("/proposals/{id}", h.GetProposal)
		r.Post("/proposals/{id}/approve", h.ApproveProposal)
		r.Post("/proposals/{id}/request-modification", h.RequestModification)
		r.Post("/proposals/{id}/ship", h.ShipProposal)
		r.Post("/proposals/{id}/confirm-delivery", h.ConfirmDeliveryProposal)

		// 4. Milestone Escrow Contracts (High-Volume)
		r.Get("/contracts", h.ListContracts)
		r.Post("/contracts", h.CreateContract)
		r.Get("/contracts/{id}", h.GetContract)
		r.Post("/milestones/{id}/submit-proof", h.SubmitMilestoneProof)
		r.Post("/milestones/{id}/approve-release", h.ApproveAndReleaseMilestone)

		// 5. Connections & Partner Requests
		r.Get("/connections", h.ListConnections)
		r.Get("/connections/requests", h.ListConnectionRequests)
		r.Post("/connections/requests", h.SendConnectionRequest)
		r.Post("/connections/requests/{id}/accept", h.AcceptConnectionRequest)

		// 6. Disputes & Arbitration Court
		r.Get("/disputes", h.ListDisputes)
		r.Post("/disputes", h.RaiseDispute)
		r.Post("/disputes/{id}/arbitrate", h.ArbitrateDispute)

		// 7. Double-Entry Ledger
		r.Get("/ledger/accounts", h.ListLedgerAccounts)
		r.Get("/ledger/journals", h.ListLedgerJournals)

		// 8. Support Ticket Desk
		r.Get("/support/tickets", h.ListSupportTickets)
		r.Post("/support/tickets", h.SubmitSupportTicket)

		// 9. Live Verification APIs
		r.Post("/mock/penny-drop", h.MockPennyDrop)
		r.Get("/mock/gst-verify/{gstin}", h.MockGSTVerify)

		// 10. Admin Executive & Owner Management APIs
		r.Get("/admin/stats", h.GetAdminStats)
		r.Get("/admin/buyers", h.ListAdminBuyers)
		r.Get("/admin/suppliers", h.ListAdminSuppliers)
		r.Get("/admin/finance", h.GetAdminFinance)
		r.Post("/admin/kyc/{userId}/verify", h.VerifyAdminKYC)
		r.Post("/admin/settlements/{id}/payout", h.ProcessAdminSettlement)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	fmt.Printf("🛡️  PayShieldX Go Backend running on port %s\n", port)
	fmt.Printf("👉 Health check: http://localhost%s/api/v1/health\n", port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
