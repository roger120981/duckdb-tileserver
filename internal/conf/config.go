package conf

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Configuration for system
var Configuration Config

func setDefaultConfig() {
	viper.SetDefault("Server.HttpHost", "0.0.0.0")
	viper.SetDefault("Server.HttpPort", 9000)
	viper.SetDefault("Server.HttpsPort", 9001)
	viper.SetDefault("Server.TlsServerCertificateFile", "")
	viper.SetDefault("Server.TlsServerPrivateKeyFile", "")
	viper.SetDefault("Server.UrlBase", "")
	viper.SetDefault("Server.BasePath", "")
	viper.SetDefault("Server.CORSOrigins", "*")
	viper.SetDefault("Server.Debug", false)
	viper.SetDefault("Server.AssetsPath", "./assets")
	viper.SetDefault("Server.ReadTimeoutSec", 5)
	viper.SetDefault("Server.WriteTimeoutSec", 30)
	viper.SetDefault("Server.DisableUi", false)

	viper.SetDefault("Database.TableIncludes", []string{})
	viper.SetDefault("Database.TableExcludes", []string{})
	viper.SetDefault("Database.FunctionIncludes", []string{"postgisftw"})
	viper.SetDefault("Database.MaxOpenConns", 25)
	viper.SetDefault("Database.MaxIdleConns", 5)
	viper.SetDefault("Database.ConnMaxLifetime", 3600) // 1 hour in seconds
	viper.SetDefault("Database.ConnMaxIdleTime", 600)  // 10 minutes in seconds

	viper.SetDefault("Paging.LimitDefault", 10)
	viper.SetDefault("Paging.LimitMax", 1000)

	viper.SetDefault("Metadata.Title", "duckdb-tileserver")
	viper.SetDefault("Metadata.Description", "DuckDB Feature Server with Spatial Extension")

	viper.SetDefault("Website.BasemapUrl", "")

	viper.SetDefault("Cache.Enabled", true)
	viper.SetDefault("Cache.MaxItems", 10000)
	viper.SetDefault("Cache.MaxMemoryMB", 1024)
	viper.SetDefault("Cache.BrowserCacheMaxAge", 3600) // 1 hour in seconds
	viper.SetDefault("Cache.DisableApi", false)
	viper.SetDefault("Cache.ApiKey", "")
}

// Config for system
type Config struct {
	Server   Server
	Paging   Paging
	Metadata Metadata
	Database Database
	Website  Website
	Cache    Cache
}

// Server config
type Server struct {
	HttpHost                 string
	HttpPort                 int
	HttpsPort                int
	TlsServerCertificateFile string
	TlsServerPrivateKeyFile  string
	UrlBase                  string
	BasePath                 string
	CORSOrigins              string
	Debug                    bool
	AssetsPath               string
	ReadTimeoutSec           int
	WriteTimeoutSec          int
	DisableUi                bool
	TransformFunctions       []string
}

// Paging config
type Paging struct {
	LimitDefault int
	LimitMax     int
}

// Database config
type Database struct {
	DatabasePath     string
	TableIncludes    []string
	TableExcludes    []string
	FunctionIncludes []string
	MaxOpenConns     int // Maximum number of open connections to the database
	MaxIdleConns     int // Maximum number of idle connections in the pool
	ConnMaxLifetime  int // Maximum lifetime of a connection in seconds
	ConnMaxIdleTime  int // Maximum idle time of a connection in seconds
}

// Metadata config
type Metadata struct {
	Title       string //`mapstructure:"METADATA_TITLE"`
	Description string
}

type Website struct {
	BasemapUrl string
}

// Cache config
type Cache struct {
	Enabled            bool
	MaxItems           int
	MaxMemoryMB        int
	BrowserCacheMaxAge int    // Browser cache max-age in seconds
	DisableApi         bool   // Disable cache management API endpoints
	ApiKey             string // API key for cache management endpoints
}

// IsHTTPSEnabled tests whether HTTPS is enabled
func (conf *Config) IsTLSEnabled() bool {
	return conf.Server.TlsServerCertificateFile != "" && conf.Server.TlsServerPrivateKeyFile != ""
}

// InitConfig initializes the configuration from the config file
func InitConfig(configFilename string, isDebug bool) {
	// --- defaults
	setDefaultConfig()

	if isDebug {
		viper.Set("Debug", true)
	}

	viper.SetEnvPrefix(AppConfig.EnvPrefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	isExplictConfigFile := configFilename != ""
	confFile := AppConfig.Name + ".toml"
	if configFilename != "" {
		viper.SetConfigFile(configFilename)
	} else {
		viper.SetConfigName(confFile)
		viper.SetConfigType("toml")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("/config")
		viper.AddConfigPath("/etc")
	}
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		_, isConfigFileNotFound := err.(viper.ConfigFileNotFoundError)
		errrConfRead := fmt.Errorf("fatal error reading config file: %s", err)
		isUseDefaultConfig := isConfigFileNotFound && !isExplictConfigFile
		if isUseDefaultConfig {
			log.Debug(errrConfRead)
		} else {
			log.Fatal(errrConfRead)
		}
	}

	log.Infof("Using config file: %s", viper.ConfigFileUsed())
	errUnM := viper.Unmarshal(&Configuration)
	if errUnM != nil {
		log.Fatal(fmt.Errorf("fatal error decoding config file: %v", errUnM))
	}

	// Read environment variable database configuration
	// It takes precedence over config file (if any)
	// A blank value is ignored
	dbconnSrc := "config file"
	if dbPath := os.Getenv("DUCKDBTS_DATABASE_PATH"); dbPath != "" {
		Configuration.Database.DatabasePath = dbPath
		dbconnSrc = "environment variable DUCKDBTS_DATABASE_PATH"
	} else if dbPath := os.Getenv("DUCKDB_PATH"); dbPath != "" {
		// Keep backward compatibility
		log.Warn("DUCKDB_PATH environment variable is deprecated, use DUCKDBTS_DATABASE_PATH instead")
		Configuration.Database.DatabasePath = dbPath
		dbconnSrc = "environment variable DUCKDB_PATH (deprecated)"
	}

	log.Infof("Using database connection info from %v", dbconnSrc)

	// sanitize the configuration
	Configuration.Server.BasePath = strings.TrimRight(Configuration.Server.BasePath, "/")
}

func DumpConfig() {
	log.Debugf("--- Configuration ---")
	//fmt.Printf("Viper: %v\n", viper.AllSettings())
	//fmt.Printf("Config: %v\n", Configuration)
	var basemapURL = Configuration.Website.BasemapUrl
	if basemapURL == "" {
		basemapURL = "*** NO URL PROVIDED ***"
	}
	log.Debugf("  BasemapUrl = %v", basemapURL)
	log.Debugf("  TableIncludes = %v", Configuration.Database.TableIncludes)
	log.Debugf("  TableExcludes = %v", Configuration.Database.TableExcludes)
	log.Debugf("  FunctionIncludes = %v", Configuration.Database.FunctionIncludes)
	log.Debugf("  TransformFunctions = %v", Configuration.Server.TransformFunctions)
}
