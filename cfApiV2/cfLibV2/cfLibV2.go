// cloudflare support library
// Author: prr, azul software
// Date: 15 Jan 2024
// Copyright (c) 2024 prr, azul software

package cfLibV2

import (
	"fmt"
	"os"
	"log"
	"time"
	"context"
	"strings"

    "github.com/cloudflare/cloudflare-go"
	json "github.com/goccy/go-json"
    yaml "github.com/goccy/go-yaml"
)

/*
type ApiObj struct {
    Api    string `yaml:"Api"`
    ApiKey string `yaml:"ApiKey"`
    ApiToken string `yaml:"ApiToken"`
	TokenId string `yaml:"TokenId"`
	TokName string `yaml:"TokenName"`
	Start time.Time `yaml:"Start"`
	Expiration time.Time `yaml:"Expire"`
	AccountId string `yaml:"AccountId"`
    Email     string `yaml:"Email"`
	YamlFile	string `yaml:"cfTokenFile"`
}

type cfApi struct {
	API *cloudflare.API
	ApiObj *ApiObj
}
*/

type CfAccount struct {
	Id string `json:"Id"`
	Name string `json:"Name"`
	Created time.Time `json:"Created"`
	Email string `json:"email:"`
}



type TokList struct {
	AccountId string `yaml:"AccountId"`
	Name string `yaml:"Name"`
	Mod string `yaml:"ModOn"`
	Toks []cfToken `yaml:"Toks"`
}

type cfToken struct {
	Id string `yaml:"Id"`
	Name string `yaml:"Name"`
	ExpTim time.Time `yaml:"Exp"`
}

type ApiPerm struct {
	Id string `yaml:"ID"`
	Name string `yaml:"Name"`
	Scopes int `yaml:"Scopes"`
}

type ApiPermList struct {
	Date time.Time `yaml:"Date"`
	ApiPerms []ApiPerm `yaml:"ApiPerm"`
}

type ZoneList struct {
	AccountId string `yaml:"AccountId"`
	Email string `yaml:"Email"`
	Modified string `yaml:"Modified"`
	ModTime time.Time
	Zones []ZoneShort `yaml:"Zones"`
}

type ZoneShort struct {
	Name string `yaml:"Name"`
	Id string `yaml:"Id"`
}

type ZoneAcme struct {
	Name string `yaml:"Name"`
	Id string `yaml:"Id"`
	AcmeId string `yaml:"AcmeId"`
	AcmeRec bool
}

type ZoneShortJson struct {
	Name string `json:"Name"`
	Id string `json:"Id"`
}

type TokOpt struct {
	Start time.Time `json:"Start"`
	End time.Time `json:"End"`
	Days int	`json:"Days"`
}

func GetTokOpt(filnam string) (opt TokOpt, err error) {

	opt.Start = time.Now().UTC().Round(time.Second)
	opt.End = time.Now().UTC().AddDate(0,0,7).Round(time.Second)

	bdat, err := os.ReadFile(filnam)
	if err != nil {return opt, fmt.Errorf("ReadFile: %v", err)}

//	fmt.Printf("bdat: %s\n", bdat)

	topt := TokOpt{}
	err = yaml.Unmarshal(bdat, &topt)
	if err != nil {return opt, fmt.Errorf("Unmarshal: %v", err)}

/*
	fmt.Printf("start: %s\n", topt.Start)
	fmt.Printf("end: %s\n", topt.End)
	fmt.Printf("days: %d\n", topt.Days)
*/
	if topt.Start == time.Date(1,time.Month(1),1, 0,0,0,0, opt.Start.Location()) {
		topt.Start = opt.Start
	} else {
		if topt.Start.Before(time.Now()) {return topt, fmt.Errorf("start time is in the past!")}
	}
	if topt.End == time.Date(1,time.Month(1),1, 0,0,0,0, opt.Start.Location()) {
		if topt.Days > 0 {topt.End = topt.Start.AddDate(0,0,topt.Days)}
	} else {
		topt.End = topt.Start.AddDate(0,0,1)
	}

	topt.End = time.Date(topt.End.Year(), topt.End.Month(), topt.End.Day(), 0,0,0,0, topt.End.Location())
	if topt.End.Before(topt.Start) {return topt, fmt.Errorf("end time is before start time")}

	return topt, nil
}


