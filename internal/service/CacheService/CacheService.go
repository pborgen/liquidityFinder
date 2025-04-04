package cacheService

import (
	"bytes"
	"compress/gzip"
	"context"

	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/pborgen/liquidityFinder/internal/compression/myGzip"
	"github.com/pborgen/liquidityFinder/internal/myConfig"
	"github.com/pborgen/liquidityFinder/internal/types"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type CacheType struct {
    Name string
    Type any
    Expiration time.Duration
    Compress bool
}

var ( CacheType_PlsPairWithHighAmountOfPls = CacheType {
        Name: "PlsPairWithHighAmountOfPls",
        Type: []types.ModelPair{},
        Expiration: 2 * time.Hour,
        Compress: true,
    }
)

var ( CacheType_NonPlsPairWithHighAmountOfPls = CacheType {
        Name: "NonPlsPairWithHighAmountOfPls",
        Type: []types.ModelPair{},
        Expiration: 2 * time.Hour,
        Compress: true,
    }
)

var ( CacheType_AllNonTaxPairs = CacheType {
        Name: "AllNonTaxPairs",
        Type: []types.ModelPair{},
        Expiration: 2 * time.Hour,
        Compress: true,
    }
)

var ( CacheType_PairService_AllPairs = CacheType {
        Name: "PairService_AllPairs",
        Type: []types.ModelPair{},
        Expiration: 2 * time.Hour,
        Compress: true,
    }
)

var ( CacheType_PairService_AllPairs_WithLimit = CacheType {
    Name: "PairService_AllPairs_WithLimit",
    Type: []types.ModelPair{},
    Expiration: 2 * time.Hour,
    Compress: true,
}
)

var ( CacheType_PairService_AllHighLiquidityPairs = CacheType {
        Name: "PairService_AllHighLiquidityPairs",
        Type: []types.ModelPair{},
        Expiration: 2 * time.Hour,
        Compress: true,
    }
)

var ( CacheType_TokenAmountService_GetByTokenAddress = CacheType {
    Name: "TokenAmountService_GetByTokenAddress",
    Type: []types.ModelTokenAmount{},
    Expiration: 30 * time.Minute,
    Compress: true,
}
)

var ( CacheType_TokenAmountService_GetByOwnerAddress = CacheType {
    Name: "TokenAmountService_GetByOwnerAddress",
    Type: []types.ModelTokenAmount{},
    Expiration: 30 * time.Minute,
    Compress: true,
}
)

var ( CacheType_TransferEventService_GetAllForAddressGroupBy = CacheType {
    Name: "TransferEventService_GetAllForAddressGroupBy",
    Type: []types.TransferEventGroupBy{},
    Expiration: 30 * time.Minute,
    Compress: true,
}
)

// DragonflyService handles all Dragonfly operations
type CacheService struct {
    client *redis.Client
}

// Config holds the configuration for DragonflyDB

var instance *CacheService
var once sync.Once

func GetInstance() *CacheService {
	once.Do(func() {
        hasConnected := false

        for !hasConnected {
            cfg := myConfig.GetInstance()
            myInstance, err := newCacheService(cfg)
            if err != nil {
                log.Error().Msgf("Error connecting to cache: %v", err)
                time.Sleep(3 * time.Second)
            } else {
                instance = myInstance
                hasConnected = true
            }
        }
	})
	return instance
}

// NewDragonflyService creates a new instance of DragonflyService
func newCacheService(cfg *myConfig.MyConfig) (*CacheService, error) {
    client := redis.NewClient(&redis.Options{
        Addr:         cfg.CacheHost + ":" + strconv.Itoa(cfg.CachePort),
        Password:     cfg.CachePassword,
        DB:           cfg.CacheDB,
        WriteTimeout: 60 * time.Second,
    })

    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if _, err := client.Ping(ctx).Result(); err != nil {
        return nil, err
    }

    return &CacheService{
        client: client,
    }, nil
}

// Close closes the Dragonfly client connection
func (s *CacheService) Close() error {
    return s.client.Close()
}


func (s *CacheService) SetObject(ctx context.Context, key string, value interface{}, cacheType CacheType ) error {
    
    originalJSON, err := json.Marshal(value)
    if err != nil {
        return err
    }

    if cacheType.Compress {
        //Compress the JSON data using gzip
        compressedJSON, err := myGzip.Compress(originalJSON)
        if err != nil {
            return err
        }
    
        return s.client.Set(ctx, key, compressedJSON, cacheType.Expiration).Err()
    } else {
        return s.client.Set(ctx, key, originalJSON, cacheType.Expiration).Err()
    }
}


func GetObject[T any](ctx context.Context, key string, cacheType CacheType)  (T, error) {
	var result T

	// Get the compressed data from Redis
	rawData, err := instance.client.Get(ctx, key).Bytes()
	if err != nil {
		if err.Error() == "redis: nil" {
			return result, nil
		} else {
			return result, err
		}
	}

    if cacheType.Compress {
        // Decompress the data
        gzipReader, err := gzip.NewReader(bytes.NewReader(rawData))
        if err != nil {
            return result, err
        }
        defer gzipReader.Close()

        // Decode the JSON data
        err = json.NewDecoder(gzipReader).Decode(&result)
        if err != nil {
            return result, err
        }
    } else {
        // Decode the JSON data
        err = json.Unmarshal(rawData, &result)
        if err != nil {
            return result, err
        }
    }


	return result, nil
}

// Delete removes a key
func (s *CacheService) Delete(ctx context.Context, keys ...string) error {
    return s.client.Del(ctx, keys...).Err()
}

// HashSet stores a hash map
func (s *CacheService) HashSet(ctx context.Context, key string, values map[string]interface{}) error {
    return s.client.HSet(ctx, key, values).Err()
}

// HashGet retrieves all fields from a hash
func (s *CacheService) HashGet(ctx context.Context, key string) (map[string]string, error) {
    return s.client.HGetAll(ctx, key).Result()
}

// ListPush adds values to a list
func (s *CacheService) ListPush(ctx context.Context, key string, values ...interface{}) error {
    return s.client.LPush(ctx, key, values...).Err()
}

// ListRange retrieves a range of elements from a list
func (s *CacheService) ListRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
    return s.client.LRange(ctx, key, start, stop).Result()
}

// SetAdd adds members to a set
func (s *CacheService) SetAdd(ctx context.Context, key string, members ...interface{}) error {
    return s.client.SAdd(ctx, key, members...).Err()
}

// SetMembers retrieves all members of a set
func (s *CacheService) SetMembers(ctx context.Context, key string) ([]string, error) {
    return s.client.SMembers(ctx, key).Result()
}

// Exists checks if keys exist
func (s *CacheService) Exists(ctx context.Context, keys ...string) (bool, error) {
    n, err := s.client.Exists(ctx, keys...).Result()
    return n > 0, err
}

// TTL gets the remaining time to live of a key
func (s *CacheService) TTL(ctx context.Context, key string) (time.Duration, error) {
    return s.client.TTL(ctx, key).Result()
}

func (s *CacheService) FlushAll(ctx context.Context) error {
    return s.client.FlushAll(ctx).Err()
}
