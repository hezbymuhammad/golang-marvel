package repository

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
        "log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	"github.com/hezbymuhammad/golang-marvel-demo/domain"
)

var (
	HttpClient = http.DefaultClient
)

type response struct {
	Data data `json:data`
}

type data struct {
	Results Characters `json:results`
}

type Characters []domain.Character

type CharacterWriteRepository struct {
	httpClient      *http.Client
	redisClient     redis.Cmdable
	api             string
	publicKey       string
	privateKey      string
	cacheExpiration time.Duration
}

func NewCharacterWriteRepository(api, publicKey, privateKey string, Conn redis.Cmdable, timeout, cacheExpiration time.Duration) domain.CharacterWriteRepository {
	httpClient := HttpClient
	httpClient.Timeout = timeout

	return &CharacterWriteRepository{
		httpClient:      httpClient,
		redisClient:     Conn,
		api:             api,
		publicKey:       publicKey,
		privateKey:      privateKey,
		cacheExpiration: cacheExpiration,
	}
}

func (r *CharacterWriteRepository) StoreByPage(ctx context.Context, page int) error {
	salt := uuid.New().String()
	hash := generateHash(salt, r.publicKey, r.privateKey)

	var pageNorm int
	if page > 0 {
		pageNorm = page - 1
	} else {
		pageNorm = 0
	}
        limit := 10
	offset := 10 * pageNorm

	req, err := http.NewRequest("GET", r.api+"/v1/public/characters", nil)
	if err != nil {
		return domain.ErrInternalServerError
	}

	q := req.URL.Query()
	q.Add("ts", salt)
	q.Add("apikey", r.publicKey)
	q.Add("hash", hash)
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	req.URL.RawQuery = q.Encode()

	res, err := r.httpClient.Do(req)
	defer res.Body.Close()
	if err != nil {
                log.Println("[ERROR][CharacterWriteRepository] StoreByPage httpClient.Do: " + err.Error())
		return domain.ErrInternalServerError
	}

	if res.StatusCode == http.StatusNotFound {
		return domain.ErrNotFound
	}

	if res.StatusCode != http.StatusOK {
		dump, _ := httputil.DumpResponse(res, true)
                nerr := fmt.Errorf("Error saving page %d: %q", page, dump)
                log.Println("[ERROR][CharacterWriteRepository] StoreByPage httpClient: " + nerr.Error())
		return domain.ErrNotFound
	}

	var rs response
	err = json.NewDecoder(res.Body).Decode(&rs)
	if err != nil {
                log.Println("[ERROR][CharacterWriteRepository] StoreByPage NewDecoder: " + err.Error())
		return domain.ErrInternalServerError
	}
	if len(rs.Data.Results) == 0 {
		return domain.ErrNotFound
	}

	IDs := getArrayFromCharacters(rs.Data.Results)
	err = r.storePage(ctx, IDs, pageNorm)
	if err != nil {
                log.Println("[INFO][CharacterWriteRepository] StoreByPage storePage: " + err.Error())
		return domain.ErrInternalServerError
	}

	err = r.storeCharacters(ctx, rs.Data.Results)
        if err != nil {
                log.Println("[WARNING][CharacterWriteRepository] StoreByPage storeCharacters: " + err.Error())
		return domain.ErrCacheKeyExists
        }

	return nil
}

func (r *CharacterWriteRepository) StoreByID(ctx context.Context, id int) error {
	salt := uuid.New().String()
	hash := generateHash(salt, r.publicKey, r.privateKey)

	url, _ := url.Parse(r.api + "/v1/public/characters/" + fmt.Sprint(id))

	q := url.Query()
	q.Set("ts", salt)
	q.Set("apikey", r.publicKey)
	q.Set("hash", hash)

	url.RawQuery = q.Encode()

	res, err := r.httpClient.Get(url.String())
	if err != nil {
                log.Println("[ERROR][CharacterWriteRepository] StoreByID httpClient: " + err.Error())
		return domain.ErrInternalServerError
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return domain.ErrNotFound
	}

	if res.StatusCode != http.StatusOK {
		dump, _ := httputil.DumpResponse(res, true)
                nerr := fmt.Errorf("Error saving ID %d: %q", id, dump)
                log.Println("[ERROR][CharacterWriteRepository] StoreByID httpClient: " + nerr.Error())
		return domain.ErrNotFound
	}

	var rs response
	err = json.NewDecoder(res.Body).Decode(&rs)
	if err != nil {
                log.Println("[ERROR][CharacterWriteRepository] StoreByID Decode: " + err.Error())
		return domain.ErrInternalServerError
	}
	if len(rs.Data.Results) == 0 {
		return domain.ErrNotFound
	}

	char := rs.Data.Results[0]
	err = r.storeCharacter(ctx, char)
        if err != nil {
                log.Println("[INFO][CharacterWriteRepository] StoreByID storeCharacter: " + err.Error())
                return err
        }
        return nil
}

func generateHash(salt, publicKey, privateKey string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(salt+privateKey+publicKey)))
}

func getArrayFromCharacters(chars []domain.Character) []int {
	var IDs []int

	for _, v := range chars {
		IDs = append(IDs, int(v.ID))
	}

	return IDs
}

func (r *CharacterWriteRepository) storePage(ctx context.Context, IDs []int, page int) error {
        key := "marvel-characters-page-"+fmt.Sprint(page)

        isExists, err := r.checkRedisKeyExists(ctx, key)
	if err != nil {
		return err
	}
        if isExists {
                return domain.ErrCacheKeyExists
        }

	json_data, err := json.Marshal(IDs)
	if err != nil {
		return err
	}
	_, err = r.redisClient.Set(ctx, key, string(json_data), r.cacheExpiration).Result()
        return err
}

func (r *CharacterWriteRepository) storeCharacters(ctx context.Context, chars []domain.Character) error {
        json_data, _ := json.Marshal(getArrayFromCharacters(chars))
        log.Println("[INFO] Caching character IDs: " + string(json_data))

	return Characters(chars).Each(10, func(c domain.Character, wg *sync.WaitGroup) error {
		err := r.storeCharacter(ctx, c)
		wg.Done()
		return err
	})
}

func (r *CharacterWriteRepository) storeCharacter(ctx context.Context, char domain.Character) error {
        key := "marvel-character-id-"+fmt.Sprint(char.ID)
	char.FetchedAt = time.Now()

        isExists, err := r.checkRedisKeyExists(ctx, key)
	if err != nil {
		return err
	}
        if isExists {
                return domain.ErrCacheKeyExists
        }

	json_data, err := json.Marshal(char)
	if err != nil {
		return err
	}

	_, err = r.redisClient.Set(ctx, key, string(json_data), r.cacheExpiration).Result()
        if err != nil {
                return err
        }
	return nil
}

func (r *CharacterWriteRepository) checkRedisKeyExists(ctx context.Context, str string) (bool, error) {
        val, err := r.redisClient.Exists(ctx, str).Result()
	if err != nil {
		return false, err
	}

        return val == 1, nil
}

func (cs Characters) Each(workers int, fn func(domain.Character, *sync.WaitGroup) error) error {
	var wg sync.WaitGroup
	err := make(chan error)

	for i, c := range cs {
		wg.Add(1)
		go func() {
			err <- fn(c, &wg)
		}()
		if i%workers == 0 {
			wg.Wait()
		}
	}
        wg.Wait()

        return <-err
}
