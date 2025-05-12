package models

type AllGuides struct {
	Student []Guide `json:"student"`
	Worker  []Guide `json:"worker"`
}

type Guide struct {
	ID        int         `json:"id,omitempty" db:"id"`
	Label     string      `json:"label" db:"title"`
	SubGuides []SubGuides `json:"themes" db:"_"`
}

type SubGuides struct {
	ID      int    `json:"id,omitempty" db:"id"`
	Label   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}

type PostThemeRequest struct {
	GuideID int `json:"guide_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
