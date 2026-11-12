// sqlitedemo: the sqlite-cs demo. It is not part of modernc.org/sqlite: it opens an in-memory database through
// database/sql, creates a table, inserts rows, and prints a query, so the native and converted outputs can be compared.
package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		fmt.Println("open:", err)
		os.Exit(1)
	}
	defer db.Close()
	for _, q := range []string{
		"create table t(id integer primary key, name text, score real)",
		"insert into t(name, score) values ('alpha', 1.5), ('beta', 2.25), ('gamma', -3)",
	} {
		if _, err := db.Exec(q); err != nil {
			fmt.Println("exec:", err)
			os.Exit(1)
		}
	}
	var version string
	if err := db.QueryRow("select sqlite_version()").Scan(&version); err != nil {
		fmt.Println("version:", err)
		os.Exit(1)
	}
	fmt.Println("sqlite", version)
	rows, err := db.Query("select id, name, score from t order by score desc")
	if err != nil {
		fmt.Println("query:", err)
		os.Exit(1)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var score float64
		if err := rows.Scan(&id, &name, &score); err != nil {
			fmt.Println("scan:", err)
			os.Exit(1)
		}
		fmt.Println(id, name, score)
	}
}
