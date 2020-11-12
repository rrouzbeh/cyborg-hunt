package click

import (
	"database/sql"
	"os"

	"github.com/prometheus/common/log"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/parser"
	"github.com/rrouzbeh/cyborg-hunt/src/winlog/pkg/prom"
)

// SysmonPrepare prepare data for insert
func SysmonPrepare(db *sql.DB) (*sql.Tx, *sql.Stmt) {
	Tx, _ := db.Begin()
	Stmt, err := Tx.Prepare("INSERT INTO cyborg.sysmon (event_date_creation, host_name, process_name, process_path, process_id, process_guid, process_current_directory, process_parent_name, process_integrity_level, process_parent_command_line, process_parent_path, process_parent_guid, process_command_line, src_port, src_ip_addr, src_ip_public, src_ip_type, src_ip_rfc, src_is_ipv6, src_host_name, src_port_name, src_ip_version, ipv6_src_addr, dst_ip_addr, dst_ip_type, dst_ip_public, dst_ip_rfc, dst_ip_version, dst_host_name, dst_port, dst_port_name, dst_asn, dst_asn_org, dst_country, dst_is_ipv6, ipv6_dst_addr, dns_query_name, dns_query_status, dns_query_results, meta_user_name_is_machine, network_protocol, tag, action, user_session_id, registry_key_value, event_type, event_id, registry_key_path, sysmon_version, file_name, file_product, fingerprint_network_community_id, thread_id, file_description, user_name, network_initiated, hash, service_state, file_version, file_previous_date_creation, user_logon_guid, file_name_original, user_domain, user_logon_id, sysmon_schema_version, file_company, file_date_creation, computer_name) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Errorf("Prepare Error: %s", err)
		os.Exit(0)
	}
	return Tx, Stmt
}

// SysmonInsert Insert event to clickhouse
func SysmonInsert(Stmt *sql.Stmt, s parser.Sysmon) (err error) {

	_, err = Stmt.Exec(s.EventDateCreation, s.HostName, s.ProcessName, s.ProcessPath, s.ProcessID, s.ProcessGUID, s.ProcessCurrentDirectory, s.ProcessParentName, s.ProcessIntegrityLevel, s.ProcessParentCommandLine, s.ProcessParentPath, s.ProcessParentGUID, s.ProcessCommandLine, s.SrcPort, s.SrcIPAddr, s.SrcIPPublic, s.SrcIPType, s.SrcIPRfc, s.SrcIsIPv6, s.SrcHostName, s.SrcPortName, s.SrcIPVersion, s.IPv6SrcAddr, s.DstIPAddr, s.DstIPType, s.DstIPPublic, s.DstIPRfc, s.DstIPVersion, s.DstHostname, s.DstPort, s.DstPortName, s.DstAsn, s.DstAsnOrg, s.DstCountry, s.DstIsIPv6, s.IPv6DstAddr, s.DNSQueryName, s.DNSQueryStatus, s.DNSQueryResults, s.MetaUsernameIsMachine, s.NetworkProtocol, s.Tag, s.Action, s.UserSessionID, s.RegistryKeyValue, s.EventType, s.EventID, s.RegistryKeyPath, s.SysmonVersion, s.FileName, s.FileProduct, s.FingerprintNetworkCommunityID, s.ThreadID, s.FileDescription, s.UserName, s.NetworkInitiated, s.Hash, s.ServiceState, s.FileVersion, s.FilePreviousDateCreation, s.UserLogonGUID, s.FileNameOriginal, s.UserDomain, s.UserLogonID, s.SysmonSchemaVersion, s.FileCompany, s.FileDateCreation, s.ComputerName)
	if err != nil {
		log.Errorf("Insert Error: %s", err)

	}
	return err
}

// SysmonCommit write events to Clickhouse
func SysmonCommit(batch []parser.Sysmon) {
	err, Db := Connect()
	if err != nil {
		log.Errorf("Error in Connecting to Clickhouse: %s", err)
	}

	Tx, Stmt := SysmonPrepare(Db)

	for _, s := range batch {
		prom.ClickhouseEventsTotal.WithLabelValues("sysmon").Inc()
		if err := SysmonInsert(Stmt, s); err != nil {
			log.Errorf("Insert Error: %s", err)
			prom.ClickhouseEventsErrors.WithLabelValues("sysmon").Inc()
		}
		prom.ClickhouseEventsSuccess.WithLabelValues("sysmon").Inc()

	}
	err = Tx.Commit()

	if err != nil {
		log.Errorf("Commit Error: %s", err)

		Db.Close()
	}

	Db.Close()
}
