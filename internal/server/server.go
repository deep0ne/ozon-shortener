package server

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"net/http"

	api "github.com/deep0ne/ozon-test/api/proto"
	"github.com/deep0ne/ozon-test/base63"
	"github.com/deep0ne/ozon-test/internal/storage"
	"github.com/deep0ne/ozon-test/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const address = "localhost:8082"

type URLShortenerServer struct {
	api.UnimplementedURLShortenerServer
	Logger *logrus.Logger
	db     storage.DB
	router *gin.Engine
}

func NewURLShortenerServer(db storage.DB) *URLShortenerServer {
	return &URLShortenerServer{
		Logger: utils.NewLogger(),
		db:     db,
	}
}

func (s *URLShortenerServer) CreateShortURL(ctx context.Context, url *api.OriginalURL) (*api.ShortenedURL, error) {

	id, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	s.Logger.Info("Adding short URL...")
	shortURL, err := s.db.AddLink(ctx, id.Int64(), url.URL)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return &api.ShortenedURL{ShortURL: shortURL}, nil
}

func (s *URLShortenerServer) GetOriginalURL(ctx context.Context, shortURL *api.ShortenedURL) (*api.OriginalURL, error) {
	id, err := base63.Decode(shortURL.ShortURL)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	fullURL, err := s.db.GetLink(ctx, id)
	if err != nil {
		s.Logger.Error(err)
		return nil, err
	}

	return &api.OriginalURL{URL: fullURL}, nil
}

func (s *URLShortenerServer) SetUpRouter() {
	router := gin.Default()

	router.POST("/shorten", s.POSTShortURL)
	router.GET("/:url", s.GETFullURL)

	s.router = router
}

func (s *URLShortenerServer) HTTPStart() error {
	s.Logger.Logln(logrus.InfoLevel, "Starting HTTP Server...")
	return s.router.Run(address)
}

func (s *URLShortenerServer) POSTShortURL(ctx *gin.Context) {
	var url api.OriginalURL
	if err := ctx.BindJSON(&url); err != nil {
		s.Logger.Error(err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	s.Logger.Logln(logrus.InfoLevel, url.URL)
	short, err := s.CreateShortURL(ctx, &url)
	if err != nil {
		s.Logger.Error(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	s.Logger.Logln(logrus.InfoLevel, "Successfully created short URL")
	ctx.JSON(http.StatusOK, fmt.Sprintf("Your Short URL: %v", short.ShortURL))
}

func (s *URLShortenerServer) GETFullURL(ctx *gin.Context) {
	url := ctx.Param("url")

	full, err := s.GetOriginalURL(ctx, &api.ShortenedURL{ShortURL: url})
	if err != nil {
		s.Logger.Error(err)
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}
	s.Logger.Logln(logrus.InfoLevel, "Redirecting...")
	ctx.Redirect(http.StatusMovedPermanently, full.URL)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
