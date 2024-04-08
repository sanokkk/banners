package models

import (
	"github.com/lib/pq"
	"time"
)

type Banner struct {
	Id        int           `json:"id" gorm:"primaryKey;autoIncrement:true"`
	TagIds    pq.Int32Array `json:"tag_ids" gorm:"type:integer[];default:'{}'"`
	FeatureId int           `json:"feature_id"`
	Content   string        `json:"content"`
	IsActive  bool          `json:"isActive" gorm:"true"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
