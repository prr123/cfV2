// creDnsRecTemplate
// Author: prr, azulsoftware
// Date: 19 Jan 2024
// copyright (c) 2024 prr, azul software
//

package main

import (
    "fmt"
    "log"
    "os"
	"strings"

    cfLib "ns/cfApiV2/cfLibV2"

    util "github.com/prr123/utility/utilLib"
    yaml "github.com/goccy/go-yaml"
//    json "github.com/goccy/go-json"
//    "github.com/cloudflare/cloudflare-go"
)


func main () {

    numArgs := len(os.Args)

    useStr := "/tpl=<template name> [/dbg]"
    helpStr := "this program creates a Dns Template file"

    if numArgs == 2 && os.Args[1] == "help" {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }


    if numArgs == 1 {
        fmt.Printf("no flags provided!\n")
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

    if dbg {log.Printf("debug -- debug on!\n")}

    tplnam := ""
    zval, ok := flagMap["tpl"]
    if !ok {
        log.Fatalf("error -- /tpl flag is missing!\n")
    }
    if zval.(string) == "none" {log.Fatalf("error -- tpl flag has no value!\n")}

	if idx := strings.Index(zval.(string), "."); idx > -1 {log.Fatalf("error -- template name has extension!\n")}
    tplnam = zval.(string)


    cfDir := os.Getenv("cfDir")
    if len(cfDir) == 0 {log.Fatalf("error -- cannot get env var cfDir!")}

	tplFilnam := cfDir +"/token/" + tplnam + ".yaml"

    log.Printf("info -- debug: %t\n", dbg)
    log.Printf("info -- tpl name:  %s\n", tplnam)
    log.Printf("info -- tpl file:  %s\n", tplFilnam)

	DnsRecList := make([]cfLib.DnsRec,2)

	fdat, err := yaml.Marshal(&DnsRecList)
	if err != nil {log.Fatalf("error -- Marshal DnsRecList: %v\n", err)}

	StartTpl := []byte("# dns template\n---\n")
	StartTpl = append(StartTpl, fdat ...)

	err = os.WriteFile(tplFilnam, StartTpl, 0666)
	if err != nil {log.Fatalf("error -- Write Tpl File: %v\n",err)}

	log.Printf("info -- success writing template file!")
}
