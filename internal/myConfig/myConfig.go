package myConfig

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/pborgen/liquidityFinder/internal/myUtil"
	"github.com/rs/zerolog/log"
)

// Config holds all configuration values
type MyConfig struct {
    BaseDir string
    
    // Cache
    CacheHost     string
    CachePort     int
    CachePassword string
    CacheDB       int
    CachePoolSize     int
    CacheMinIdleConns int
    CacheMaxRetries   int

    // Use Local DB
    UseLocalDB bool

    // Postgres
    PostgresHost string
    PostgresPort int
    PostgresUser string
    PostgresPassword string
    PostgresDB string
    PostgresSSLMode string

    // Blockchain Client
    BlockchainClientUrlHttp string
    BlockchainClientPublicUrlHttp string
    BlockchainClientUrlWs string


    // Moralis
    MoralisApiKey string
    MoralisBaseUrl string

    // Transfer Event Gather
    TransferEventGatherBatchSize int
    TokenAmountModelInsertBatchSize int
    TokenAmountServiceBatchSize uint64
    TokenAmountModelBatchSize int
    TokenAmountModelWorkerCount int
    SeedPw string
}

var instance *MyConfig
var once sync.Once

func GetInstance() *MyConfig {
	once.Do(func() {
		setup()
	})
	return instance
}

func GetInstanceRefresh() *MyConfig {
	setup()
	
	return instance
}

func setup() {
    var envFilePath string

    dockerContainerEnvFile := "/opt/.env"
    // Check if we are running in a container
    if myUtil.FileExists(dockerContainerEnvFile) {
        log.Info().Msg("Running in a container")
        envFilePath = dockerContainerEnvFile
    } else {
        log.Info().Msg("Running in a non-container environment")
        envFilePath = os.Getenv("BASE_DIR") + "/.env"
    }
    
    log.Info().Msg("Loading config from " + envFilePath )

    myConfig, err := instance.load(envFilePath)
    if err != nil {
        panic(err)
    }

    log.Info().Msg("Config loaded")

    instance = myConfig
}

// Load reads the configuration from environment variables
func (c *MyConfig) load(envFile string) (*MyConfig, error) {
    // Load .env file if it exists
    if envFile != "" {
        err := godotenv.Load(envFile)
        if err != nil && !os.IsNotExist(err) {
            return nil, fmt.Errorf("error loading .env file: %w", err)
        }
    }

    config := &MyConfig{}

    config.BaseDir = getEnvString("BASE_DIR", "")
    // Cache configuration
    config.CacheHost = getEnvString("CACHE_HOST", "localhost")
    config.CachePort = getEnvInt("CACHE_PORT", 6379)
    config.CachePassword = getEnvString("CACHE_PASSWORD", "")
    config.CacheDB = getEnvInt("CACHE_DB", 0)

    config.UseLocalDB = getEnvBool("USE_LOCAL_DB", false)

    // Postgres configuration
    config.PostgresHost = getEnvString("POSTGRES_HOST", "localhost")
    config.PostgresPort = getEnvInt("POSTGRES_PORT", 5432)
    config.PostgresUser = getEnvString("POSTGRES_USER", "postgres")
    config.PostgresPassword = getEnvString("POSTGRES_PASSWORD", "")
    config.PostgresDB = getEnvString("POSTGRES_DB", "postgres")
    config.PostgresSSLMode = getEnvString("POSTGRES_SSL_MODE", "disable")

    // Blockchain Client
    config.BlockchainClientUrlHttp = getEnvString("BLOCKCHAIN_CLIENT_URL_HTTP", "")
    config.BlockchainClientPublicUrlHttp = getEnvString("BLOCKCHAIN_CLIENT_PUBLIC_URL_HTTP", "")
    config.BlockchainClientUrlWs = getEnvString("BLOCKCHAIN_CLIENT_URL_WS", "")

    // Moralis
    config.MoralisApiKey = getEnvString("MORALIS_API_KEY", "")
    config.MoralisBaseUrl = getEnvString("MORALIS_BASE_URL", "")

    config.TransferEventGatherBatchSize = getEnvInt("TRANSFER_EVENT_GATHER_BATCH_SIZE", 10)
    config.TokenAmountModelInsertBatchSize = getEnvInt("TOKEN_AMOUNT_MODEL_INSERT_BATCH_SIZE", 1000)
    config.TokenAmountServiceBatchSize = getEnvUint64("TOKEN_AMOUNT_SERVICE_BATCH_SIZE", 1000)
    config.TokenAmountModelBatchSize = getEnvInt("TOKEN_AMOUNT_MODEL_BATCH_SIZE", 100)
    config.TokenAmountModelWorkerCount = getEnvInt("TOKEN_AMOUNT_MODEL_WORKER_COUNT", 5)
    return config, nil
}

func (c *MyConfig) GetIsDevMode() bool {
    return getEnvBool("IS_DEV", true)
}

func (c *MyConfig) GetBaseDir() string {
    return os.Getenv("BASE_DIR")
}


// GetDSN returns the PostgreSQL connection string
func (c *MyConfig) GetDSN() string {
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        c.CacheHost, c.CachePort, c.CachePassword, c.CacheDB)
}

func (c *MyConfig) GetMoralisApiKey() string {
    return os.Getenv("MORALIS_API_KEY")
}

func (c *MyConfig) GetMoralisBaseUrl() string {
    return os.Getenv("MORALIS_BASE_URL")
}

func (c *MyConfig) GetTransferEventGatherBatchSize() int {
    return c.TransferEventGatherBatchSize
}


// Helper functions to get environment variables with default values
func getEnvString(key string, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value, exists := os.LookupEnv(key); exists {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
    if value, exists := os.LookupEnv(key); exists {
        if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvUint64(key string, defaultValue uint64) uint64 {
    if value, exists := os.LookupEnv(key); exists {
        if uintValue, err := strconv.ParseUint(value, 10, 64); err == nil {
            return uintValue
        }
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if value, exists := os.LookupEnv(key); exists {
        boolValue, err := strconv.ParseBool(value)
        if err == nil {
            return boolValue
        }
    }
    return defaultValue
}