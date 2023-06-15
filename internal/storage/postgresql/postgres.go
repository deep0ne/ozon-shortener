package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/deep0ne/ozon-test/base63"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dsn = "host=postgres user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"

type Link struct {
	FullURL  string `gorm:"unique"`
	ShortUrl int64  `gorm:"unique"`
}

type PostgreSQL struct {
	DB *gorm.DB
}

func NewPostgreSQL() (*PostgreSQL, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Link{})

	return &PostgreSQL{
		DB: db,
	}, nil
}

func (p *PostgreSQL) AddLink(ctx context.Context, id int64, url string) (string, error) {

	var link Link
	result := p.DB.Where(&Link{FullURL: url}).Attrs(Link{ShortUrl: id}).FirstOrCreate(&link)

	if result.RowsAffected == 0 {
		short := base63.Encode(link.ShortUrl)
		return "", errors.New(fmt.Sprintf("short URL for this URL was already generated: %v", short))
	}

	shortURL := base63.Encode(id)
	return shortURL, nil
}

func (p *PostgreSQL) GetLink(ctx context.Context, id int64) (string, error) {
	var link Link
	if err := p.DB.First(&link, "short_url = ?", id).Error; err != nil {
		return "", errors.New("there is no URL for such short URL")
	}
	return link.FullURL, nil
}
