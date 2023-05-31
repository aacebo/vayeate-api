package main

import (
	"fmt"
	"net/http"
	"time"
	"vayeate-api/api"
	"vayeate-api/common"
	"vayeate-api/logger"
	"vayeate-api/node"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func main() {
	log := logger.New("vayeate")
	n, err := node.New(
		common.GetEnv("VAYEATE_SOCKET_PORT", "6789"),
		common.GetEnv("VAYEATE_PEER_PORT", "9876"),
		common.GetEnv("VAYEATE_USERNAME", "admin"),
		common.GetEnv("VAYEATE_PASSWORD", "admin"),
	)

	if err != nil {
		log.Error(err)
		return
	}

	if eaddr := common.GetEnv("VAYEATE_ENTRY_ADDRESS", ""); eaddr != "" {
		err = n.Discover(eaddr)

		if err != nil {
			log.Error(err)
			return
		}
	}

	defer n.Close()
	n.Listen()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(cors.AllowAll().Handler)
	r.Mount("/", api.NewRouter(n))

	err = http.ListenAndServe(fmt.Sprintf(":%s", common.GetEnv("VAYEATE_API_PORT", "3000")), r)

	if err != nil {
		log.Error(err)
	}
}
