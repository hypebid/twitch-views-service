package db

import "github.com/hypebid/go-micro-template/internal/config"

func PingDB(c *config.Config) error {
	ping := c.Psql.DB.Raw("SELECT * FROM information_schema.information_schema_catalog_name;")
	return ping.Error
}
