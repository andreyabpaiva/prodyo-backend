package api

import (
	"net/http"
	"prodyo-backend/cmd/internal/config"
	"prodyo-backend/cmd/internal/repositories"
)

func main() {
	cfg := config.Load()

	db := repositories.NewDB(cfg.DSN())
	defer db.Close()

	http.ListenAndServe(":8080", nil)
}
