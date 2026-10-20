package movies

import "github.com/labstack/echo/v5"

func RegisterRoutes(e *echo.Group, moviesHandler *Handler) {
	grp := e.Group("/movies")

	grp.GET("/search", moviesHandler.SearchMovies)
	grp.GET("/suggest", moviesHandler.GetSuggestions)
	grp.GET("/:id", moviesHandler.GetMovieById)
}
