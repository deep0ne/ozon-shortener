package server

import (
	"context"
	"log"
	"os/exec"
	"testing"
	"time"

	api "github.com/deep0ne/ozon-test/api/proto"
	"github.com/deep0ne/ozon-test/base63"
	"github.com/deep0ne/ozon-test/internal/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestURLShortenerServer(t *testing.T) {
	cmd := exec.Command("make", "-C", "../../", "redis")
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
	time.Sleep(2 * time.Second)
	db := memory.NewRedisDB()

	server := NewURLShortenerServer(db)
	testCases := []struct {
		originalURL string
	}{
		{"www.google.com"},
		{"www.vk.com"},
		{"www.instagram.com"},
		{"www.yahoo.com"},
		{"www.yandex.com"},
		{"www.mail.ru"},
		{"www.ozon.com"},
		{"www.wb.com"},
		{"www.amocrm.com"},
		{"www.avito.com"},
		{"www.kaspersky.com"},
		{"www.sberbank.com"},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.originalURL, func(t *testing.T) {
			originalURL := &api.OriginalURL{
				URL: tc.originalURL,
			}

			shortenedURL, err := server.CreateShortURL(context.Background(), originalURL)
			assert.NoError(t, err)

			createdID, err := base63.Decode(shortenedURL.ShortURL)
			assert.NoError(t, err)

			retrievedOriginalURL, err := server.GetOriginalURL(context.Background(), shortenedURL)
			assert.NoError(t, err)
			assert.Equal(t, originalURL.URL, retrievedOriginalURL.URL)

			short := base63.Encode(createdID)
			assert.Equal(t, shortenedURL.ShortURL, short)

		})
	}
}
