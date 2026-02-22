package requests

import "encoding/json"

type IRequestCreateStore struct {
	Name           string          `json:"name" binding:"required"`
	Description    string          `json:"description"`
	Logo           string          `json:"logo"`
	Coordinate     string          `json:"coordinate"`
	Address        string          `json:"address" binding:"required"`
	Phone          string          `json:"phone" binding:"required"`
	IsActive       bool            `json:"is_active"`
	Category       []string        `json:"category"`
	Contacts       json.RawMessage `json:"contacts"`
	TermAndService string          `json:"term_and_service"`
}
