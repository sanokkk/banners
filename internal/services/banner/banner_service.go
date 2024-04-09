package banner

import (
	"banner-service/internal/domain/models"
	"banner-service/internal/lib/slice"
	"banner-service/internal/services"
	"banner-service/internal/storage"
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log/slog"
)

type Service struct {
	bannersRepo
	bannersCacheRepo
	log *slog.Logger
}

func New(repo bannersRepo, cacheRepo bannersCacheRepo, logger *slog.Logger) *Service {
	return &Service{repo, cacheRepo, logger}
}

type bannersRepo interface {
	GetUserBanner(
		ctx context.Context,
		tagId int,
		featureId int,
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
		tagIds pq.Int32Array,
		featureId int,
		content string,
		isActive bool,
	) (bannerId int, err error)

	UpdateBanner(
		ctx context.Context,
		id int,
		tagIds pq.Int32Array,
		featureId int,
		content string,
		isActive bool,
	) error

	DeleteBanner(
		ctx context.Context,
		id int,
	) error
}

type bannersCacheRepo interface {
	GetBannerFromCache(
		tagId int,
		featureId int,
	) (*models.Banner, error)

	CreateCacheBanner(
		tagIds pq.Int32Array,
		featureId int,
		content string,
		isActive bool,
	) error

	UpdateCacheBanner(
		featureId int,
		tagIds pq.Int32Array,
		banner models.Banner,
	)

	GetCacheBanners() []models.Banner
}

func (s *Service) GetUserBanner(
	ctx context.Context,
	tagId int,
	featureId int,
	useLastRevision bool,
) (*models.Banner, error) {
	const op = "services:banner:GetUserBanner"

	var banner *models.Banner
	var err error

	if useLastRevision {
		banner, err = s.bannersRepo.GetUserBanner(ctx, tagId, featureId)
	} else {
		banner, err = s.bannersCacheRepo.GetBannerFromCache(tagId, featureId)
	}

	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, services.ErrServiceNotFound
		}

		return nil, fmt.Errorf("%s:%w", op, err)
	}
	return banner, nil
}

func (s *Service) GetBanners(
	ctx context.Context,
	tagId *int,
	featureId *int,
	limit *int,
	offset *int,
) ([]models.Banner, error) {
	const op = "services:banner:GetBanners"

	banners, err := s.bannersRepo.GetBanners(ctx, tagId, featureId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	return banners, nil
}

func (s *Service) CreateBanner(
	ctx context.Context,
	tagIds pq.Int32Array,
	featureId int,
	content string,
	isActive bool,
) (bannerId int, err error) {
	const op = "services:banner:CreateBanner"

	logger := s.log.With(slog.String("op", op))

	id, err := s.bannersRepo.CreateBanner(ctx, tagIds, featureId, content, isActive)
	if err != nil {

		logger.Warn(err.Error(),
			slog.Int("featureId", featureId),
			slog.String("tagIds", fmt.Sprintf("%v", tagIds)))

		return -1, err
	}

	if err := s.bannersCacheRepo.CreateCacheBanner(tagIds, featureId, content, isActive); err != nil {
		logger.Warn("error while inserting to cache", slog.String("error", err.Error()))
	}

	return id, nil
}

func (s *Service) UpdateBanner(
	ctx context.Context,
	id int,
	tagIds pq.Int32Array,
	featureId int,
	content string,
	isActive bool,
) error {
	const op = "services:banner:UpdateBanner"
	logger := s.log.With(slog.String("operation", op))

	err := s.bannersRepo.UpdateBanner(ctx, id, tagIds, featureId, content, isActive)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			logger.Warn("banner not found",
				slog.Int("featureId", featureId),
				slog.String("tagIds", fmt.Sprintf("%v", tagIds)))

			return services.ErrServiceNotFound
		}
		logger.Warn(err.Error(),
			slog.Int("featureId", featureId),
			slog.String("tagIds", fmt.Sprintf("%v", tagIds)))

		return err
	}

	return nil
}

func (s *Service) DeleteBanner(
	ctx context.Context,
	id int,
) error {
	const op = "services:banner:DeleteBanner"
	logger := s.log.With(slog.String("operation", op))

	err := s.bannersRepo.DeleteBanner(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			logger.Warn("banner not found",
				slog.Int("bannerId", id))

			return services.ErrServiceNotFound
		}
		logger.Warn(err.Error(),
			slog.Int("bannerId", id))

		return err
	}

	return nil
}

func (s *Service) UpdateCache() {
	const op = "services:banner:UpdateCache"
	logger := s.log.With(slog.String("op", op))

	logger.Info("starting updating cache")

	bannersFromCache := s.bannersCacheRepo.GetCacheBanners()
	if len(bannersFromCache) == 0 {
		return
	}

	ctx := context.Background()

	for _, cached := range bannersFromCache {
		tagsArray := slice.ConvertToIntSlice(cached.TagIds)

		bannerFromDb, err := s.bannersRepo.GetUserBanner(ctx, tagsArray[0], cached.FeatureId)
		if err != nil {
			logger.Warn("error while getting banner to update cache", slog.String("err", err.Error()))
			continue
		}

		s.bannersCacheRepo.UpdateCacheBanner(cached.FeatureId, cached.TagIds, *bannerFromDb)
	}
}
