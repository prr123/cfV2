// createTokenFile.go
// Author: prr, azul software
// Date 12 Jan 2024
// copyright (c) 2024 prr, azul software
//
// program that creates a token json file
//


package main

import (
	"context"
	"fmt"
	"log"
	"os"
//	"time"
	"strings"

	cfLib "ns/cfApiV2/cfLibV2"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
//    json "github.com/goccy/go-json"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)

	useStr := "./createTokenFile /token=val /tokName=name /vfy [/dbg]"
	helpStr := "this program retrieves a token with name tokNam and creates a token json file"

	if numArgs == 2 && os.Args[1] == "help" {
		fmt.Printf("help: %s\n", helpStr)
		fmt.Printf("usage: %s\n",useStr)
		os.Exit(1)
	}

    if numArgs == 1 {
		fmt.Printf("no flags provided!\n")
		fmt.Printf("usage: %s\n", useStr)
		os.Exit(1)
	}

    flags := []string{"dbg", "vfy", "token", "tokName"}
    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {
        log.Fatalf("error parseFlags: %v\n",err)
    }

	token := ""
	tokval, ok := flagMap["token"]
	if !ok {
		log.Fatalf("/token flag is missing!")
	}
	if tokval.(string) == "none" {
        log.Fatalf("token value not provided!")
    }
	token = tokval.(string)


	tokNam := ""
	toknval, ok := flagMap["tokName"]
	if !ok {
		log.Fatalf("/tokName flag is missing!")
	}
	if toknval.(string) == "none" {
        log.Fatalf("token name not provided!")
    }
	tokNam = toknval.(string)

    cfDir := ""
	tokFilnam := tokNam
	if idx:= strings.Index(tokNam, ".json"); idx == -1 {tokFilnam = tokNam + ".json"}

    nTokFilnam := cfDir + "token/" + tokFilnam

	dbg:= false
	_, ok = flagMap["dbg"]
	if ok {
		dbg = true
	}

	vfy:= false
	_, ok = flagMap["vfy"]
	if ok {
		vfy = true
	}

	log.Printf("info -- token value: %s\n", token)
	log.Printf("info -- token name:  %s\n", tokNam)
	log.Printf("info -- token file:  %s\n", nTokFilnam)
	log.Printf("info -- verify %t\n", vfy)
	log.Printf("info -- debug  %t\n", dbg)

	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {log.Fatalf("error -- api client: %v\n", err)}

	ctx := context.Background()

    tokResp, err := api.VerifyAPIToken(ctx)
    if err != nil {log.Fatalf("error -- VerifyApiToken: %v\n", err)}

    if dbg {cfLib.PrintTokResp(&tokResp)}

    if tokResp.Status != "active" {log.Fatalf("invalid status returned! %s", tokResp.Status)}

	tokId := tokResp.ID

    key := "bc7d4aef3af4c6968e641c656e1771e00a0df"
	napi, err := cloudflare.New(key, "azulsoftwarevlc@gmail.com")
    if err != nil {log.Fatalf("error generating napi obj: %v\n",err)}
	if dbg {log.Printf("info -- created api from key!\n")}

	NewTok, err := napi.GetAPIToken(ctx, tokId)
	if err != nil {log.Fatalf("error -- cannot get ApiToken!\n")}
	if dbg {log.Printf("info -- received apitok\n"); cfLib.PrintToken(NewTok);}

	if dbg {
		log.Printf("debug -- retrieved token!\n")
		cfLib.PrintToken(NewTok)
	}
	NewTok.Value = token
	if dbg {
		log.Printf("debug -- assigned token!\n")
		cfLib.PrintToken(NewTok)
	}
    err = cfLib.CreateTokFile(nTokFilnam, &NewTok, dbg)
    if err != nil {log.Fatalf("error -- CreateTokFile: %v", err) }

	Tok, err := cfLib.ReadTokenFile(nTokFilnam, vfy)
    if err != nil {log.Fatalf("error -- ReadTokenFile: %v\n", err)}

	if Tok.ID != NewTok.ID {log.Fatalf("error -- token id mismatch!")}
	if Tok.Name != NewTok.Name {log.Fatalf("error -- token name mismatch!")}
	if Tok.Value != NewTok.Value {log.Fatalf("error -- token value mismatch!")}

	log.Printf("info -- created token file!\n")
}

