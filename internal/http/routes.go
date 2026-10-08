package http

import "movie-db/internal/http/handlers/movies"

func RegisterRoutes(srv *Server, moviesHandler *movies.Handler) {
	v1 := srv.echo.Group("/v1")

	movies.RegisterRoutes(v1, moviesHandler)
}
