package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/prom"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/util"
)

func EventHandler(msg []byte) Sysmon {
	s := &Sysmon{}
	e := Winevent{}

	var (
		provider  string
		eventData json.RawMessage
		eventID   uint32
	)
	err := json.Unmarshal(msg, &e)
	if err != nil {
		fmt.Println("Error:", err)

	}
	prom.KafkaEventsTotal.WithLabelValues("winlog").Inc()

	if len(e.SourceName) > 0 {
		provider = e.SourceName
		eventData = e.EventData
		eventID = e.EventID

	} else {
		provider = e.Event.Provider
		eventData = e.Winlog.EventData
		eventID = e.Winlog.EventID
	}
	if provider == "Microsoft-Windows-Sysmon" {

		s.SysmonParser(eventData, eventID, e.Host.Name)

	}
	return *s

}

func (s *Sysmon) SysmonParser(msg []byte, eventid uint32, hostname string) {

	// endTime := util.MakeTimestamp() - startTime
	err := json.Unmarshal(msg, &s)
	if err != nil {
		fmt.Println("Error in decoding sysmon json:", err)

	}
	prom.KafkaEventsTotal.WithLabelValues("sysmon").Inc()
	s.HostName = hostname
	s.EventID = eventid
	s.EventDateCreation, _ = time.Parse("2006-01-02 15:04:05.000", s.UtcTime)
	s.FileDateCreation, _ = time.Parse("2006-01-02 15:04:05.000", s.FileDateCreationStr)
	s.FilePreviousDateCreation, _ = time.Parse("2006-01-02 15:04:05.000", s.FilePreviousDateCreationStr)
	if s.ProcessParentPath != "" {
		s.ProcessParentName = util.RemoveBefor(s.ProcessParentPath, "\\")
	}
	// Split Process Name
	if s.ProcessPath != "" {
		s.ProcessName = util.RemoveBefor(s.ProcessPath, "\\")
	}
	// Split User_domain and UserName
	if s.User != "" {
		s.UserDomain = util.RemoveAfter(s.User, "\\")
		s.UserName = util.RemoveBefor(s.User, "\\")
	}
	if strings.Contains(s.UserName, "\\$$") {
		s.MetaUsernameIsMachine = "true"
	} else {
		s.MetaUsernameIsMachine = "false"
	}
	switch eventid {
	case 1:
		s.Action = "processcreate"
	case 2:
		s.Action = "filecreatetime"
	case 3:
		s.Action = "networkconnect"
		if s.SourceIP != "" {
			if s.SrcIsIPv6 == "false" {

				s.FingerprintNetworkCommunityID = util.CommunityIDHash(s.SourceIP, s.DstIP, s.SrcPort, s.DstPort, s.NetworkProtocol)
				s.SrcIPAddr = s.SourceIP
				s.SrcIPRfc, s.SrcIPPublic, s.SrcIPType, s.SrcIPVersion = util.IPProcess(s.SrcIPAddr)

			} else {
				s.IPv6SrcAddr = s.SourceIP
				s.SrcIPRfc, s.SrcIPPublic, s.SrcIPType, s.SrcIPVersion = util.IPProcess(s.IPv6SrcAddr)
			}
			if s.DstIP != "" {
				if s.DstIsIPv6 == "false" {
					s.DstIPAddr = s.DstIP
					s.DstIPRfc, s.DstIPPublic, s.DstIPType, s.DstIPVersion = util.IPProcess(s.DstIPAddr)

					if err != nil {
						fmt.Println("Error in public ip check: ", err)
					}
					if s.DstIPPublic == "true" {
						s.DstCountry = util.Geocountry(s.DstIPAddr)
						s.DstAsn, s.DstAsnOrg = util.GeoASN(s.DstIPAddr)
					}

				} else {
					s.IPv6DstAddr = s.DstIP
				}
			}
		}
	case 4:
		s.Action = "sysmonservicestatechanged"
	case 6:
		s.Action = "driverload"
	case 7:
		s.Action = "moduleload"
	case 8:
		s.Action = "createremotethread"
	case 9:
		s.Action = "rawaccessread"
	case 10:
		s.Action = "processaccess"
	case 11:
		s.Action = "filecreate"
	case 12, 13, 14:
		s.Action = "registryevent"
	case 15:
		s.Action = "filecreatestreamhash"
	case 16:
		s.Action = "sysmonconfigstatechanged"
	case 17, 18:
		s.Action = "pipeevent"
	case 19, 20, 21:
		s.Action = "wmievent"
	case 22:
		s.Action = "dnsquery"
	}

}
