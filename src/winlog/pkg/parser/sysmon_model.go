package parser

import (
	"encoding/json"
	"time"
)

type Sysmon struct {
	EventID                       uint32
	ComputerName                  string
	Tag                           string `json:"RuleName"`
	HostName                      string `json:"host_name"`
	UtcTime                       string `json:"UtcTime"`
	EventDateCreation             time.Time
	ProcessGUID                   string `json:"ProcessGuid"`
	ProcessID                     string `json:"ProcessId"`
	ProcessPath                   string `json:"Image"`
	ProcessName                   string
	FileVersion                   string `json:"FileVersion"`
	FileDescription               string `json:"Description"`
	FileProduct                   string `json:"Product"`
	FileCompany                   string `json:"Company"`
	FileNameOriginal              string `json:"OriginalFileName"`
	ProcessCommandLine            string `json:"CommandLine"`
	FileCurrentDirectory          string `json:"CurrentDirectory"`
	User                          string `json:"User"`
	UserLogonGUID                 string `json:"LogonGuid"`
	UserLogonID                   string `json:"LogonId"`
	UserSessionID                 string `json:"TerminalSessionId"`
	ProcessIntegrityLevel         string `json:"IntegrityLevel"`
	Hash                          string `json:"Hashes"`
	ProcessParentGUID             string `json:"ParentProcessGuid"`
	ProcessParentID               string `json:"ParentProcessId"`
	ProcessParentName             string
	ProcessParentPath             string `json:"ParentImage"`
	ProcessParentCommandLine      string `json:"ParentCommandLine"`
	ProcessCurrentDirectory       string `json:"CurrentDirectory"`
	FileName                      string `json:"TargetFilename"`
	DstHostname                   string `json:"DestinationHostname"`
	DstPort                       string `json:"DestinationPort"`
	DstPortName                   string `json:"DestinationPortName"`
	DstIsIPv6                     string `json:"DestinationIsIpv6"`
	NetworkInitiated              string `json:"Initiated"`
	NetworkProtocol               string `json:"Protocol"`
	SrcHostName                   string `json:"SourceHostname"`
	SrcPort                       string `json:"SourcePort"`
	SrcPortName                   string `json:"SourcePortName"`
	SrcIsIPv6                     string `json:"SourceIsIpv6"`
	SrcIPAddr                     string
	SrcIPRfc                      string
	SrcIPPublic                   string
	SrcIPType                     string
	SrcIPVersion                  string
	IPv6SrcAddr                   string
	SourceIP                      string `json:"SourceIp"`
	DstIP                         string `json:"DestinationIp"`
	DstIPAddr                     string
	DstIPRfc                      string
	DstIPPublic                   string
	DstIPType                     string
	DstIPVersion                  string
	DstCountry                    string
	DstAsn                        uint32
	DstAsnOrg                     string
	IPv6DstAddr                   string
	ServiceState                  string `json:"State"`
	SysmonVersion                 string `json:"Version"`
	SysmonSchemaVersion           string `json:"SchemaVersion"`
	DriverLoaded                  string `json:"ImageLoaded"`
	Signature                     string `json:"Signature"`
	SignatureStatus               string `json:"SignatureStatus"`
	Signed                        string `json:"Signed"`
	ThreadNewID                   string `json:"NewThreadId"`
	ThreadStartAddress            string `json:"StartAddress"`
	ThreadStartFunction           string `json:"StartFunction"`
	ThreadStartModule             string `json:"StartModule"`
	DeviceName                    string `json:"Device"`
	ProcessCallTrace              string `json:"CallTrace"`
	ProcessGrantedAccess          string `json:"GrantedAccess"`
	ThreadID                      string `json:"SourceThreadId"`
	EventType                     string `json:"EventType"`
	RegistryKeyPath               string `json:"TargetObject"`
	RegistryKeyValue              string `json:"Details"`
	RegistryKeyNewName            string `json:"NewName"`
	SysmonHash                    string `json:"hash"`
	SysmonConfiguration           string `json:"Configuration"`
	PipeName                      string `json:"PipeName"`
	WmiOperation                  string `json:"Operation"`
	WmiNamespace                  string `json:"EventNamespace"`
	WmiFilterName                 string `json:"Name"`
	WmiQuery                      string `json:"Query"`
	WmiConsumerType               string `json:"Type"`
	WmiConsumerDestination        string `json:"Destination"`
	WmiConsumerPath               string `json:"Consumer"`
	WmiFilterPath                 string `json:"Filter"`
	DNSQueryName                  string `json:"QueryName"`
	DNSQueryStatus                string `json:"QueryStatus"`
	DNSQueryResults               string `json:"QueryResults"`
	FileDateCreation              time.Time
	FileDateCreationStr           string `json:"CreationUtcTime"`
	FilePreviousDateCreation      time.Time
	FilePreviousDateCreationStr   string `json:"PreviousCreationUtcTime"`
	Action                        string
	MetaUsernameIsMachine         string
	UserDomain                    string
	UserName                      string
	FingerprintNetworkCommunityID string
}
type Winevent struct {
	EventDateCreation time.Time
	Message           string `json:"message"`
	EventHash         string

	ProviderName string `json:"provider_name"`
	SourceName   string `json:"source_name"`
	Timestamp    string `json:"@timestamp"`
	Winlog       struct {
		EventData    json.RawMessage `json:"event_data"`
		ComputerName string          `json:"computer_name"`
		EventID      uint32          `json:"event_id"`
	} `json:"winlog"`
	EventID      uint32          `json:"event_id"`
	EventData    json.RawMessage `json:"event_data"`
	ComputerName string          `json:"computer_name"`

	Metadata struct {
		Beat    string `json:"beat"`
		Version string `json:"version"`
	} `json:"@metadata"`
	Event struct {
		Provider string `json:"provider"`
	} `json:"event"`
	Host struct {
		Name string `json:"name"`
	} `json:"Host"`
	Index uint
}
