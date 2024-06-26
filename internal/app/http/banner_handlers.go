package http

import (
	"banner-service/internal/app/http/requests"
	"banner-service/internal/lib/validation"
	"banner-service/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

func (s *Server) handleGetBanner(c *gin.Context) {
	const op = "http:handleGetBanner"
	logger := s.Log.With("op", op)

	var req requests.GetBannerRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		s.respondWithError(c, http.StatusBadRequest, validation.ErrValidation)
		logger.Warn(err.Error())

		return
	}

	banners, err := s.Service.GetBanners(c, req.TagId, req.FeatureId, req.Limit, req.Offset)
	if err != nil {
		s.respondWithError(c, http.StatusInternalServerError, err)
		logger.Warn(err.Error())

		return
	}

	c.IndentedJSON(200, banners)
}

func (s *Server) handleCreateBanner(c *gin.Context) {
	const op = "http:handleCreateBanner"
	logger := s.Log.With("op", op)

	var req requests.CreateBannerRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		s.respondWithError(c, http.StatusBadRequest, validation.ErrValidation)
		logger.Warn(err.Error())

		return
	}

	if err := validation.ValidateCreateBannerRequest(req); err != nil {
		s.respondWithError(c, http.StatusBadRequest, validation.ErrValidation)
		logger.Warn(err.Error())

		return
	}

	bannerId, err := s.Service.CreateBanner(c, req.TagIds, req.FeatureId, req.Content, req.IsActive)

	if err != nil {
		s.respondWithError(c, http.StatusInternalServerError, err)
		logger.Warn(err.Error())

		return
	}

	c.JSON(http.StatusCreated, requests.CreateBannerResponse{BannerId: bannerId})
}

func (s *Server) handleUpdateBanner(c *gin.Context) {
	const op = "http:handleUpdateBanner"
	logger := s.Log.With("op", op)

	var req requests.CreateBannerRequest
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		s.respondWithError(c, http.StatusBadRequest, validation.ErrValidation)
		logger.Warn(err.Error())

		return
	}

	if err := validation.ValidateCreateBannerRequest(req); err != nil {
		s.respondWithError(c, http.StatusBadRequest, validation.ErrValidation)
		logger.Warn(err.Error())

		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.respondWithError(c, http.StatusBadRequest, validation.ErrValidation)
		logger.Warn(err.Error())

		return
	}

	if err := s.Service.UpdateBanner(c, id, req.TagIds, req.FeatureId, req.Content, req.IsActive); err != nil {
		if errors.Is(err, services.ErrServiceNotFound) {
			s.respondWithError(c, http.StatusNotFound, err)
			logger.Warn(err.Error())

			return
		}
		s.respondWithError(c, http.StatusInternalServerError, err)
		logger.Warn(err.Error())

		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) handleDeleteBanner(c *gin.Context) {
	const op = "http:handleUpdateBanner"
	logger := s.Log.With("op", op)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.respondWithError(c, http.StatusBadRequest, validation.ErrValidation)
		logger.Warn(err.Error())

		return
	}

	if err := s.Service.DeleteBanner(c, id); err != nil {
		if errors.Is(err, services.ErrServiceNotFound) {
			s.respondWithError(c, http.StatusNotFound, err)
			logger.Warn(err.Error())

			return
		}
		s.respondWithError(c, http.StatusInternalServerError, err)
		logger.Warn(err.Error())

		return
	}

	c.Status(http.StatusOK)
}
