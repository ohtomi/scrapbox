package command

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

const (
	initializeSQL = `
  create table related_page (host text, project text, page text, related_page text, tag_list text, first_url text, unique(host, project, page, related_page));
  create table local_cache (host text, project text, page text, cached_at integer, unique(host, project, page));
  `
)

func initializeDB() error {
	directory := path.Join(scrapboxHome, "db")
	filepath := path.Join(directory, "data.db")

	if _, err := os.Stat(filepath); err == nil {
		return err
	}

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(initializeSQL)
	if err != nil {
		return err
	}

	return nil
}

func querySQL(statement string, parameters []interface{}, handler func(*sql.Rows) error) error {

	initializeDB()

	directory := path.Join(scrapboxHome, "db")
	filepath := path.Join(directory, "data.db")

	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(parameters...)
	if err != nil {
		return err
	}
	defer rows.Close()

	err = handler(rows)
	if err != nil {
		return err
	}

	return nil
}

func execSQL(statement string, parameters []interface{}) error {

	initializeDB()

	directory := path.Join(scrapboxHome, "db")
	filepath := path.Join(directory, "data.db")

	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()
	stmt.Exec(parameters...)

	return nil
}
