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
	"tradeshield-backend/internal/middleware"
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
		log.Println("📥 Syncing all accounts, proposals, contracts, buyers, suppliers, and settlements into PostgreSQL...")
		for _, u := range repository.DB.Users {
			_ = pg.SaveUser(u)
		}
		for _, pr := range repository.DB.Proposals {
			_ = pg.SaveProposal(pr)
		}
		for _, c := range repository.DB.Contracts {
			_ = pg.SaveContract(c)
		}
		for _, b := range repository.DB.Buyers {
			_ = pg.SaveBuyer(b)
		}
		for _, sp := range repository.DB.Suppliers {
			_ = pg.SaveSupplier(sp)
		}
		for _, st := range repository.DB.Settlements {
			_ = pg.SaveSettlement(st)
		}
	}

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://payshieldx.in",
			"https://www.payshieldx.in",
			"https://tradeshield-frontend.vercel.app",
			"http://localhost:5173",
			"http://localhost:3000",
			"http://127.0.0.1:5173",
		},
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

	// Top-level aliases for direct compliance with integration specs
	r.Post("/api/create-order", h.CreateOrder)
	r.Post("/api/verify-payment", h.VerifyPayment)
	r.Get("/documents/*", h.GetDocumentPreview)

	r.Route("/api/v1", func(r chi.Router) {
		// Public Endpoints (no auth required)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"status":"ok","service":"PayShieldX Trade Protection API","version":"2.4.2-verified-production"}`))
		})

		// 1. Auth & Public Directory
		r.Post("/auth/login", h.Login)
		r.Post("/auth/register", h.Register)
		r.Post("/auth/send-otp", h.SendOTP)
		r.Post("/auth/verify-otp", h.VerifyOTP)
		r.Post("/auth/reset-password", h.ResetPassword)
		r.Get("/directory", h.ListDirectory)

		// 2. Escrow Metrics & Directory Lookups
		r.Get("/escrow/summary", h.GetSummary)
		r.Get("/suppliers", h.ListAdminSuppliers)

		// 3. Proposal Public View & PDF Render
		r.Get("/proposals/{id}", h.GetProposal)
		r.Get("/proposals/{id}/pdf", h.GetProposalPDFHTML)

		// 4. Integrations & Documents
		r.Post("/mock/penny-drop", h.MockPennyDrop)
		r.Get("/mock/gst-verify/{gstin}", h.MockGSTVerify)
		r.Post("/create-order", h.CreateOrder)
		r.Post("/verify-payment", h.VerifyPayment)
		r.Post("/gst/lookup", h.GSTLookup)
		r.Get("/documents/*", h.GetDocumentPreview)

		// Authenticated Routes (Requires valid JWT Bearer Token)
		r.Group(func(auth chi.Router) {
			auth.Use(middleware.AuthMiddleware)

			// Proposals Management
			auth.Get("/proposals", h.ListProposals)
			auth.Post("/proposals", h.CreateProposal)
			auth.Post("/proposals/{id}/approve", h.ApproveProposal)
			auth.Post("/proposals/{id}/request-modification", h.RequestModification)
			auth.Post("/proposals/{id}/ship", h.ShipProposal)
			auth.Post("/proposals/{id}/confirm-delivery", h.ConfirmDeliveryProposal)

			// Milestone Escrow Contracts
			auth.Get("/contracts", h.ListContracts)
			auth.Post("/contracts", h.CreateContract)
			auth.Get("/contracts/{id}", h.GetContract)
			auth.Post("/milestones/{id}/submit-proof", h.SubmitMilestoneProof)
			auth.Post("/milestones/{id}/approve-release", h.ApproveAndReleaseMilestone)

			// Connections & Networking
			auth.Get("/connections", h.ListConnections)
			auth.Get("/connections/requests", h.ListConnectionRequests)
			auth.Post("/connections/requests", h.SendConnectionRequest)
			auth.Post("/connections/requests/{id}/accept", h.AcceptConnectionRequest)

			// Disputes (User raising & listing)
			auth.Get("/disputes", h.ListDisputes)
			auth.Post("/disputes", h.RaiseDispute)

			// Support Tickets
			auth.Get("/support/tickets", h.ListSupportTickets)
			auth.Post("/support/tickets", h.SubmitSupportTicket)

			// Admin Restricted Routes (Requires role: "admin")
			auth.Group(func(admin chi.Router) {
				admin.Use(middleware.RequireRole("admin"))

				admin.Get("/admin/stats", h.GetAdminStats)
				admin.Get("/admin/buyers", h.ListAdminBuyers)
				admin.Get("/admin/suppliers", h.ListAdminSuppliers)
				admin.Get("/admin/finance", h.GetAdminFinance)
				admin.Post("/admin/kyc/{userId}/verify", h.VerifyAdminKYC)
				admin.Post("/admin/settlements/{id}/payout", h.ProcessAdminSettlement)

				admin.Get("/ledger/accounts", h.ListLedgerAccounts)
				admin.Get("/ledger/journals", h.ListLedgerJournals)

				admin.Post("/disputes/{id}/arbitrate", h.ArbitrateDispute)
			})
		})
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
