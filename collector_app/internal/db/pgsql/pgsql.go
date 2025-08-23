package pgsql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mbydanov/collector/collector_app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PgSQLInterface interface {
	Db() *gorm.DB
	GetDB() *sql.DB
	GetError() error
	AutoMigrate()
}

type PgSQL struct {
	db  *gorm.DB
	err error
}

func NewPgSQL(connStr *ConnStringBuilder, schemaName *string) PgSQLInterface {

	var schemaName_ string = ""
	if schemaName != nil {
		schemaName_ = *schemaName
	}

	db, err := gorm.Open(postgres.Open(connStr.GetString()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: schemaName_ + ".",
		},
	})
	if err != nil {
		return &PgSQL{db: nil, err: err}
	}

	return &PgSQL{db: db, err: err}
}

func (pg *PgSQL) GetError() error {
	return pg.err
}

func (pg *PgSQL) Db() *gorm.DB {
	return pg.db
}

func (pg *PgSQL) GetDB() *sql.DB {
	db, err := pg.db.DB()
	if err != nil {
		pg.err = err
	}
	return db
}

func (pg *PgSQL) AutoMigrate() {
	err := pg.db.AutoMigrate(&models.CoinStat{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}
}

type ConnStringBuilder struct {
	connString string
	err        error
}

func NewConnStr() *ConnStringBuilder {
	return &ConnStringBuilder{}
}

func (c *ConnStringBuilder) GetString() string {
	return c.connString
}

func (c *ConnStringBuilder) GetError() error {
	return c.err
}

func (c *ConnStringBuilder) Host(host string) *ConnStringBuilder {
	if c.err == nil {
		if host != "" {
			c.connString += " host=" + host
		} else if host == "" {
			c.err = fmt.Errorf("Host is empty")
		}
	}
	return c
}

func (c *ConnStringBuilder) Port(port string) *ConnStringBuilder {
	if c.err == nil {
		if port != "" {
			c.connString += " port=" + port
		} else if port == "" {
			c.err = fmt.Errorf("Port is empty")
		}
	}
	return c
}

func (c *ConnStringBuilder) User(user string) *ConnStringBuilder {
	if c.err == nil {
		if user != "" {
			c.connString += " user=" + user
		} else if user == "" {
			c.err = fmt.Errorf("User is empty")
		}
	}
	return c
}

func (c *ConnStringBuilder) Password(password string) *ConnStringBuilder {
	if c.err == nil {
		if password != "" {
			c.connString += " password=" + password
		} else if password == "" {
			c.err = fmt.Errorf("Password is empty")
		}
	}
	return c
}

func (c *ConnStringBuilder) Dbname(dbname string) *ConnStringBuilder {
	if c.err == nil {
		if dbname != "" {
			c.connString += " dbname=" + dbname
		} else if dbname == "" {
			c.err = fmt.Errorf("Dbname is empty")
		}
	}
	return c
}

func (c *ConnStringBuilder) Sslmode(sslmode string) *ConnStringBuilder {
	if c.err == nil {
		if sslmode != "" {
			c.connString += " sslmode=" + sslmode
		} else if sslmode == "" {
			c.err = fmt.Errorf("Sslmode is empty")
		}
	}
	return c
}
