package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
	_ "url-shorter/app/docs"
	"url-shorter/app/internal/adapters/db/redis"
	"url-shorter/app/internal/config"
	http2 "url-shorter/app/internal/controllers/v1/http"
	"url-shorter/app/internal/domain/service/reduction"
	r "url-shorter/app/pkg/client/redis"
)

// TODO: Test

type App struct {
	cfg *config.Config

	router     *gin.Engine
	httpServer *http.Server
}

func NewApp(ctx context.Context, config *config.Config) (App, error) {
	log.Println("router initializing")
	router := gin.Default()

	log.Println("swagger docs initializing")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	iDb, _ := strconv.Atoi(config.Redis.Database)

	redisConfig := r.NewRedisConfig(config.Redis.Host, config.Redis.Port, iDb)
	client := r.NewClient(ctx, redisConfig)
	storage := redis.NewReductionStorage(config, client)
	service := reduction.NewReductionService(storage)
	handler := http2.NewReductionHandler(service)
	handler.Register(router)

	return App{
		cfg:    config,
		router: router,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, _ := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startHTTP()
	})

	return grp.Wait()
}

func (a *App) startHTTP() error {
	log.Println("HTTP Server initializing")

	log.Println(fmt.Sprintf("%s:%d", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	if err != nil {
		log.Fatalln("failed to create listener")
	}

	a.httpServer = &http.Server{
		Handler:      a.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Println("server shutdown")
		default:
			log.Fatal(err)
		}
	}

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return err
}
