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

type Server struct {
	ctx        context.Context
	idleConsCh chan struct{}
	store      store.Store

	Address string
}

func NewServer(ctx context.Context, address string, store store.Store) *Server {
	return &Server{
		ctx:        ctx,
		idleConsCh: make(chan struct{}),
		store:      store,
		Address:    address,
	}
}

func (s *Server) basicHandler() chi.Router {
	router := chi.NewRouter()

	router.Post("/coupons", func(writer http.ResponseWriter, request *http.Request) {
		coupon := new(models.Coupon)
		if err := json.NewDecoder(request.Body).Decode(coupon); err != nil {
			fmt.Fprintf(writer, "Unknown err: %v", err)
			return
		}

		s.store.Create(request.Context(), coupon)
	})

	router.Get("/coupons", func(writer http.ResponseWriter, request *http.Request) {
		coupons, err := s.store.All(request.Context())
		if err != nil {
			fmt.Fprintf(writer, "Unknown err: %v", err)
			return
		}

		render.JSON(writer, request, coupons)
	})

	router.Get("/coupons/{id}", func(writer http.ResponseWriter, request *http.Request) {
		idStr := chi.URLParam(request, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(writer, "Unknown err: %v", err)
			return
		}

		coupon, err := s.store.ByID(request.Context(), id)
		if err != nil {
			fmt.Fprintf(writer, "Unknown err: %v", err)
			return
		}

		render.JSON(writer, request, coupon)
	})

	router.Put("/coupons", func(w http.ResponseWriter, r *http.Request) {
		coupon := new(models.Coupon)
		if err := json.NewDecoder(r.Body).Decode(coupon); err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Update(r.Context(), coupon)
	})

	router.Delete("/coupons/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(w, "Unknown err: %v", err)
			return
		}

		s.store.Delete(r.Context(), id)
	})

	return router
}

func (s *Server) Run() error {
	srv := &http.Server{
		Addr:         s.Address,
		Handler:      s.basicHandler(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 30,
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
	close(s.idleConsCh)
}

func (s *Server) WaitForGracefulTermination() {
	<-s.idleConsCh
}
