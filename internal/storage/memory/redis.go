package memory

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/deep0ne/ozon-test/base63"
	"github.com/go-redis/redis"
)

const expirationTime = time.Hour * 24 * 180

type RedisDB struct {
	Client *redis.Client
}

func NewRedisDB() *RedisDB {
	return &RedisDB{Client: redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})}
}

func (r *RedisDB) AddLink(ctx context.Context, id int64, url string) (string, error) {
	idKey := strconv.FormatInt(id, 10)
	_, err := r.Client.Get(idKey).Result()
	if err == nil {
		return "", errors.New("such short URL was already generated")
	}

	err = r.Client.Get(url).Err()
	if err == nil {
		return "", errors.New("short url for that URL was already generated")
	}

	shortURL := base63.Encode(id)
	err = r.Client.Set(idKey, url, expirationTime).Err()
	if err != nil {
		return "", err
	}

	err = r.Client.Set(url, idKey, expirationTime).Err()
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (r *RedisDB) GetLink(ctx context.Context, id int64) (string, error) {
	idKey := strconv.FormatInt(id, 10)
	fullURL, err := r.Client.Get(idKey).Result()
	if err != nil {
		return "", errors.New("there is no url for such short URL")
	}
	return fullURL, nil
}
