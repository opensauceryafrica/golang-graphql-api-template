package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

const (
	SessionExpirationTime = 60 * time.Minute
)

// Redis Access Layer
type RAL struct {
	client *redis.Client
}

var Ral *RAL

// connectRedis connects to a redis instance defined by
// the REDIS_URL in the environment variable
func connectRedis(url string) (*redis.Client, error) {
	opt, parseErr := redis.ParseURL(url)
	if parseErr != nil {
		return nil, errors.Wrap(parseErr, "Unable to parse redis url")
	}

	client := redis.NewClient(opt)
	_, pingErr := client.Ping(context.Background()).Result()
	return client, errors.Wrap(pingErr, "Unable to ping redis")
}

// New creates a new DAL instance at startup
func New(url string) (*RAL, error) {
	client, err := connectRedis(url)
	if err != nil {
		return nil, err
	}

	Ral = &RAL{
		client: client,
	}

	return Ral, nil
}

// Get
func (d RAL) Get(token string) ([]byte, error) {

	return d.client.Get(context.Background(), token).Bytes()
}

// Set
func (d RAL) Set(token string, session interface{}, expiration time.Duration) error {

	_, err := d.client.Set(context.Background(), token, session, expiration).Result()

	return err
}

// Del
func (d RAL) Del(key string) error {

	_, err := d.client.Del(context.Background(), key).Result()

	return err
}
