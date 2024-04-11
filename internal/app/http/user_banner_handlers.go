package http

import (
	"banner-service/internal/app/http/requests"
	"banner-service/internal/lib/validation"
	"banner-service/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func (s *Server) handleGetUserBanner(c *gin.Context) {
	const op = "http:handleGetUserBanner"
	logger := s.Log.With(slog.String("operation", op))

	var requestModel requests.GetUserBannerRequest
	if err := c.ShouldBindQuery(&requestModel); err != nil {
		logger.Warn(err.Error())
		s.respondWithError(c, http.StatusBadRequest, validation.ErrValidation)

		return
	}

	if err := validation.ValidateGetUserBannerRequest(requestModel); err != nil {
		s.respondWithError(c, http.StatusBadRequest, err)
		logger.Warn(err.Error())

		return
	}

	banner, err := s.Service.GetUserBanner(
		c,
		requestModel.TagId,
		requestModel.FeatureId,
		requestModel.UseLastRevision)

	if err != nil {
		if errors.Is(err, services.ErrServiceNotFound) {
			s.respondWithError(c, http.StatusNotFound, err)
			logger.Warn(err.Error())

			return
		}
		s.respondWithError(c, http.StatusInternalServerError, err)
		logger.Warn(err.Error())

		return
	}

	c.IndentedJSON(200, banner.Content)
}
