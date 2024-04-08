package postgres

import (
	"banner-service/internal/domain/models"
	"banner-service/internal/storage"
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

const (
	limitDefault  = 100
	offsetDefault = 0
)

type Storage struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetUserBanner(
	ctx context.Context,
	tagId int,
	featureId int,
) (*models.Banner, error) {
	var resultBanner models.Banner

	db := s.db.WithContext(ctx)

	tx := db.
		Where("tag_ids @> ?", pq.Int32Array{int32(tagId)}).
		Where("feature_id=?", featureId).
		Where("is_active=?", true).
		First(&resultBanner)

	if tx.RowsAffected == 0 {
		return nil, storage.ErrNotFound
	}
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &resultBanner, nil
}

func (s *Storage) GetBanners(
	ctx context.Context,
	tagId *int,
	featureId *int,
	limit *int,
	offset *int,
) ([]models.Banner, error) {
	const op = "storage:postgres:GetBanners"
	db := s.db.WithContext(ctx)

	var result []models.Banner
	limitNum, offsetNum := defineLimitAndOffset(limit, offset)
	var queryResult *gorm.DB

	if tagId != nil && featureId != nil {
		queryResult = db.
			Where("TagId=?", *tagId).
			Where("FeatureId=?", *featureId).
			Limit(limitNum).Offset(offsetNum).
			Find(result)
	} else if tagId != nil {
		queryResult = db.Where("TagId=?", *tagId).Limit(limitNum).Offset(offsetNum).Find(result)
	} else if featureId != nil {
		queryResult = db.Where("FeatureId=?", *featureId).Limit(limitNum).Offset(offsetNum).Find(result)
	} else {
		queryResult = db.Limit(limitNum).Offset(offsetNum).Find(result)
	}

	if queryResult.Error != nil {
		return nil, fmt.Errorf("%s: %w", op, queryResult.Error)
	}

	return result, nil
}

func (s *Storage) CreateBanner(
	ctx context.Context,
	tagIds pq.Int32Array,
	featureId int,
	content string,
	isActive bool,
) (bannerId int, err error) {
	const op = "storage:postgres:CreateBanner"

	db := s.db.WithContext(ctx)

	bannerModel := models.Banner{
		FeatureId: featureId,
		TagIds:    tagIds,
		Content:   content,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := db.Create(&bannerModel)
	if result.Error != nil {
		return -1, fmt.Errorf("%s: %w", op, result.Error)
	}

	return bannerModel.Id, nil

}

func (s *Storage) UpdateBanner(
	ctx context.Context,
	id int,
	tagIds pq.Int32Array,
	featureId int,
	content string,
	isActive bool,
) error {
	const op = "storage:postgres:UpdateBanner"

	db := s.db.WithContext(ctx)

	var bannerToUpdate models.Banner
	result := db.Where("Id=?", id).First(&bannerToUpdate)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return storage.ErrNotFound
		}

		return fmt.Errorf("%s: %w", op, result.Error)
	}

	bannerToUpdate.UpdatedAt = time.Now()
	bannerToUpdate.TagIds = tagIds
	bannerToUpdate.FeatureId = featureId
	bannerToUpdate.IsActive = isActive
	bannerToUpdate.Content = content

	resultUpdate := db.Save(&bannerToUpdate)
	if resultUpdate.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return storage.ErrNotFound
		}

		return fmt.Errorf("%s: %w", op, result.Error)
	}

	return nil
}

func (s *Storage) DeleteBanner(
	ctx context.Context,
	id int,
) error {
	const op = "storage:postgres:DeleteBanner"

	db := s.db.WithContext(ctx)

	deleteResult := db.Delete(&models.Banner{}, id)
	if deleteResult.Error != nil {
		if errors.Is(deleteResult.Error, gorm.ErrRecordNotFound) {
			return storage.ErrNotFound
		}

		return fmt.Errorf("%s: %w", op, deleteResult.Error)
	}

	return nil
}

func defineLimitAndOffset(limitLink *int, offsetLink *int) (int, int) {
	var limit, offset = limitDefault, offsetDefault
	if limitLink != nil {
		limit = *limitLink
	}
	if offsetLink != nil {
		offset = *offsetLink
	}

	return limit, offset
}
