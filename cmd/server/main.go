package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"sync"

	api "github.com/deep0ne/ozon-test/api/proto"
	"github.com/deep0ne/ozon-test/internal/server"
	"github.com/deep0ne/ozon-test/internal/storage"
	"github.com/deep0ne/ozon-test/internal/storage/memory"
	"github.com/deep0ne/ozon-test/internal/storage/postgresql"
	"github.com/deep0ne/ozon-test/internal/storage/redisdb"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func flagParsing() (bool, bool) {
	var inMemory, Redis bool
	flag.BoolVar(&inMemory, "m", false, "Usage: pass -m if you want to be your db in memory")
	flag.BoolVar(&Redis, "r", false, "Usage: pass -r if you want to initialize redis")
	flag.Parse()
	return inMemory, Redis
}

func run() error {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}

	inMemory, Redis := flagParsing()
	if inMemory && Redis {
		return errors.New("you need to pass either -m or -r to initialize database. See usage.")
	}
	var db storage.DB

	if Redis {
		db = redisdb.NewRedisDB()
	} else if inMemory {
		db = memory.NewInMemoryDB()
	} else {
		db, err = postgresql.NewPostgreSQL()
		if err != nil {
			return err
		}
	}

	grpcServer := grpc.NewServer()
	if err != nil {
		return err
	}
	shortener := server.NewURLShortenerServer(db)
	api.RegisterURLShortenerServer(grpcServer, shortener)

	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error, 2)

	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			errChan <- err
		}
		wg.Done()
	}()

	shortener.SetUpRouter()
	go func() {
		log.Println("Starting http server")
		err = shortener.HTTPStart()
		if err != nil {
			errChan <- err
		}
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
