package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASS")
	dbname   = os.Getenv("DB_NAME")
)

func connect() *sql.DB {
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	db, err := sql.Open("mysql", mysqlInfo)

	fmt.Println("connecting database")

	if err != nil {
		panic(err)
	}

	fmt.Println("database ok")

	return db
}

func GetRow(query string) ([]interface{}, int) {
	db := connect()
	defer db.Close()

	rows, err := db.Query(query)

	var r []interface{}

	if err != nil {
		log.Panicln(err)
		return r, -1
	}

	columns, err := rows.Columns()

	if err != nil {
		log.Panicln(err)
		return r, -1
	}

	result := make([]interface{}, len(columns))

	dest := make([]interface{}, len(columns))
	for i, _ := range columns {
		dest[i] = &columns[i]
	}
	for rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			log.Fatal(err)
		}

		for i, raw := range columns {
			result[i] = raw
		}
	}

	return result, 1
}

func GetRows(query string) ([]map[string]interface{}, int) {
	db := connect()
	defer db.Close()

	hsMap := make([]map[string]interface{}, 0, 0)

	rows, err := db.Query(query)

	if err != nil {
		log.Panicln(err)
		return hsMap, -1
	}

	columns, err := rows.Columns()

	if err != nil {
		log.Panicln(err)
		return hsMap, -1
	}

	array := make(map[string]interface{})

	dest := make([]interface{}, len(columns))
	for i, _ := range columns {
		dest[i] = &columns[i]
	}
	for rows.Next() {
		err := rows.Scan(dest...)
		if err != nil {
			log.Panic(err)
		}

		for i, raw := range columns {
			array[columns[i]] = raw
		}

		hsMap = append(hsMap, array)

	}

	return hsMap, 1
}

func InsertRow(query string) int64 {
	db := connect()
	defer db.Close()

	exc, err := db.Exec(query)

	if err != nil {
		log.Panicln(err)
		return -1
	}

	result, err := exc.LastInsertId()

	if err != nil {
		log.Panicln(err)
		return -1
	}

	return result
}

func UDRow(query string) int64 {
	db := connect()
	defer db.Close()

	exc, err := db.Exec(query)

	if err != nil {
		log.Panicln(err)
		return -1
	}

	result, err := exc.RowsAffected()

	if err != nil {
		log.Panicln(err)
		return -1
	}

	return result
}
