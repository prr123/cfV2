// cloudflare support library
// Author: prr, azul software
// Date: 29 March 2023
// Copyright 2023 prr, azul software

package cfLib

import (
	"fmt"
	"os"
	"log"
	"time"
	"context"
	"strings"


    yaml "github.com/goccy/go-yaml"
    "github.com/cloudflare/cloudflare-go"
	json "github.com/goccy/go-json"
)

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

func InitCfLib(yamlFilnam string) (apiObjRef *ApiObj, err error) {
	var apiObj ApiObj

    // open file and decode
    buf, err := os.ReadFile(yamlFilnam)
    if err != nil {
        return nil, fmt.Errorf("cannot open yaml File: os.ReadFile: %v\n", err)
    }

//    fmt.Printf("buf [%d]:\n%s\n", len(buf), string(buf))

    if err := yaml.Unmarshal(buf, &apiObj); err !=nil {
		return nil, fmt.Errorf("error Unmarshalling Yaml File: %v\n", err)
    }

	if apiObj.Api != "cloudflare" {
		return nil, fmt.Errorf("Api is not cloudflare!")
	}

	apiObj.YamlFile = yamlFilnam

	return &apiObj, nil
}

//  function that initiates the cloudflare api with a yaml api file containing a token
func InitCfApi(apiFilnam string) (cfapi *cfApi, err error) {

	cfDir := os.Getenv("Cloudflare")
	if len(cfDir) == 0 {return nil, fmt.Errorf("could not get env: Cloudflare\n")}

	yamlFilnam := cfDir + "/token/cfZones.yaml"


	if len(apiFilnam) > 0 {
		yamlFilnam = cfDir + "/token/" + apiFilnam
	}

	apiObj, err := InitCfLib(yamlFilnam)
    if err != nil {return nil, fmt.Errorf("cfLib.InitCfLib: %v\n", err)}

	api, err := cloudflare.NewWithAPIToken(apiObj.ApiToken)
	if err != nil {return nil, fmt.Errorf("NewWithAPIToken: %v\n", err)}

	cfApiObj := &cfApi{API: api, ApiObj: apiObj,}

	return cfApiObj, nil
}

func CreateTokFile(filnam string, token string, dbg bool) (err error){

	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {return fmt.Errorf("NewWithAPIToken: %v\n", err)}
	if dbg {log.Printf("cf api returned")}

	if len(filnam) ==0 {return fmt.Errorf("no filnam provided!")}
	idx := strings.Index(filnam, ".yaml")
	if idx == -1 { filnam += ".yaml"}

	cfDir := os.Getenv("Cloudflare")
	if len(cfDir) == 0 {return fmt.Errorf("could not get env: Cloudflare\n")}

	// todo: check whether token dir exists

	rdTokFilnam := cfDir + "/token/cfTokRead.yaml"

	tokFilnam := cfDir + "/token/" + filnam
	if dbg {log.Printf("token filnam: %s\n", tokFilnam)}

    ctx := context.Background()

	tokResp, err := api.VerifyAPIToken(ctx)
	if err != nil {return fmt.Errorf("VerifyApiToken: %v", err)}

	if dbg {PrintTokResp(&tokResp)}

	if tokResp.Status != "active" {return fmt.Errorf("invalid status returned! %s", tokResp.Status)}


	apiObj := ApiObj{
		Api: "Cloudflare",
		ApiKey: "na",
		ApiToken: token,
		TokenId: tokResp.ID,
		TokName: "",
		Start: tokResp.NotBefore,
		Expiration: tokResp.ExpiresOn,
		AccountId: "d0e0781201c0536742831e308ce406fb",
		Email: "azulsoftwarevlc@gmail.com",
		YamlFile: filnam,
	}

	rdTokDat, err := os.ReadFile(rdTokFilnam)
	if err != nil {return fmt.Errorf("readFile rdToken: %v", err)}

	tokApiObj := ApiObj{}
	err = yaml.Unmarshal(rdTokDat, &tokApiObj)
	if err != nil {return fmt.Errorf("Unmarshal reTokDat: %v", err)}

	tokApi, err := cloudflare.NewWithAPIToken(tokApiObj.ApiToken)
	if err != nil {return fmt.Errorf("NewWithAPIToken tokApi: %v\n", err)}

	tok, err := tokApi.GetAPIToken(ctx, tokResp.ID)
	if err != nil {return fmt.Errorf("GetApiToken: %v", err)}

	apiObj.TokName = tok.Name
//ww
	if dbg {PrintToken(tok)}
	outfil, err := os.Create(tokFilnam)
	if err != nil {return fmt.Errorf("os.Create: %v", err)}

	_, err = outfil.Write([]byte("---\n"))
	if err != nil {return fmt.Errorf("header outfil.Write: %v", err)}

	yamlData, err := yaml.Marshal(&apiObj)
	if err != nil {return fmt.Errorf("yaml.Marshal: %v", err)}

	_, err = outfil.Write(yamlData)
	if err != nil {return fmt.Errorf("yamlData os=utfil.Write: %v", err)}

	return nil
}