func CreateTokFile(tokFilnam string, Token *cloudflare.APIToken, dbg bool) (err error){

    if dbg {log.Printf("token filnam: %s\n", tokFilnam)}

    if dbg {PrintToken(*Token)}

    bdat, err := json.Marshal(*Token)
    if err != nil {return fmt.Errorf("json encode: %v\n", err)}

    if dbg {fmt.Printf("json: %s\n", string(bdat))}

    tokfil, err := os.Create(tokFilnam)
    if err != nil {return fmt.Errorf("failure creating token file: %v", err)}
    defer tokfil.Close()

    _, err = tokfil.Write(bdat)
    if err !=nil {return fmt.Errorf("token write: %v\n", err)}

    return nil
}

func ReadDnsRecord(zoneId string, token string, dnspar *cloudflare.ListDNSRecordsParams)(dnsrec []cloudflare.DNSRecord, err error) {

	// token  is Token.Value
	if len(token) == 0 {return dnsrec, fmt.Errorf("no token supplied")}
	if len(zoneId) == 0 {return dnsrec, fmt.Errorf("no zoneId supplied")}
   	api, err := cloudflare.NewWithAPIToken(token)
    if err != nil {return dnsrec, fmt.Errorf("initiating api obj: %v",err)}

    ctx := context.Background()

    rc := cloudflare.ResourceContainer{
        Level: cloudflare.ZoneRouteLevel,
        Identifier: zoneId,
    }
	// https://pkg.go.dev/github.com/cloudflare/cloudflare-go#ListDNSRecordsParams
    DnsPars:=cloudflare.ListDNSRecordsParams{}
	if dnspar != nil { DnsPars = *dnspar }

    dnsRecs, _, err := api.ListDNSRecords(ctx, &rc, DnsPars)
    if err != nil {return dnsrec, fmt.Errorf("api.ListDNSRecords: %v", err)}

	return dnsRecs, nil
}


func AddDnsRecord(param *cloudflare.CreateDNSRecordParams, token string)(recId string, err error) {

/*

type CreateDNSRecordParams struct {
	CreatedOn  time.Time   `json:"created_on,omitempty" url:"created_on,omitempty"`
	ModifiedOn time.Time   `json:"modified_on,omitempty" url:"modified_on,omitempty"`
	Type       string      `json:"type,omitempty" url:"type,omitempty"`
	Name       string      `json:"name,omitempty" url:"name,omitempty"`
	Content    string      `json:"content,omitempty" url:"content,omitempty"`
	Meta       interface{} `json:"meta,omitempty"`
	Data       interface{} `json:"data,omitempty"` // data returned by: SRV, LOC
	ID         string      `json:"id,omitempty"`
	ZoneID     string      `json:"zone_id,omitempty"`
	ZoneName   string      `json:"zone_name,omitempty"`
	Priority   *uint16     `json:"priority,omitempty"`
	TTL        int         `json:"ttl,omitempty"`
	Proxied    *bool       `json:"proxied,omitempty" url:"proxied,omitempty"`
	Proxiable  bool        `json:"proxiable,omitempty"`
	Locked     bool        `json:"locked,omitempty"`
	Comment    string      `json:"comment,omitempty" url:"comment,omitempty"` // to the server, there's no difference between "no comment" and "empty comment"
	Tags       []string    `json:"tags,omitempty"`
}

*/

	if len(token) == 0 {return "", fmt.Errorf("no token supplied")}
//	if len(zoneId) == 0 {return "", fmt.Errorf("no zoneId supplied")}

	if len(param.Type) == 0 {return "", fmt.Errorf("param Type missing!")}
	if len(param.Name) == 0 {return "", fmt.Errorf("param Name missing!")}
	if len(param.Content) == 0 {return "", fmt.Errorf("param Content missing!")}
	if len(param.ZoneID) == 0 {return "", fmt.Errorf("param ZoneID missing!")}


	// token  is Token.Value
   	api, err := cloudflare.NewWithAPIToken(token)
    if err != nil {return "", fmt.Errorf("initiating api obj: %v",err)}
    ctx := context.Background()


    rc := cloudflare.ResourceContainer{
        Level: cloudflare.ZoneRouteLevel,
//        Identifier: zoneId,
        Identifier: param.ZoneID,
    }

	dnsRec, err := api.CreateDNSRecord(ctx, &rc, *param)
	if err != nil {return "", fmt.Errorf("CreateDnsREc: %v", err)}

	return dnsRec.ID, nil
}

