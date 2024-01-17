// listDnsRec.go
// Author: prr, azul software
// Date 12 Jan 2024
// copyright (c) 2024 prr, azul software
//
// usage: listDnsRec
// print Dns records of a domain
//

package main

import (
	"fmt"
	"log"
	"os"
	"context"
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

	useStr := "/zone=domain [/dbg]"
	helpStr := "this program prints all Dns records of a domain"

	if numArgs == 2 && os.Args[1] == "help" {
		fmt.Printf("help: %s\n", helpStr)
		fmt.Printf("usage: %s %s\n", os.Args[0], useStr)
		os.Exit(1)
	}


    if numArgs == 1 {
		fmt.Printf("no flags provided!")
		fmt.Printf("usage: %s %s\n", os.Args[0], useStr)
		os.Exit(1)
	}

    flags := []string{"dbg", "zone"}
    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {
        log.Fatalf("error -- parseFlags: %v\n",err)
    }

	dbg:= false
	_, ok := flagMap["dbg"]
	if ok {
		dbg = true
	}

	if dbg {log.Printf("debug -- debug on!\n")}

	zone := ""
    zval, ok := flagMap["zone"]
    if !ok {
        log.Fatalf("error -- /zone flag is missing!\n")
    }
	if zval.(string) == "none" {log.Fatalf("error -- zone flag has no value!\n")}

	zone = zval.(string)

	log.Printf("info -- debug: %t\n", dbg)
	log.Printf("info -- zone:  %s\n", zone)

	// read zone files
	zoneFilnam := "zones/ZoneList.yaml"
	zones, err := cfLib.ReadAcmeDns(zoneFilnam)
	if err != nil {log.Fatalf("error -- reading acm dns file: %v\n", err)}
	if dbg {log.Printf("debug -- read zone file\n")}

	// check zone against files
	zoneIdx := -1
	for i:=0; i< len(zones); i++ {
		if zones[i].Name == zone {
			zoneIdx = i
			break
		}
	}

	if zoneIdx == -1 {log.Fatalf("error -- zone: %s not found in zone file!\n", zone)}
	if dbg {log.Printf("debug -- zone found in zone file\n")}
	zoneId := zones[zoneIdx].Id
	// get Dns records for zone

	// first need to get token
	vfy := false
	rdTokFilnam := "token/DnsRead.json"
    DnsTok, err := cfLib.ReadTokenFile(rdTokFilnam, vfy)
    if err != nil {log.Fatalf("error -- ReadTokenFile: %v\n", err)}
	if dbg {log.Printf("debug -- read Dns Token\n")}
	fmt.Printf("DnsToken: %s\n", DnsTok.Value)
	// instantiate the api
    api, err := cloudflare.NewWithAPIToken(DnsTok.Value)
    if err != nil {log.Fatalf("error -- initiating api obj: %v\n",err)}

    ctx := context.Background()

    rc := cloudflare.ResourceContainer{
        Level: cloudflare.ZoneRouteLevel,
        Identifier: zoneId,
    }

    DnsPars:=cloudflare.ListDNSRecordsParams{}

    dnsRecs, _, err := api.ListDNSRecords(ctx, &rc, DnsPars)
    if err != nil {log.Fatalf("error -- api.ListDNSRecords: %v\n", err)}

	cfLib.PrintDnsRecs(&dnsRecs)
	log.Printf("info -- success\n")
}
