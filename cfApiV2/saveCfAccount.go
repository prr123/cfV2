// saveCfAccount.go
// Author: prr, azul software
// Date 19 Jan 2024
// copyright (c) 2024 prr, azul software
//
// reads a yaml file cfDir/accountInfo.yaml
// and generates a json file: cfDir/cfAccount.json
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
    yaml "github.com/goccy/go-yaml"
    json "github.com/goccy/go-json"
	"github.com/cloudflare/cloudflare-go"
)

type AcntInfo struct {
	Name string `yaml:"Name"`
	Email string `yaml:"Email"`
}

func main() {

    numArgs := len(os.Args)

	useStr := "[/dbg]"
	helpStr := "this program retrieves the cloudflare account info."

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
		log.Fatalf("error -- ParseFlags: %v\n",err)
    }

    cfDir := os.Getenv("cfDir")
	if len(cfDir) == 0 {log.Fatalf("error -- cannot get env var cfDir!\n")}


	dbg:= false
	_, ok := flagMap["dbg"]
	if ok {	dbg = true}

    cfAcntFilnam := cfDir + "/cfAccount.json"
	acntFilnam := cfDir + "/accountInfo.yaml"

	log.Printf("info -- account file: %s\n", acntFilnam)
	log.Printf("info -- cf account file: %s\n", cfAcntFilnam)
	log.Printf("info -- debug: %t\n", dbg)


	actdat, err := os.ReadFile(acntFilnam)
	if err != nil {log.Fatalf("error -- ReadFile: %v\n", err)}

	Acnt:= AcntInfo{}
	err = yaml.Unmarshal(actdat, &Acnt)
	if err != nil {log.Fatalf("error -- Unmarshal Acnt: %v\n", err)}

	if dbg {
		log.Printf("debug -- account name:  %s\n", Acnt.Name)
		log.Printf("debug -- account email: %s\n", Acnt.Email)
	}

	key := os.Getenv("cfApi")
	if len(key) == 0 {log.Fatalf("error -- not a valid key!")}

    api, err := cloudflare.New(key, Acnt.Email)
//    api, err := cloudflare.New(key, "azulsoftwarevlc@gmail.com")
    if err != nil {log.Fatalf("error generating api obj: %v\n",err)}

	log.Printf("info -- generated api!\n")

	// Most API calls require a Context
	ctx := context.Background()

/*
	pagOpt := cloudflare.PaginationOptions{
		Page: 1,
		PerPage: 20,
	}
*/
	acntParam := cloudflare.AccountsListParams{
		Name: Acnt.Name,
//		PaginationOptions: pagOpt,
	}
//	log.Printf("acntParam: %v\n", acntParam)

	acntList, _, err := api.Accounts(ctx, acntParam)

	if err != nil {log.Fatalf("error -- getting accounts: %v\n", err)}

	if len(acntList) == 0 {log.Fatalf("error -- no account found!\n")}
	if len(acntList) >1 {log.Fatalf("error -- multiple accounts found!\n")}

//	CfAcnt := acntList[0]
	CfAcnt := cfLib.CfAccount {
		Id: acntList[0].ID,
		Name: Acnt.Name,
		Email: Acnt.Email,
		Created: acntList[0].CreatedOn,
	}

	bdat, err := json.Marshal(CfAcnt)
	if err != nil {log.Fatalf("error -- marshal account info!\n")}
//	fmt.Printf("bdat:\n%s\n", bdat)

	err = os.WriteFile(cfAcntFilnam, bdat, 0666)
	if err != nil {log.Fatalf("error -- write yaml account file!\n")}

	log.Printf("info -- success\n")
}

