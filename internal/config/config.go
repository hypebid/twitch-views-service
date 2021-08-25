package config

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Constants struct {
	Port            string
	LogLevel        string
	DatabaseURL     string
	ServiceName     string
	ReleaseDate     string
	ReleaseSlug     string
	ReleaseVersion  string
	HashSecret      string
	MetadataKeyList string
	MetadataHashKey string
}

type PsqlInstance struct {
	DB *gorm.DB
}

type Config struct {
	Constants
	Psql PsqlInstance
	Log  *logrus.Logger
}

var ServiceConfig Config

func NewServiceConfig() (*Config, error) {
	c := Config{}
	// Load constants
	c.Constants = Constants{
		Port:            os.Getenv("PORT"),
		LogLevel:        os.Getenv("LOG_LEVEL"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		ServiceName:     os.Getenv("SERVICE_NAME"),
		ReleaseDate:     os.Getenv("RELEASE_DATE"),
		ReleaseSlug:     os.Getenv("RELEASE_SLUG"),
		ReleaseVersion:  os.Getenv("RELEASE_VERSION"),
		HashSecret:      os.Getenv("HASH_SECRET"),
		MetadataKeyList: os.Getenv("METADATA_KEY_LIST"),
		MetadataHashKey: os.Getenv("METADATA_HASH_KEY"),
	}

	// create logger
	c.Log = logrus.New()
	c.Log.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	c.Log.SetOutput(os.Stdout)
	logLvl, err := logrus.ParseLevel(c.Constants.LogLevel)
	if err != nil {
		logLvl = 4
	}
	c.Log.SetLevel(logLvl)

	// Gorm setup
	c.Log.Println(fmt.Sprintf("Connecting to db: %v", getDSN(c.Constants.DatabaseURL)))
	database, err := sql.Open("postgres", getDSN(c.Constants.DatabaseURL))
	if err != nil {
		c.Log.Error("sql connection error")
		return nil, err
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}), &gorm.Config{})
	if err != nil {
		c.Log.Error("gorm connection error")
		return nil, err
	}
	c.Psql.DB = gormDB

	return &c, nil
}

func getDSN(url string) string {
	var host string
	var user string
	var password string
	var dbname string
	var port string

	s1 := strings.Split(url, "://")
	s2 := strings.Split(s1[1], ":")
	user = s2[0]
	s3 := strings.Split(s2[1], "@")
	password = s3[0]
	host = s3[1]
	s4 := strings.Split(s2[2], "/")
	port = s4[0]
	dbname = s4[1]

	return fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", host, user, password, dbname, port)
}
