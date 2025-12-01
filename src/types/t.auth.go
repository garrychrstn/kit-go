package types

import (
	"encoding/json"
)

type IRequestLogin struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

type IPartStore struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description" binding:"required"`
	Map            string          `json:"map" binding:"required"`
	Coordinate     string          `json:"coordinate"`
	Address        string          `json:"address" binding:"required"`
	Phone          string          `json:"phone" binding:"required"`
	Category       []string        `json:"category"`
	Contacts       json.RawMessage `json:"contacts"`
	TermAndService string          `json:"term_and_service" binding:"required"`
}

type IPartUser struct {
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	Name            string `json:"name" binding:"required"`
}

type IRequestSetup struct {
	Store IPartStore `json:"store" binding:"required"`
	User  IPartUser  `json:"user" binding:"required"`
}
