package limiter

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr string) *RedisStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStorage{client: rdb}
}

func (r *RedisStorage) Increment(ctx context.Context, key string, limit int) (int, error) {
	log.Printf("Incrementando a chave: %s", key)

	var val int64
	var err error

	// Carregar o valor de BLOCK_TIME_SECONDS do arquivo .env
	blockTimeSeconds, _ := strconv.Atoi(os.Getenv("BLOCK_TIME_SECONDS"))
	expiration := time.Duration(blockTimeSeconds) * time.Second
	// Verificar se a chave já existe
	val, err = r.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		// Chave não existe, inicializar com 1
		_, err = r.client.Set(ctx, key, 1, expiration).Result()
		if err != nil {
			log.Printf("Erro ao inicializar a chave: %s, erro: %v", key, err)
			return 0, err
		}
		val = 1
	} else if err != nil {
		log.Printf("Erro ao obter valor da chave: %s, erro: %v", key, err)
		return 0, err
	} else {
		// Chave existe, incrementar o valor
		val, err = r.client.Incr(ctx, key).Result()
		if err != nil {
			log.Printf("Erro ao incrementar a chave: %s, erro: %v", key, err)
			return 0, err
		}
	}
	// Redefinir o tempo de expiração
	_, err = r.client.Expire(ctx, key, expiration).Result()
	if err != nil {
		log.Printf("Erro ao definir expiração da chave: %s, erro: %v", key, err)
		return 0, err
	}

	getKeyExpiration(ctx, r.client, key)

	log.Printf("Valor da chave %s após incremento: %d", key, val)
	return int(val), nil
}

func (r *RedisStorage) Block(ctx context.Context, key string, duration time.Duration) error {
	return r.client.Set(ctx, key+":blocked", 1, duration).Err()
}

func (r *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	log.Printf("Verificando se a chave está bloqueada: %s", key+":blocked")
	val, err := r.client.Get(ctx, key+":blocked").Result()

	if err == redis.Nil {
		log.Printf("Chave não encontrada: %s", key+":blocked")
		return false, nil
	}

	if err != nil {
		log.Printf("Erro ao verificar a chave: %s, erro: %v", key+":blocked", err)
		return false, err
	}
	log.Printf("Valor da chave %s: %s", key+":blocked", val)
	return val == "1", err
}

func getKeyExpiration(ctx context.Context, rdb *redis.Client, key string) {
	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		log.Fatalf("Erro ao obter o tempo de expiração da chave: %v", err)
	}

	switch {
	case ttl == -1:
		fmt.Printf("A chave %s não tem tempo de expiração.\n", key)
	case ttl == -2:
		fmt.Printf("A chave %s não existe.\n", key)
	default:
		fmt.Printf("A chave %s expira em %v.\n", key, ttl)
	}
}
