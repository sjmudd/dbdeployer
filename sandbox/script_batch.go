// DBDeployer - The MySQL Sandbox
// Copyright © 2024 The dbdeployer authors
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

package sandbox

import (
	"fmt"

	"github.com/datacharmer/dbdeployer/common"
	"github.com/datacharmer/dbdeployer/defaults"
)

// ScriptBatch holds scripts to be used in a sandbox directory
type ScriptBatch struct {
	tc         TemplateCollection
	logger     *defaults.Logger
	sandboxDir string
	data       common.StringMap
	scripts    []Script
}

// String prints out the representation of the ScriptBatch
func (sb ScriptBatch) String() string {
	return fmt.Sprintf("ScriptBatch<tc: %+v, logger: %+v, sandbox: %q, data: %+v, scripts: %+v>",
		sb.tc,
		sb.logger,
		sb.sandboxDir,
		sb.data,
		sb.scripts)
}

// FIXME: use the runtime to figure this out automatically
func (sb ScriptBatch) WriteScripts(hint string) error {
	fmt.Printf("ScriptBatch.WriteScripts(%q): sandboxDir=%q\n", hint, sb.sandboxDir)
	for _, scriptDef := range sb.scripts {
		//		if scriptDef.scriptName == "initialize_nodes" {
		//			fmt.Printf("ERROR: WriteScripts(%v): writeScript(%q,?,?,?,?,data=%+v,?) IS GOING TO FAIL!\n",
		//				hint,
		//				scriptDef.scriptName,
		//				sb.data)
		//		}
		if err := writeScript(
			sb.logger,
			sb.tc,
			scriptDef.scriptName,
			scriptDef.templateName,
			sb.sandboxDir,
			sb.data,
			scriptDef.makeExecutable,
		); err != nil {
			fmt.Printf("ERROR: WriteScripts(%v): writeScript(?,?,%q,%q,%q,len=%d,%v) failed: %v\n",
				hint,
				scriptDef.scriptName,
				scriptDef.templateName,
				sb.sandboxDir,
				len(sb.data),
				scriptDef.makeExecutable,
				err)
			return err
		}
	}
	fmt.Printf("ScriptBatch.WriteScripts() completes\n")
	return nil
}

func writeScript(
	logger *defaults.Logger,
	tempVar TemplateCollection,
	scriptName,
	templateName,
	directory string,
	data common.StringMap,
	makeExecutable bool) error {

	//	fmt.Printf("- writeScript(?,?, scriptName: %q, templateName: %q, directory: %q, len(data): %v)\n",
	//		scriptName,
	//		templateName,
	//		directory,
	//		len(data),
	//	)

	if directory == "" {
		return fmt.Errorf("writeScript (%s): missing directory", scriptName)
	}
	if _, ok := tempVar[templateName]; !ok {
		return fmt.Errorf("writeScript (%s): template %s not found", scriptName, templateName)
	}
	template := tempVar[templateName].Contents
	template = common.TrimmedLines(template)
	data["TemplateName"] = templateName
	var err error
	text, err := common.SafeTemplateFill(templateName, template, data)
	if err != nil {
		return err
	}
	executableStatus := ""
	if makeExecutable {
		err = writeExec(scriptName, text, directory)
		executableStatus = " executable"
	} else {
		_, err = writeRegularFile(scriptName, text, directory)
	}
	if err != nil {
		return err
	}
	if logger != nil {
		logger.Printf("Creating %s script '%s/%s' using template '%s'\n", executableStatus, common.ReplaceLiteralHome(directory), scriptName, templateName)
	}
	return nil
}
