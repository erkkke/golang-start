package resources

import (
	"encoding/json"
	"fmt"
	"github.com/erkkke/golang-start/project/internal/cache"
	"github.com/erkkke/golang-start/project/internal/models"
	"github.com/erkkke/golang-start/project/internal/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"net/http"
	"strconv"
)

type CategoriesResource struct {
	store store.Store
	cache cache.Cache
}

func NewCategoriesResource(store store.Store, cache cache.Cache) *CategoriesResource {
	return &CategoriesResource{
		store: store,
		cache: cache,
	}
}

func (cr *CategoriesResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", cr.CreateCategory)
	r.Get("/", cr.AllCategories)
	r.Get("/{id}", cr.ById)
	r.Put("/", cr.UpdateCategory)
	r.Delete("/{id}", cr.DeleteCategory)

	return r
}

func (cr *CategoriesResource) CreateCategory(w http.ResponseWriter, r *http.Request) {
	category := new(models.Category)
	if err := json.NewDecoder(r.Body).Decode(category); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err := cr.store.Categories().Create(r.Context(), category); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "BD err: %v", err)
		return
	}

	// В идеале надо пройтись по всем буквам и по всем словам
	if err := cr.cache.DeleteAll(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Cache err: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cr *CategoriesResource) AllCategories(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	filter := &models.CategoriesFilter{}

	searchQuery := queryValues.Get("query")
	if searchQuery != "" {
		categoriesFromCache, err := cr.cache.Categories().Get(r.Context(), searchQuery)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Cache err: %v", err)
			return
		}
		if categoriesFromCache != nil {
			render.JSON(w, r, categoriesFromCache)
			return
		}

		filter.Query = &searchQuery
	}

	categories, err := cr.store.Categories().All(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	if searchQuery != "" && len(categories) > 0{
		err = cr.cache.Categories().Set(r.Context(), searchQuery, categories)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Cache err: %v", err)
			return
		}
	}

	render.JSON(w, r, categories)
}

func (cr *CategoriesResource) ById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	category, err := cr.store.Categories().ByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}

	render.JSON(w, r, category)
}

func (cr *CategoriesResource) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	category := new(models.Category)
	if err := json.NewDecoder(r.Body).Decode(category); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	err := validation.ValidateStruct(category,
		validation.Field(&category.Id, validation.Required),
		validation.Field(&category.Name, validation.Required))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = cr.store.Categories().Update(r.Context(), category); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
}

func (cr *CategoriesResource) DeleteCategory (w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unknown err: %v", err)
		return
	}

	if err = cr.store.Categories().Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "DB err: %v", err)
		return
	}
}
