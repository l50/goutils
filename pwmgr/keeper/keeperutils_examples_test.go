package keeper_test

import (
	"log"

	"github.com/l50/goutils/v2/pwmgr/keeper"
)

func ExampleKeeper_CommanderInstalled() {
	k := keeper.Keeper{}
	if !k.CommanderInstalled() {
		log.Fatal("keeper commander is not installed.")
	}
}

func ExampleKeeper_LoggedIn() {
	k := keeper.Keeper{}
	if !k.LoggedIn() {
		log.Fatal("not logged into keeper vault.")
	}
}

func ExampleKeeper_RetrieveRecord() {
	k := keeper.Keeper{}
	record, err := k.RetrieveRecord("1234abcd")
	if err != nil {
		log.Fatalf("failed to retrieve record: %v", err)
	}
	log.Printf("retrieved record: %+v\n", record)
}

func ExampleKeeper_SearchRecords() {
	k := keeper.Keeper{}
	uid, err := k.SearchRecords("search term")
	if err != nil {
		log.Fatalf("failed to search records: %v", err)
	}
	log.Printf("found matching record with UID: %s\n", uid)
}
