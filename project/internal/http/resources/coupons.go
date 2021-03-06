package resources

import (
	"encoding/json"
	"fmt"
	"github.com/erkkke/golang-start/project/internal/cache"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/erkkke/golang-start/project/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type CouponsResource struct {
	store store.Store
	cache cache.Cache
}

func NewCouponsResource(store store.Store, cache cache.Cache) *CouponsResource {
	return &CouponsResource{
		store: store,
		cache: cache,
	}
}

func (cr *CouponsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", cr.CreateCoupon)
	r.Get("/", cr.AllCoupons)
	r.Get("/{id}", cr.ByID)
	r.Put("/", cr.UpdateCoupon)
	r.Delete("/{id}", cr.DeleteCoupon)

	return r
}

func (cr *CouponsResource) CreateCoupon(w http.ResponseWriter, r *http.Request) {
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
	if err = cr.cache.DeleteAll(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Cache err: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cr *CouponsResource) AllCoupons(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.NameFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		couponsFromCache, err := cr.cache.Coupons().Get(r.Context(), searchQuery)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Cache err: %v", err)
			return
		}

		if couponsFromCache != nil {
			render.JSON(w, r, couponsFromCache)
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
		err = cr.cache.Coupons().Set(r.Context(), searchQuery, coupons)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Cache err: %v", err)
			return
		}
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

	coupon, err := cr.store.Coupons().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	render.JSON(w, r, coupon)
}

func (cr *CouponsResource) UpdateCoupon(w http.ResponseWriter, r *http.Request) {
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
}

