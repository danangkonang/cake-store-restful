package model

import "time"

type ProductResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Rating      int    `json:"rating"`
	Image       string `json:"image"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ProductPostRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Rating      int       `json:"rating" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductUpdateRequest struct {
	Id          int       `json:"id" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Rating      int       `json:"rating" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductDeleteRequest struct {
	Id int `json:"id" validate:"required"`
}
