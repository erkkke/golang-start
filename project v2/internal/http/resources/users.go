package resources

import (
	"encoding/json"
	"fmt"
	"github.com/erkkke/golang-start/project/internal/message_broker"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/erkkke/golang-start/project/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"strconv"
)

type UsersResource struct {
	store  store.Store
	broker message_broker.MessageBroker
	cache  *lru.TwoQueueCache
}

func NewUsersResource(store store.Store, broker message_broker.MessageBroker, cache *lru.TwoQueueCache) *UsersResource {
	return &UsersResource{
		store:  store,
		broker: broker,
		cache:  cache,
	}
}

func (ur *UsersResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/registration", ur.CreateUser)
	r.Post("/login", ur.LoginUser)
	r.Get("/", ur.AllUsers)
	r.Put("/", ur.UpdateUser)
	r.Delete("/{id}", ur.DeleteUser)

	return r
}

func (ur *UsersResource) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := user.BeforeCreating(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	err := ur.store.Users().Create(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}

	// В идеале надо пройтись по всем буквам и по всем словам
	if err = ur.broker.Cache().Purge(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ur *UsersResource) LoginUser(w http.ResponseWriter, r *http.Request) {
	req := &struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	u, err := ur.store.Users().FindByEmail(r.Context(), req.Email)
	if err != nil || !u.ComparePassword(req.Password) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Incorrect email or password")
		return
	}
}

func (ur *UsersResource) AllUsers(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := new(models.NameFilter)

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		fromCache, ok := ur.cache.Get(searchQuery)
		if ok {
			render.JSON(w, r, fromCache)
			return
		}

		filter.Query = &searchQuery
	}

	users, err := ur.store.Users().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	if searchQuery != "" && len(users) > 0 {
		ur.cache.Add(searchQuery, users)
	}

	render.JSON(w, r, users)
}

func (ur *UsersResource) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := ur.store.Users().Update(r.Context(), user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	if err := ur.broker.Cache().Remove(user.ID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}
}

func (ur *UsersResource) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = ur.store.Users().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	if err = ur.broker.Cache().Remove(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Broker err: %v", err)
		return
	}
}
