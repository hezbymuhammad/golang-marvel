package main

import (
	"context"
	"log"
	"time"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
	redis "github.com/go-redis/redis/v8"

	characterRepository "github.com/hezbymuhammad/golang-marvel-demo/model/character/repository"
	characterUsecase "github.com/hezbymuhammad/golang-marvel-demo/model/character/usecase"
	characterHttpDelivery "github.com/hezbymuhammad/golang-marvel-demo/model/character/delivery/http"
)

func init() {
	viper.SetConfigFile(`config/common.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
        apiUrl := viper.GetString(`marvel_api.url`)
        publicKey := viper.GetString(`marvel_api.public_key`)
        privateKey := viper.GetString(`marvel_api.private_key`)
        httpMarvelApiTimeout := time.Duration(viper.GetInt(`marvel_api.timeout_in_sec`)) * time.Second
        cacheExpiration := time.Duration(viper.GetInt(`cache_expiration_in_sec`)) * time.Second
        httpTimeout := time.Duration(viper.GetInt(`server.timeout_in_sec`)) * time.Second
        redisHost := viper.GetString(`redis.host`)
        redisPort := viper.GetString(`redis.port`)

        e := echo.New()
	redisConn := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})


        crRead := characterRepository.NewCharacterReadRepository(redisConn)
        crWrite := characterRepository.NewCharacterWriteRepository(
                apiUrl,
                publicKey,
                privateKey,
                redisConn,
                httpMarvelApiTimeout,
                cacheExpiration,
        )
        cu := characterUsecase.NewCharacterUsecase(
                crRead,
                crWrite,
                httpTimeout,
        )
        characterHttpDelivery.NewCharacterHandler(e, cu)

        log.Println("[INFO] Warming up cache for several seconds")
        for i := 0; i <= 15; i++ {
                crWrite.StoreByPage(context.Background(), i)
        }

        log.Fatal(e.Start(viper.GetString("server.address")))
}
