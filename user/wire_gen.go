// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"cloud.google.com/go/pubsub"
	"database/sql"
	"github.com/StevanoZ/dv-shared/middleware"
	"github.com/StevanoZ/dv-shared/pubsub"
	"github.com/StevanoZ/dv-shared/s3"
	"github.com/StevanoZ/dv-shared/service"
	"github.com/StevanoZ/dv-shared/token"
	"github.com/StevanoZ/dv-shared/utils"
	"github.com/StevanoZ/dv-user/app"
	"github.com/StevanoZ/dv-user/db/user/sqlc"
	"github.com/StevanoZ/dv-user/handler"
	"github.com/StevanoZ/dv-user/service"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

// Injectors from injector.go:

func InitializedApp(route *chi.Mux, DB *sql.DB, config *shrd_utils.BaseConfig) (app.Server, error) {
	userRepo := user_db.NewUserRepo(DB)
	client, err := s3_client.Init(config)
	if err != nil {
		return nil, err
	}
	presignClient := s3_client.PreSignClient(client)
	fileSvc := s3_client.NewS3Client(client, presignClient, config)
	pubsubClient, err := pubsub_client.NewGooglePubSub(config)
	if err != nil {
		return nil, err
	}
	pubSubClient := pubsub_client.NewPubSubClient(config, pubsubClient)
	redisClient := shrd_service.NewRedisClient(config)
	cacheSvc := shrd_service.NewCacheSvc(config, redisClient)
	maker, err := shrd_token.NewPasetoMaker(config)
	if err != nil {
		return nil, err
	}
	userSvc := service.NewUserSvc(userRepo, fileSvc, pubSubClient, cacheSvc, maker, config)
	authMiddleware := shrd_middleware.NewAuthMiddleware(maker)
	userHandler := handler.NewUserHandler(userSvc, authMiddleware)
	server := app.NewServer(route, config, userHandler)
	return server, nil
}

// injector.go:

var fileset = wire.NewSet(wire.Bind(new(s3_client.S3Client), new(*s3.Client)), wire.Bind(new(s3_client.S3PreSign), new(*s3.PresignClient)), s3_client.Init, s3_client.PreSignClient, s3_client.NewS3Client)

var userSet = wire.NewSet(user_db.NewUserRepo, service.NewUserSvc, shrd_middleware.NewAuthMiddleware, handler.NewUserHandler)

var tokenSet = wire.NewSet(shrd_token.NewPasetoMaker)

var pubSubSet = wire.NewSet(wire.Bind(new(pubsub_client.GooglePubSub), new(*pubsub.Client)), pubsub_client.NewGooglePubSub, pubsub_client.NewPubSubClient)

var cacheSet = wire.NewSet(wire.Bind(new(shrd_service.RedisClient), new(*redis.Client)), shrd_service.NewRedisClient, shrd_service.NewCacheSvc)
