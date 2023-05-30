package api

import (
	"net/http"
	"time"
	"vayeate-api/api/peers"
	"vayeate-api/api/sockets"
	"vayeate-api/node"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func NewRouter(n *node.Node) *chi.Mux {
	now := time.Now()
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]interface{}{
			"id":         n.ID,
			"started_at": now.Unix(),
		})
	})

	r.Route("/v1", func(r chi.Router) {
		sockets.NewRouter(r, n)
		peers.NewRouter(r, n)
	})

	return r
}
