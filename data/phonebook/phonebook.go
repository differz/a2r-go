package phonebook

import (
	"database/sql"
	"log"

	".."
	"../../config"
)

// Get all phone numbers and corresponding project numbers into table
// format: "a2r:1111@tcp(opensun:3306)/asterisk?charset=utf8"
func Get() {
	var qRow = data.QueryRow{
		Columns:   []string{"id", "num", "cid", "projectID"},
		Types:     []string{"", "*", "*", ""},
		Values:    make([]interface{}, 4),
		ValuePtrs: make([]interface{}, 4)}

	selectText := qRow.Init()

	cfg := config.New()

	connectionString := "" + cfg.HelpDeskUser + ":" + cfg.HelpDeskPassword + "@tcp(" + cfg.HelpDeskServer + ":3306)/" + cfg.HelpDeskDatabase + "?charset=utf8"
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	id := 0

	// SQL query text
	sqlText := "SELECT " + selectText + " FROM cidresolve WHERE id > ?"

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

	pt := data.NewProjectTable()

	// Take result
	for rows.Next() {
		rows.Scan(qRow.ValuePtrs...)
		num := qRow.GetString("num")
		cid := qRow.GetString("cid")
		pid := qRow.GetInt64("projectID")
		if pid != 0 {
			n := new(data.ProjectTable)
			n.Num = num
			n.CID = cid
			n.ProjectID = pid
			*pt = append(*pt, *n)
		}
	}
}
