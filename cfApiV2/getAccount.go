// getAccount.go
// Author: prr, azul software
// Date 16 Jan 2024
// copyright (c) 2024 prr, azul software
//
// usage getAccount
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
    yaml "github.com/goccy/go-yaml"
//    json "github.com/goccy/go-json"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)

	useStr := "[/dbg]"
	helpStr := "this program retireves the cloudflare account info."

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
		log.Fatalf("ParseFlags: %v\n",err)
    }

    cfDir := ""
    acntFilnam := cfDir + "zones/accountInfo.yaml"

	dbg:= false
	_, ok := flagMap["dbg"]
	if ok {
		dbg = true
	}

	if dbg {
		log.Printf("debug -- accountfile:      %s\n",acntFilnam)
		log.Printf("debug -- debug: %t\n", dbg)
	}

    key := "bc7d4aef3af4c6968e641c656e1771e00a0df"
    api, err := cloudflare.New(key, "azulsoftwarevlc@gmail.com")
    if err != nil {log.Fatalf("error generating api obj: %v\n",err)}

	// Most API calls require a Context
	ctx := context.Background()

	acntParam := cloudflare.AccountsListParams{
		Name: "azulsoftware",
	}
	acntList, _, err := api.Accounts(ctx, acntParam)

	if err != nil {log.Fatalf("error -- getting accounts!\n")}

	if len(acntList) == 0 {log.Fatalf("error -- no account found!\n")}
	if len(acntList) >1 {log.Fatalf("error -- multiple accounta found!\n")}

	Acnt := acntList[0]

	bdat, err := yaml.Marshal(Acnt)
	if err != nil {log.Fatalf("error -- mashal account info!\n")}
	fmt.Printf("bdat:\n%s\n", bdat)

	err = os.WriteFile(acntFilnam, bdat, 0666)
	if err != nil {log.Fatalf("error -- write yaml account file!\n")}

//	log.Printf("acntList: %v\n", acntList)
//	log.Printf("info: %v\n", info)
	log.Printf("info -- success\n")
}

