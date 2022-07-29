package cors

import (
	"net/http"

	"github.com/rs/cors"
)

func CORSEnabled(handler http.Handler) http.Handler {
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
		},
	})
	return cors.Handler(handler)
}
