package peers

import (
	"vayeate-api/node"

	"github.com/go-chi/chi/v5"
)

func NewRouter(r chi.Router, n *node.Node) {
	r.Get("/peers", Get(n))
}
