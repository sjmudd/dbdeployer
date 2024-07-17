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
)

// Templates for replication

var (
	//go:embed templates/replication/initialize_slaves.gotxt
	initializeSlavesTemplate string

	//go:embed templates/replication/semi_sync_start.gotxt
	semiSyncStartTemplate string

	//go:embed templates/replication/start_all.gotxt
	startAllTemplate string

	//go:embed templates/replication/restart_all.gotxt
	restartAllTemplate string

	//go:embed templates/replication/exec_all.gotxt
	execAllTemplate string

	//go:embed templates/replication/use_all.gotxt
	useAllTemplate string

	//go:embed templates/replication/metadata_all.gotxt
	metadataAllTemplate string

	//go:embed templates/replication/use_all_admin.gotxt
	useAllAdminTemplate string

	//go:embed templates/replication/use_all_slaves.gotxt
	useAllSlavesTemplate string

	//go:embed templates/replication/use_all_masters.gotxt
	useAllMastersTemplate string

	//go:embed templates/replication/exec_all_slaves.gotxt
	execAllSlavesTemplate string

	//go:embed templates/replication/exec_all_masters.gotxt
	execAllMastersTemplate string

	//go:embed templates/replication/wipe_and_restart_all.gotxt
	wipeAndRestartAllTemplate string

	//go:embed templates/replication/stop_all.gotxt
	stopAllTemplate string

	//go:embed templates/replication/send_kill_all.gotxt
	sendKillAllTemplate string

	//go:embed templates/replication/clear_all.gotxt
	clearAllTemplate string

	//go:embed templates/replication/status_all.gotxt
	statusAllTemplate string

	//go:embed templates/replication/test_sb_all.gotxt
	testSbAllTemplate string

	//go:embed templates/replication/check_slaves.gotxt
	checkSlavesTemplate string

	//go:embed templates/replication/master.gotxt
	masterTemplate string

	//go:embed templates/replication/master_admin.gotxt
	masterAdminTemplate string

	//go:embed templates/replication/slave.gotxt
	slaveTemplate string

	//go:embed templates/replication/slave_admin.gotxt
	slaveAdminTemplate string

	//go:embed templates/replication/test_replication.gotxt
	testReplicationTemplate string

	//go:embed templates/replication/multi_source.gotxt
	multiSourceTemplate string

	//go:embed templates/replication/multi_source_use_slaves.gotxt
	multiSourceUseSlavesTemplate string

	//go:embed templates/replication/multi_source_use_masters.gotxt
	multiSourceUseMastersTemplate string

	//go:embed templates/replication/multi_source_exec_slaves.gotxt
	multiSourceExecSlavesTemplate string

	//go:embed templates/replication/multi_source_exec_masters.gotxt
	multiSourceExecMastersTemplate string

	//go:embed templates/replication/check_multi_source.gotxt
	checkMultiSourceTemplate string

	//go:embed templates/replication/multi_source_test.gotxt
	multiSourceTestTemplate string

	//go:embed templates/replication/repl_replicate_from.gotxt
	replicateFromReplTemplate string

	//go:embed templates/replication/repl_sysbench.gotxt
	sysbenchReplTemplate string

	//go:embed templates/replication/repl_sysbench_ready.gotxt
	sysbenchReadyReplTemplate string

	ReplicationTemplates = TemplateCollection{
		TmplInitializeSlaves: TemplateDesc{
			Description: "Initialize slaves after deployment",
			Notes:       "Can also be run after calling './clear_all'",
			Contents:    initializeSlavesTemplate,
		},
		TmplSemiSyncStart: TemplateDesc{
			Description: "Starts semi synch replication ",
			Notes:       "",
			Contents:    semiSyncStartTemplate,
		},
		TmplStartAll: TemplateDesc{
			Description: "Starts nodes in replication order (with optional mysqld arguments)",
			Notes:       "",
			Contents:    startAllTemplate,
		},
		TmplRestartAll: TemplateDesc{
			Description: "stops all nodes and restarts them (with optional mysqld arguments)",
			Notes:       "",
			Contents:    restartAllTemplate,
		},
		TmplUseAll: TemplateDesc{
			Description: "Execute a query for all nodes",
			Notes:       "",
			Contents:    useAllTemplate,
		},
		TmplExecAll: TemplateDesc{
			Description: "Execute a command in all nodes",
			Notes:       "",
			Contents:    execAllTemplate,
		},
		TmplMetadataAll: TemplateDesc{
			Description: "Execute a metadata query for all nodes",
			Notes:       "",
			Contents:    metadataAllTemplate,
		},
		TmplUseAllAdmin: TemplateDesc{
			Description: "Execute a query (as admin user) for all nodes",
			Notes:       "",
			Contents:    useAllAdminTemplate,
		},
		TmplUseAllSlaves: TemplateDesc{
			Description: "Execute a query for all slaves",
			Notes:       "master-slave topology",
			Contents:    useAllSlavesTemplate,
		},
		TmplUseAllMasters: TemplateDesc{
			Description: "Execute a query for all masters",
			Notes:       "master-slave topology",
			Contents:    useAllMastersTemplate,
		},
		TmplExecAllSlaves: TemplateDesc{
			Description: "Execute a command in all slave nodes",
			Notes:       "master-slave topology",
			Contents:    execAllSlavesTemplate,
		},
		TmplExecAllMasters: TemplateDesc{
			Description: "Execute a command in all master nodes",
			Notes:       "master-slave topology",
			Contents:    execAllMastersTemplate,
		},
		TmplStopAll: TemplateDesc{
			Description: "Stops all nodes in reverse replication order",
			Notes:       "",
			Contents:    stopAllTemplate,
		},
		TmplSendKillAll: TemplateDesc{
			Description: "Send kill signal to all nodes",
			Notes:       "",
			Contents:    sendKillAllTemplate,
		},
		TmplClearAll: TemplateDesc{
			Description: "Remove data from all nodes",
			Notes:       "",
			Contents:    clearAllTemplate,
		},
		TmplStatusAll: TemplateDesc{
			Description: "Show status of all nodes",
			Notes:       "",
			Contents:    statusAllTemplate,
		},
		TmplTestSbAll: TemplateDesc{
			Description: "Run sb test on all nodes",
			Notes:       "",
			Contents:    testSbAllTemplate,
		},
		TmplTestReplication: TemplateDesc{
			Description: "Tests replication flow",
			Notes:       "",
			Contents:    testReplicationTemplate,
		},
		TmplCheckSlaves: TemplateDesc{
			Description: "Checks replication status in master and slaves",
			Notes:       "",
			Contents:    checkSlavesTemplate,
		},
		TmplMaster: TemplateDesc{
			Description: "Runs the MySQL client for the master",
			Notes:       "",
			Contents:    masterTemplate,
		},
		TmplMasterAdmin: TemplateDesc{
			Description: "Runs the MySQL client for the master as admin user",
			Notes:       "",
			Contents:    masterAdminTemplate,
		},
		TmplSlave: TemplateDesc{
			Description: "Runs the MySQL client for a slave",
			Notes:       "",
			Contents:    slaveTemplate,
		},
		TmplSlaveAdmin: TemplateDesc{
			Description: "Runs the MySQL client for a slave as admin_user",
			Notes:       "",
			Contents:    slaveAdminTemplate,
		},
		TmplMultiSource: TemplateDesc{
			Description: "Initializes nodes for multi-source replication",
			Notes:       "fan-in and all-masters",
			Contents:    multiSourceTemplate,
		},
		TmplMultiSourceUseSlaves: TemplateDesc{
			Description: "Runs a query for all slave nodes",
			Notes:       "group replication and multi-source topologies",
			Contents:    multiSourceUseSlavesTemplate,
		},
		TmplMultiSourceUseMasters: TemplateDesc{
			Description: "Runs a query for all master nodes",
			Notes:       "group replication and multi-source topologies",
			Contents:    multiSourceUseMastersTemplate,
		},
		TmplMultiSourceExecSlaves: TemplateDesc{
			Description: "Runs a command in each slave node",
			Notes:       "group replication and multi-source topologies",
			Contents:    multiSourceExecSlavesTemplate,
		},
		TmplMultiSourceExecMasters: TemplateDesc{
			Description: "Runs a command in each slave node",
			Notes:       "group replication and multi-source topologies",
			Contents:    multiSourceExecMastersTemplate,
		},
		TmplWipeAndRestartAll: TemplateDesc{
			Description: "clears the databases and restarts them all",
			Notes:       "group replication and multi-source topologies",
			Contents:    wipeAndRestartAllTemplate,
		},
		TmplMultiSourceTest: TemplateDesc{
			Description: "Test replication flow for multi-source replication",
			Notes:       "fan-in and all-masters",
			Contents:    multiSourceTestTemplate,
		},
		TmplCheckMultiSource: TemplateDesc{
			Description: "checks replication status for multi-source replication",
			Notes:       "fan-in and all-masters",
			Contents:    checkMultiSourceTemplate,
		},
		TmplReplReplicateFrom: TemplateDesc{
			Description: "use replicate_from script from the master",
			Notes:       "",
			Contents:    replicateFromReplTemplate,
		},
		TmplReplSysbench: TemplateDesc{
			Description: "use sysbench script from the master",
			Notes:       "",
			Contents:    sysbenchReplTemplate,
		},
		TmplReplSysbenchReady: TemplateDesc{
			Description: "use sysbench_ready script from the master",
			Notes:       "",
			Contents:    sysbenchReadyReplTemplate,
		},
	}
)
