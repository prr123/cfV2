// get testKey

package main

import (
    "context"
    "fmt"
    "log"
//  "os"

    "github.com/cloudflare/cloudflare-go"
)

func main() {
    // Construct a new API object using a global API key
//  api, err := cloudflare.New(os.Getenv("CLOUDFLARE_API_KEY"), os.Getenv("CLOUDFLARE_API_EMAIL"))
    // alternatively, you can use a scoped API token
	key := "bc7d4aef3af4c6968e641c656e1771e00a0df"
    api, err := cloudflare.New(key, "azulsoftwarevlc@gmail.com")
    if err != nil {
        log.Fatalf("get api obj: %v\n",err)
    }

    // Most API calls require a Context
    ctx := context.Background()

    fmt.Printf("received key\ncontext: %v\n", ctx )


	fmt.Printf("api:\n%v\n",api)
	printApi(api)
}

func printApi(api *cloudflare.API) {
	fmt.Printf("APIKey: %s\n", api.APIKey)
	fmt.Printf("APIToken: %s\n", api.APIToken)
	fmt.Printf("Debug: %t\n", api.Debug)
}
