package peers

import (
	"net/http"
	"vayeate-api/common"
	"vayeate-api/node"
	"vayeate-api/peer"

	"github.com/go-chi/render"
)

func Get(n *node.Node) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, common.SliceMap(n.GetPeers(), func(p *peer.Peer) any {
			return p.JSON()
		}))
	}
}
