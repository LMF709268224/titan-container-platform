package dao

import (
	"database/sql"
	"fmt"
	"time"

	"titan-container-platform/config"

	_ "github.com/go-sql-driver/mysql"
	logging "github.com/ipfs/go-log/v2"
	"github.com/jmoiron/sqlx"
)

var log = logging.Logger("db")

// mDB reference to database
var mDB *sqlx.DB

const (
	maxOpenConnections = 60
	connMaxLifetime    = 120
	maxIdleConnections = 30
	connMaxIdleTime    = 20
)

const (
	userInfoTable     = "users"
	orderInfoTable    = "orders"
	userClaimsTable   = "user_claims"
	hourlyQuotasTable = "hourly_quotas"
)

// ErrNoRow is returned when no matching row is found in the database.
var ErrNoRow = fmt.Errorf("no matching row found")

// Init initializes the database connection with the provided configuration.
func Init(cfg *config.Config) error {
	if cfg.DatabaseURL == "" {
		return fmt.Errorf("database url not setup")
	}

	db, err := sqlx.Connect("mysql", cfg.DatabaseURL)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(maxOpenConnections)
	db.SetConnMaxLifetime(connMaxLifetime * time.Second)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	mDB = db

	initTables()

	return nil
}

// initTables initializes data tables.
func initTables() error {
	doExec()

	// init table
	tx, err := mDB.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		err = tx.Rollback()
		if err != nil && err != sql.ErrTxDone {
			log.Errorf("InitTables Rollback err:%s", err.Error())
		}
	}()

	// Execute table creation statements
	tx.MustExec(fmt.Sprintf(cUserTable, userInfoTable))
	tx.MustExec(fmt.Sprintf(cOrderTable, orderInfoTable))
	tx.MustExec(fmt.Sprintf(cUserClaimsTable, userClaimsTable))
	tx.MustExec(fmt.Sprintf(cHourlyQuotasTable, hourlyQuotasTable))

	return tx.Commit()
}

func doExec() {
	// _, err := DB.Exec(fmt.Sprintf("ALTER TABLE %s CHANGE area_id area_id       VARCHAR(256)   DEFAULT ''", onlineCountTable))
	// if err != nil {
	// 	log.Errorf("InitTables doExec err:%s", err.Error())
	// }
	// _, err := DB.Exec(fmt.Sprintf("ALTER TABLE %s DROP COLUMN cpu_cores ;", projectInfoTable))
	// if err != nil {
	// 	log.Errorf("InitTables doExec err:%s", err.Error())
	// }
	// _, err := mDB.Exec(fmt.Sprintf("ALTER TABLE %s ADD price        INT           DEFAULT 0", orderInfoTable))
	// if err != nil {
	// 	log.Errorf("InitTables doExec err:%s", err.Error())
	// }
}
