// updToken.go
// Author: prr, azul software
// Date 12 Jan 2024
// copyright (c) 2024 prr, azul software
//
// usage updToken
// updated 16 Jan
// - logs started wiith error, debug, info
//
// 18 Jan
// read environmental var cfDir
//

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"strings"

	cfLib "ns/cfApiV2/cfLibV2"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
//    json "github.com/goccy/go-json"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)

	useStr := "./updToken /token=[dnsread|dnswrite|zones] [/opt=opt.yaml] [/vfy] [/dbg]"
	helpStr := "this program creates a new  token with the ability to change Dns records or list zones"

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

    flags := []string{"token","opt", "vfy", "dbg"}
    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {
        log.Fatalf("error -- parseFlags: %v\n",err)
    }

	optFilnam := ""
	optval, ok := flagMap["opt"]
	if ok {
		if optval.(string) == "none" {
			log.Fatalf("error -- /opt flag value is missing!\n")
		}
		if idx := strings.Index(optval.(string), "."); idx > -1 {
			log.Fatalf("error -- opt value has extension!\n")
		}
		optFilnam = optval.(string) + ".yaml"
	}


	token := ""
	tokval, ok := flagMap["token"]
	if !ok {
		log.Fatalf("error -- /token flag is missing!\n")
	}
	switch tokval.(string) {
		case "none":
			log.Fatalf("error -- token value not provided!\n")
		case "dnsread", "dnswrite", "zones":
			token, _ = tokval.(string)
		default:
			log.Fatalf("error -- invalid token: %s!\n", tokval.(string))
	}

	dbg:= false
	_, ok = flagMap["dbg"]
	if ok {	dbg = true}

	vfy:= false
	_, ok = flagMap["vfy"]
	if ok {vfy = true}

    cfDir := ""
	cfDir = os.Getenv("cfDir")
	if len(cfDir) == 0 {log.Fatalf("error -- cannot read env var cfDir!")}

	tokFilnam := ""
	switch token {
		case "dnsread":
			tokFilnam = "DnsRead.json"
		case "dnswrite":
			tokFilnam = "DnsWrite.json"
		case "zones":
			tokFilnam = "ZonesList.json"
		default:
			log.Fatalf("error -- invalid token: %s!\n", token)
	}

    nTokFilnam := cfDir + "/token/" + tokFilnam
	nOptFilnam :=""

//.Format("2005-12-30T01:02:03Z")
	timopt := cfLib.TokOpt{
		Start: time.Now().UTC().Round(time.Second),
		End: time.Now().UTC().AddDate(0,0,7).Round(time.Second),
	}

	timOpt := cfLib.TokOpt{}
	timOpt.Start = timopt.Start
//	timOpt.Start = time.Date(timopt.Start.Year(), timopt.Start.Month(), timopt.Start.Day(), 0,0,0,0, timopt.Start.Location())
	timOpt.End = time.Date(timopt.End.Year(), timopt.End.Month(), timopt.End.Day(), 0,0,0,0, timopt.End.Location())

	if len(optFilnam) > 0 {
		nOptFilnam = cfDir + "/token/" + optFilnam
		timOpt, err = cfLib.GetTokOpt(nOptFilnam)
		if err != nil {log.Fatalf("error -- GetTokOpt: %v\n", err)}
	}

	log.Printf("info -- token file: %s\n", nTokFilnam)
	log.Printf("info -- token:      %s\n", token)
	log.Printf("info -- opt file:   %s\n", nOptFilnam)
	log.Printf("info -- opt:        %s\n", optFilnam)
	log.Printf("info -- debug:      %t\n", dbg)
	log.Printf("info -- vfy:        %t\n", vfy)
	log.Printf("debug -- start: %s\n", timOpt.Start)
	log.Printf("debug -- end:   %s\n", timOpt.End)

	os.Exit(0)
	key := os.Getenv("cfApi")
	if len(key) == 0 {log.Fatalf("error -- not a valid key!")}
	if dbg {log.Printf("debug -- key: %s\n", key)}

    api, err := cloudflare.New(key, "azulsoftwarevlc@gmail.com")
    if err != nil {log.Fatalf("error -- generating api obj: %v\n",err)}

	// Most API calls require a Context
	ctx := context.Background()
