package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

var (
	DBPath  = os.Getenv("ZSH_HISTORY_FILE")
	Table   = "history"
	Columns = []string{
		"id", "date", "dir", "command", "status", "host",
	}
)

var (
	QueryList = fmt.Sprintf("select * from %s", Table)
)

type DBHandler struct {
	dbMap *gorp.DbMap
}

type Record struct {
	ID        int    `db:"id"`
	DateTime  string `db:"date"`
	Directory string `db:"dir"`
	Command   string `db:"command"`
	Status    int    `db:"status"`
	Hostname  string `db:"host"`
}

type Records []Record

func NewDBHandler() *DBHandler {
	return &DBHandler{
		dbMap: initDb(),
	}
}

func newHistory(cmd string, status int) Record {
	return Record{
		DateTime: time.Now().Format("2006-01-02 15:04:05"),
		Directory: func() string {
			dir, err := os.Getwd()
			if err != nil {
				return ""
			}
			return dir
		}(),
		Command: cmd,
		Status:  status,
		Hostname: func() string {
			host, err := os.Hostname()
			if err != nil {
				return ""
			}
			return host
		}(),
	}
}

func (db *DBHandler) Query(query string) (Records, error) {
	var rs Records
	_, err := db.dbMap.Select(&rs, query)
	return rs, err
}

func (db *DBHandler) Columns(cols []string) (Records, error) {
	return db.Query(fmt.Sprintf("select %s from history", cols[0]))
}

func (db *DBHandler) QueryList() (Records, error) {
	return db.Query(QueryList)
}

func (db *DBHandler) Insert(cmd string, status int) error {
	h := newHistory(cmd, status)
	return db.dbMap.Insert(&h)
}

func initDb() *gorp.DbMap {
	// if DBPath == "" {
	// 	fmt.Fprintf(os.Stderr, "Please set ZSH_HISTORY_FILE\n")
	// 	return nil
	// }
	//
	// if _, err := os.Stat(DBPath); os.IsNotExist(err) {
	// 	fmt.Fprintf(os.Stderr, "%s: no such db file\n", DBPath)
	// 	return nil
	// }

	// connect to db using standard Go database/sql API
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		return nil
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'journal' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Record{}, Table).SetKeys(true, "ID")

	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		return nil
	}

	return dbmap
}