func AddDnsChalRecord(zoneId, val, token string)(recId string, err error) {

	if len(token) == 0 {return "", fmt.Errorf("no token supplied")}
	if len(zoneId) == 0 {return "", fmt.Errorf("no zoneId supplied")}

	// token  is Token.Value
   	api, err := cloudflare.NewWithAPIToken(token)
    if err != nil {return "", fmt.Errorf("initiating api obj: %v",err)}

    ctx := context.Background()

    param := cloudflare.CreateDNSRecordParams{
        CreatedOn: time.Now(),
        Type: "TXT",
        Name: "_acme-challenge",
        Content: val,
        TTL: 30000,
        Comment: "acme challenge record",
    }

    rc := cloudflare.ResourceContainer{
        Level: cloudflare.ZoneRouteLevel,
        Identifier: zoneId,
    }

	dnsRec, err := api.CreateDNSRecord(ctx, &rc, param)
	if err != nil {return "", fmt.Errorf("CreateDnsREc: %v", err)}

	return dnsRec.ID, nil
}

func UpdDnsRecord(param *cloudflare.UpdateDNSRecordParams, token string)(recId string, err error) {

/*

type UpdateDNSRecordParams struct {
	Type     string      `json:"type,omitempty"`
	Name     string      `json:"name,omitempty"`
	Content  string      `json:"content,omitempty"`
	Data     interface{} `json:"data,omitempty"` // data for: SRV, LOC
	ID       string      `json:"-"`
	Priority *uint16     `json:"priority,omitempty"`
	TTL      int         `json:"ttl,omitempty"`
	Proxied  *bool       `json:"proxied,omitempty"`
	Comment  *string     `json:"comment,omitempty"` // nil will keep the current comment, while StringPtr("") will empty it
	Tags     []string    `json:"tags"`
}

*/
	if len(token) == 0 {return "", fmt.Errorf("no token supplied")}
//	if len(zoneId) == 0 {return "", fmt.Errorf("no zoneId supplied")}

	if len(param.Type) == 0 {return "", fmt.Errorf("param Type missing!")}
	if len(param.Name) == 0 {return "", fmt.Errorf("param Name missing!")}
	if len(param.Content) == 0 {return "", fmt.Errorf("param Content missing!")}
	if len(param.ID) == 0 {return "", fmt.Errorf("param ID missing!")}


	// token  is Token.Value
   	api, err := cloudflare.NewWithAPIToken(token)
    if err != nil {return "", fmt.Errorf("initiating api obj: %v",err)}

    ctx := context.Background()


    rc := cloudflare.ResourceContainer{
        Level: cloudflare.ZoneRouteLevel,
//        Identifier: zoneId,
        Identifier: param.ID,
    }

	dnsRec, err := api.UpdateDNSRecord(ctx, &rc, *param)
	if err != nil {return "", fmt.Errorf("UpdateDnsRec: %v", err)}

	return dnsRec.ID, nil
}


