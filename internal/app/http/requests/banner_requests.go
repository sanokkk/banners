package requests

import "github.com/lib/pq"

type GetBannerRequest struct {
	TagId     *int `json:"tag_id"`
	FeatureId *int `json:"feature_id"`
	Limit     *int `json:"limit"`
	Offset    *int `json:"offset"`
}

type UpdateBannerRequest struct {
	TagIds    pq.Int32Array `json:"tag_ids" validate:"required"`
	FeatureId int           `json:"feature_id" validate:"required"`
	Content   string        `json:"content" validate:"required"`
	IsActive  bool          `json:"is_active" validate:"required"`
}

type CreateBannerRequest struct {
	TagIds    pq.Int32Array `json:"tag_ids" validate:"required"`
	FeatureId int           `json:"feature_id" validate:"required"`
	Content   string        `json:"content" validate:"required"`
	IsActive  bool          `json:"is_active" validate:"required"`
}
