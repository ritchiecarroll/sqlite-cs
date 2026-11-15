// sqlitedemo-file: the sqlite-cs demo over a database FILE. It is not part of modernc.org/sqlite. For each journal mode
// (SQLite's default, then WAL) it creates a database file in a fresh temporary directory, creates a table, inserts rows
// and closes the database; then it reopens the same file, checks it, queries it, closes it and removes it. It prints
// only lines that do not depend on the machine (no paths, no times), so the native and converted outputs can be compared.
package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func fail(step string, err error) {
	fmt.Println(step+":", err)
	os.Exit(1)
}

// open opens the database file through database/sql with a single connection, so every statement of a phase runs on
// the same SQLite connection.
func open(path string) *sql.DB {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		fail("open", err)
	}
	db.SetMaxOpenConns(1)
	return db
}

// create is phase one: a new database file, the table and its rows, then a close that writes everything to the file.
func create(path, mode string) {
	db := open(path)
	if mode == "wal" {
		var got string
		if err := db.QueryRow("pragma journal_mode=wal").Scan(&got); err != nil {
			fail("journal_mode", err)
		}
		fmt.Printf("[%s] set journal_mode=%s\n", mode, got)
	}
	for _, q := range []string{
		"create table t(id integer primary key, name text, score real)",
		"insert into t(name, score) values ('alpha', 1.5), ('beta', 2.25), ('gamma', -3)",
		"create index t_name on t(name)",
	} {
		if _, err := db.Exec(q); err != nil {
			fail("exec", err)
		}
	}
	if err := db.Close(); err != nil {
		fail("close", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		fail("stat", err)
	}
	fmt.Printf("[%s] created: file present=%t, nonempty=%t\n", mode, true, info.Size() > 0)
}

// reopen is phase two: the same file opened again, checked, queried and closed.
func reopen(path, mode string) {
	db := open(path)
	var journal, integrity string
	var pages, count int
	if err := db.QueryRow("pragma journal_mode").Scan(&journal); err != nil {
		fail("journal_mode", err)
	}
	if err := db.QueryRow("pragma integrity_check").Scan(&integrity); err != nil {
		fail("integrity_check", err)
	}
	if err := db.QueryRow("pragma page_count").Scan(&pages); err != nil {
		fail("page_count", err)
	}
	if err := db.QueryRow("select count(*) from t").Scan(&count); err != nil {
		fail("count", err)
	}
	fmt.Printf("[%s] reopened: journal_mode=%s integrity=%s pages=%d rows=%d\n", mode, journal, integrity, pages, count)
	rows, err := db.Query("select id, name, score from t order by score desc")
	if err != nil {
		fail("query", err)
	}
	for rows.Next() {
		var id int
		var name string
		var score float64
		if err := rows.Scan(&id, &name, &score); err != nil {
			fail("scan", err)
		}
		fmt.Printf("[%s] %d %s %v\n", mode, id, name, score)
	}
	if err := rows.Err(); err != nil {
		fail("rows", err)
	}
	rows.Close()
	var name string
	if err := db.QueryRow("select name from t where name >= 'b' order by name limit 1").Scan(&name); err != nil {
		fail("index", err)
	}
	fmt.Printf("[%s] first name from b: %s\n", mode, name)
	if err := db.Close(); err != nil {
		fail("close", err)
	}
}

func main() {
	dir, err := os.MkdirTemp("", "sqlitecs-demo-file-")
	if err != nil {
		fail("tempdir", err)
	}
	defer os.RemoveAll(dir)

	var version string
	db := open(":memory:")
	if err := db.QueryRow("select sqlite_version()").Scan(&version); err != nil {
		fail("version", err)
	}
	db.Close()
	fmt.Println("sqlite", version)

	for _, mode := range []string{"delete", "wal"} {
		path := filepath.Join(dir, "demo-"+mode+".db")
		create(path, mode)
		reopen(path, mode)
		for _, suffix := range []string{"", "-wal", "-shm"} {
			if err := os.Remove(path + suffix); err != nil && !os.IsNotExist(err) {
				fail("remove", err)
			}
		}
		_, err := os.Stat(path)
		fmt.Printf("[%s] removed: %t\n", mode, os.IsNotExist(err))
	}
}
