package movies

import (
	"log/slog"
	"movie-db/internal/helpers"
	"movie-db/internal/http/httphelpers"
	"movie-db/internal/service"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc    *service.Service
	logger *slog.Logger
}

func New(svc *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

type GetMoviesRequest struct {
	Id int64 `param:"id" validate:"required,min=1,max=2147483647"`
}

type GetSuggestionsRequest struct {
	Query string `query:"q" validate:"required,min=2,max=255"`
}

type SearchMoviesRequest struct {
	Type           string                 `query:"type" validate:"omitempty,oneof=movie series"`
	YearFrom       int32                  `query:"yearFrom" validate:"omitempty,min=1,max=9999"`
	YearTo         int32                  `query:"yearTo" validate:"omitempty,min=1,max=9999,gtefield=YearFrom"`
	GenreIds       httphelpers.Int32Slice `query:"genreIds" validate:"omitempty,dive,min=1"`
	ImdbRatingFrom float32                `query:"imdbRatingFrom" validate:"omitempty,min=1.0,max=10.0"`
	ImdbRatingTo   float32                `query:"imdbRatingTo" validate:"omitempty,min=1.0,max=10.0,gtefield=ImdbRatingFrom"`
	OrderBy        string                 `query:"orderBy" validate:"omitempty,oneof=imdb_rating year_from"`
	OrderDir       string                 `query:"orderDir" validate:"omitempty,oneof=asc desc"`
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
//	@Param			id	path		int64	true	"Movie ID"	Minimum(1)	Maximum(2147483647)	Example(12163543)
//	@Success		200	{array}		internal.Movie
//	@Failure		400	{object}	http.ErrorResponse
//	@Failure		404	{object}	http.ErrorResponse
//	@Failure		422	{object}	http.ErrorResponse
//	@Failure		500	{object}	http.ErrorResponse
//	@Router			/movies/{id} [get]
func (s *Handler) GetMovieById(c *echo.Context) error {
	var req GetMoviesRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	movie, err := s.svc.GetMovieById(c.Request().Context(), req.Id)
	if err != nil {
		s.logger.Error("Failed to get movie",
			slog.Int64("id", req.Id),
			slog.Any("error", err),
		)

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
//	@Param			q	query		string	true	"Search string"	MinLength(2)	MaxLength(255)	Example("game of thrones")
//	@Success		200	{array}		internal.MovieSuggestion
//	@Failure		400	{object}	http.ErrorResponse
//	@Failure		404	{object}	http.ErrorResponse
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
		s.logger.Error("Failed to get movie suggestions",
			slog.String("query", req.Query),
			slog.Any("error", err),
		)

		return echo.ErrInternalServerError.Wrap(err)
	}

	if len(movies) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, movies)
}

// SearchMovies
//
//	@Summary		Search movies
//	@Description	Search movies by criteria
//	@Tags			movies
//	@Produce		json
//	@Param			type			query		string	false	"Movie type to search for"																					Enums(movie, series)			Example(movie)
//	@Param			yearFrom		query		int32	false	"Lower bound of the year range for searching movies"														Minimum(1)						Maximum(9999)	Example(1982)
//	@Param			yearTo			query		int32	false	"Upper bound of the year range for searching movies. Should be greater or equal to yearFrom"				Minimum(1)						Maximum(9999)	Example(2024)
//	@Param			genreIds		query		[]int32	false	"IDs of movie genres to search for"																			CollectionFormat(csv)			Example(12,15)
//	@Param			imdbRatingFrom	query		float32	false	"Lower bound of the IMDB rating range for searching movies"													Minimum(1.0)					Maximum(10.0)			Example(6.0)
//	@Param			imdbRatingTo	query		float32	false	"Upper bound of the IMDB rating range for searching movies. Should be greater or equal to imdbRatingFrom"	Minimum(1.0)					Maximum(10.0)			Example(9.3)
//	@Param			orderBy			query		string	false	"Sorting key"																								Enums(imdb_rating, year_from)	Default(imdb_rating)	Example(year_from)
//	@Param			orderDir		query		string	false	"Sorting direction"																							Enums(asc, desc)				Default(desc)			Example(asc)
//	@Success		200				{array}		internal.MovieSuggestion
//	@Failure		400				{object}	http.ErrorResponse
//	@Failure		404				{object}	http.ErrorResponse
//	@Failure		422				{object}	http.ErrorResponse
//	@Failure		500				{object}	http.ErrorResponse
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
		GenreIds:       helpers.CopySlice(req.GenreIds),
		ImdbRatingFrom: req.ImdbRatingFrom,
		ImdbRatingTo:   req.ImdbRatingTo,
		OrderBy:        req.OrderBy,
		OrderDir:       req.OrderDir,
	})
	if err != nil {
		s.logger.Error("Failed to get movie suggestions",
			slog.Any("type", req.Type),
			slog.Any("genreIds", req.GenreIds),
			slog.Float64("imdbRatingFrom", float64(req.ImdbRatingFrom)),
			slog.Float64("imdbRatingTo", float64(req.ImdbRatingTo)),
			slog.Int64("yearFrom", int64(req.YearFrom)),
			slog.Int64("yearTo", int64(req.YearTo)),
			slog.String("orderBy", req.OrderBy),
			slog.String("orderDir", req.OrderDir),
			slog.Any("error", err),
		)

		return echo.ErrInternalServerError.Wrap(err)
	}

	if len(movies) == 0 {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, movies)
}
