package storage

import (
	"github.com/MohGanji/braindump/pkg/models"
)

type Store interface {
	Add(note *models.Note) error
	Get(id string) (*models.Note, error)
	GetByTitle(category, title string) (*models.Note, error)
	List(category string) ([]*models.Note, error)
	Update(note *models.Note) error
	Delete(id string) error
	Search(query string, category string, tags []string) ([]*models.Note, error)
	GetCategories() ([]string, error)
	GetTags() ([]string, error)
	Close() error
}
