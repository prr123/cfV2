// delDnsRecord.go
// Author: prr, azul software
// Date 19 Jan 2024
// copyright (c) 2024 prr, azul software
//
// usage: delDnsRecord /tpl=record [/save] [/dbg]
// program looks for template files in folder cfDir/dns
// print Dns records of a domain
//

package main

import (
	"fmt"
	"log"
	"os"
//	"context"
	"time"
	"strings"

	cfLib "ns/cfApiV2/cfLibV2"

    util "github.com/prr123/utility/utilLib"
    yaml "github.com/goccy/go-yaml"
//    json "github.com/goccy/go-json"
//	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)

	useStr := "/tpl=<tpl file> [/dbg]"
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

    flags := []string{"dbg", "tpl"}
    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {
        log.Fatalf("error -- parseFlags: %v\n",err)
    }

	dbg:= false
	_, ok := flagMap["dbg"]
	if ok {dbg = true}

	rec := ""
    zval, ok := flagMap["tpl"]
    if !ok {
        log.Fatalf("error -- /rec flag is missing!\n")
    }
	if zval.(string) == "none" {log.Fatalf("error -- rec flag has no value!\n")}

	rec = zval.(string)
	if idx := strings.Index(rec, "."); idx > -1 {log.Fatalf("error -- tpl value has a period!\n")}

	cfDir := os.Getenv("cfDir")
	tplFilnam := cfDir + "/dns/" + rec + ".yaml"

	if len(cfDir) == 0 {log.Fatalf("error -- cannot get env var cfDir!")}

	log.Printf("info -- debug: %t\n", dbg)
	log.Printf("info -- template: %s\n", rec)
	log.Printf("info -- tpl file: %s\n", tplFilnam)


	// get Dns Record Info
	RecList, err := ReadDnsTpl(tplFilnam)
	if err != nil {log.Fatalf("error -- reading Dns template file: %v\n", err)}

	for i:=0; i< len(RecList); i++ {
		PrintDnsRec(RecList[i])
	}

	// some basic checks
	zoneId := RecList[0].ZoneId
	if len(zoneId) == 0 {log.Fatalf("error -- no zone id value in DnsTemplate file!\n")}
	if len(RecList[0].ID) == 0 {log.Fatalf("error -- Dns record already has no Id!")}

	log.Printf("info -- record verification passed!\n")

	// instantiate the api
	vfy := true
	wrTokFilnam := cfDir + "/token/DnsWrite.json"
    DnsTok, err := cfLib.ReadTokenFile(wrTokFilnam, vfy)
    if err != nil {log.Fatalf("error -- cannot read TokenFile: %v\n", err)}
	log.Printf("info -- obtained Token\n")

    err = cfLib.DelDnsRecord(RecList[0], DnsTok.Value)
    if err != nil {log.Fatalf("error -- DelDnsRecords: %v\n", err)}
	log.Printf("info -- deleted dns record!\n")

	err = os.Remove(tplFilnam)
	if err != nil {log.Fatalf("error -- del tpl file: %v\n", err)}
	log.Printf("info -- deleted tpl file!\n")

	log.Printf("info -- success\n")
}


func ReadDnsTpl(filnam string) (RecList []cfLib.DnsRec, err error) {

	//open file
	fdat, err := os.ReadFile(filnam)
	if err != nil {return RecList, fmt.Errorf("Read File: %v", err)}

	// Unmarshal
	err = yaml.Unmarshal(fdat, &RecList)
	if err != nil {return RecList, fmt.Errorf("yaml Unmarshal: %v", err)}

	return RecList, nil
}

func WriteDnsTpl(filnam string, RecList []cfLib.DnsRec) (err error) {


	// Unmarshal
	fdat, err := yaml.Marshal(RecList)
	if err != nil {return fmt.Errorf("yaml Marshal: %v", err)}

	//open file
	err = os.WriteFile(filnam, fdat, 0666)
	if err != nil {return fmt.Errorf("Write File: %v", err)}

	return nil
}




func PrintDnsRec (rec cfLib.DnsRec) {

	fmt.Printf("************ Dns Record *************\n")
	fmt.Printf("Record Id:   %s\n", rec.ID)
	fmt.Printf("Record Type: %s\n", rec.Type)
	fmt.Printf("Record Name: %s\n", rec.Name)
	fmt.Printf("Record Content: %s\n", rec.Content)
	fmt.Printf("Record TTL:  %d\n", rec.TTL)
	fmt.Printf("CreatedOn:  %s\n", rec.CreatedOn.Format(time.RFC1123))
	fmt.Printf("ModifiedOn: %s\n", rec.ModifiedOn.Format(time.RFC1123))
	fmt.Printf("Zone Name:  %s\n", rec.Zone)
	fmt.Printf("Zone Id:    %s\n", rec.ZoneId)
	fmt.Printf("********** End Dns Record ***********\n")
}