//	apiobj := apiObj.ApiObj

	if dbg {fmt.Println("********************************************")}

	permGroups := make([]cloudflare.APITokenPermissionGroups, 1)

	switch token {
		case "dnsread":
			permGroup :=cloudflare.APITokenPermissionGroups {
				ID: "82e64a83756745bbbb1c9c2701bf816b",
				Name: "DnsRead",
				Scopes: nil,
			}

			permGroups[0] = permGroup

		case "dnswrite":
			permGroup :=cloudflare.APITokenPermissionGroups {
				ID: "4755a26eedb94da69e1066d98aa820be",
				Name: "DnsWrite",
				Scopes: nil,
			}

			permGroups[0] = permGroup

		case "zones":
			scopes := make([]string, 1)
			scopes[0] = "com.cloudflare.api.account.zone"
			permGroup :=cloudflare.APITokenPermissionGroups {
				ID: "c8fed203ed3043cba015a93ad1616f1f",
				Name: "Zone Read",
				Scopes: scopes,
			}

			permGroups[0] = permGroup

		default:
			log.Fatalf("error -- invalid token!")
	}

	res := make(map[string]interface{})
	res["com.cloudflare.api.account.d0e0781201c0536742831e308ce406fb"] = "*"

	policy := cloudflare.APITokenPolicies{
			Effect: "allow",
			Resources: res,
			PermissionGroups: permGroups,
		}

	policies := make([]cloudflare.APITokenPolicies, 1)
	policies[0] = policy

	startTime := time.Now().UTC().Round(time.Second)
//.Format("2005-12-30T01:02:03Z")
	endTime := time.Now().UTC().AddDate(0,2,0).Round(time.Second)
//.Format("2005-12-30T01:02:03Z")

	ipList := make([]string, 1)
	ipList[0]="89.116.30.49/32"

	ipCond := cloudflare.APITokenRequestIPCondition {
		In: ipList,
	}

	cond := cloudflare.APITokenCondition{}
	cond.RequestIP = &ipCond

	// first we need to retrieve account
	tok:= cloudflare.APIToken{
		NotBefore: &startTime,
		ExpiresOn: &endTime,
		Policies: policies,
		Condition: &cond,
	}

	switch token {
	 case "dnsread":
		tok.Name = "DnsRead"
	 case "dnswrite":
		tok.Name = "DnsWrite"
	case "zones":
		tok.Name = "ZonesRead"
	default:
		log.Fatalf("error -- invalid token: %s\n", token)
	}


	NewTok, err := api.CreateAPIToken(ctx, tok)
	if err != nil {log.Fatalf("error -- CreateApiToken: %v\n", err)}

	if dbg {cfLib.PrintToken(NewTok)}

    err = cfLib.CreateTokFile(nTokFilnam, &NewTok, dbg)
    if err != nil {log.Fatalf("error -- CreateTokFile: %v", err) }

	if !vfy {os.Exit(0)}

	// test token

    napi, err := cloudflare.NewWithAPIToken(NewTok.Value)
    if err != nil {log.Fatalf("error -- initiating api obj: %v\n",err)}

    tokResp, err := napi.VerifyAPIToken(ctx)
    if err != nil {log.Fatalf("error -- VerifyApiToken: %v", err)}

    if dbg {cfLib.PrintTokResp(&tokResp)}

    if tokResp.Status != "active" {log.Fatalf("error -- invalid status returned! %s\n", tokResp.Status)}

	log.Printf("info -- success creating Dns Token!")
}

