package handler

import (
	"chi-sqlx/database/entity"
	"chi-sqlx/service"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

type productHandler struct {
	ctx     context.Context
	service *service.ProductService
}

func NewProductController(service *service.ProductService) *productHandler {
	return &productHandler{
		ctx:     context.Background(),
		service: service,
	}
}

func toStoreProduct(p entity.ProductReq) *entity.Product {
	return &entity.Product{
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
	}
}

func toProductRes(p *entity.Product) entity.ProductRes {
	return entity.ProductRes{
		ID:           p.ID,
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
		DeletedAt:    p.DeletedAt,
		Name:         p.Name,
		Image:        p.Image,
		Category:     p.Category,
		Description:  p.Description,
		Rating:       p.Rating,
		NumReviews:   p.NumReviews,
		Price:        p.Price,
		CountInStock: p.CountInStock,
	}
}

func toTimePtr(t time.Time) time.Time {
	return t
}

func patchProductReq(product *entity.Product, p entity.ProductReq) {
	if p.Name != "" {
		product.Name = p.Name
	}
	if p.Image != "" {
		product.Image = p.Image
	}
	if p.Category != "" {
		product.Category = p.Category
	}
	if p.Description != "" {
		product.Description = p.Description
	}
	if p.Rating != 0 {
		product.Rating = p.Rating
	}
	if p.NumReviews != 0 {
		product.NumReviews = p.NumReviews
	}
	if p.Price != 0 {
		product.Price = p.Price
	}
	if p.CountInStock != 0 {
		product.CountInStock = p.CountInStock
	}
	product.UpdatedAt = toTimePtr(time.Now())
}

func (h *productHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	var p entity.ProductReq
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	fmt.Println("test", p, toStoreProduct(p))

	product, err := h.service.CreateProduct(h.ctx, toStoreProduct(p))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error creating product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *productHandler) getProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProduct(h.ctx, i)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error getting product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *productHandler) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(h.ctx)
	if err != nil {
		http.Error(w, "error listing product", http.StatusInternalServerError)
		return
	}

	var res []entity.ProductRes
	for _, p := range products {
		res = append(res, toProductRes(&p))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (h *productHandler) updateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing ID", http.StatusBadRequest)
		return
	}

	var p entity.ProductReq
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "error decoding request body", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProduct(h.ctx, i)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error getting product", http.StatusInternalServerError)
		return
	}

	// patch our product request
	patchProductReq(product, p)

	updated, err := h.service.UpdateProduct(h.ctx, product)
	if err != nil {
		http.Error(w, "error creating product", http.StatusInternalServerError)
		return
	}

	res := toProductRes(updated)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (h *productHandler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, "error parsing ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteProduct(h.ctx, i); err != nil {
		http.Error(w, "error getting product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
