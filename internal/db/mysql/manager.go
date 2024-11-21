package mysql

import (
	"database/sql"
	"fmt"
	"time"
	"vpngigabot/internal/config"
	"vpngigabot/internal/db"
	"vpngigabot/internal/db/mysql/storages"

	_ "github.com/go-sql-driver/mysql"
)

// Manager to manage DB connection and storage
type Manager struct {
	conn *sql.DB
}

// Label for errors inside package
const label = "storage.manager"

// Connect to MySQL database
func (d *Manager) Connect(cfg config.StorageConfig) error {
	var err error
	const op = label + ".Connect"

	d.conn, err = sql.Open(
		"mysql",
		cfg.User+":"+cfg.Pass+"@tcp("+cfg.Host+":"+cfg.Port+")/"+cfg.Database+"?parseTime=true",
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// ping to make sure the connection is established
	for i := 0; i < 10; i++ {
		if err = d.conn.Ping(); err == nil {
			break
		}
		fmt.Println(err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Setup additional settings if present
	if cfg.MaxOpenConn != 0 {
		d.conn.SetMaxOpenConns(cfg.MaxOpenConn)
	}
	if cfg.MaxIdleConn != 0 {
		d.conn.SetMaxIdleConns(cfg.MaxIdleConn)
	}
	if cfg.ConnMaxLife != 0 {
		d.conn.SetConnMaxLifetime(cfg.ConnMaxLife)
	}

	return nil
}

// Disconnect from database
func (d *Manager) Disconnect() error {
	return d.conn.Close()
}

// TODO split connection management and storage factory
// TODO Implement lazy initialization and objects storing except the object creation

func (d *Manager) NewLinkStorage() db.Links {
	return storages.NewLinkStorage(d.conn)
}

func (d *Manager) NewUsersStorage() db.Users {
	return storages.NewUsersStorage(d.conn)
}
func (d *Manager) NewPaymentHistoryStorage() db.PaymentHistory {
	return storages.NewPaymentHistoryStorage(d.conn)
}

func (d *Manager) NewIpStoriesStorage() db.IpStories {
	return storages.NewIpStoriesStorage(d.conn)
}

func (d *Manager) NewPayLinkStorage() db.PayLink {
	return storages.NewPayLinkStorage(d.conn)
}
