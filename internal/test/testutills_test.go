package transport

import (
	"UserService/internal/config"
	"UserService/internal/middleware"
	"UserService/internal/model"
	"UserService/internal/model/db"
	"UserService/internal/repository"
	"UserService/internal/service"
	"UserService/internal/transport/user"
	"UserService/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"testing"
	"time"
)

var cfg *config.Config

func newTxManager(t *testing.T, database *sqlx.DB) *db.TxManager {
	t.Helper()
	return db.NewTxManager(database, cfg.MainDatabase.MaxTransactionRetries)
}

func newRedisUserRepository(t *testing.T, client *redis.Client) *repository.RedisUserRepository {
	t.Helper()
	return repository.NewRedisUserRepository(client)
}

func newUserRepository(t *testing.T, db *sqlx.DB) *repository.UserRepository {
	t.Helper()
	return repository.NewUserRepository(db)
}

func newUserService(t *testing.T, db *sqlx.DB, redisClient *redis.Client) *service.UserService {
	t.Helper()
	return service.NewUserService(newUserRepository(t, db), newRedisUserRepository(t, redisClient), newTxManager(t, db))
}

func newUserController(t *testing.T, db *sqlx.DB, redisClient *redis.Client) *user.UserController {
	t.Helper()
	return user.NewUserController(newUserService(t, db, redisClient))
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
}

func createUserModel() *user.CreateUser {
	return &user.CreateUser{
		FirstName:     "Илья",
		LastName:      utils.StringPtr("Дмитриев"),
		MiddleName:    utils.StringPtr("Дмитриевич"),
		Email:         "test" + uuid.NewString() + "@ya.ru",
		Phone:         utils.StringPtr("79" + generateRandNumbersAsString(9, 10)),
		Notifications: true,
		Type:          model.Parent,
		ChildName:     utils.StringPtr("Алеша"),
		ChildAge:      utils.IntPtr(18),
	}
}

func updateUserModel() *user.UpdateUser {
	return &user.UpdateUser{
		FirstName:     "John",
		LastName:      "Doe",
		MiddleName:    "A",
		Email:         "test" + uuid.NewString() + "@doe.com",
		Phone:         "79" + generateRandNumbersAsString(9, 10),
		Notifications: false,
		Type:          model.ContentCreator,
		ChildName:     "Robin",
		ChildAge:      1,
		CreatedAt:     time.Now(),
	}
}

func patchUserModel() *user.PatchUser {
	return &user.PatchUser{
		FirstName:     utils.StringPtr("John"),
		LastName:      utils.StringPtr("Doe"),
		MiddleName:    utils.StringPtr("A"),
		Email:         utils.StringPtr("test" + uuid.NewString() + "@doe.com"),
		Phone:         utils.StringPtr("79" + generateRandNumbersAsString(9, 10)),
		Notifications: utils.BoolPtr(false),
		Type:          utils.StringPtr(model.ContentCreator),
		ChildName:     utils.StringPtr("Robin"),
		ChildAge:      utils.IntPtr(1),
		CreatedAt:     utils.TimePtr(time.Now()),
	}
}

func generateRandNumbersAsString(count int, rightBoundExclude int) string {
	if count <= 0 {
		count = 1
	}

	var str string
	for i := 0; i < count; i++ {
		str += strconv.Itoa(rand.IntN(rightBoundExclude))
	}
	return str
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

type TestControllers struct {
	userController *user.UserController
}

func configureEnvironment(t *testing.T) (chi.Router, *TestControllers) {
	t.Helper()

	setUpConfig()
	db := connectDB(t)
	redis := connectRedis(t)
	t.Cleanup(func() { _ = db.Close() })

	controller := newUserController(t, db, redis)

	r := chi.NewRouter()
	r.Use(middleware.RequestId)
	controller.RegisterRoutes(r)
	regCustomValidators()

	return r, &TestControllers{userController: controller}
}