func VerifyCFToken (tokFilnam string) (err error) {

	var apiObj ApiObj

	yamlFilnam := ""

	cfDir := os.Getenv("Cloudflare")
	if len(cfDir) == 0 {return fmt.Errorf("could not get env: Cloudflare\n")}

	if len(tokFilnam) == 0 {
	    yamlFilnam = cfDir + "/token/cfZones.yaml"
	} else {
		yamlFilnam = cfDir + "/token/" + tokFilnam
	}

    buf, err := os.ReadFile(yamlFilnam)
    if err != nil {
        return fmt.Errorf("cannot open yaml File: os.ReadFile: %v\n", err)
    }

//    fmt.Printf("buf [%d]:\n%s\n", len(buf), string(buf))

    if err := yaml.Unmarshal(buf, &apiObj); err !=nil {
		return fmt.Errorf("error Unmarshalling Yaml File: %v\n", err)
    }

	if apiObj.Api != "cloudflare" {
		return fmt.Errorf("Api is not cloudflare!")
	}

	api, err := cloudflare.NewWithAPIToken(apiObj.ApiToken)
	if err != nil {return fmt.Errorf("NewWithAPIToken: %v\n", err)}

    ctx := context.Background()

	tokResp, err := api.VerifyAPIToken(ctx)
	if err != nil {return fmt.Errorf("VerifyApiToken: %v", err)}

	fmt.Printf("*** Token Response ****\n")
	fmt.Printf("ID:     %s\n", tokResp.ID)
	fmt.Printf("Status: %s\n", tokResp.Status)
	fmt.Printf("No Before: %s\n", tokResp.NotBefore.Format(time.RFC1123))
	fmt.Printf("Expires:   %s\n", tokResp.ExpiresOn.Format(time.RFC1123))

	return nil
}


func (cfapi *cfApi) ListDnsRecords (zoneId string) (dnsList *[]cloudflare.DNSRecord, err error) {

	if cfapi == nil {return nil, fmt.Errorf("cfApi is nil!")}
	if cfapi.API == nil {return nil, fmt.Errorf("cfApi.api is nil!")}

	if len(zoneId) == 0 {return nil, fmt.Errorf("no zoneId provided!")}

	api := cfapi.API

    ctx := context.Background()

    rc := cloudflare.ResourceContainer{
   		Level: cloudflare.ZoneRouteLevel,
		Identifier: zoneId,
	}

	DnsPars:=cloudflare.ListDNSRecordsParams{}

	dnsRecs, _, err := api.ListDNSRecords(ctx, &rc, DnsPars)
	if err != nil {return nil, fmt.Errorf("api.ListDNSRecords: %v\n", err)}
//    PrintDnsRecs(&dnsRecs)

	return &dnsRecs, nil
}

