package http

import (
	"banner-service/internal/config"
	"banner-service/internal/domain/models"
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"strconv"
	"time"
)

type Server struct {
	Log     *slog.Logger
	timeout time.Duration
	port    int
	Service BannerService
}

// New initialize Server instance
func New(log *slog.Logger, cfg config.Config, service BannerService) *Server {
	return &Server{Log: log, timeout: cfg.ApiConfig.Timeout, port: cfg.ApiConfig.Port, Service: service}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s *Server) respondWithError(c *gin.Context, code int, err error) {
	errResponse := errorResponse{Error: err.Error()}

	c.IndentedJSON(code, errResponse)
}

type BannerService interface {
	GetUserBanner(
		ctx context.Context,
		tagId int,
		featureId int,
		useLastRevision bool,
	) (*models.Banner, error)

	GetBanners(
		ctx context.Context,
		tagId *int,
		featureId *int,
		limit *int,
		offset *int,
	) ([]models.Banner, error)

	CreateBanner(
		ctx context.Context,
		tagIds []int,
		featureId int,
		content string,
		isActive bool,
	) (bannerId int, err error)

	UpdateBanner(
		ctx context.Context,
		id int,
		tagIds []int,
		featureId int,
		content string,
		isActive bool,
	) error

	DeleteBanner(
		ctx context.Context,
		id int,
	) error

	UpdateCache()
}

func (s *Server) MustServe() {
	const op = "http:MustServe"

	router := gin.Default()
	router.Use(s.timeoutMiddleware())

	groupUserBanner := router.Group("/user_banner")
	s.groupUserBanner(groupUserBanner)

	groupBanner := router.Group("/banner")
	s.groupBanner(groupBanner)

	s.Log.Info("starting http server", slog.Int("port", s.port))
	if err := router.Run(":" + strconv.Itoa(s.port)); err != nil {
		panic(err)
	}
}

// groupUserBanner groups methods for /user_banner route
func (s *Server) groupUserBanner(group *gin.RouterGroup) {
	group.Use(s.baseAuthMiddleware)

	group.GET("/user_banner", s.handleGetUserBanner)
}

// groupBanner groups methods for /banner route
func (s *Server) groupBanner(group *gin.RouterGroup) {
	group.Use(s.adminAuthMiddleware)

	group.GET("", s.handleGetBanner)
	group.POST("", s.handleCreateBanner)
	group.PATCH("/:id", s.handleUpdateBanner)
}
