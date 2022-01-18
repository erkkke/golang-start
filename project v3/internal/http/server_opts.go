package http

import (
	"github.com/erkkke/golang-start/project/internal/message_broker"
	"github.com/erkkke/golang-start/project/internal/store"
	"github.com/erkkke/golang-start/project/pkg/auth"
	lru "github.com/hashicorp/golang-lru"
)

type ServerOption func(srv *Server)

func WithAddress(address string) ServerOption {
	return func(srv *Server) {
		srv.Address = address
	}
}

func WithStore(store store.Store) ServerOption {
	return func(srv *Server) {
		srv.store = store
	}
}

func WithCache(cache *lru.TwoQueueCache) ServerOption {
	return func(srv *Server) {
		srv.cache = cache
	}
}

func WithBroker(broker message_broker.MessageBroker) ServerOption {
	return func(srv *Server) {
		srv.broker = broker
	}
}

func WithTokenManager(tokenManager auth.TokenManger) ServerOption {
	return func(srv *Server) {
		srv.tokenManager = tokenManager
	}
}
