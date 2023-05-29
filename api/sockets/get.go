package sockets

import (
	"net/http"
	"vayeate-api/node"

	"github.com/go-chi/render"
)

func Get(n *node.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, n.GetSockets())
	}
}
