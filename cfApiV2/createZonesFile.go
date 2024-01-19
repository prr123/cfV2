// createZonesFile
// program that creates a zone file.
// the program uses a DNS token json file for authentication
//
// Author: prr, azulsoftware
// Date: 13 Jan 2024
// copyright (c) 2024 prr, azulsoftware
//

package main

import (
	"os"
	"log"
	"fmt"
//	"strings"
	"context"

    cfLib "ns/cfApiV2/cfLibV2"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
    json "github.com/goccy/go-json"
    "github.com/cloudflare/cloudflare-go"

)


func main() {

    numArgs := len(os.Args)

    useStr := "./creZoneFile [/dbg]"
    helpStr := "this program creates a zone file"

    if numArgs == 2 && os.Args[1] == "help" {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage: %s\n",useStr)
        os.Exit(1)
    }

    if numArgs == 1 {
        fmt.Printf("no flags provided!")
        fmt.Printf(useStr)
        os.Exit(1)
    }

    flags := []string{"dbg"}
    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {
        log.Fatalf("ParseFlags: %v\n",err)
    }

    dbg:= false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

    cfDir := os.Getenv("cfDir")
    if len(cfDir) == 0 {log.Fatalf("error -- could not get env var: cfDir\n")}

//    rdTokFilnam := cfDir + "token/DNsToken.yaml"

    rdTokFilnam := cfDir + "/token/ZonesList.json"

	zoneFilnam := cfDir + "/zones/ZoneList.yaml"

	log.Printf("info -- token file: %s\n", rdTokFilnam)
	log.Printf("info -- zone file:  %s\n", zoneFilnam)
	log.Printf("info -- debug:  %t\n", dbg)

	bdat, err := os.ReadFile(rdTokFilnam)
	if err != nil {log.Fatalf("error -- cannot read token file: %v\n", err)}

//	if dbg {fmt.Printf("debug -- json: %s\n***********\n", string(bdat))}

	ZoneTok := cloudflare.APIToken {}
	err = json.Unmarshal(bdat, &ZoneTok)
	if err != nil {log.Fatalf("error -- cannot unmarshal token file: %v\n", err)}

    if dbg {
		fmt.Println("********** DNS Token ********")
		cfLib.PrintToken(ZoneTok)
		fmt.Println("******** End DNS Token ******")
	}

	log.Println("info -- success decoding token")


    api, err := cloudflare.NewWithAPIToken(ZoneTok.Value)
    if err != nil {log.Fatalf("error -- getting cf api obj: %v\n",err)}

	ctx := context.Background()

	resp, err :=api.VerifyAPIToken(ctx)
	if err != nil {log.Fatalf("error -- cannot verify token: %v\n", err)}
	if dbg {cfLib.PrintTokResp(&resp)}

	fmt.Printf("Resp Status: %s\n", resp.Status)
	if !(resp.Status == "active") {log.Fatalf("error -- token status is not active: %s\n", resp.Status)}

	if dbg {log.Println("debug -- verified token!")}

    zones, err := api.ListZones(ctx)
    if err != nil {log.Fatalf("error -- cannot get zones: %v\n", err)}

	log.Printf("info -- found %d zones\n", len(zones))

    acmeZones := make([]cfLib.ZoneAcme, len(zones))
	for i:=0; i< len(zones); i++ {
		acmeZones[i].Name = zones[i].Name
		acmeZones[i].Id = zones[i].ID
		acmeZones[i].Select = false
	}

	zoneFil, err := os.Create(zoneFilnam)
	if err !=nil {log.Fatalf("error -- could not create zone file: %v!\n", err)}

	err = cfLib.SaveAcmeDns(acmeZones[:], zoneFil)
	if err != nil {log.Fatalf("error -- cfLib.SaveZonesYaml: %v\n", err)}
	log.Printf("info -- success listAcmeDomains created Acme Domain File")

	log.Printf("info -- wrote zone file!")
}
