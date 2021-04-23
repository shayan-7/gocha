package internal

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/shayan-7/gocha/internal/config"
	"github.com/shayan-7/gocha/internal/data"
	"github.com/shayan-7/gocha/internal/handlers"
	"github.com/shayan-7/gocha/internal/services"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type app struct {
	engine *gin.Engine
	db     *gorm.DB
	cache  *redis.Client
	queue  string
}

func NewApp(db *gorm.DB, cacheAddr, queue string, conf *config.Config) *app {
	if conf != nil {
		return &app{
			engine: gin.Default(),
			db:     db,
			cache:  NewRedis(conf.RedisAddr),
			queue:  conf.RedisChan,
		}
	} else {
		return &app{
			engine: gin.Default(),
			db:     db,
			cache:  NewRedis(cacheAddr),
			queue:  queue,
		}
	}
}

func NewRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       1,
	})
}

func InitDatabase(user, pass, addr, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, pass, addr, dbname)
	log.Println(dsn)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (a *app) Launch(listenAddr string) error {
	// blocking call
	err := a.setupRoutes()
	if err != nil {
		return fmt.Errorf("unable to setup routes: %v", err)
	}
	return a.engine.Run(listenAddr)
}

func (a *app) setupRoutes() error {
	orderService := services.NewOrderService(a.db, a.cache, a.queue)
	orderHandler := handlers.NewOrderHandler(orderService)

	apiGroup := a.engine.Group("/api")
	orderAPI := apiGroup.Group("/order")
	orderAPI.POST("/", orderHandler.PostOrderHandler)

	return nil
}

func (a *app) Subscribe() {
	psNewMessage := a.cache.Subscribe(a.queue)
	for {
		msg, _ := psNewMessage.ReceiveMessage()
		log.Println("Received message:", msg.Payload)

		order := &data.Order{}
		json.Unmarshal([]byte(msg.Payload), order)
		fmt.Println(order)
		result := a.db.Create(&order)
		log.Println(result.Error)
	}
}
