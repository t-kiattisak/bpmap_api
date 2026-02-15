package main

import (
	"fmt"

	"pbmap_api/src/internal/database"
	"pbmap_api/src/internal/delivery/http"
	"pbmap_api/src/internal/delivery/http/v1"
	"pbmap_api/src/internal/repository"
	"pbmap_api/src/internal/usecase"
	"pbmap_api/src/internal/worker"
	"pbmap_api/src/pkg/auth"
	"pbmap_api/src/pkg/config"
	"pbmap_api/src/pkg/redis"
	"pbmap_api/src/pkg/validator"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()
	db := config.NewDatabase(cfg)

	if err := database.Migrate(db); err != nil {
		panic(err)
	}

	cleanupJobs := worker.StartBackgroundJobs(cfg)
	defer cleanupJobs()

	fcmRepo, err := repository.NewFCMRepo(cfg)
	if err != nil {
		fmt.Printf("Warning: Failed to initialize FCM Repository: %v\n", err)
	}

	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		fmt.Printf("Warning: Failed to connect to Redis: %v\n", err)
	} else {
		fmt.Println("Successfully connected to Redis")
	}

	v := validator.New()

	userRepo := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, deviceRepo)

	alarmUsecase := usecase.NewAlarmUsecase(fcmRepo)
	notificationUsecase := usecase.NewNotificationUsecase(fcmRepo)

	tokenRepo := repository.NewTokenRepository(redisClient)
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	sessionRepo := repository.NewSessionRepository(db)
	tm := repository.NewTransactionManager(db)
	authUsecase := usecase.NewAuthService(userUsecase, tokenRepo, sessionRepo, tm, jwtService, cfg)

	alarmHandler := v1.NewAlarmHandler(alarmUsecase, v)
	authHandler := v1.NewAuthHandler(authUsecase, v)
	userHandler := v1.NewUserHandler(userUsecase, v, jwtService)
	notificationHandler := v1.NewNotificationHandler(notificationUsecase, v)

	handlers := &http.Handlers{
		Alarm:        alarmHandler,
		Auth:         authHandler,
		User:         userHandler,
		Notification: notificationHandler,
	}

	app := http.Router(handlers, jwtService, tokenRepo)

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}
