package mysql

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// newResolver return new resolver
func newResolver(config *Config) (res *dbresolver.DBResolver) {
	c := dbresolver.Config{
		Sources:  []gorm.Dialector{},
		Replicas: []gorm.Dialector{},
		Policy:   dbresolver.RandomPolicy{},
	}

	c.Sources = append(c.Sources, mysql.Open(config.Master))
	for _, item := range config.Slaves {
		c.Replicas = append(c.Replicas, mysql.Open(item))
	}

	res = dbresolver.Register(c)
	res.SetConnMaxIdleTime(config.IdleTimeout * time.Millisecond)
	res.SetConnMaxLifetime(config.MaxConnLifeTime * time.Millisecond)
	res.SetMaxIdleConns(config.MaxIdle)
	res.SetMaxOpenConns(config.MaxActive)

	return
}
