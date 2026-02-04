package models

import (
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID       string            `json:"id"`
	Category string            `json:"category"`
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	Tags     []string          `json:"tags,omitempty"`
	Created  time.Time         `json:"created"`
	Updated  time.Time         `json:"updated"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func NewNote(category, title, content string, tags []string) *Note {
	now := time.Now()
	return &Note{
		ID:       uuid.New().String(),
		Category: category,
		Title:    title,
		Content:  content,
		Tags:     tags,
		Created:  now,
		Updated:  now,
		Metadata: map[string]string{
			"created_by": "agent",
			"source":     "braindump",
		},
	}
}