func DelDnsRecord(zoneId, recId, token string) (err error) {

	if len(token) == 0 {return fmt.Errorf("no token supplied")}
	if len(zoneId) == 0 {return fmt.Errorf("no zoneId supplied")}
	// token  is Token.Value
   	api, err := cloudflare.NewWithAPIToken(token)
    if err != nil {return fmt.Errorf("initiating api obj: %v",err)}
    ctx := context.Background()

    rc := cloudflare.ResourceContainer{
        Level: cloudflare.ZoneRouteLevel,
        Identifier: zoneId,
    }

    err = api.DeleteDNSRecord(ctx, &rc, recId)
    if err != nil {return fmt.Errorf("DeleteDnsRecord: %v", err)}

    return nil
}

func FindAcmeChalRecord(dnsrec  []cloudflare.DNSRecord)(recId string, err error) {

	if len(dnsrec) == 0 {return recId, fmt.Errorf("no dns records supplied!")}
	recId = ""
	for i:=0; i< len(dnsrec); i++ {
		if idx := strings.Index(dnsrec[i].Name, "_acme-challenge"); idx > -1 {
			recId = dnsrec[i].ID
			break
		}
	}
	if len(recId) == 0 {return recId, fmt.Errorf("no Acme Chal Record found!")}
	return recId, nil
}

/*
func (cfapi *cfApi) DelDnsChalRecord (zone ZoneAcme) (err error) {

	api := cfapi.API

    ctx := context.Background()

    var rc cloudflare.ResourceContainer
    //domains
    rc.Level = cloudflare.ZoneRouteLevel
    rc.Identifier = zone.Id
	recId := zone.AcmeId

	err = api.DeleteDNSRecord(ctx, &rc, recId)
	if err != nil {return fmt.Errorf("DeleteDnsRecord: %v", err)}

	return nil
}
*/

func SaveZonesJson(zones []cloudflare.Zone, outfil *os.File)(err error) {

	if outfil == nil { return fmt.Errorf("no file provided!")}

	jsonData, err := json.Marshal(zones)
	if err != nil {return fmt.Errorf("json.Marshal: %v", err)}

	_, err = outfil.Write(jsonData)
	if err != nil {return fmt.Errorf("jsonData os.Write: %v", err)}
	return nil
}

func SaveZonesYaml(zones []cloudflare.Zone, outfil *os.File)(err error) {

	if outfil == nil { return fmt.Errorf("no file provided!")}
	yamlData, err := yaml.Marshal(zones)
	if err != nil {return fmt.Errorf("yaml.Marshal: %v", err)}

	_, err = outfil.Write(yamlData)
	if err != nil {return fmt.Errorf("yamlData os.Write: %v", err)}
	return nil
}

func SaveZonesShortJson(zones []ZoneShort, outfil *os.File)(err error) {

	if outfil == nil { return fmt.Errorf("no file provided!")}

	jsonData, err := json.Marshal(zones)
	if err != nil {return fmt.Errorf("json.Marshal: %v", err)}

	_, err = outfil.Write(jsonData)
	if err != nil {return fmt.Errorf("jsonData os.Write: %v", err)}
	return nil
}

func SaveZonesShortYaml(zones []ZoneShort, outfil *os.File)(err error) {

	if outfil == nil { return fmt.Errorf("no file provided!")}

	yamlData, err := yaml.Marshal(zones)
	if err != nil {return fmt.Errorf("yaml.Marshal: %v", err)}

	_, err = outfil.WriteString("---\n")
	if err != nil {return fmt.Errorf("yamlData os.WriteString: %v", err)}

	_, err = outfil.Write(yamlData)
	if err != nil {return fmt.Errorf("yamlData os.Write: %v", err)}
	return nil
}

func SaveZoneShortFile(zoneList *ZoneList, outfil *os.File)(err error) {

	if outfil == nil { return fmt.Errorf("no file provided!")}

	yamlData, err := yaml.Marshal(zoneList)
	if err != nil {return fmt.Errorf("yaml.Marshal: %v", err)}

	_, err = outfil.WriteString("---\n")
	if err != nil {return fmt.Errorf("yamlData os.WriteString: %v", err)}

	_, err = outfil.Write(yamlData)
	if err != nil {return fmt.Errorf("yamlData os.Write: %v", err)}
	return nil
}

