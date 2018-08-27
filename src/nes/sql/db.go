package sql

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path"
)

var db *sql.DB
var counter = 0
var buffer bytes.Buffer
var addressMap = map[uint16]int{}

func Connect(filename string) error {
	dir, _ := path.Split(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	_db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}

	db = _db

	buffer.WriteString("INSERT OR REPLACE INTO address (address, isOpCode) VALUES ")

	return err
}

func Close() {
	db.Close()
}

func InitTable() {
	stmt := `create table if not exists address ( address text not null primary key, isOpCode integer);`
	Exec(stmt)
}

func Insert(addr uint16, nOpCode int) {
	counter += 1

	// dup check
	if _, ok := addressMap[addr]; ok {
		return
	}
	addressMap[addr] = 1

	buffer.WriteString(fmt.Sprintf("( \"%x\" , %d )", addr, nOpCode))

	if counter > 1000000 {
		counter = 0
		buffer.WriteString(";")
		str := buffer.String()
		buffer.Reset()
		buffer.WriteString("INSERT OR REPLACE INTO address (address, isOpCode) VALUES ")

		go func() {
			Exec(str)
		}()

	} else {
		buffer.WriteString(",")
	}
}

func Exec(sqlStmt string) {

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

}

func __nouse_Query(query string) {

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		/*
		   var id int
		   var name string
		   err = rows.Scan(&id, &name)
		   if err != nil {
		       log.Fatal(err)
		   }
		   fmt.Println(id, name)
		*/
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

}
