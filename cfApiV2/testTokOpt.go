package main

import (
	"os"
	"log"
	"time"

	cfLib "ns/cfApiV2/cfLibV2"
)

func main() {

	cfDir := ""
    cfDir = os.Getenv("cfDir")
    if len(cfDir) == 0 {log.Fatalf("error -- cannot read env var cfDir!")}

	nOptFilnam := cfDir + "/token/tokopt.yaml"
	log.Printf("info -- file: %s\n", nOptFilnam)

	timOpt, err := cfLib.GetTokOpt(nOptFilnam)
	if err != nil {log.Fatalf("error -- GetTokOpt: %v\n", err)}

	log.Printf("start: %s\n", timOpt.Start.Format(time.RFC1123))
	log.Printf("end:   %s\n", timOpt.End.Format(time.RFC1123))
}