func SaveAcmeDns(zones []ZoneAcme, outfil *os.File)(err error) {

	if outfil == nil { return fmt.Errorf("no file provided!")}

	yamlData, err := yaml.Marshal(zones)
	if err != nil {return fmt.Errorf("yaml.Marshal: %v", err)}

	_, err = outfil.WriteString("---\n")
	if err != nil {return fmt.Errorf("yamlData os.WriteString: %v", err)}

	_, err = outfil.Write(yamlData)
	if err != nil {return fmt.Errorf("yamlData os.Write: %v", err)}
	return nil
}

func ReadAcmeDns(filnam string)(zones []ZoneAcme, err error) {

	bdat, err := os.ReadFile(filnam)
	if err != nil {return zones, fmt.Errorf("reading file: %v", err)}

	err = yaml.Unmarshal(bdat, &zones)
	if err != nil {return zones, fmt.Errorf("unmarshal: %v", err)}
	return zones, nil
}

func ReadZonesShortYaml(infil *os.File)(zoneListObj *[]ZoneShort, err error) {

	var zonesShort []ZoneShort

	if infil == nil { return nil, fmt.Errorf("no file provided!")}

	info, err := infil.Stat()
	if err != nil {return nil, fmt.Errorf("info.Stat: %v", err)}

	size := info.Size()

	inBuf := make([]byte, int(size))

	_, err = infil.Read(inBuf)
	if err != nil {return nil, fmt.Errorf("infil.Read: %v", err)}

	err = yaml.Unmarshal(inBuf, &zonesShort)
	if err != nil {return nil, fmt.Errorf("yaml.Unmarshal: %v", err)}

	return &zonesShort, nil
}


// read acme file
func ReadAcmeZones(inFilNam string)(zoneListObj *[]ZoneAcme, err error) {

	var zones []ZoneAcme

	inBuf, err := os.ReadFile(inFilNam)
	if err != nil {return nil, fmt.Errorf("os.ReadFile: %v", err)}

	err = yaml.Unmarshal(inBuf, &zones)
	if err != nil {return nil, fmt.Errorf("yaml.Unmarshal: %v", err)}

	return &zones, nil
}

func ReadZoneShortFile(inFilNam string)(zoneList *ZoneList, err error) {

	var zonelist ZoneList

	inBuf, err := os.ReadFile(inFilNam)
	if err != nil {return nil, fmt.Errorf("os.ReadFile: %v", err)}

	err = yaml.Unmarshal(inBuf, &zonelist)
	if err != nil {return nil, fmt.Errorf("yaml.Unmarshal: %v", err)}

	return &zonelist, nil

}

func ReadTokenFile(rdTokFilnam string, vfy bool)(tok cloudflare.APIToken, err error) {

	bdat, err := os.ReadFile(rdTokFilnam)
    if err != nil {return tok, fmt.Errorf("Open Token File: %v\n", err)}

//    if dbg {log.Printf(debug -- ReadRTokenFile: 
//	fmt.Printf("json: %s\n***********\n", string(bdat))

    DnsTok := cloudflare.APIToken {}
    err = json.Unmarshal(bdat, &DnsTok)
    if err != nil {return tok, fmt.Errorf("unMarshal: %v\n", err)}

	if len(DnsTok.Value) == 0 {return DnsTok, fmt.Errorf("no token value!")}
    log.Printf("info -- success reading Token File!\n")

	if !vfy {return DnsTok, nil}

	// verify
	api, err := cloudflare.NewWithAPIToken(DnsTok.Value)
    if err != nil {return tok, fmt.Errorf("initializing api obj: %v\n",err)}

    ctx := context.Background()

    resp, err :=api.VerifyAPIToken(ctx)
    if err != nil {return tok, fmt.Errorf("verify: %v\n", err)}

    if !(resp.Status == "active") {return tok, fmt.Errorf("token status is not active: %s\n", resp.Status)}

	return DnsTok, nil
}

