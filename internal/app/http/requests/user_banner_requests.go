package requests

type GetUserBannerRequest struct {
	TagId           int  `json:"tag_id" validate:"gte=0" form:"tag_id"`
	FeatureId       int  `json:"feature_id" validate:"gte=0" form:"feature_id"`
	UseLastRevision bool `json:"use_last_revision,default:false" form:"use_last_revision"`
}
