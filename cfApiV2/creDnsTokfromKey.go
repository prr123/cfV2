// creDnsTokfromKey.go
// Author: prr, azul software
// Date 12 Jan 2024
// copyright (c) 2024 prr, azul software
//
// usage creDnsToken
//

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"strings"

	cfLib "ns/cfApiV2/cfLib"

    util "github.com/prr123/utility/utilLib"
//    yaml "github.com/goccy/go-yaml"
    json "github.com/goccy/go-json"
	"github.com/cloudflare/cloudflare-go"
)

func main() {

    numArgs := len(os.Args)

	useStr := "./creDnsTokfromKey /out=file [/dbg]"
	helpStr := "this program creates a new  token with the ability to change Dns records."

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

	flags := []string{"out", "dbg"}
	flagMap, err := util.ParseFlags(os.Args, flags)
	if err != nil {
		log.Fatalf("ParseFlags: %v\n",err)
    }

	val, ok := flagMap["out"]
	if !ok {
		log.Fatalf("/out flag is missing!")
	}
	tokFilnam, _ := val.(string)
	if tokFilnam == "none" {
		log.Fatalf("token file name not provided!")
	}

    if idx := strings.Index(tokFilnam, ".json"); idx == -1 {
        tokFilnam += ".json"
    }

    cfDir := ""
    newTokenFilnam := cfDir + "token/" + tokFilnam

	dbg:= false
	_, ok = flagMap["dbg"]
	if ok {
		dbg = true
	}

	if dbg {
		log.Printf("New token file:      %s\n", newTokenFilnam)
		log.Printf("debug: %t\n", dbg)
	}

    key := "bc7d4aef3af4c6968e641c656e1771e00a0df"
    api, err := cloudflare.New(key, "azulsoftwarevlc@gmail.com")
    if err != nil {log.Fatalf("error generating api obj: %v\n",err)}

	// Most API calls require a Context
	ctx := context.Background()
//	apiobj := apiObj.ApiObj

	fmt.Println("********************************************")

	permGroup :=cloudflare.APITokenPermissionGroups {
		ID: "4755a26eedb94da69e1066d98aa820be",
		Name: "DNS Write",
		Scopes: nil,
	}

	permGroups := make([]cloudflare.APITokenPermissionGroups, 1)
	permGroups[0] = permGroup

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


	// first we need to retrieve account
	tok:= cloudflare.APIToken{
		Name: "TestDnsChange",
		NotBefore: &startTime,
		ExpiresOn: &endTime,
		Policies: policies,
	}

	NewTok, err := api.CreateAPIToken(ctx, tok)
	if err != nil {log.Fatalf("CreateApiToken: %v\n", err)}

	if dbg {cfLib.PrintToken(NewTok)}

    err = CreateTokFile(newTokenFilnam, &NewTok, dbg)
    if err != nil {log.Fatalf("CreateTokFile: %v", err) }

/*
    tokResp, err := api.VerifyAPIToken(ctx)
    if err != nil {return fmt.Errorf("VerifyApiToken: %v", err)}

    if dbg {cfLib.PrintTokResp(&tokResp)}

    if tokResp.Status != "active" {return fmt.Errorf("invalid status returned! %s", tokResp.Status)}
*/
	log.Printf("success creating Dns Token!")
}

func CreateTokFile(filnam string, Token *cloudflare.APIToken, dbg bool) (err error){

    if len(filnam) ==0 {return fmt.Errorf("no filnam provided!")}
    idx := strings.Index(filnam, ".json")
    if idx == -1 { filnam += ".json"}

//    cfDir := os.Getenv("Cloudflare")
//    if len(cfDir) == 0 {return fmt.Errorf("could not get env: Cloudflare\n")}

	cfDir := ""
//    rdTokFilnam := cfDir + "token/DNsToken.yaml"

    tokFilnam := cfDir + "token/" + filnam
    if dbg {log.Printf("token filnam: %s\n", tokFilnam)}

	if dbg {cfLib.PrintToken(*Token)}

	bdat, err := json.Marshal(*Token)
	if err != nil {return fmt.Errorf("json encode: %v\n", err)}

	if dbg {fmt.Printf("json: %s\n", string(bdat))}
//    ctx := context.Background()

	tokfil, err := os.Create(tokFilnam)
	if err != nil {return fmt.Errorf("create token file: %v", err)}
	defer tokfil.Close()

	_, err = tokfil.Write(bdat)
	if err !=nil {return fmt.Errorf("token write: %v\n", err)}

    return nil
}

