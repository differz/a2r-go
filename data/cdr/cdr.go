package cdr

import (
	"database/sql"
	"fmt"
	"log"

	".."
	"../../config"
	"../../redminemini"
	"../processed"
)

// Get info from asterisk format
// redminemini.New(pID2, "2017-03-05", 1, "test time entry")
func Get() (data.In, error) {
	var qRow = data.QueryRow{
		Columns:   []string{"id", "dcontext", "calldate", "uniqueid", "src", "dst", "billsec"},
		Values:    make([]interface{}, 7),
		ValuePtrs: make([]interface{}, 7)}

	selectText := qRow.Init()

	cfg := config.New()

	db, err := sql.Open("postgres", "host="+cfg.AsteriskServer+" user="+cfg.AsteriskUser+" password="+cfg.AsteriskPassword+" dbname="+cfg.AsteriskDatabase+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	//#LastStore need rewrite on New()
	// Get last processed row
	lastStore := processed.Store{}
	lastStore.Get()
	//if debug {
	fmt.Println(lastStore)

	//id := 29749
	id, _ := lastStore.Last()

	// SQL query text
	sqlText := "SELECT " + selectText + " FROM cdr WHERE billsec > 0 AND id > $1 ORDER BY id"

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

	var newID int64
	var newUID string

	pt := data.NewProjectTable()

	// Take result
	for rows.Next() {
		rows.Scan(qRow.ValuePtrs...)

		clnum := qRow.ClientNumber()
		newID = qRow.GetValue("id").(int64)
		newUID = qRow.GetValue("uniqueid").(string)

		for index := 0; index < len(*pt); index++ {

			curnum := (*pt)[index].Num
			if clnum == curnum {
				cID2 := (*pt)[index].CID
				pID2 := int((*pt)[index].ProjectID)

				calldate := qRow.GetValue("calldate")
				calldate2 := fmt.Sprintf("%s", calldate)
				calldate3 := string([]rune(calldate2)[0:10])
				calldate4 := string([]rune(calldate2)[0:19])
				_ = calldate4

				billsec := int(qRow.GetValue("billsec").(int64))

				billhour := 0.01
				billhour = float64(billsec) / 3600
				if billhour < 0.07 {
					billhour = 0.07
				}

				// if debug
				fmt.Println(calldate3)

				messageFormat := "☎" + " " + curnum + " | " + cID2 + " (" + calldate4 + ")"
				timeEnt := redminemini.NewTimeEntry(pID2, calldate3, billhour, messageFormat)
				timeEnt.CreateTimeEntry()
				_ = timeEnt
				break
			}
		}

	}

	if newID > id {
		lastStore.Set(newID, newUID)
	}

	return nil, nil
}
