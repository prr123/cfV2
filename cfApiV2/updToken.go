// updToken.go
// Author: prr, azul software
// Date 12 Jan 2024
// copyright (c) 2024 prr, azul software
//
// usage updToken
// updated 16 Jan
// - logs started wiith error, debug, info
//

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
//	"strings"

	cfLib "ns/cfApiV2/cfLibV2"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
//    json "github.com/goccy/go-json"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)

	useStr := "./updToken /token=[dns|zones] [/dbg]"
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

    flags := []string{"token","dbg"}
    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {
        log.Fatalf("error parseFlags: %v\n",err)
    }

	token := ""
	tokval, ok := flagMap["token"]
	if !ok {
		log.Fatalf("/token flag is missing!")
	}
	switch tokval.(string) {
		case "none":
			log.Fatalf("token value not provided!")
		case "dns", "zones":
			token, _ = tokval.(string)
		default:
			log.Fatalf("invalid token: %s!", tokval.(string))
	}

	dbg:= false
	_, ok = flagMap["dbg"]
	if ok {
		dbg = true
	}


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

    nTokFilnam := cfDir + "token/" + tokFilnam

	if dbg {
		log.Printf("token file: %s\n", nTokFilnam)
		log.Printf("token:      %s\n", token)
	}

    key := "bc7d4aef3af4c6968e641c656e1771e00a0df"
    api, err := cloudflare.New(key, "azulsoftwarevlc@gmail.com")
    if err != nil {log.Fatalf("error generating api obj: %v\n",err)}

	// Most API calls require a Context
	ctx := context.Background()
//	apiobj := apiObj.ApiObj

	if dbg {fmt.Println("********************************************")}

	permGroups := make([]cloudflare.APITokenPermissionGroups, 1)

	switch token {
		case "dns": 
			permGroup :=cloudflare.APITokenPermissionGroups {
				ID: "4755a26eedb94da69e1066d98aa820be",
				Name: "DNS Write",
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
//	actStr := "com.cloudflare.api.account." + apiObj.ApiObj.AccountId
//	if dbg {fmt.Printf("Account: %s\n", actStr)}

	res := make(map[string]interface{})
	res["com.cloudflare.api.account.d0e0781201c0536742831e308ce406fb"] = "*"
//	res[actStr] = "*"

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
		Name: "Test",
		NotBefore: &startTime,
		ExpiresOn: &endTime,
		Policies: policies,
		Condition: &cond,
	}

	switch token {
	 case "dns":
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

/*
    tokResp, err := api.VerifyAPIToken(ctx)
    if err != nil {return fmt.Errorf("VerifyApiToken: %v", err)}

    if dbg {cfLib.PrintTokResp(&tokResp)}

    if tokResp.Status != "active" {return fmt.Errorf("invalid status returned! %s", tokResp.Status)}
*/
	log.Printf("info -- success creating Dns Token!")
}

