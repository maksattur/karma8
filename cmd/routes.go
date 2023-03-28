package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maksattur/karma8/internal/server"
)

func registerRoutes(router fiber.Router, api *server.API) {
	router.Put("/upload", api.UploadFile)
	router.Get("/download/:file_name", api.DownloadFile)
}
