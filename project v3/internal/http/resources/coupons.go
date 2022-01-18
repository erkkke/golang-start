package resources

import (
	"encoding/json"
	"fmt"
	"github.com/erkkke/golang-start/project/internal/message_broker"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/erkkke/golang-start/project/internal/store"
	"github.com/erkkke/golang-start/project/pkg"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"strconv"
)

type CouponsResource struct {
	store  store.Store
	broker message_broker.MessageBroker
	cache  *lru.TwoQueueCache
}

func NewCouponsResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *CouponsResource {
	return &CouponsResource{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (cr *CouponsResource) Routes(auth func(handler http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", cr.AllCoupons)
	r.Get("/{id}", cr.ByID)

	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Post("/", cr.CreateCoupon)
		r.Put("/", cr.UpdateCoupon)
		r.Delete("/{id}", cr.DeleteCoupon)
	})

	return r
}

func (cr *CouponsResource) CreateCoupon(w http.ResponseWriter, r *http.Request) {
	if !pkg.IsUserAdmin(r.Context(), w) {
		return
	}

	coupon := new(models.Coupon)

	if err := json.NewDecoder(r.Body).Decode(coupon); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := cr.store.Coupons().Create(r.Context(), coupon)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	// В идеале надо пройтись по всем буквам и по всем словам
	if err = cr.broker.Cache().Purge(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cr *CouponsResource) AllCoupons(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.NameFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		fromCache, ok := cr.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, fromCache)
			return
		}

		filter.Query = &searchQuery
	}

	coupons, err := cr.store.Coupons().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if searchQuery != "" && len(coupons) > 0 {
		cr.cache.Add(searchQuery, coupons)
	}

	render.JSON(w, r, coupons)
}

func (cr *CouponsResource) ByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	coupon, err := cr.store.Coupons().ById(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	render.JSON(w, r, coupon)
}

func (cr *CouponsResource) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
	if !pkg.IsUserAdmin(r.Context(), w) {
		return
	}

	coupon := new(models.Coupon)
	if err := json.NewDecoder(r.Body).Decode(coupon); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := cr.store.Coupons().Update(r.Context(), coupon)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if err = cr.broker.Cache().Remove(coupon.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}
}

func (cr *CouponsResource) DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err = cr.store.Coupons().Delete(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if err = cr.broker.Cache().Remove(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}
}
