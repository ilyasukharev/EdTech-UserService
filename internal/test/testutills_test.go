package transport

import (
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
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"testing"
	"time"
)

var cfg *config.Config

func newTxManager(t *testing.T, database *sqlx.DB) *db.TxManager {
	t.Helper()
	return db.NewTxManager(database, cfg.MainDatabase.MaxTransactionRetries)
}

func newRedisUserRepository(t *testing.T, client *redis.Client) *repository.UserRedisRepository {
	t.Helper()
	return repository.NewRedisUserRepository(client)
}

func newUserRepository(t *testing.T, db *sqlx.DB) *repository.UserRepository {
	t.Helper()
	return repository.NewUserRepository(db)
}

func newReferralRepository(t *testing.T, db *sqlx.DB) *repository.ReferralsRepository {
	t.Helper()
	return repository.NewReferralsRepository(db)
}

func newChildrenRepository(t *testing.T, db *sqlx.DB) *repository.ChildrenRepository {
	t.Helper()
	return repository.NewChildrenRepository(db)
}

func newUserService(t *testing.T, db *sqlx.DB, redisClient *redis.Client) *service.UserService {
	t.Helper()
	return service.NewUserService(newUserRepository(t, db), newRedisUserRepository(t, redisClient), newTxManager(t, db))
}

func newReferralService(t *testing.T, db *sqlx.DB) *service.ReferralService {
	t.Helper()
	return service.NewReferralService(newTxManager(t, db), newReferralRepository(t, db))
}

func newChildrenService(t *testing.T, db *sqlx.DB) *service.ChildrenService {
	t.Helper()
	return service.NewChildrenService(newTxManager(t, db), newChildrenRepository(t, db))
}

func newAuthService(t *testing.T, redisClient *redis.Client) *service.AuthService {
	t.Helper()
	return service.NewAuthService(newRedisUserRepository(t, redisClient))
}

func newUserController(t *testing.T, db *sqlx.DB, redisClient *redis.Client) *user.UserController {
	t.Helper()
	return user.NewUserController(newUserService(t, db, redisClient))
}

func newReferralController(t *testing.T, db *sqlx.DB) *referral.ReferralController {
	t.Helper()
	return referral.NewReferralController(newReferralService(t, db))
}

func newChildrenController(t *testing.T, db *sqlx.DB) *child.ChildrenController {
	t.Helper()
	return child.NewChildrenController(newChildrenService(t, db))
}

func newAuthController(t *testing.T, redisClient *redis.Client) *auth.AuthController {
	t.Helper()
	return auth.NewAuthController(newAuthService(t, redisClient))
}

func connectDB(t *testing.T) *sqlx.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		log.Fatal("environment variable TEST_DSN is not set")
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func connectRedis(t *testing.T) *redis.Client {
	t.Helper()

	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisDatabase.Address,
		Password: cfg.RedisDatabase.Password,
		DB:       cfg.RedisDatabase.Database,
	})
}

func assertPtrEqual[T comparable](t *testing.T, field string, got, want *T) {
	t.Helper()

	if got == nil && want == nil {
		return
	}
	if got == nil || want == nil {
		t.Fatalf("%s: expected %v, got %v", field, want, got)
	}
	if *got != *want {
		t.Fatalf("%s: expected %v, got %v", field, *want, *got)
	}
}

func assertTimePtrEqual(t *testing.T, field string, got, want *time.Time) {
	t.Helper()

	if got == nil && want == nil {
		return
	}
	if got == nil || want == nil {
		t.Fatalf("%s: expected %v, got %v", field, want, got)
	}

	g := got.UTC().Truncate(time.Microsecond)
	w := want.UTC().Truncate(time.Microsecond)

	if !g.Equal(w) {
		t.Fatalf("%s: expected %v, got %v", field, w, g)
	}
}

func regCustomValidators() {
	model.Validator = validator.New()
	_ = model.Validator.RegisterValidation("user_type", user.CheckTypeValid)
	_ = model.Validator.RegisterValidation("gender", child.CheckGenderValid)
}

func setUpConfig() {
	cfg = &config.Config{
		MainDatabase: &config.MainDatabase{
			MaxTransactionRetries: 1,
		},
		RedisDatabase: &config.RedisDatabase{
			Address:  "localhost:6379",
			Password: os.Getenv("REDIS_PASSWORD"),
			Database: 0,
		},
	}
}

type TestData struct {
	db          *sqlx.DB
	redisClient *redis.Client

	userController     *user.UserController
	referralController *referral.ReferralController
	childrenController *child.ChildrenController
	authController     *auth.AuthController
}

func configure(t *testing.T) (chi.Router, *TestData) {
	t.Helper()

	setUpConfig()
	db := connectDB(t)
	redis := connectRedis(t)
	t.Cleanup(func() {
		_ = db.Close()
		_ = redis.Close()
	})

	userController := newUserController(t, db, redis)
	referralController := newReferralController(t, db)
	childrenController := newChildrenController(t, db)
	authController := newAuthController(t, redis)

	r := chi.NewRouter()
	r.Use(middleware.RequestId)
	userController.RegisterRoutes(r)
	referralController.RegisterRoutes(r)
	childrenController.RegisterRoutes(r)
	authController.RegisterRoutes(r)
	regCustomValidators()

	return r, &TestData{
		db:                 db,
		redisClient:        redis,
		userController:     userController,
		referralController: referralController,
		childrenController: childrenController,
		authController:     authController,
	}
}
