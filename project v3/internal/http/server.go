package http

import (
	"context"
	"github.com/erkkke/golang-start/project/internal/http/resources"
	"github.com/erkkke/golang-start/project/internal/message_broker"
	"github.com/erkkke/golang-start/project/internal/store"
	"github.com/erkkke/golang-start/project/pkg/auth"
	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"net/http"
	"time"
)

const (
	readTimeout         = 5 * time.Second
	writeTimeout        = 30 * time.Second
)

type Server struct {
	ctx               context.Context
	idleConnectionsCh chan struct{}
	store             store.Store
	cache             *lru.TwoQueueCache
	broker            message_broker.MessageBroker
	tokenManager      auth.TokenManger

	Address string
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	srv := &Server{
		ctx:               ctx,
		idleConnectionsCh: make(chan struct{}),
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (s *Server) basicHandler() http.Handler {
	r := chi.NewRouter()

	// CATEGORIES HANDLER
	categoriesResource := resources.NewCategoriesResource(s.store, s.broker, s.cache)
	r.Mount("/categories", categoriesResource.Routes(s.userIdentity))

	// COUPONS HANDLER
	couponsResource := resources.NewCouponsResource(s.store, s.broker, s.cache)
	r.Mount("/coupons", couponsResource.Routes(s.userIdentity))

	// USERS HANDLER
	usersResource := resources.NewUsersResource(s.store, s.broker, s.cache)
	r.Mount("/users", usersResource.Routes(s.userIdentity))

	// Authentication
	authResource := resources.NewAuthResource(s.store, s.cache, s.tokenManager)
	r.Mount("/auth", authResource.Routes())

	// Order
	orderResource := resources.NewOrdersResource(s.store)
	r.Mount("/orders", orderResource.Routes(s.userIdentity))

	return r
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
