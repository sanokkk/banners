package requests

type GetUserBannerRequest struct {
	TagId           int   `json:"tag_id" validate:"required"`
	FeatureId       int   `json:"feature_id" validate:"required"`
	UseLastRevision *bool `json:"use_last_revision,omitempty"`
}
