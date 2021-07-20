package main

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/labstack/echo"
	"github.com/spf13/viper"
	redis "github.com/go-redis/redis/v8"

	characterRepository "github.com/hezbymuhammad/golang-marvel-demo/model/article/repository"
	characterUsecase "github.com/hezbymuhammad/golang-marvel-demo/model/article/usecase"
	characterHttpDelivery "github.com/hezbymuhammad/golang-marvel-demo/model/article/delivery/http"
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
        httpMarvelApiTimeout := viper.GetInt(`marvel_api.timeout_in_sec`)
        cacheExpiration := viper.GetInt(`cache_expiration_in_sec`)
        httpTimeout := viper.GetInt(`server.timeout_in_sec`)
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
        chd := characterHttpDelivery.NewCharacterHandler(e, cu)

        log.Fatal(e.Start(viper.GetString("server.address")))
}
