// DBDeployer - The MySQL Sandbox
// Copyright © 2006-2021 Giuseppe Maxia
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

//go:build go1.16
// +build go1.16

package sandbox

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/sjmudd/dbdeployer/common"
	"github.com/sjmudd/dbdeployer/defaults"
	"github.com/sjmudd/dbdeployer/globals"
)

type TemplateDesc struct {
	TemplateInFile bool
	Description    string
	Notes          string
	Contents       string
}

type TemplateCollection map[string]TemplateDesc
type AllTemplateCollection map[string]TemplateCollection

// templates for single sandbox

var (
	//go:embed templates/single/init_db.gotxt
	initDbTemplate string

	//go:embed templates/single/start.gotxt
	startTemplate string

	//go:embed templates/single/wipe_and_restart.gotxt
	wipeAndRestartTemplate string

	//go:embed templates/single/use.gotxt
	useTemplate string

	//go:embed templates/single/use_admin.gotxt
	useAdminTemplate string

	//go:embed templates/single/sysbench.gotxt
	sysbenchTemplate string

	//go:embed templates/single/sysbench_ready.gotxt
	sysbenchReadyTemplate string

	//go:embed templates/single/mysqlsh.gotxt
	mysqlshTemplate string

	//go:embed templates/single/stop.gotxt
	stopTemplate string

	//go:embed templates/single/clear.gotxt
	clearTemplate string

	//go:embed templates/single/my_cnf.gotxt
	myCnfTemplate string

	//go:embed templates/single/send_kill.gotxt
	sendKillTemplate string

	//go:embed templates/single/status.gotxt
	statusTemplate string

	//go:embed templates/single/restart.gotxt
	restartTemplate string

	//go:embed templates/single/load_grants.gotxt
	loadGrantsTemplate string

	//go:embed templates/single/grants5x.gotxt
	grantsTemplate5x string

	//go:embed templates/single/grants57.gotxt
	grantsTemplate57 string

	//go:embed templates/single/task_user_grants.gotxt
	grantsTaskUserTemplate string

	//go:embed templates/single/grants8x.gotxt
	grantsTemplate8x string

	//go:embed templates/single/add_option.gotxt
	addOptionTemplate string

	//go:embed templates/single/show_log.gotxt
	showLogTemplate string

	//go:embed templates/single/show_binlog.gotxt
	showBinlogTemplate string

	//go:embed templates/single/my.gotxt
	myTemplate string

	//go:embed templates/single/show_relaylog.gotxt
	showRelaylogTemplate string

	//go:embed templates/single/test_sb.gotxt
	testSbTemplate string

	//go:embed templates/single/replication_options.gotxt
	replicationOptions string

	//go:embed templates/single/semisync_master_options.gotxt
	semisyncMasterOptions string

	//go:embed templates/single/semisync_slave_options.gotxt
	semisyncSlaveOptions string

	//go:embed templates/single/repl_crash_safe_options.gotxt
	replCrashSafeOptions string

	//go:embed templates/single/gtid_options_56.gotxt
	gtidOptions56 string

	//go:embed templates/single/gtid_options_57.gotxt
	gtidOptions57 string

	//go:embed templates/single/expose_dd_tables.gotxt
	exposeDdTables string

	//go:embed templates/single/sb_locked.gotxt
	sbLockedTemplate string

	//go:embed templates/mock/no_op_mock.gotxt
	noOpMockTemplate string

	//go:embed templates/mock/mysqld_safe_mock.gotxt
	mysqldSafeMockTemplate string

	//go:embed templates/mock/tidb_mock.gotxt
	tidbMockTemplate string

	//go:embed templates/single/after_start.gotxt
	afterStartTemplate string

	//go:embed templates/single/sb_include.gotxt
	sbIncludeTemplate string

	//go:embed templates/single/connection_info_sql.gotxt
	connectionInfoSql string

	//go:embed templates/single/connection_info_json.gotxt
	ConnectionInfoJson string

	//go:embed templates/single/connection_info_super_json.gotxt
	ConnectionInfoSuperJson string

	//go:embed templates/single/connection_info_conf.gotxt
	ConnectionInfoConf string

	//go:embed templates/single/connection_info_super_conf.gotxt
	ConnectionInfoSuperConf string

	//go:embed templates/single/clone_connection_sql.gotxt
	cloneConnectionSql string

	//go:embed templates/single/clone_from.gotxt
	cloneFromTemplate string

	//go:embed templates/single/replicate_from.gotxt
	replicateFromTemplate string

	//go:embed templates/single/metadata.gotxt
	metadataTemplate string

	SingleTemplates = TemplateCollection{
		TmplCopyright: TemplateDesc{
			Description: "Copyright for every sandbox script",
			Notes:       "",
			Contents:    globals.ShellScriptCopyright,
		},
		TmplReplicationOptions: TemplateDesc{
			Description: "Replication options for my.cnf",
			Notes:       "",
			Contents:    replicationOptions,
		},
		TmplSemisyncMasterOptions: TemplateDesc{
			Description: "master semi-synch options for my.cnf",
			Notes:       "",
			Contents:    semisyncMasterOptions,
		},
		TmplSemisyncSlaveOptions: TemplateDesc{
			Description: "slave semi-synch options for my.cnf",
			Notes:       "",
			Contents:    semisyncSlaveOptions,
		},
		TmplGtidOptions56: TemplateDesc{
			Description: "GTID options for my.cnf 5.6.x",
			Notes:       "",
			Contents:    gtidOptions56,
		},
		TmplGtidOptions57: TemplateDesc{
			Description: "GTID options for my.cnf 5.7.x and 8.0",
			Notes:       "",
			Contents:    gtidOptions57,
		},
		TmplReplCrashSafeOptions: TemplateDesc{
			Description: "Replication crash safe options",
			Notes:       "",
			Contents:    replCrashSafeOptions,
		},
		TmplExposeDdTables: TemplateDesc{
			Description: "Commands needed to enable data dictionary table usage",
			Notes:       "",
			Contents:    exposeDdTables,
		},
		TmplInitDb: TemplateDesc{
			Description: "Initialization template for the database",
			Notes:       "This should normally run only once",
			Contents:    initDbTemplate,
		},
		TmplStart: TemplateDesc{
			Description: "starts the database in a single sandbox (with optional mysqld arguments)",
			Notes:       "",
			Contents:    startTemplate,
		},
		TmplUse: TemplateDesc{
			Description: "Invokes the MySQL client with the appropriate options",
			Notes:       "",
			Contents:    useTemplate,
		},
		TmplUseAdmin: TemplateDesc{
			Description: "Invokes the MySQL client as admin",
			Notes:       "For MySQL 8.0.14+",
			Contents:    useAdminTemplate,
		},
		TmplSysbench: TemplateDesc{
			Description: "Invokes the sysbench tool with custom defined options",
			Notes:       "Requires sysbench to be installed",
			Contents:    sysbenchTemplate,
		},
		TmplSysbenchReady: TemplateDesc{
			Description: "Invokes the sysbench tool with predefined actions",
			Notes:       "Requires sysbench to be installed",
			Contents:    sysbenchReadyTemplate,
		},
		TmplMysqlsh: TemplateDesc{
			Description: "Invokes the MySQL shell with an appropriate URI",
			Notes:       "",
			Contents:    mysqlshTemplate,
		},
		TmplStop: TemplateDesc{
			Description: "Stops a database in a single sandbox",
			Notes:       "",
			Contents:    stopTemplate,
		},
		TmplClear: TemplateDesc{
			Description: "Remove all data from a single sandbox",
			Notes:       "",
			Contents:    clearTemplate,
		},
		TmplMyCnf: TemplateDesc{
			Description: "Default options file for a sandbox",
			Notes:       "",
			Contents:    myCnfTemplate,
		},
		TmplStatus: TemplateDesc{
			Description: "Shows the status of a single sandbox",
			Notes:       "",
			Contents:    statusTemplate,
		},
		TmplRestart: TemplateDesc{
			Description: "Restarts the database (with optional mysqld arguments)",
			Notes:       "",
			Contents:    restartTemplate,
		},
		TmplSendKill: TemplateDesc{
			Description: "Sends a kill signal to the database",
			Notes:       "",
			Contents:    sendKillTemplate,
		},
		TmplLoadGrants: TemplateDesc{
			Description: "Loads the grants defined for the sandbox",
			Notes:       "",
			Contents:    loadGrantsTemplate,
		},
		TmplGrants5x: TemplateDesc{
			Description: "Grants for sandboxes up to 5.6",
			Notes:       "",
			Contents:    grantsTemplate5x,
		},
		TmplGrants57: TemplateDesc{
			Description: "Grants for sandboxes from 5.7+",
			Notes:       "",
			Contents:    grantsTemplate57,
		},
		TmplGrants8x: TemplateDesc{
			Description: "Grants for sandboxes from 8.0+",
			Notes:       "",
			Contents:    grantsTemplate8x,
		},
		TmplTaskUserGrants: TemplateDesc{
			Description: "Grants for task user (8.0+)",
			Notes:       "",
			Contents:    grantsTaskUserTemplate,
		},
		TmplMy: TemplateDesc{
			Description: "Prefix script to run every my* command line tool",
			Notes:       "",
			Contents:    myTemplate,
		},
		TmplAddOption: TemplateDesc{
			Description: "Adds options to the my.sandbox.cnf file and restarts",
			Notes:       "",
			Contents:    addOptionTemplate,
		},
		TmplShowLog: TemplateDesc{
			Description: "Shows error log or custom log",
			Notes:       "",
			Contents:    showLogTemplate,
		},
		TmplShowBinlog: TemplateDesc{
			Description: "Shows a binlog for a single sandbox",
			Notes:       "",
			Contents:    showBinlogTemplate,
		},
		TmplShowRelaylog: TemplateDesc{
			Description: "Show the relaylog for a single sandbox",
			Notes:       "",
			Contents:    showRelaylogTemplate,
		},
		TmplTestSb: TemplateDesc{
			Description: "Tests basic sandbox functionality",
			Notes:       "",
			Contents:    testSbTemplate,
		},
		TmplSbLocked: TemplateDesc{
			Description: "locked sandbox script",
			Notes:       "This script is replacing 'clear' when the sandbox is locked",
			Contents:    sbLockedTemplate,
		},
		TmplAfterStart: TemplateDesc{
			Description: "commands to run after the database started",
			Notes:       "This script does nothing. You can change it and reuse through --use-template",
			Contents:    afterStartTemplate,
		},
		TmplSbInclude: TemplateDesc{
			Description: "Common variables and routines for sandboxes scripts",
			Notes:       "",
			Contents:    sbIncludeTemplate,
		},
		TmplConnectionInfoSql: TemplateDesc{
			Description: "connection info to replicate from this sandbox",
			Notes:       "",
			Contents:    connectionInfoSql,
		},
		TmplConnectionInfoConf: TemplateDesc{
			Description: "connection info to replicate from this sandbox (.conf)",
			Notes:       "",
			Contents:    ConnectionInfoConf,
		},
		TmplConnectionInfoSuperConf: TemplateDesc{
			Description: "connection info use this sandbox as super user (.conf)",
			Notes:       "",
			Contents:    ConnectionInfoSuperConf,
		},
		TmplConnectionInfoJson: TemplateDesc{
			Description: "connection info to replicate from this sandbox (.json)",
			Notes:       "",
			Contents:    ConnectionInfoJson,
		},
		TmplConnectionInfoSuperJson: TemplateDesc{
			Description: "connection info to use this sandbox as super user (.json)",
			Notes:       "",
			Contents:    ConnectionInfoSuperJson,
		},
		TmplReplicateFrom: TemplateDesc{
			Description: "starts replication from another sandbox",
			Notes:       "",
			Contents:    replicateFromTemplate,
		},
		TmplCloneConnectionSql: TemplateDesc{
			Description: "connection info to clone from this sandbox",
			Notes:       "",
			Contents:    cloneConnectionSql,
		},
		TmplCloneFrom: TemplateDesc{
			Description: "clone from another sandbox",
			Notes:       "",
			Contents:    cloneFromTemplate,
		},
		TmplMetadata: TemplateDesc{
			Description: "shows data about the sandbox",
			Notes:       "",
			Contents:    metadataTemplate,
		},
		TmplWipeAndRestart: TemplateDesc{
			Description: "wipe the database and re-create it",
			Notes:       "",
			Contents:    wipeAndRestartTemplate,
		},
	}
	MockTemplates = TemplateCollection{
		TmplNoOpMock: TemplateDesc{
			Description: "mock script that does nothing",
			Notes:       "Used for internal tests",
			Contents:    noOpMockTemplate,
		},
		TmplMysqldSafeMock: TemplateDesc{
			Description: "mock script for mysqld_safe",
			Notes:       "Used for internal tests",
			Contents:    mysqldSafeMockTemplate,
		},
		TmplTidbMock: TemplateDesc{
			Description: "mock script for tidb-server",
			Notes:       "Used for internal tests",
			Contents:    tidbMockTemplate,
		},
	}

	AllTemplates = AllTemplateCollection{
		"mock":        MockTemplates,
		"single":      SingleTemplates,
		"tidb":        TidbTemplates,
		"import":      ImportTemplates,
		"multiple":    MultipleTemplates,
		"replication": ReplicationTemplates,
		"group":       GroupTemplates,
		"pxc":         PxcTemplates,
		"ndb":         NdbTemplates,
	}
)

func FillMockTemplates() error {
	data := defaults.DefaultsToMap()
	for name, template := range MockTemplates {
		tempString, err := common.SafeTemplateFill(name, template.Contents, data)
		if err != nil {
			return fmt.Errorf("error initializing mock template %s", name)
		}
		MockTemplates[name] = TemplateDesc{
			Description:    template.Description,
			TemplateInFile: template.TemplateInFile,
			Notes:          template.Notes,
			Contents:       tempString,
		}
	}
	return nil
}

func init() {
	// The command "dbdeployer defaults template show templateName"
	// depends on the template names being unique across all collections.
	// This initialisation routine will ensure that there are no duplicates.
	var seen = make(map[string]bool)
	for collName, coll := range AllTemplates {
		for name := range coll {
			_, found := seen[name]
			if found {
				// name already exists:
				fmt.Printf("Duplicate template %s found in %s\n", name, collName)
				os.Exit(1)
			}
			seen[name] = true
			if len(coll[name].Contents) == 0 {
				// The template is empty
				fmt.Printf("Template '%s' is empty\n", name)
				os.Exit(1)
			}
		}
	}
}
