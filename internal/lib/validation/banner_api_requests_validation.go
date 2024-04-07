package validation

import (
	"banner-service/internal/app/http/requests"
	"errors"
	"github.com/go-playground/validator/v10"
)

var ErrValidation = errors.New("Некорректные данные")

var validatorObject = validator.Validate{}

func ValidateGetUserBannerRequest(req requests.GetUserBannerRequest) error {
	if err := validatorObject.Struct(req); err != nil {
		return ErrValidation
	}

	return nil
}

func ValidateCreateBannerRequest(req requests.CreateBannerRequest) error {
	if err := validatorObject.Struct(req); err != nil {
		return ErrValidation
	}

	return nil
}

func ValidateUpdateBannerRequest(req requests.UpdateBannerRequest) error {
	if err := validatorObject.Struct(req); err != nil {
		return ErrValidation
	}

	return nil
}
