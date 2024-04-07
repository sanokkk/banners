package local

import (
	"banner-service/internal/domain/models"
	"banner-service/internal/storage"
	"sync"
	"time"
)

type LocalStorage struct {
	sync.RWMutex
	Banners map[LocalStorageKey]*models.Banner
}

type LocalStorageKey struct {
	TagId     int
	FeatureId int
}

func New() *LocalStorage {
	return &LocalStorage{Banners: make(map[LocalStorageKey]*models.Banner)}
}

// todo add update background method
func (s *LocalStorage) GetBannerFromCache(
	tagId int,
	featureId int,
) (*models.Banner, error) {
	key := LocalStorageKey{TagId: tagId, FeatureId: featureId}

	s.RWMutex.Lock()

	banner, exists := s.Banners[key]

	s.RWMutex.Unlock()

	if !exists {
		return nil, storage.ErrNotFound
	}

	return banner, nil
}

func (s *LocalStorage) CreateCacheBanner(
	tagIds []int,
	featureId int,
	content string,
	isActive bool,
) error {
	banner := &models.Banner{
		TagIds:    tagIds,
		FeatureId: featureId,
		Content:   content,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	for _, tagId := range tagIds {
		if success := s.RWMutex.TryLock(); success {
			key := LocalStorageKey{FeatureId: featureId, TagId: tagId}
			s.Banners[key] = banner
			s.RWMutex.Unlock()
		} else {
			return storage.ErrInsertInCache
		}
	}

	return nil
}

func (s *LocalStorage) UpdateCacheBanner(featureId int, tagIds []int, banner models.Banner) {
	for _, tagId := range tagIds {
		key := LocalStorageKey{FeatureId: featureId, TagId: tagId}

		s.RWMutex.Lock()
		s.Banners[key] = &banner
		s.RWMutex.Unlock()
	}
}

func (s *LocalStorage) GetCacheBanners() []models.Banner {
	result := make([]models.Banner, 0)
	s.RWMutex.Lock()
	for _, banner := range s.Banners {
		result = append(result, *banner)
	}
	s.RWMutex.Unlock()

	return result
}