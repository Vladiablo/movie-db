package movies

import (
	"log/slog"
	"movie-db/internal/service"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	svc    *service.Service
	logger *slog.Logger
}

// TODO: separate search endpoint by year, genres, rating filter & sort
// TODO: sort by rating, votes, runtime
// TODO: with pagination
// TODO: separate tv_episodes into different table

type GetMoviesRequest struct {
	Id int64 `param:"id" validate:"required,min=1,max=2147483647"`
}

type GetSuggestionsRequest struct {
	Query string `query:"q" validate:"required,min=2,max=255"`
}

func New(svc *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
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

func (s *Handler) GetSuggestions(c *echo.Context) error {
	var req GetSuggestionsRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	movie, err := s.svc.GetSuggestions(c.Request().Context(), req.Query)
	if err != nil {
		s.logger.Error("Failed to get movie suggestions",
			slog.String("query", req.Query),
			slog.Any("error", err),
		)

		return echo.ErrInternalServerError.Wrap(err)
	}

	if movie == nil {
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, movie)
}
