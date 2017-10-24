package processed

import (
	"database/sql"
	"log"

	"../../config"
	// MySQL driver
	_ "github.com/Go-SQL-Driver/MySQL"
)

// Store record UID
type Store struct {
	id       int64
	uniqueID string
}

var lastStore Store

// Last takes data from lastStore variable
func (p *Store) Last() (int64, string) {
	return p.id, p.uniqueID
}

// Get data from first row
func (p *Store) Get() error {

	cfg := config.New()

	connectionString := "" + cfg.HelpDeskUser + ":" + cfg.HelpDeskPassword + "@tcp(" + cfg.HelpDeskServer + ":3306)/" + cfg.HelpDeskDatabase + "?charset=utf8"
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	id := 1

	// SQL query text
	sqlText := "SELECT lastId, lastUniqueId FROM processed WHERE id = ? LIMIT 1"

	// Prepare statement
	stmt, err := db.Prepare(sqlText)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Run query
	rows, err := stmt.Query(id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Take result
	for rows.Next() {
		rows.Scan(&p.id, &p.uniqueID)
	}

	return nil
}

// Set data to first row
func (p *Store) Set(lastID int64, lastUniqueID string) error {

	cfg := config.New()

	connectionString := "" + cfg.HelpDeskUser + ":" + cfg.HelpDeskPassword + "@tcp(" + cfg.HelpDeskServer + ":3306)/" + cfg.HelpDeskDatabase + "?charset=utf8"
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	id := 1

	// SQL query text
	sqlText := "UPDATE processed SET lastId=?, lastUniqueId=? WHERE id = ? LIMIT 1"

	// Prepare statement
	stmt, err := db.Prepare(sqlText)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Run query
	_, err = stmt.Query(lastID, lastUniqueID, id)
	if err != nil {
		log.Fatal(err)
	}

	p.id = lastID
	p.uniqueID = lastUniqueID

	return nil
}
