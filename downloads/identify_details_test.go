// DBDeployer - The MySQL Sandbox
// Copyright © 2024 Simon J Mudd <sjmudd@pobox.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package downloads

import (
	"testing"
)

// try to identify the flavour, version, arch, os based on the filename
func TestIdentify(t *testing.T) {

	tests := []struct {
		filename             string
		expectedFlavour      string
		expectedVersion      string
		expectedShortVersion string
		expectedArch         string
		expectedOS           string
	}{
		{"mysql-8.0.34-macos13-arm64.tar.gz", "mysql", "8.0.34", "8.0", "arm64", "Darwin"},
		{"mysql-8.0.37-macos14-arm64.tar.gz", "mysql", "8.0.37", "8.0", "arm64", "Darwin"},
		{"mysql-8.2.0-macos13-arm64.tar.gz", "mysql", "8.2.0", "8.2", "arm64", "Darwin"},
		{"mysql-8.4.0-macos14-arm64.tar.gz", "mysql", "8.4.0", "8.4", "arm64", "Darwin"},
		{"mysql-cluster-8.0.34-macos13-arm64.tar.gz", "ndb", "8.0.34", "8.0", "arm64", "Darwin"},
		{"mysql-shell-8.0.33-linux-glibc2.12-x86-64bit.tar.gz", "shell", "8.0.33", "8.0", "x86-64", "linux"},
		{"mysql-shell-8.0.33-macos13-arm64.tar.gz", "shell", "8.0.33", "8.0", "arm64", "Darwin"},
		{"mysql-shell-8.0.33-macos13-x86-64bit.tar.gz", "shell", "8.0.33", "8.0", "x86-64", "Darwin"},
		{"mysql-shell-8.0.33-linux-glibc2.12-x86-64bit.tar.gz", "shell", "8.0.33", "8.0", "x86-64", "linux"},
		{"mysql-8.0.33-macos13-arm64.tar.gz", "mysql", "8.0.33", "8.0", "arm64", "Darwin"},
		{"mysql-8.0.33-macos13-x86_64.tar.gz", "mysql", "8.0.33", "8.0", "x86_64", "Darwin"},
		{"mysql-cluster-8.0.33-macos13-arm64.tar.gz", "ndb", "8.0.33", "8.0", "arm64", "Darwin"},
		{"mysql-cluster-8.0.33-macos13-x86_64.tar.gz", "ndb", "8.0.33", "8.0", "x86_64", "Darwin"},
		{"mysql-cluster-8.0.33-linux-glibc2.12-x86_64.tar.xz", "ndb", "8.0.33", "8.0", "x86_64", "linux"},
		{"mysql-8.0.33-linux-glibc2.17-x86_64-minimal.tar.xz", "mysql", "8.0.33", "8.0", "x86_64", "linux"},
		{"mysql-8.0.33-linux-glibc2.17-aarch64-minimal.tar.xz", "mysql", "8.0.33", "8.0", "aarch64", "linux"},
		{"mysql-8.0.33-linux-glibc2.28-aarch64.tar.gz", "mysql", "8.0.33", "8.0", "aarch64", "linux"},
		{"mysql-8.0.33-linux-glibc2.28-x86_64.tar.gz", "mysql", "8.0.33", "8.0", "x86_64", "linux"},
	}

	for _, test := range tests {
		flavour, err := identifyFlavour(test.filename)
		if flavour != test.expectedFlavour {
			t.Errorf("identifyFlavour(%v) failed. Expected: %v. Got: %v, %v",
				test.filename,
				test.expectedFlavour,
				flavour,
				err,
			)
		}

		version, shortVersion, err := identifyVersion(test.filename)
		if version != test.expectedVersion || shortVersion != test.expectedShortVersion {
			t.Errorf("identifyVersion(%v) failed. Expected: %v, %v. Got: %v, %v, %v",
				test.filename,
				test.expectedVersion,
				test.expectedShortVersion,
				version,
				shortVersion,
				err,
			)
		}

		arch, err := identifyArchitecture(test.filename)
		if arch != test.expectedArch {
			t.Errorf("identifyArchitecture(%v) failed. Expected: %v. Got: %v, %v",
				test.filename,
				test.expectedArch,
				arch,
				err,
			)
		}

		os, err := identifyOS(test.filename)
		if os != test.expectedOS {
			t.Errorf("identifyOS(%v) failed. Expected: %v. Got: %v, %v",
				test.filename,
				test.expectedOS,
				os,
				err,
			)
		}
	}
}
