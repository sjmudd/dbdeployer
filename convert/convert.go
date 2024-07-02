// DBDeployer - The MySQL Sandbox
// Copyright © 2024 Simon J Mudd <sjmudd@pobox.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Routines for handling the move from old replication terminology
// - master / slave
// to new replication terminology
// - source / replica

package convert

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/datacharmer/dbdeployer/common"
)

// intPart returns the numeric integer part of a version string to be used in comparisons
// - return 0 if the conversion to an integer fails
func intPart(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0 // nothing better to return
	}
	return i
}

// Version is expected to a 3 digit version such as 8.4.0, 8.0.37,
// 5.7.56, ... used for version comparisons
type Version string

// Less returns true if v < other
func (v Version) Less(other Version) bool {
	v1List := strings.Split(string(v), ".")
	v2List := strings.Split(string(other), ".")

	for i := range v1List {
		if i > len(v2List)-1 {
			return false
		}
		if intPart(v1List[i]) < intPart(v2List[i]) {
			return true
		}
	}
	return false
}

// Equal returns true if the "." split parts of the versions match
func (v Version) Equal(other Version) bool {
	v1List := strings.Split(string(v), ".")
	v2List := strings.Split(string(other), ".")

	if len(v1List) != len(v2List) {
		return false
	}

	for i := range v1List {
		if intPart(v1List[i]) == intPart(v2List[i]) {
			return true
		}
	}
	return false
}

// Map MySQL 8.4+ terms to the pre-8.4 equivalents
// - related to replication settings.
// - this is for:
//   - commands
//   - column names
//   - /etc/my.cnf or similar server settings
// - SQL commands are expected to be in UPPER CASE
// - variable settings are expected to be in lower case.
// - columns values must match the value provided by MySQL (and
//   some are mixed case)
// - ensure all values are sorted alphabetically

type OldAndNewNames struct {
	OldName string
	NewName string
}
type OldAndNewNamesMap map[string]OldAndNewNames

