// delDnsToken
// program that reads a token file and deletes the token and its file
//
// Author: prr, azul software
// Date: 16. Jan 2024
// copyright (c) 2024 prr, azul software
//

package main


import (
    "fmt"
    "log"
    "os"
//    "time"
    "context"
//    "strings"

//    cfLib "ns/cfApiV2/cfLib"


    util "github.com/prr123/utility/utilLib"
    json "github.com/goccy/go-json"
    "github.com/cloudflare/cloudflare-go"
)



func main() {

    numArgs := len(os.Args)

    useStr := "./delToken /token=[dns|zones] [/dbg]"
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

    flags := []string{"token", "dbg"}
    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {
        log.Fatalf("error -- ParseFlags: %v\n",err)
    }

	token := ""
    val, ok := flagMap["token"]
    if !ok {
        log.Fatalf("error -- /token flag is missing!")
    }
    if val.(string) == "none" {
        log.Fatalf("error -- token value not provided!")
    }
    token, _ = val.(string)

    cfDir := ""
    tokFilnam := ""
    switch token {
        case "dns":
            tokFilnam = "DnsWrite.json"
        case "zones":
            tokFilnam = "ZonesList.json"
        default:
            log.Fatalf("error -- invalid token: %s!\n", token)
    }

    nTokFilnam := cfDir + "token/" + tokFilnam

    dbg:= false
    _, ok = flagMap["dbg"]
    if ok {
        dbg = true
    }


//    cfDir := os.Getenv("Cloudflare")
//    if len(cfDir) == 0 {return fmt.Errorf("could not get env: Cloudflare\n")}


    if dbg {
		log.Printf("debug -- token: %s\n", token)
        log.Printf("debug -- New token file: %s\n", nTokFilnam)
        log.Printf("debug -- debug: %t\n", dbg)
	}


	bdat, err := os.ReadFile(nTokFilnam)
	if err != nil {log.Fatalf("error -- opening token file: %v\n", err)}

	DnsTok := cloudflare.APIToken {}
	err = json.Unmarshal(bdat, &DnsTok)
	if err != nil {log.Fatalf("error -- unMarshal token string: %v\n", err)}

	key:= os.Getenv("cfApi")
	if len(key) == 0 {log.Fatalf("error -- no valid key!")}

	api, err := cloudflare.New(key, "azulsoftwarevlc@gmail.com")
	if err != nil {log.Fatalf("error -- generating api obj: %v\n",err)}

    // Most API calls require a Context
    ctx := context.Background()

	err = api.DeleteAPIToken(ctx, DnsTok.ID)
	if err != nil {log.Fatalf("error -- deleting token: %v\n", err)}

	err = os.Remove(nTokFilnam)
	if err != nil {log.Fatalf("error -- deleting token file: %v\n", err)}

	log.Println("info -- success deleting token!")
}