func (cfapi *cfApi) AddDnsRecord (zoneId string, dnsPar *cloudflare.CreateDNSRecordParams) (dnsRec *cloudflare.DNSRecord, err error) {

    // Most API calls require a Context

	if cfapi == nil {return nil, fmt.Errorf("cfApi is nil!")}
	if cfapi.API == nil {return nil, fmt.Errorf("cfApi.api is nil!")}

	api := cfapi.API

    ctx := context.Background()

    rc := cloudflare.ResourceContainer{
   		Level: cloudflare.ZoneRouteLevel,
		Identifier: zoneId,
	}

    dnsrec, err := api.CreateDNSRecord(ctx, &rc, *dnsPar)
    if err != nil { return nil, fmt.Errorf("cfApi.CreateDNSRecord: %v\n", err)
    }

    return &dnsrec, nil
}

// methods that deletes a DnsRec from a zone with zoneId
func (cfapi *cfApi) DelDnsRec (zoneId, recId string) (err error) {

    api := cfapi.API

    ctx := context.Background()

    rc := cloudflare.ResourceContainer{
   		Level: cloudflare.ZoneRouteLevel,
		Identifier: zoneId,
	}

    err = api.DeleteDNSRecord(ctx, &rc, recId)
    if err != nil {return fmt.Errorf("DeleteDnsRecord: %v", err)}

    return nil
}


// function that creates DNS Challenge record
func (cfapi *cfApi) AddDnsChalRecord (zoneId string, val string) (recId string, err error) {

    // Most API calls require a Context

	if cfapi == nil {return "", fmt.Errorf("cfApi is nil!")}
	if cfapi.API == nil {return "", fmt.Errorf("cfApi.api is nil!")}

	api := cfapi.API

    ctx := context.Background()

    // try to create DNS Record
    dnsPar := cloudflare.CreateDNSRecordParams{
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

    dnsRec, err := api.CreateDNSRecord(ctx, &rc, dnsPar)
    if err != nil { return "", fmt.Errorf("cfApi.AddDNSRecord: %v\n", err)
    }

	recId = dnsRec.ID

    fmt.Printf("success creating Dns Record!\n")
    PrintDnsRec(&dnsRec)

    return recId, nil
}

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

// method that verifies tokens
func (cfapi *cfApi) VerifyToken() (tokInfo *cloudflare.APITokenVerifyBody, err error) {

	api:= cfapi.API
	if api == nil {return nil, fmt.Errorf("no valid api!")}

	res , err := api.VerifyAPIToken(context.Background())
	if err != nil {return &res, fmt.Errorf("VerifyAPIToken %v", err)}

	return &res, nil
}

func PrintTokResp(res *cloudflare.APITokenVerifyBody) {
	fmt.Printf("************ Token Verification ******\n")
	fmt.Printf("ID:       %s\n", res.ID)
	fmt.Printf("Status:   %s\n", res.Status)
	fmt.Printf("NotBefore: %s\n",res.NotBefore.Format(time.RFC1123))
	fmt.Printf("Expires:   %s\n", res.ExpiresOn.Format(time.RFC1123))
}


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

/*
func ReadCsrFil(inFilNam string)(csrDatList *CsrList, err error) {

	//todo check for yaml extension
    bytData, err := os.ReadFile(inFilNam)
    if err != nil {
        return nil, fmt.Errorf("os.ReadFile: %v\n",err)
    }

    csrList := &CsrList{}
    err = yaml.Unmarshal(bytData, csrList)
    if err != nil {
        return nil, fmt.Errorf("yaml Unmarshal: %v\n", err)
    }

//    PrintCsr(CsrList)
	return csrList, nil
}

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
type ApiObj struct {
    Api    string `yaml:"Api"`
    ApiKey string `yaml:"ApiKey"`
    ApiToken string `yaml:"ApiToken"`
	TokenId string `yaml:"token Id"`
	TokName string `yaml:"Token Name"`
	Start time.Time `yaml:"Start"`
	Expiration time.Time `yaml:"Expire"`
	AccountId string `yaml:"AccountId"`
    Email     string `yaml:"Email"`
	YamlFile	string `yaml:"cfToken File"`
}
*/

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

