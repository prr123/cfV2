package main

import (
	"context"
	"fmt"
	"log"
//	"os"

	"github.com/cloudflare/cloudflare-go"
)

func main() {
	// Construct a new API object using a global API key
//	api, err := cloudflare.New(os.Getenv("CLOUDFLARE_API_KEY"), os.Getenv("CLOUDFLARE_API_EMAIL"))
	// alternatively, you can use a scoped API token
	api, err := cloudflare.NewWithAPIToken("xtTY17B-67nnMSiC0YzNhrvrnkh_PAE5J8e0X13B")
	if err != nil {
		log.Fatalf("get api obj: %v\n",err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch user details on the account
	u, err := api.UserDetails(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Print user details
	fmt.Printf("user: %v\n",u)
}
