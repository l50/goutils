package keeper_test

import (
	"log"

	"github.com/l50/goutils/v2/keeper"
)

func ExampleCommanderInstalled() {
	if !keeper.CommanderInstalled() {
		log.Fatal("keeper commander is not installed.")
	}
}

func ExampleLoggedIn() {
	if !keeper.LoggedIn() {
		log.Fatal("not logged into keeper vault.")
	}
}

func ExampleRetrieveRecord() {
	record, err := keeper.RetrieveRecord("1234abcd")
	if err != nil {
		log.Fatalf("failed to retrieve record: %v", err)
	}
	log.Printf("retrieved record: %+v\n", record)
}

func ExampleSearchRecords() {
	uid, err := keeper.SearchRecords("search term")
	if err != nil {
		log.Fatalf("failed to search records: %v", err)
	}
	log.Printf("found matching record with UID: %s\n", uid)
}
