package resources

import (
	"encoding/json"
	"fmt"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/erkkke/golang-start/project/internal/store"
	"github.com/erkkke/golang-start/project/pkg"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
)

type OrdersResource struct {
	store store.Store
}

func NewOrdersResource(store store.Store) *OrdersResource {
	return &OrdersResource{store: store}
}

func (o *OrdersResource) Routes(auth func(handler http.Handler) http.Handler) chi.Router {
	r := chi.NewRouter()
	r.Group(func(r chi.Router) {
		r.Use(auth)
		r.Post("/create", o.CreateOrder)
		r.Get("/all", o.AllOrders)
		r.Get("/", o.AllUserOrders)
		r.Post("/", o.ChangeOrderStatus)
		r.Get("/{id}", o.ById)
	})

	return r
}

func (o *OrdersResource) CreateOrder(w http.ResponseWriter, r *http.Request) {
	order := new(models.Order)
	if err := json.NewDecoder(r.Body).Decode(order); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}
	order.UserId = r.Context().Value(pkg.CtxKeyUser).(*models.AuthorizedUserInfo).Id

	if err := o.store.Orders().Create(r.Context(), order); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (o *OrdersResource) AllOrders(w http.ResponseWriter, r *http.Request) {
	if !pkg.IsUserAdmin(r.Context(), w) {
		return
	}

	queryValues := r.URL.Query()
	filter := new(models.NameFilter)

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		filter.Query = &searchQuery
	}

	orders, err := o.store.Orders().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	render.JSON(w, r, orders)
}

func (o *OrdersResource) AllUserOrders(w http.ResponseWriter, r *http.Request) {
	userInfo := r.Context().Value(pkg.CtxKeyUser).(*models.AuthorizedUserInfo)

	orders, err := o.store.Orders().AllOfUsers(r.Context(), userInfo.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	render.JSON(w, r, orders)
}

func (o *OrdersResource) ById(w http.ResponseWriter, r *http.Request) {
	if userInfo := r.Context().Value(pkg.CtxKeyUser).(*models.AuthorizedUserInfo); userInfo.Role != models.Admin {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "err: Insufficient rights to access data")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	order, err := o.store.Orders().ById(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	render.JSON(w, r, order)
}

func (o *OrdersResource) ChangeOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderStatus := new(models.OrderStatusDTO)

	if err := json.NewDecoder(r.Body).Decode(orderStatus); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := o.store.Orders().ChangeStatus(r.Context(), orderStatus); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}