var (
	ToPreMySQL84Values = OldAndNewNamesMap{
		// commands
		"ChangeMasterCmd":     {"CHANGE MASTER TO", "CHANGE REPLICATION SOURCE TO"},
		"PurgeMasterLogsCmd":  {"PURGE MASTER LOGS", "PURGE BINARY LOGS"},
		"ResetMasterCmd":      {"RESET MASTER", "RESET BINARY LOGS AND GTIDS"},
		"ResetSlaveCmd":       {"RESET SLAVE", "RESET REPLICA"},
		"ShowMasterStatusCmd": {"SHOW MASTER STATUS", "SHOW BINARY LOG STATUS"},
		"ShowMasterLogsCmd":   {"SHOW MASTER LOGS", "SHOW BINARY LOGS"},
		"ShowSlaveStatusCmd":  {"SHOW SLAVE STATUS", "SHOW REPLICA STATUS"},
		"ShowSlaveHostsCmd":   {"SHOW SLAVE HOSTS", "SHOW REPLICAS"},
		"StartSlaveCmd":       {"START SLAVE", "START REPLICA"},
		"StopSlaveCmd":        {"STOP SLAVE", "STOP REPLICA"},
		// terms
		"GET_MASTER_PUBLIC_KEY":         {"GET_MASTER_PUBLIC_KEY", "GET_SOURCE_PUBLIC_KEY"},
		"MASTER_AUTO_POSITION":          {"MASTER_AUTO_POSITION", "SOURCE_AUTO_POSITION"},
		"MASTER_BIND":                   {"MASTER_BIND", "SOURCE_BIND"},
		"MASTER_COMPRESSION_ALGORITHMS": {"MASTER_COMPRESSION_ALGORITHMS", "SOURCE_COMPRESSION_ALGORITHMS"},
		"MASTER_CONNECT_RETRY":          {"MASTER_CONNECT_RETRY", "SOURCE_CONNECT_RETRY"},
		"MASTER_DELAY":                  {"MASTER_DELAY", "SOURCE_DELAY"},
		"MASTER_HEARTBEAT_PERIOD":       {"MASTER_HEARTBEAT_PERIOD", "SOURCE_HEARTBEAT_PERIOD"},
		"MASTER_HOST":                   {"MASTER_HOST", "SOURCE_HOST"},
		"MASTER_LOG_FILE":               {"MASTER_LOG_FILE", "SOURCE_LOG_FILE"},
		"MASTER_LOG_POS":                {"MASTER_LOG_POS", "SOURCE_LOG_POS"},
		"MASTER_PASSWORD":               {"MASTER_PASSWORD", "SOURCE_PASSWORD"},
		"MASTER_PORT":                   {"MASTER_PORT", "SOURCE_PORT"},
		"MASTER_PUBLIC_KEY_PATH":        {"MASTER_PUBLIC_KEY_PATH", "SOURCE_PUBLIC_KEY_PATH"},
		"MASTER_RETRY_COUNT":            {"MASTER_RETRY_COUNT", "SOURCE_RETRY_COUNT"},
		"MASTER_SSL":                    {"MASTER_SSL", "SOURCE_SSL"},
		"MASTER_SSL_CA":                 {"MASTER_SSL_CA", "SOURCE_SSL_CA"},
		"MASTER_SSL_CAPATH":             {"MASTER_SSL_CAPATH", "SOURCE_SSL_CAPATH"},
		"MASTER_SSL_CERT":               {"MASTER_SSL_CERT", "SOURCE_SSL_CERT"},
		"MASTER_SSL_CIPHER":             {"MASTER_SSL_CIPHER", "SOURCE_SSL_CIPHER"},
		"MASTER_SSL_CRL":                {"MASTER_SSL_CRL", "SOURCE_SSL_CRL"},
		"MASTER_SSL_CRLPATH":            {"MASTER_SSL_CRLPATH", "SOURCE_SSL_CRLPATH"},
		"MASTER_SSL_KEY":                {"MASTER_SSL_KEY", "SOURCE_SSL_KEY"},
		"MASTER_SSL_VERIFY_SERVER_CERT": {"MASTER_SSL_VERIFY_SERVER_CERT", "SOURCE_SSL_VERIFY_SERVER_CERT"},
		"MASTER_TLS_CIPHERSUITES":       {"MASTER_TLS_CIPHERSUITES", "SOURCE_TLS_CIPHERSUITES"},
		"MASTER_TLS_VERSION":            {"MASTER_TLS_VERSION", "SOURCE_TLS_VERSION"},
		"MASTER_USER":                   {"MASTER_USER", "SOURCE_USER"},
		"MASTER_ZSTD_COMPRESSION_LEVEL": {"MASTER_ZSTD_COMPRESSION_LEVEL", "SOURCE_ZSTD_COMPRESSION_LEVEL"},
		// variables
		"rpl_semi_sync_master":         {"rpl_semi_sync_master", "rpl_semi_sync_source"},
		"rpl_semi_sync_master_enabled": {"rpl_semi_sync_master_enabled", "rpl_semi_sync_source_enabled"},
		"rpl_semi_sync_slave":          {"rpl_semi_sync_slave", "rpl_semi_sync_replica"},
		"rpl_semi_sync_slave_enabled":  {"rpl_semi_sync_slave_enabled", "rpl_semi_sync_replica_enabled"},
		"semisync_master":              {"semisync_master", "semisync_source"},
		"semisync_slave":               {"semisync_slave", "semisync_replica"},
		// replication fields
		"Master_Log_File":   {"Master_Log_File", "Source_Log_File"},
		"Master_Log_Pos":    {"Master_Log_Pos", "Source_Log_Pos"},
		"Slave_IO_Running":  {"Slave_IO_Running", "Replica_IO_Running"},
		"Slave_SQL_Running": {"Slave_SQL_Running", "Replica_SQL_Running"},
		// replication functions
		"MASTER_POS_WAIT": {"MASTER_POS_WAIT", "SOURCE_POS_WAIT"},

		"log-slave-updates": {"log-slave-updates", "log-replica-updates"},
	}
)

// MySQLNamesByVersion returns a map of K,V MySQL 8.4 or non-8.4
// settings based on the provided version string
func MySQLNamesByVersion(version string) common.StringMap {
	var (
		isMySQL84 = isMySQL84Version(version)
		names     = make(common.StringMap)
	)

	for key, value := range ToPreMySQL84Values {
		if isMySQL84 {
			names[key] = value.NewName
		} else {
			names[key] = value.OldName
		}
	}
	return names
}

// isMySQL84Version is true if the version provided is an 8.4 (compatible version)
func isMySQL84Version(version string) bool {
	if Version(version).Less(Version("8.4.0")) {
		return false
	}
	if Version(version).Less(Version("10.0.0")) {
		return true
	}
	return false
}

// Given the version and key lookup the value and return the new or old version account to the version.
// - if the key is not fuond return a variable to indicate this.
func VersionedValue(version string, key string) (string, error) {
	isMySQL84 := isMySQL84Version(version)

	for k, value := range ToPreMySQL84Values {
		if k == key {
			if isMySQL84 {
				return value.NewName, nil
			} else {
				return value.OldName, nil
			}
		}
	}
	return "", fmt.Errorf("VersionedValue(%v,%v): key not found", version, key)
}
