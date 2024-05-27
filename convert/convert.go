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
	"strconv"
	"strings"
)

// convert a part of a version to a numberic value for comparison
// - return 0 if we get an error
func intPart(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0 // nothing better to return
	}
	return i
}

type Version string

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
// - this is for commands and also for field names
// - commands are expected to be in upper case
// - configuration settings are expected to be in lower case
// - ensure all values are sorted alphabetically
var PreMySQL84Terms = map[string]string{
	"CHANGE REPLICATION SOURCE TO":  "CHANGE MASTER TO",
	"GET_SOURCE_PUBLIC_KEY":         "GET_MASTER_PUBLIC_KEY",
	"PURGE BINARY LOGS":             "PURGE MASTER LOGS",
	"RESET BINARY LOGS AND GTIDS":   "RESET MASTER",
	"RESET REPLICA":                 "RESET SLAVE",
	"SHOW BINARY LOG STATUS":        "SHOW MASTER STATUS",
	"SHOW BINARY LOGS":              "SHOW MASTER LOGS",
	"SHOW REPLICA STATUS":           "SHOW SLAVE STATUS",
	"SHOW REPLICAS":                 "SHOW SLAVE HOSTS",
	"SOURCE_AUTO_POSITION":          "MASTER_AUTO_POSITION",
	"SOURCE_BIND":                   "MASTER_BIND",
	"SOURCE_COMPRESSION_ALGORITHMS": "MASTER_COMPRESSION_ALGORITHMS",
	"SOURCE_CONNECT_RETRY":          "MASTER_CONNECT_RETRY",
	"SOURCE_DELAY":                  "MASTER_DELAY",
	"SOURCE_HEARTBEAT_PERIOD":       "MASTER_HEARTBEAT_PERIOD",
	"SOURCE_HOST":                   "MASTER_HOST",
	"SOURCE_LOG_FILE":               "MASTER_LOG_FILE",
	"SOURCE_LOG_POS":                "MASTER_LOG_POS",
	"SOURCE_PASSWORD":               "MASTER_PASSWORD",
	"SOURCE_PORT":                   "MASTER_PORT",
	"SOURCE_PUBLIC_KEY_PATH":        "MASTER_PUBLIC_KEY_PATH",
	"SOURCE_RETRY_COUNT":            "MASTER_RETRY_COUNT",
	"SOURCE_SSL":                    "MASTER_SSL",
	"SOURCE_SSL_CA":                 "MASTER_SSL_CA",
	"SOURCE_SSL_CAPATH":             "MASTER_SSL_CAPATH",
	"SOURCE_SSL_CERT":               "MASTER_SSL_CERT",
	"SOURCE_SSL_CIPHER":             "MASTER_SSL_CIPHER",
	"SOURCE_SSL_CRL":                "MASTER_SSL_CRL",
	"SOURCE_SSL_CRLPATH":            "MASTER_SSL_CRLPATH",
	"SOURCE_SSL_KEY":                "MASTER_SSL_KEY",
	"SOURCE_SSL_VERIFY_SERVER_CERT": "MASTER_SSL_VERIFY_SERVER_CERT",
	"SOURCE_TLS_CIPHERSUITES":       "MASTER_TLS_CIPHERSUITES",
	"SOURCE_TLS_VERSION":            "MASTER_TLS_VERSION",
	"SOURCE_USER":                   "MASTER_USER",
	"SOURCE_ZSTD_COMPRESSION_LEVEL": "MASTER_ZSTD_COMPRESSION_LEVEL",
	"START REPLICA":                 "START SLAVE",
	"STOP REPLICA":                  "STOP SLAVE",
	"rpl_semi_sync_slave_enabled":   "rpl_semi_sync_slave_enabled",
	"rpl_semi_sync_source_enabled":  "rpl_semi_sync_master_enabled",
}

// Translate replication statements or settings if they are not for 8.4 to the older (original) format
func OldValue(version string, command string) string {
	if !mysql84Version(version) {
		return command
	}

	if replacement, ok := PreMySQL84Terms[strings.ToUpper(command)]; ok {
		return replacement
	}

	return command
}

// mysql84Version is true if the version provided is an 8.4 (compatible version)
func mysql84Version(version string) bool {
	if Version(version).Less(Version("8.4.0")) {
		return false
	}
	if Version(version).Less(Version("10.0.0")) {
		return true
	}
	return false
}

// generic84Pattern will take a version and old and new strings and
// return the appropriate string based on 8.4 vs other versions
//   - for the moment we assume a range of 8.4.0 >= v < 10.0.0 where the
//     new string will be returned
func generic84Pattern(version string, oldName string, newName string) string {
	if Version(version).Less(Version("8.4.0")) {
		return oldName
	}
	if Version(version).Less(Version("10.0.0")) {
		return newName
	}
	return oldName
}

func SourceAutoPosition(version string) string {
	return generic84Pattern(version, "MASTER_AUTO_POSITION", "SOURCE_AUTO_POSITION")
}

func GetSourcePublicKey(version string) string {
	return generic84Pattern(version, "GET_MASTER_PUBLIC_KEY", "GET_SOURCE_PUBLIC_KEY")
}

func ChangeReplicationSourceTo(version string) string {
	return generic84Pattern(version, "CHANGE MASTER TO", "CHANGE REPLICATION SOURCE TO")
}

func ShowReplicaStatus(version string) string {
	return generic84Pattern(version, "SHOW SLAVE STATUS", "SHOW REPLICA STATUS")
}

func ShowBinaryLogStatus(version string) string {
	return generic84Pattern(version, "SHOW MASTER STATUS", "SHOW BINARY LOG STATUS")
}

func StartReplica(version string) string {
	return generic84Pattern(version, "START SLAVE", "START REPLICA")
}
