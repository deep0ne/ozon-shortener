package memory

import (
	"context"
	"errors"
	"fmt"

	"github.com/deep0ne/ozon-test/base63"
)

type InMemoryDB struct {
	DB map[string]string
}

func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		DB: make(map[string]string),
	}
}

func (m *InMemoryDB) AddLink(ctx context.Context, id int64, url string) (string, error) {

	short := base63.Encode(id)
	if _, ok := m.DB[url]; ok {
		return "", errors.New(fmt.Sprintf("short URL for this URL was already generated: %v", m.DB[url]))
	}
	if _, ok := m.DB[short]; ok {
		return "", errors.New("algorithm generated short url that already exists, try again")
	}
	m.DB[short] = url
	m.DB[url] = short
	return short, nil
}

func (m *InMemoryDB) GetLink(ctx context.Context, id int64) (string, error) {
	short := base63.Encode(id)
	if fullURL, ok := m.DB[short]; ok {
		return fullURL, nil
	} else {
		return "", errors.New("there is no url for such short URL")
	}
}
