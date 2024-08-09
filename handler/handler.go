package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

var r *chi.Mux

func ProductHandler(handler *productHandler) {
	r = chi.NewRouter()

	r.Route("/product", func(r chi.Router) {
		r.Get("/", handler.listProducts)
		r.Post("/", handler.createProduct)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.getProduct)
			r.Patch("/", handler.updateProduct)
			r.Delete("/", handler.deleteProduct)
		})
	})
}

func Start(addr string) error {
	return http.ListenAndServe(addr, r)
}
