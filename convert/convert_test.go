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

// Routines for handling the move from old replication terminology
// - master / slave
// to new replication terminology
// - source / replica

package convert

import (
	"testing"
)

func TestVersionLess(t *testing.T) {
	tests := []struct {
		one      string
		two      string
		expected bool
	}{
		{"5.6.0", "5.6.1", true},
		{"5.6.0", "5.7.0", true},
		{"5.7.36", "8.0.0", true},
		{"8.0.0", "8.0.1", true},
		{"5.10.0", "5.2.0", false},
		{"5.6.0", "5.10.0", true},
		{"8.0.10", "8.0.9", false},
	}

	for _, test := range tests {
		got := Version(test.one).Less(Version(test.two))
		if got != test.expected {
			t.Errorf("Version(%q).Less(%q) returned %v, expected: %v", test.one, test.two, got, test.expected)
		}
	}
}

func TestVersionEqual(t *testing.T) {
	tests := []struct {
		one      string
		two      string
		expected bool
	}{
		{"5", "5", true},
		{"5.5", "5.5", true},
		{"8.0.37", "8.0.37", true},
		{"8.4.0", "8.4.0", true},
		{"5", "5.6", false},
		{"5", "5.6.0", false},
		{"5.6", "5.6.0", false},
		{"8.0", "8.0.37", false},
	}

	for _, test := range tests {
		got := Version(test.one).Equal(Version(test.two))
		if got != test.expected {
			t.Errorf("Version(%q).Equal(%q) returned %v, expected: %v", test.one, test.two, got, test.expected)
		}
	}
}
