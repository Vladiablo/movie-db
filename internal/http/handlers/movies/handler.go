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

type GetMoviesRequest struct {
	Id int64 `param:"id" validate:"required,gte=1,lte=2147483647"`
}

func New(svc *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

func (s *Handler) GetMovieById(c *echo.Context) error {
	var req GetMoviesRequest

	if err := c.Bind(&req); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "validation failed").Wrap(err)
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
