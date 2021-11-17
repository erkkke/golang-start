package resources

import (
	"encoding/json"
	"fmt"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/erkkke/golang-start/project/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	lru "github.com/hashicorp/golang-lru"

	"net/http"
	"strconv"
)

type UsersResource struct {
	store store.Store
	cache *lru.TwoQueueCache
}

func NewUsersResource(store store.Store, cache *lru.TwoQueueCache) *UsersResource {
	return &UsersResource{
		store: store,
		cache: cache,
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

	err := ur.store.Users().Create(r.Context(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
		return
	}
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
	users, err := ur.store.Users().All(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
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
}
