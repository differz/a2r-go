package main

import (
	"fmt"
	"log"

	"../../config"
	"../../data"
	"../../data/phonebook"
	"../../manager"
	_ "github.com/Go-SQL-Driver/MySQL"
	_ "github.com/lib/pq"
)

func main() {

	debug := false

	// 1) Get config
	/////////////////////////////////////////////////////////////////////////////////////////
	cfg := config.New()
	err := cfg.Init()

	if err != nil {
		log.Fatal(err)
	}

	if debug {
		fmt.Println(cfg)
	}

	// 2) Get phones & projects
	/////////////////////////////////////////////////////////////////////////////////////////
	phonebook.Get()

	if debug {
		pt := data.NewProjectTable()
		for index := 0; index < len(*pt); index++ {
			num1 := (*pt)[index].Num
			pID1 := (*pt)[index].ProjectID
			fmt.Println(num1, pID1)
		}
	}

	// 3) Get info via manager CID or CDR records
	/////////////////////////////////////////////////////////////////////////////////////////

	manager := manager.New()

	store, err := manager.Create()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	// ^
	// | cdr.Get()

	if false {
		fmt.Println(manager)
		fmt.Println(store)
	}

}
