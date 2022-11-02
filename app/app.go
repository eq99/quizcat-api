package app

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"quizcat/conf"
)

type App struct {
	DB     *gorm.DB
	Fiber  *fiber.App
	Cache  *redis.Client
	Logger *log.Logger
}

var app *App

func init() {
	app = InitApp()
}

func InitApp() *App {
	// init database
	dsn := fmt.Sprintf("user=%s password=%s  host=%s port=%s dbname=%s  sslmode=disable", conf.Conf().Get("POSTGRES_USER"), conf.Conf().Get("POSTGRES_PASSWORD"), conf.Conf().Get("POSTGRES_HOST"), conf.Conf().Get("POSTGRES_PORT"), conf.Conf().Get("POSTGRES_DB"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// init redis cache
	rdb := redis.NewClient(&redis.Options{
		Addr: conf.Conf().GetString("REDIS_URI"),
	})

	// init fiber app
	myfiber := fiber.New()

	// init log
	logger := log.Default()
	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//  return App
	return &App{
		DB:     db,
		Cache:  rdb,
		Fiber:  myfiber,
		Logger: logger,
	}
}

// get app instance
func GetApp() *App {
	return app
}

// A shortcut for externl package to get DB
func DB() *gorm.DB {
	return app.DB
}

func Fiber() *fiber.App {
	return app.Fiber
}

// A shortcut for externl package to get redis Cache
func Cache() *redis.Client {
	return app.Cache
}

func Log() *log.Logger {
	return app.Logger
}

func (app *App) Run() {
	port := conf.Conf().Get("PORT")
	host := conf.Conf().Get("HOST")

	// start web server
	app.Logger.Fatal(app.Fiber.Listen(fmt.Sprintf("%s:%s", host, port)))
}