func SaveTokList(outFilnam string, tokList []cloudflare.APIToken) (err error){

	var tokSav TokList

	cfTokList := make([]cfToken, len(tokList))

	tokSav.Toks = cfTokList

    for i:=0; i<len(tokList); i++ {
//        tok := tokList[i]
		cfTokList[i].Id = tokList[i].ID
		cfTokList[i].Name = tokList[i].Name
		cfTokList[i].ExpTim = *(tokList[i].ExpiresOn)
//        fmt.Printf("  [%d]: %-20s| %-30s| %-5s| %-10s %-20s\n",i+1, tok.ID, tok.Name, tok.Value, tok.Status, tok.ExpiresOn.Format(time.RFC1123))
	}

    yamlData, err := yaml.Marshal(&tokSav)
    if err != nil {return fmt.Errorf("yaml.Marshal: %v", err)}

   	err = os.WriteFile(outFilnam,yamlData, 0666)
    if err != nil {return fmt.Errorf("yamlData os.Write: %v", err)}

    return nil
}

/*
func PrintCsr(csrlist *CsrList) {

    fmt.Println("******** Csr List *********")
    fmt.Printf("template: %s\n", csrlist.Template)
	numDom := len(csrlist.Domains)
	fmt.Printf("domains: %d\n", numDom)
	for i:=0; i< numDom; i++ {
		csrdat := csrlist.Domains[i]
	    fmt.Printf("  domain:   %s\n", csrdat.Domain)
    	fmt.Printf("  email:    %s\n", csrdat.Email)
	    fmt.Printf("  name:\n")
    	nam:= csrdat.Name
    	fmt.Printf("    CommonName:   %s\n", nam.CommonName)
    	fmt.Printf("    Country:      %s\n", nam.Country)
    	fmt.Printf("    Province:     %s\n", nam.Province)
    	fmt.Printf("    Locality:     %s\n", nam.Locality)
    	fmt.Printf("    Organisation: %s\n", nam.Organisation)
    	fmt.Printf("    OrgUnit:      %s\n", nam.OrganisationUnit)
	}

    fmt.Println("******** End Csr List *******")

}
*/

func PrintTokResp(res *cloudflare.APITokenVerifyBody) {
	fmt.Printf("************ Token Verification ******\n")
	fmt.Printf("ID:       %s\n", res.ID)
	fmt.Printf("Status:   %s\n", res.Status)
	fmt.Printf("NotBefore: %s\n",res.NotBefore.Format(time.RFC1123))
	fmt.Printf("Expires:   %s\n", res.ExpiresOn.Format(time.RFC1123))
}


func PrintZones(zones []cloudflare.Zone) {

    fmt.Printf("************** Zones/Domains [%d] *************\n", len(zones))

    for i:=0; i< len(zones); i++ {
        zone := zones[i]
        fmt.Printf("%d %-20s %s\n",i+1, zone.Name, zone.ID)
    }
}

func PrintZoneList(zoneList *ZoneList){

    fmt.Printf("************** ZoneList *************\n")

	fmt.Printf("AccountId: %s\n", zoneList.AccountId)
	fmt.Printf("Email:     %s\n", zoneList.Email)
	fmt.Printf("Modified:  %s\n", zoneList.ModTime.Format(time.RFC1123))
	zonesLen := len((*zoneList).Zones)
	fmt.Printf("*** Zones[%d]: ***\n", zonesLen)
    for i:=0; i< zonesLen; i++ {
        zone :=(* zoneList).Zones[i]
        fmt.Printf("   %d %-20s %s\n",i+1, zone.Name, zone.Id)
    }
}

/*
func PrintApiObj (apiObj *ApiObj) {

    fmt.Println("***************** Api Obj ******************")
    fmt.Printf("API:       %s\n", apiObj.Api)
    fmt.Printf("APIKey:    %s\n", apiObj.ApiKey)
    fmt.Printf("ApiToken:  %s\n", apiObj.ApiToken)
    fmt.Printf("TokenId:   %s\n", apiObj.TokenId)
    fmt.Printf("TokName:   %s\n", apiObj.TokName)
    fmt.Printf("Start:     %s\n", apiObj.Start.Format(time.RFC1123))
    fmt.Printf("Expiry:    %s\n", apiObj.Expiration.Format(time.RFC1123))
    fmt.Printf("AccountId: %s\n", apiObj.AccountId)
    fmt.Printf("Email:     %s\n", apiObj.Email)
    fmt.Printf("YamlFile:  %s\n", apiObj.YamlFile)
    fmt.Println("********************************************")
}
*/

