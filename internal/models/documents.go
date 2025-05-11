package models

type (
	PostDocumentRequest struct {
		Documents []File `json:"documents"`
		Type      string
	}
)
