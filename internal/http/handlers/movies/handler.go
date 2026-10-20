package movies

import (
	"movie-db/internal/http/httphelpers"
	"movie-db/internal/service"
	"net/http"
	"slices"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

type GetMovieRequest struct {
	Id int64 `param:"id" validate:"required,min=1,max=2147483647" example:"12163543"` // Movie ID
}

type GetSuggestionsRequest struct {
	Query string `query:"q" json:"q" validate:"required,min=2,max=255" example:"game of thrones"` // Search string
}

type SearchMoviesRequest struct {
	Type           string                 `query:"type" validate:"omitempty,oneof=movie series" example:"movie"`                                       // Movie type to search for
	YearFrom       int32                  `query:"yearFrom" validate:"omitempty,min=1,max=9999" example:"1982"`                                        // Lower bound of the year range for searching movies
	YearTo         int32                  `query:"yearTo" validate:"omitempty,min=1,max=9999,gtefield=YearFrom" example:"2024"`                        // Upper bound of the year range for searching movies. Should be greater or equal to yearFro
	GenreIds       httphelpers.Int32Slice `query:"genreIds" validate:"omitempty,dive,min=1" example:"12,15"`                                           // IDs of movie genres to search for
	ImdbRatingFrom float32                `query:"imdbRatingFrom" validate:"omitempty,min=1.0,max=10.0" example:"6.0"`                                 // Lower bound of the IMDB rating range for searching movies
	ImdbRatingTo   float32                `query:"imdbRatingTo" validate:"omitempty,min=1.0,max=10.0,gtefield=ImdbRatingFrom" example:"9.3"`           // Upper bound of the IMDB rating range for searching movies. Should be greater or equal to imdbRatingFrom
	OrderBy        string                 `query:"orderBy" validate:"omitempty,oneof=imdb_rating year_from" example:"year_from" default:"imdb_rating"` // Sorting key
	OrderDir       string                 `query:"orderDir" validate:"omitempty,oneof=asc desc" example:"asc" default:"desc"`                          // Sorting direction
	PersonId       int32                  `query:"personId" validate:"omitempty,min=1" example:"15694579"`                                             // ID of a person who should be part of the movies's crew
}

func bindAndValidate(c *echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "validation failed").Wrap(err)
	}

	return nil
}

// GetMovieById
//
//	@Summary		Get movie
//	@Description	Get movie by its ID
//	@Tags			movies
//	@Produce		json
//	@Param			_	path		GetMovieRequest	true	"path params"
//	@Success		200	{array}		internal.Movie
//	@Failure		400	{object}	http.ErrorResponse
//	@Failure		404	{object}	http.ErrorResponse
//	@Failure		422	{object}	http.ErrorResponse
//	@Failure		500	{object}	http.ErrorResponse
//	@Router			/movies/{id} [get]
func (s *Handler) GetMovieById(c *echo.Context) error {
	var req GetMovieRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	movie, err := s.svc.GetMovieById(c.Request().Context(), req.Id)
	if err != nil {
		return echo.ErrInternalServerError.Wrap(err)
	}

	if movie == nil {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, movie)
}

// GetSuggestions
//
//	@Summary		Get movie suggestions
//	@Description	Suggest movies with title similar to search string
//	@Tags			movies
//	@Produce		json
//	@Param			_	query		GetSuggestionsRequest	true	"query params"
//	@Success		200	{array}		internal.MovieBase
//	@Failure		400	{object}	http.ErrorResponse
//	@Failure		422	{object}	http.ErrorResponse
//	@Failure		500	{object}	http.ErrorResponse
//	@Router			/movies/suggest [get]
//
// TODO: order by rating
func (s *Handler) GetSuggestions(c *echo.Context) error {
	var req GetSuggestionsRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	movies, err := s.svc.GetSuggestions(c.Request().Context(), req.Query)
	if err != nil {
		return echo.ErrInternalServerError.Wrap(err)
	}

	return c.JSON(http.StatusOK, movies)
}

// SearchMovies
//
//	@Summary		Search movies
//	@Description	Search movies by criteria
//	@Tags			movies
//	@Produce		json
//	@Param			_	query		SearchMoviesRequest	true	"query params"
//	@Success		200	{array}		internal.MovieBase
//	@Failure		400	{object}	http.ErrorResponse
//	@Failure		422	{object}	http.ErrorResponse
//	@Failure		500	{object}	http.ErrorResponse
//	@Router			/movies/search [get]
//
// TODO: with pagination
func (s *Handler) SearchMovies(c *echo.Context) error {
	var req SearchMoviesRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	movies, err := s.svc.SearchMovies(c.Request().Context(), &service.MoviesSearchParams{
		Type:           req.Type,
		YearFrom:       req.YearFrom,
		YearTo:         req.YearTo,
		GenreIds:       slices.Clone(req.GenreIds),
		ImdbRatingFrom: req.ImdbRatingFrom,
		ImdbRatingTo:   req.ImdbRatingTo,
		OrderBy:        req.OrderBy,
		OrderDir:       req.OrderDir,
		PersonId:       req.PersonId,
	})

	if err != nil {
		return echo.ErrInternalServerError.Wrap(err)
	}

	return c.JSON(http.StatusOK, movies)
}
