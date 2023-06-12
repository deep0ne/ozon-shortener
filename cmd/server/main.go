package main

import (
	"flag"
	"log"
	"net"
	"os/exec"
	"sync"
	"time"

	api "github.com/deep0ne/ozon-test/api/proto"
	"github.com/deep0ne/ozon-test/internal/server"
	"github.com/deep0ne/ozon-test/internal/storage"
	"github.com/deep0ne/ozon-test/internal/storage/memory"
	"github.com/deep0ne/ozon-test/internal/storage/postgresql"
	"google.golang.org/grpc"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func flagParsing() bool {
	var inMemory bool
	flag.BoolVar(&inMemory, "m", false, "Usage: pass -m if you want to be your db in memory")
	flag.Parse()
	return inMemory
}

func run() error {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		return err
	}

	inMemory := flagParsing()
	var db storage.DB

	if inMemory {
		log.Println("Initializing Redis... Please wait a few seconds...")
		cmd := exec.Command("make", "-C", "../../", "redis")
		if err := cmd.Run(); err != nil {
			log.Println(err)
			return err
		}
		time.Sleep(5 * time.Second)
		db = memory.NewRedisDB()
	} else {
		log.Println("Initializing Postgres... Please wait a few seconds...")
		cmd := exec.Command("make", "-C", "../../", "postgres")
		if err := cmd.Run(); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
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
