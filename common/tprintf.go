// DBDeployer - The MySQL Sandbox
// Copyright © 2006-2019 Giuseppe Maxia
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

package common

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/datacharmer/dbdeployer/globals"
)

// StringMap defines the map of variable types, for brevity
type StringMap map[string]interface{}

// Add returns a StringMap of the original map with the sm2 map values added
// - if there are duplicate/overlapping entries a warning will be printed.
func (sm StringMap) Add(sm2 StringMap) StringMap {
	m := make(StringMap)
	for k, v := range sm {
		m[k] = v
	}
	for k, v := range sm2 {
		if v, ok := m[k]; ok {
			fmt.Printf("WARNING: StringMap.Add(): key: %q already defined in StringMap with valuee: <%v>. Will overwrite with: <%v>\n",
				k,
				sm[k],
				v)
		}
		m[k] = v
	}
	return m
}

// Given a multi-line string, this function removes leading
// spaces from every line.
// It also removes the first line, if it is empty
func TrimmedLines(s string) string {
	// matches the start of the text followed by an EOL
	re := regexp.MustCompile(`(?m)\A\s*$`)
	s = re.ReplaceAllString(s, globals.EmptyString)

	re = regexp.MustCompile(`(?m)^\t\t`)
	s = re.ReplaceAllString(s, globals.EmptyString)
	return s
}

func GetBashPath(wantedValue string) (string, error) {

	var shellPath string
	defaultValue := globals.ShellPathValue

	if wantedValue != "" {
		shellPath = wantedValue
	} else {
		shellPath = os.Getenv("SHELL_PATH")
		if shellPath == "" {
			shellPath = defaultValue
		}
	}
	if !ExecExists(shellPath) {
		return "", fmt.Errorf("executable %s does not exist", shellPath)
	}
	out, err := RunCmdCtrlWithArgs(shellPath, []string{"-c", "echo $BASH_VERSION"}, true)
	if err != nil {
		return "", fmt.Errorf("error checking BASH_VERSION for %s: %s", shellPath, err)
	}
	if IsEmptyOrBlank(out) {
		return "", fmt.Errorf("executable '%s' does not appear to be a Bash interpreter", shellPath)
	}
	return shellPath, nil
}

// Returns true if a StringMap (or an inner StringMap) contains a given key
func hasKey(sm StringMap, wantedKey string) bool {
	for key, value := range sm {
		if key == wantedKey {
			return true
		}
		valType := reflect.TypeOf(value)
		if valType == reflect.TypeOf(StringMap{}) {
			innerSm := value.(StringMap)
			return hasKey(innerSm, wantedKey)
		}
		if valType == reflect.TypeOf([]StringMap{}) {
			innerSm := value.([]StringMap)
			for _, ism := range innerSm {
				if hasKey(ism, wantedKey) {
					return true
				}
			}
		}
	}
	return false
}

// Returns a unique, sorted list of all variables defined in a template string
// FIXME: does not handle NESTED variables
func GetVarsFromTemplate(tmpl string) []string {
	var (
		varList = make([]string, 0)
		varMap  = make(map[string]struct{})
	)

	reTemplateVar := regexp.MustCompile(`\{\{\.([^{]+)\}\}`)
	captureList := reTemplateVar.FindAllStringSubmatch(tmpl, -1)

	for _, capture := range captureList {
		varMap[capture[1]] = struct{}{}
	}
	// push map keys to list (unsorted)
	for k := range varMap {
		varList = append(varList, k)
	}
	// sort the keys
	sort.Strings(varList)

	return varList
}

func listToCommaSeparatedString(list []string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return strings.Join(list, ", ")
	}
}

// checkAllParametersProvided checks if all parameters defined in the template are provided in data
// and returns an error indicating which entries are missing if appropriate.
func checkAllParametersProvided(template_name, tmpl string, data StringMap) error {
	var extraData string
	fields := make(map[string]struct{})

	for _, capture := range GetVarsFromTemplate(tmpl) {
		if !hasKey(data, capture) {
			fields[capture] = struct{}{}
		}
	}

	if len(fields) > 0 {
		quoted := make([]string, 0)

		for field := range fields {
			quoted = append(quoted, fmt.Sprintf("%q", field))
		}

		switch len(quoted) {
		case 1:
			extraData = fmt.Sprintf("required data for field %v was not populated", listToCommaSeparatedString(quoted))
		default:
			extraData = fmt.Sprintf("required data for %d fields were not populated: %s", len(fields), listToCommaSeparatedString(quoted))
		}

		return fmt.Errorf("template %q: %v",
			template_name,
			extraData,
		)
	}

	return nil
}

// SafeTemplateFill passed template string is formatted using its operands and returns the resulting string.
// It checks that the data was safely initialized
func SafeTemplateFill(template_name, tmpl string, data StringMap) (string, error) {

	// Adds shell path, timestamp, version info, and empty engine clause if one was not provided
	timestamp := time.Now()
	shellPath, shellPathExists := data["ShellPath"]
	_, engineClauseExists := data["EngineClause"]
	_, timeStampExists := data["DateTime"]
	_, versionExists := data["AppVersion"]
	if !timeStampExists {
		data["DateTime"] = timestamp.Format(time.UnixDate)
	}
	if !versionExists {
		data["AppVersion"] = VersionDef
	}
	if !engineClauseExists {
		data["EngineClause"] = ""
	}
	if shellPathExists {
		if shellPath == "" {
			shellPathExists = false
		}
	}
	if !shellPathExists {
		data["ShellPath"] = globals.ShellPathValue
	}

	// Checks that all data was initialized
	// This check is especially useful when introducing new templates

	if err := checkAllParametersProvided(template_name, tmpl, data); err != nil {
		return globals.EmptyString, err
	}

	// Creates a template
	processTemplate := template.Must(template.New("tmp").Parse(tmpl))
	buf := &bytes.Buffer{}

	// If an error occurs, returns an empty string
	if err := processTemplate.Execute(buf, data); err != nil {
		return globals.EmptyString, err
	}

	// Returns the populated template
	return buf.String(), nil
}
