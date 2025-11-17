package storage

import (
    "os"
)

type Config struct {
    PostgresDSN string
    RedisAddr   string
}

func Load() Config {
    return Config{
        PostgresDSN: os.Getenv("POSTGRES_DSN"),
        RedisAddr:   os.Getenv("REDIS_ADDR"),
    }
}

type DB struct{}
type Cache struct{}

func Init(cfg Config) (*DB, *Cache, error) {
    return &DB{}, &Cache{}, nil
}

