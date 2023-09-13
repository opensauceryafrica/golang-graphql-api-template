package redis

import (
	"blacheapi/config"
	"bytes"
	"context"
	"encoding/gob"
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
func connectRedis(cfg *config.Config) (*redis.Client, error) {
	opt, parseErr := redis.ParseURL(cfg.RedisURL)
	if parseErr != nil {
		return nil, errors.Wrap(parseErr, "Unable to parse redis url")
	}

	client := redis.NewClient(opt)
	_, pingErr := client.Ping(context.Background()).Result()
	return client, errors.Wrap(pingErr, "Unable to ping redis")
}

// New creates a new DAL instance at startup
func New(cfg *config.Config) (*RAL, error) {
	client, err := connectRedis(cfg)
	if err != nil {
		return nil, err
	}

	Ral = &RAL{
		client: client,
	}

	return Ral, nil
}

// GetSession fetches the session information for redis
// with key token
func (d RAL) GetSession(token string) (*Session, error) {

	cmd := d.client.Get(context.Background(), token)

	cmdb, err := cmd.Bytes()
	if err != nil {
		return &Session{}, err
	}

	b := bytes.NewReader(cmdb)

	var res Session

	if decodeErr := gob.NewDecoder(b).Decode(&res); decodeErr != nil {
		return &Session{}, err
	}

	return &res, nil
}
