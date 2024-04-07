package models

import (
	"time"
)

type Banner struct {
	Id        int       `json:"id" gorm:"primaryKey;autoIncrement:true"`
	TagIds    []int     `json:"tag_ids"`
	FeatureId int       `json:"feature_id"`
	Content   string    `json:"content"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
