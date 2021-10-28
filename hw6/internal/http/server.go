package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/erkkke/golang-start/hw6/internal/models"
	"github.com/erkkke/golang-start/hw6/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 30 * time.Second
)

type Server struct {
	ctx               context.Context
	idleConnectionsCh chan struct{}
	store             store.Store
	Address           string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:               ctx,
		idleConnectionsCh: make(chan struct{}),
		store:             store,
		Address:           address,
	}
}

func (s *Server) basicHandler() chi.Router {
	router := chi.NewRouter()

	router.Post("/coupons", func(w http.ResponseWriter, r *http.Request) {
		coupon := new(models.Coupon)
		if err := json.NewDecoder(r.Body).Decode(coupon); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Coupons().Create(r.Context(), coupon)
	})

	router.Get("/coupons", func(w http.ResponseWriter, r *http.Request) {
		coupons, err := s.store.Coupons().All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, coupons)
	})

	router.Get("/coupons/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		coupon, err := s.store.Coupons().ByID(r.Context(), id)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, coupon)
	})

	router.Put("/coupons", func(w http.ResponseWriter, r *http.Request) {
		coupon := new(models.Coupon)
		if err := json.NewDecoder(r.Body).Decode(coupon); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Coupons().Update(r.Context(), coupon)
	})

	router.Delete("/coupons/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Coupons().Delete(r.Context(), id)
	})

	router.Post("/registration", func(w http.ResponseWriter, r *http.Request) {
		user := new(models.User)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		err := s.store.Users().Create(r.Context(), user)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	})

	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		req := new(request)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		u, err := s.store.Users().ByEmail(r.Context(), req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			fmt.Fprintln(w, "Incorrect email or password")
			return
		}

	})

	router.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := s.store.Users().All(r.Context())
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		render.JSON(w, r, users)
	})

	return router
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
	go s.ListenCtxForGT(srv)

	log.Println("[HTTP] Server running on", s.Address)
	return srv.ListenAndServe()
}

func (s *Server) ListenCtxForGT(srv *http.Server) {
	<-s.ctx.Done()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[HTTP] Got err while shutting down^ %v", err)
	}

	log.Println("[HTTP] Processed all idle connections")
	close(s.idleConnectionsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConnectionsCh
}
