package requests

import "encoding/json"

type IRequestCreateStore struct {
	Name           string          `json:"name" binding:"required"`
	Phone          string          `json:"phone" binding:"required"`
	Coordinate     string          `json:"coordinate"`
	Category       []string        `json:"category" binding:"required,min=1"`
	Description    string          `json:"description"`
	Contacts       json.RawMessage `json:"contacts"`
	Logo           string          `json:"logo"`
	IsActive       bool            `json:"is_active"`
	TermAndService string          `json:"tos"`
}
