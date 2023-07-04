/*
 * @Descripttion:
 * @Author:
 * @Date: 2023-05-05 07:19:42
 * @LastEditTime: 2023-05-08 07:28:14
 */
package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("conf")
	viper.AddConfigPath("../chatapp/config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(viper.Get("database"))

}

func InitMySQL() *gorm.DB {
	// 自定义打印日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: time.Second, //slow sql
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	DB, _ = gorm.Open(mysql.Open("root:password@tcp(127.0.0.1:3306)/chat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{Logger: newLogger})

	return DB
}

func InitRedis() {

	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

}

const (
	PublishKey = "websokcet"
)

//Publish 发布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error
	err = Red.Publish(ctx, channel, msg).Err()
	return err
}

//Subscribe 订阅Redis的消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	var err error
	sub := Red.Subscribe(ctx, channel)

	msg, err := sub.ReceiveMessage(ctx)
	fmt.Println("subscribe...", msg.Payload)
	return msg.Payload, err
}
