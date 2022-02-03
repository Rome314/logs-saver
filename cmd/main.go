package main

import (
	"context"

	"emperror.dev/emperror"
	"github.com/gin-gonic/gin"
	"github.com/rome314/idkb-events/cmd/config"
	"github.com/rome314/idkb-events/internal/events"
	eventsRepository "github.com/rome314/idkb-events/internal/events/repository"
	eventsWeb "github.com/rome314/idkb-events/internal/events/web"
	"github.com/rome314/idkb-events/pkg/connections"
	"github.com/rome314/idkb-events/pkg/logging"
)

func main() {

	logger := logging.GetLogger("main", "")

	logger.Info("Preparing config...")
	cfg := config.GetConfig()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	logger.Info("Preparing DB connections...")
	logger.Info("	Preparing Postgres connection...")
	pgConn, err := connections.GetPostgresDatabase(ctx, cfg.PgConnString)
	emperror.Panic(err)

	logger.Info("	Preparing Redis connection...")
	redisConn, err := connections.GetRedisConnection(ctx, connections.RedisConfig{
		Address:  cfg.Redis.Address,
		Password: cfg.Redis.Password,
		Db:       cfg.Redis.Db,
	})
	emperror.Panic(err)

	logger.Info("Configuring internal modules...")
	pubSub, err := connections.GetRedisPubSub(ctx, redisConn.Connection)
	emperror.Panic(err)

	eventsRepo := eventsRepository.NewPostgres(logging.GetLogger("events", "repository"), pgConn)
	eventsUC := events.NewUseCase(logging.GetLogger("events", "use_case"), eventsRepo, pubSub.Sub, cfg.App)

	logger.Info("Running main listener...")
	err = eventsUC.Run(ctx)
	emperror.Panic(err)

	eventsGin := eventsWeb.NewGinDelivery(logging.GetLogger("events", "gin"), pubSub.Pub, cfg.App.EventsTopic)

	logger.Info("Configuring web server...")
	router := gin.Default()

	apiGroup := router.Group("/api")

	eventsGin.SetEndpoints(apiGroup)

	logger.Info("Running...")
	err = router.Run("0.0.0.0:" + cfg.ServerPort)
	emperror.Panic(err)

}
