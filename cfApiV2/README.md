# cfApiV2

## updToken

usage: ./updToken /token=[dns|zones] [/dbg]

This program can create tow types of cloudflare API tokens:
	- write dns records in a zone (domain)
	- read zone records managed by cloudflare

## delToken

usage: ./delToken /token=[dns|zones] [/dbg]

This program deletes the cloudflare API token and the token file.

## listToken

usage: ./listToken [/dbg]

A program that lists all cloudflare API tokens for the account.

## creZoneFile

usage: /creZoneFile

This program creates a file named "zones/ZoneList.yaml".
The yaml file is a list of all zone records managed by cloudflare

## listPermGroups

usage: /listPermGroups

This program lists all Permission Groups of cloudflare tokens
PermGroups are necessary for the cration of tokens.

## createTokenFile

program that creates a json token file in the token directory.
Inputs are the token value and the token file name.

## readTokenFile

usage: ./readTokenFile /token=[zones|dns] /vfy

program reads token file
/vfy flag verifies token and if token is not valid removes the token file in the token folder.


## cfLibV2

### CreateTokFile
func CreateTokFile(tokFilnam string, Token *cloudflare.APIToken, dbg bool) (err error)

### ReadDnsRecord
func ReadDnsRec(zoneId string, token string, dnspar cloudflare.ListDNSRecordsParams)(dnsrec []cloudflare.DnsRecord, err error)

function that retrieves all dns records of a zone with id zoneId.
dnspar can be specified to reduce the search universe.

### AddDnsRecord
func AddDnsRecord(param *cloudflare.CreateDNSRecordParams, token string)(recId string, err error)

### AddDnsChalRecord
func AddDnsChalRecord(zoneId, val, token string)(recId string, err error)

### UpdDnsRecord
func UpdDnsRecord(param *cloudflare.UpdateDNSRecordParams, token string)(recId string, err error)

### DelDnsRecord
func DelDnsRecord (zoneId, recId, token string) (err error)
function that deletes a dns record from the zone with id zoneId.

### FindAcmeChalRecord
func FindAcmeChalRecord(dnsrec []cloudflare.DNSRecord)(recId string, err error)
function that retrieves all dns records with the name that includes the phrase <_acme-challenge>

### DelDnsChalRecord
func (cfapi *cfApi) DelDnsChalRecord (zone ZoneAcme) (err error)
tobe deleted

### SaveZonesJson
func SaveZonesJson(zones []cloudflare.Zone, outfil *os.File)(err error)

### SaveZonesYaml
func SaveZonesYaml(zones []cloudflare.Zone, outfil *os.File)(err error)

### SaveZonesShortJson
func SaveZonesShortJson(zones []ZoneShort, outfil *os.File)(err error)

### SaveAcmeDns
func SaveAcmeDns(zones []ZoneAcme, outfil *os.File)(err error)

### ReadZonesShortYaml
func ReadZonesShortYaml(infil *os.File)(zoneListObj *[]ZoneShort, err error)

### ReadAcmeZones
func ReadAcmeZones(inFilNam string)(zoneListObj *[]ZoneAcme, err error

### ReadZoneShortFile
func ReadZoneShortFile(inFilNam string)(zoneList *ZoneList, err error)

### SaveTokList
func SaveTokList(outFilnam string, tokList []cloudflare.APIToken) (err error)

### PrintTokResp
func PrintTokResp(res *cloudflare.APITokenVerifyBody)

token response is an object that is returned from the cloudflare api method VerifyAPIToken

### PrintZones
func PrintZones(zones []cloudflare.Zone)

Creates a list of zones with the fields xzone name and zone id.

### PrintZoneList
func PrintZoneList(zoneList *ZoneList)

prints a zonelist 
A zonelist is an object that contains fields for:
	- Accountid
	- account email
	- modified time
plus a slice of all zones (see above)

### PrintUserInfo
func PrintUserInfo (u *cloudflare.User)

function that prints out all field of a cloudflare user struct (cloudflare.User)
cloudflare api: func (api *API) UserDetails(ctx context.Context) (User, error)

### PrintResInfo
func PrintResInfo(res *cloudflare.ResultInfo)

function that prints out the cloudflare ResultInfo struct (cloudflare.ResultInfo)
ResultInfo contains metadata about the Response.

### PrintDnsRecs
func PrintDnsRecs(recs *[]cloudflare.DNSRecord)

function that prints a list of DNS records
function that prints out the cloudflare DNSRecord fields (cloudflare.DNSRec)

func (api *API) CreateDNSRecord(ctx context.Context, rc *ResourceContainer, params CreateDNSRecordParams) (DNSRecord, error)


func ZoneIdentifier(id string) *ResourceContainer

### PrintAcmeZones
func PrintAcmeZones(zones []ZoneAcme)

### PrintAccount
func PrintAccount(act *cloudflare.Account)

### PrintTokList
func PrintTokList(tokList []cloudflare.APIToken)

### PrintToken
func PrintToken(tok cloudflare.APIToken)

function that prints out all fields of a cloudflare APItoken
