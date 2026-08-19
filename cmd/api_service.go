package cmd

import (
	_ "UserService/docs"
	"UserService/internal/config"
	"UserService/internal/middleware"
	"UserService/internal/model"
	"UserService/internal/model/db"
	"UserService/internal/repository"
	"UserService/internal/service"
	"UserService/internal/transport/auth"
	"UserService/internal/transport/child"
	"UserService/internal/transport/referral"
	"UserService/internal/transport/user"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartApiService @title UserService API
// @version 1.0
// @description User-Service
// @host localhost:8082
// @BasePath /
func StartApiService() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	cfg := config.Load()
	psql := initDb(cfg)
	redis := model.ConnectRedis(ctx, cfg.RedisDatabase)
	stopResources := func() {
		_ = psql.Close()
		_ = redis.Close()
	}
	registerValidators()

	root := chi.NewRouter()
	setupMiddlewares(root)
	registerRoutesAndSchedulers(ctx, psql, redis, root, cfg)

	if config.IsDevStand() {
		root.Get("/swagger/*", httpSwagger.WrapHandler)
	}

	server := http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: root,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Server shutdowned with reason: %s", err)
			cancel()
		}
	}()

	ctx1, _ := context.WithTimeout(context.Background(), time.Minute*5)
	<-signals
	_ = server.Shutdown(ctx1)
	stopResources()
}

func initDb(cfg *config.Config) *sqlx.DB {
	database := db.ConnectPsql(cfg.MainDatabase)
	if err := db.StartMigration(cfg.MainDatabase); err != nil {
		log.Fatalf("Error starting migration: %s", err)
	}

	log.Println("Migration finished")
	return database
}

func registerRoutesAndSchedulers(ctx context.Context, psql *sqlx.DB, redis *redis.Client, root chi.Router, cfg *config.Config) {
	txManager := db.NewTxManager(psql, cfg.MainDatabase.MaxTransactionRetries)

	redisUserRepo := repository.NewRedisUserRepository(redis)

	userRepo := repository.NewUserRepository(psql)
	childrenRepo := repository.NewChildrenRepository(psql)
	referralRepo := repository.NewReferralsRepository(psql)

	userService := service.NewUserService(userRepo, redisUserRepo, txManager)
	childrenService := service.NewChildrenService(txManager, childrenRepo)
	referralService := service.NewReferralService(txManager, referralRepo)
	authService := service.NewAuthService(redisUserRepo)

	userController := user.NewUserController(userService)
	childrenController := child.NewChildrenController(childrenService)
	referralController := referral.NewReferralController(referralService)
	authController := auth.NewAuthController(authService)

	userController.RegisterRoutes(root)
	childrenController.RegisterRoutes(root)
	referralController.RegisterRoutes(root)
	authController.RegisterRoutes(root)

	go func() {
		//paymentCheckerScheduler.StartPaymentStatusChecker()
	}()
}

func setupMiddlewares(r chi.Router) {
	r.Use(middleware.RequestId)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

func registerValidators() {
	model.Validator = validator.New()
	_ = model.Validator.RegisterValidation("user_type", user.CheckTypeValid)
	_ = model.Validator.RegisterValidation("gender", child.CheckGenderValid)
}
