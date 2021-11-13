package withcache

import (
	"context"
	"fmt"
	"github.com/ppdraga/go-shortener/internal/shortener/link/datatype"
	"github.com/ppdraga/go-shortener/internal/shortener/link/withdb"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type cachedRepo struct {
	repo  *withdb.WithDB
	cache *cache.Cache
}

func (r *cachedRepo) FindLink(shortLink string) (*datatype.Link, error) {
	key := fmt.Sprintf("short_link:%s", shortLink)

	var link datatype.Link

	err := r.cache.Get(context.Background(), key, &link)

	switch err {
	case nil:
		return &link, nil

	case cache.ErrCacheMiss:
		dbLink, dbErr := r.repo.FindLink(shortLink)
		if dbErr != nil {
			return nil, dbErr
		}
		err = r.cache.Set(&cache.Item{
			Ctx:   context.Background(),
			Key:   key,
			Value: *dbLink,
			TTL:   time.Hour,
		})
		if err != nil {
			return nil, err
		}
		return dbLink, nil
	}
	return nil, err
}

func (r *cachedRepo) WriteLink(external *datatype.Link) (error, int64) {
	return r.repo.WriteLink(external)
}

func (r *cachedRepo) ReadLink(id int64) (*datatype.Link, error) {
	key := fmt.Sprintf("link:%d", id)

	var link datatype.Link

	err := r.cache.Get(context.Background(), key, &link)

	switch err {
	case nil:
		return &link, nil

	case cache.ErrCacheMiss:
		dbLink, dbErr := r.repo.ReadLink(id)
		if dbErr != nil {
			return nil, dbErr
		}
		err = r.cache.Set(&cache.Item{
			Ctx:   context.Background(),
			Key:   key,
			Value: *dbLink,
			TTL:   time.Hour,
		})
		if err != nil {
			return nil, err
		}
		return dbLink, nil
	}
	return nil, err
}

func NewCachedRepo(repo *withdb.WithDB) *cachedRepo {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rCache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &cachedRepo{
		repo:  repo,
		cache: rCache,
	}
}