// https://github.com/cloudflare/cloudflare-go/blob/0d05fc09483641dde8abb4c64cf2f6016f590d79/user.go#L12
func PrintUserInfo (u *cloudflare.User) {

    var actTyp string

    fmt.Println("************** User Info **************")
    fmt.Printf("First Name:  %s\n", u.FirstName)
    fmt.Printf("Last Name:   %s\n", u.LastName)
    fmt.Printf("Email:       %s\n", u.Email)
    fmt.Printf("ID:          %s\n", u.ID)
    fmt.Printf("Country:     %s\n", u.Country)
    fmt.Printf("Zip Code:    %s\n", u.Zipcode)
    fmt.Printf("Phone:       %s\n", u.Telephone)
    fmt.Printf("2FA:         %t\n", u.TwoFA)
    timStr := (u.CreatedOn).Format("02 Jan 06 15:04 MST")
    fmt.Printf("Created:     %s\n", timStr)
    timStr = (u.ModifiedOn).Format("02 Jan 06 15:04 MST")
    fmt.Printf("Modified:    %s\n", timStr)
    fmt.Printf("ApiKey:      %s\n", u.APIKey)
    if len(u.Accounts) == 1 {
        act := u.Accounts[0]
        actTyp = act.Type
        if len(actTyp) == 0 {actTyp = "-"}
        fmt.Printf("account ID: %s Name: %s Type: %s\n", act.ID, act.Name, actTyp)
    } else {
        fmt.Printf("Accounts [%d]:\n", len(u.Accounts))
        fmt.Printf("Nu ID  Name  Type\n")
        for i:=0; i< len(u.Accounts); i++ {
            act := u.Accounts[i]
            actTyp = act.Type
            if len(actTyp) == 0 {actTyp = "-"}
            fmt.Printf("%d: %s %s %s\n", i+1, act.ID, act.Name, actTyp)
        }
    }
    fmt.Println("********** End User Info **************")
}

func PrintResInfo(res *cloudflare.ResultInfo) {

    fmt.Println("************** ResultInfo **************")
    fmt.Printf("Page:       %d\n", res.Page)
    fmt.Printf("PerPage:    %d\n", res.PerPage)
    fmt.Printf("TotalPages: %d\n", res.TotalPages)
    fmt.Printf("Count:      %d\n", res.Count)
    fmt.Printf("Total:      %d\n", res.Total)
    fmt.Println("********** End ResultInfo **************")
}

func PrintDnsRecs(recs *[]cloudflare.DNSRecord) {
    fmt.Printf("************** DNS Records: %d *************\n", len(*recs))
    fmt.Println("number           ID          type      name             value/ content")
    for i:=0; i< len(*recs); i++ {
		rec := (*recs)[i]
        fmt.Printf("Record[%d]: %-15s %-3s %s %s\n", i+1, rec.ID, rec.Type, rec.Name, rec.Content)
    }
    fmt.Printf("************** End DNS Records **************\n")
}

func PrintDnsRec(rec *cloudflare.DNSRecord) {
    fmt.Printf("************* DNS Record  ***************\n", )
	fmt.Printf("ID: %s Type: %-3s Name: %s Value: %s\n", rec.ID, rec.Type, rec.Name, rec.Content)
    fmt.Printf("************* End DNS Record ************\n")
}

func PrintAcmeZones(zones []ZoneAcme) {
	fmt.Printf("*********** Acme Zones: %d ***************\n", len(zones))
	for i:=0; i< len(zones); i++ {
		zone := zones[i]
		fmt.Printf("Zone [%d] Id: %s Name: %s Acme Record Id: %s\n", i+1, zone.Id, zone.Name, zone.AcmeId)
	}
	fmt.Printf("*********** End Acme Zones ***************\n")
}

