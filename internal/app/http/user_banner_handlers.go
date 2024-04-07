package http

import (
	"banner-service/internal/app/http/requests"
	"banner-service/internal/lib/validation"
	"banner-service/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// todo: return dto instead of banner model
func (s *Server) handleGetUserBanner(c *gin.Context) {
	var requestModel requests.GetUserBannerRequest
	if err := c.ShouldBindQuery(&requestModel); err != nil {
		s.respondWithError(c, http.StatusBadRequest, validation.ErrValidation)
	}

	if err := validation.ValidateGetUserBannerRequest(requestModel); err != nil {
		s.respondWithError(c, http.StatusBadRequest, err)
	}

	var useLastRevisionParameter bool
	if requestModel.UseLastRevision == nil {
		useLastRevisionParameter = false
	} else {
		useLastRevisionParameter = *requestModel.UseLastRevision
	}

	banner, err := s.Service.GetUserBanner(
		c,
		requestModel.TagId,
		requestModel.FeatureId,
		useLastRevisionParameter)

	if err != nil {
		if errors.Is(err, services.ErrServiceNotFound) {
			s.respondWithError(c, http.StatusNotFound, err)

			return
		}
		s.respondWithError(c, http.StatusInternalServerError, err)

		return
	}

	c.IndentedJSON(200, banner.Content)
}
