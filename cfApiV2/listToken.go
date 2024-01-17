// listToken.go
// Author: prr, azul software
// Date 12 Jan 2024
// copyright (c) 2024 prr, azul software
//
// usage updToken /
//


package main

import (
	"context"
	"fmt"
	"log"
	"os"
//	"time"
//	"strings"

	cfLib "ns/cfApiV2/cfLibV2"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
//    json "github.com/goccy/go-json"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)

	useStr := "./listToken [/dbg]"
	helpStr := "this program lists all cloudflare API tokens"

	if numArgs == 2 && os.Args[1] == "help" {
		fmt.Printf("help: %s\n", helpStr)
		fmt.Printf("usage: %s\n",useStr)
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
	if ok {
		dbg = true
	}

	if dbg {log.Printf("debug -- debug on!\n")}

	key := os.Getenv("cfApi")
	if len(key) ==0 {log.Fatalf("error -- no a valid key!")}

    api, err := cloudflare.New(key, "azulsoftwarevlc@gmail.com")
    if err != nil {log.Fatalf("error -- generating api obj: %v\n",err)}

	// Most API calls require a Context
	ctx := context.Background()
//	apiobj := apiObj.ApiObj

	if dbg {fmt.Println("********************************************")}

	TokList, err := api.APITokens(ctx)
	if err != nil {log.Fatalf("error -- cannot get ApiTokens: %v\n", err)}

	cfLib.PrintTokList(TokList)
}

func PrintTok(tokList []cloudflare.APIToken) {

	fmt.Printf("************ Tokens: %d ***************\n", len(tokList))



	fmt.Printf("********** End Token List *************\n", len(tokList))
}
