package handler

import (
	"encoding/json"
	"net/http"
	"FoodStore-AdvProg2/domain"
	"FoodStore-AdvProg2/usecase"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	UC *usecase.ProductUseCase
}

func NewProductHandler(uc *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{UC: uc}
}

func (h *ProductHandler) respondJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *ProductHandler) parseID(r *http.Request) string {
	return mux.Vars(r)["id"]
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	_ = json.NewDecoder(r.Body).Decode(&p)
	if err := h.UC.Create(p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.respondJSON(w, p, http.StatusCreated)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	p, err := h.UC.GetByID(h.parseID(r))
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	h.respondJSON(w, p, http.StatusOK)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	_ = json.NewDecoder(r.Body).Decode(&p)
	if err := h.UC.Update(h.parseID(r), p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.respondJSON(w, p, http.StatusOK)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.UC.Delete(h.parseID(r)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
    products, err := h.UC.List()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    h.respondJSON(w, products, http.StatusOK)
}