package main

import (
	"gopkg.in/macaron.v1"

	h "github.com/liangchenye/update-service/cmd/server/handler"
)

// SetRouters is the Dockyard Updater Server Router Definition
func SetRouters(m *macaron.Macaron) {
	// Web API
	m.Get("/", h.IndexMetaV1Handler)

	// App Discovery
	m.Group("/app", func() {
		m.Group("/v1", func() {
			m.Group("/:namespace/:repository", func() {
				// List files
				m.Get("/", h.AppListFileV1Handler)
				// Get pub key of the whole repo
				m.Get("/pubkey", h.AppGetPublicKeyV1Handler)
				// Get meta data of the whole repo
				m.Get("/meta", h.AppGetMetaV1Handler)
				// Get meta signature data of the whole repo
				m.Get("/metasign", h.AppGetMetaSignV1Handler)
				// Get file data of a certain app
				m.Get("/blob/:name", h.AppGetFileV1Handler)
				// Add file to the repo
				m.Post("/:name", h.AppPostFileV1Handler)
			})
		})
	})

}