func PrintAccount(act *cloudflare.Account) {

	fmt.Println("****** Account Info *****")
	fmt.Printf("Id:    %s\n", act.ID)
	fmt.Printf("Name: %s\n", act.Name)
	fmt.Printf("Type: %s\n", act.Type)
	t := act.CreatedOn
	fmt.Printf("CreatedOn: %s\n", t.Format(time.RFC1123))
	fmt.Printf("2Fa: %t\n",act.Settings.EnforceTwoFactor)
}

func PrintTokList(tokList []cloudflare.APIToken) {

    fmt.Printf("************ Token List [%d] **************\n", len(tokList))
    fmt.Printf("   seq     ID        Name        Value      Status  Exp \n")
    for i:=0; i<len(tokList); i++ {
        tok := tokList[i]
        fmt.Printf("  [%d]: %-20s| %-30s| %-5s| %-10s %-20s\n",i+1, tok.ID, tok.Name, tok.Value, tok.Status, tok.ExpiresOn.Format(time.RFC1123))
	}
    for i:=0; i<len(tokList); i++ {
        tok := tokList[i]
		fmt.Printf("**** detail token: %d ******\n", i+1)
		PrintToken(tok)
    }
}

func PrintToken(tok cloudflare.APIToken) {

	fmt.Printf("  Id:     %s\n", tok.ID)
	fmt.Printf("  Name:   %s\n", tok.Name)
	fmt.Printf("  Value:  %s\n", tok.Value)
	fmt.Printf("  Status: %s\n", tok.Status)
	timStr := "NA"
	if tok.NotBefore != nil {timStr = tok.NotBefore.Format(time.RFC1123)}
	fmt.Printf("  Start:  %s\n", timStr)
	timStr = "NA"
	if tok.ExpiresOn != nil {timStr = tok.ExpiresOn.Format(time.RFC1123)}
	fmt.Printf("  Expiration: %s\n", timStr)
	timStr = "NA"
	if tok.ModifiedOn != nil {timStr = tok.ModifiedOn.Format(time.RFC1123)}
	fmt.Printf("  Modified:   %s\n", timStr)
	fmt.Printf("  Policies: %d\n", len(tok.Policies))
	for j:=0; j< len(tok.Policies); j++ {
		pol := tok.Policies[j]
		fmt.Printf("  ***** Policy %d ****\n", j+1)
		fmt.Printf("    ID:     %s\n", pol.ID)
		fmt.Printf("    Effect: %s\n", pol.Effect)
		fmt.Printf("    Resources: %d\n", len(pol.Resources))
		for k,v := range pol.Resources {
			fmt.Printf("       key: %s val: %v\n",k , v)
		}
		fmt.Printf("    PermGroups: %d\n", len(pol.PermissionGroups))
		for k:=0; k<len(pol.PermissionGroups); k++ {
			pgrp := pol.PermissionGroups[k]
			fmt.Printf("        pgrp[%d]: %s %s %d\n", k+1, pgrp.ID, pgrp.Name, len(pgrp.Scopes))
			for l:=1; l<len(pgrp.Scopes); l++ {
				fmt.Printf("         scope[%d]: %s\n", l+1, pgrp.Scopes[l])
			}
		}
	}
	cond := tok.Condition
	if cond == nil {return}
	ipCond := cond.RequestIP
	if ipCond == nil {return}

	if len(ipCond.In) > 0 {
		fmt.Printf("  **** Conditions In:\n", )
		for j:=0; j< len(ipCond.In); j++ {
			fmt.Printf("    %d: %s\n", j+1, ipCond.In[j])
		}
	}

	if len(ipCond.NotIn) > 0 {
		fmt.Printf("  **** Conditions NotIn:\n", )
		for j:=0; j< len(ipCond.NotIn); j++ {
			fmt.Printf("     %d: %s\n", j+1, ipCond.NotIn[j])
		}
	}
}

