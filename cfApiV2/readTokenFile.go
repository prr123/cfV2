// readTokenFile
// program that reads a DNS token json file
//
// Author: prr, azulsoftware
// Date: 13 Jan 2024
// copyright (c) 2024 prr, azulsoftware
//
// updated 16 Jan 2024 from readDnsTok
//

package main

import (
	"os"
	"log"
	"fmt"
//	"strings"
	"context"

    cfLib "ns/cfApiV2/cfLib"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
    json "github.com/goccy/go-json"
    "github.com/cloudflare/cloudflare-go"

)


func main() {

    numArgs := len(os.Args)

    useStr := "./readTokenFile /token=[zones|dns] [/vfy] [/dbg]"
    helpStr := "this program reads a Token File"

    if numArgs == 2 && os.Args[1] == "help" {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage: %s\n",useStr)
        os.Exit(1)
    }

    if numArgs == 1 {
        fmt.Printf("no flags provided!")
        fmt.Printf("usage: %s\n", useStr)
        os.Exit(1)
    }

    flags := []string{"token", "vfy", "dbg"}
    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {
        log.Fatalf("error -- ParseFlags: %v\n",err)
    }

	token :=""
    val, ok := flagMap["token"]
    if !ok {log.Fatalf("error -- /token flag is missing!")}
    if val.(string) == "none" {
        log.Fatalf("token value is not provided!")
    }
	token= val.(string)

/*
    switch token {
     case "dns":
        tok.Name = "DnsWrite"
    case "zones":
        tok.Name = "ZonesRead"
    default:
        log.Fatalf("error -- invalid token: %s\n", token)
    }
*/

    vfy:= false
    _, ok = flagMap["vfy"]
    if ok {vfy = true}

    dbg:= false
    _, ok = flagMap["dbg"]
    if ok {dbg = true}

//    cfDir := os.Getenv("Cloudflare")
//    if len(cfDir) == 0 {return fmt.Errorf("could not get env: Cloudflare\n")}

//    rdTokFilnam := cfDir + "token/DNsToken.yaml"
    cfDir := ""
    tokFilnam := ""
    switch token {
        case "dns":
            tokFilnam = "DnsWrite.json"
        case "zones":
            tokFilnam = "ZonesList.json"
        default:
            log.Fatalf("invalid token: %s!\n", token)
    }

    rdTokFilnam := cfDir + "token/" + tokFilnam

    if dbg {
		log.Printf("debug -- token:  %s\n", token)
        log.Printf("debug -- token file: %s\n", rdTokFilnam)
        log.Printf("debug -- verify: %t\n", vfy)
        log.Printf("debug -- debug:  %t\n", dbg)
    }

	bdat, err := os.ReadFile(rdTokFilnam)
	if err != nil {log.Fatalf("Open: %v\n", err)}

	if dbg {fmt.Printf("json: %s\n***********\n", string(bdat))}

	DnsTok := cloudflare.APIToken {}
	err = json.Unmarshal(bdat, &DnsTok)
	if err != nil {log.Fatalf("error -- unMarshal: %v\n", err)}

    if dbg {
		fmt.Println("********** DNS Token ********")
		cfLib.PrintToken(DnsTok)
		fmt.Println("******** End DNS Token ******")
	}
	log.Printf("info -- success reading Token File!\n")

	if !vfy {os.Exit(0)}

    api, err := cloudflare.NewWithAPIToken(DnsTok.Value)
    if err != nil {
        log.Fatalf("error -- get cf api obj: %v\n",err)
    }

	ctx := context.Background()

	resp, err :=api.VerifyAPIToken(ctx)
	if err != nil {
		log.Printf("info -- VerifyApi: %v\n", err)
		err2 := os.Remove(rdTokFilnam)
		if err2 != nil {log.Fatalf("error -- could not remove token file: %v\n", err2)}
		log.Printf("info -- removed invalid token file\n")
		os.Exit(1)
	}

	if dbg {cfLib.PrintTokResp(&resp)}

	log.Printf("info -- Verification Resp Status: %s\n", resp.Status)
	if !(resp.Status == "active") {log.Fatalf("error -- token status is not active: %s\n", resp.Status)}

	log.Printf("info -- verified token!\n")
}
