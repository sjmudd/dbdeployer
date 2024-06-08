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
	scripts    []ScriptDef
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

// AppendScript appens a script to scripts with the given parameters
func (sb *ScriptBatch) AppendScript(scriptName string, templateName string, makeExecutable bool) {
	sb.scripts = append(sb.scripts, ScriptDef{
		scriptName:     scriptName,
		templateName:   templateName,
		makeExecutable: makeExecutable,
	})
}
