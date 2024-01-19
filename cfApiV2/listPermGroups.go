// listPermGroups.go
// Author: prr, azul software
// Date 15 Jan 2024
// copyright (c) 2024 prr, azul software
//
// usage listPermGroups /
// program that lists all available API token permission groups
//

package main

import (
	"context"
	"fmt"
	"log"
	"os"
//	"time"
//	"strings"

//	cfLib "ns/cfApiV2/cfLib"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
//    json "github.com/goccy/go-json"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)

	useStr := "[/dbg]"
	helpStr := "this program creates a new  token with the ability to change Dns records or list zones"

	if numArgs == 2 && os.Args[1] == "help" {
		fmt.Printf("help: %s\n", helpStr)
		fmt.Printf("usage: %s %s\n",os.Args[0], useStr)
		os.Exit(1)
	}

/*
    if numArgs == 1 {
		fmt.Printf("no flags provided!")
		fmt.Printf(useStr)
		os.Exit(1)
	}
*/
    flags := []string{"dbg"}
    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {
        log.Fatalf("error parseFlags: %v\n",err)
    }

    dbg:= false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

	if dbg {log.Printf("debug on\n")}

	key:= os.Getenv("cfApi")
	if len(key) == 0 {log.Fatalf("error -- invalid key!")}

    api, err := cloudflare.New(key, "azulsoftwarevlc@gmail.com")
    if err != nil {log.Fatalf("error generating api obj: %v\n", err)}

	// Most API calls require a Context
	ctx := context.Background()
//	apiobj := apiObj.ApiObj

	ApiPermGroups, err := api.ListAPITokensPermissionGroups(ctx)
	if err != nil {log.Fatalf("error -- ListApiTokensPermGroups: %v\n", err)}

	log.Printf("ApiPermGroups: %d\n", len(ApiPermGroups))
	PrintPermGroups(ApiPermGroups)
	log.Printf("success listPermGroups!")
}

func PrintPermGroups(permGroups []cloudflare.APITokenPermissionGroups) {

	num := len(permGroups)
	fmt.Printf("************* APITokenPermissionGroups: %d ***************\n", num)
	for i:=0; i< num; i++ {
		fmt.Printf("  *** group: %d ***\n", i)
		fmt.Printf("  ID:   %s\n", permGroups[i].ID)
		fmt.Printf("  Name: %s\n", permGroups[i].Name)
		fmt.Printf("  Scopes [%d]\n", len(permGroups[i].Scopes))
		for j:=0; j<  len(permGroups[i].Scopes); j++ {
			fmt.Printf("    Scope[%d]: %s\n", j, permGroups[i].Scopes[j])
		}
	}

	fmt.Printf("************ End APITokenPermissionGroups ****************\n")

}
