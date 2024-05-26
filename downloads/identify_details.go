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
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"strings"
)

// UrlDetails contains the details taken from a tarball Url
// - size is discovered by getting the details from the webserver
// - checksum needs to be discovered by downloading the file
type UrlDetails struct {
	Architecture string
	Checksum     string
	Flavour      string
	Minimal      bool
	OS           string
	ShortVersion string
	Size         int64
	Version      string
}

// Given a URL attempt to identify the details of the expected MySQL like binary
// - size and checksum can only be calculated later if the other details are recognised.
func GetDetailsFromUrl(url string) (UrlDetails, error) {
	// FIXME(sm): not sure if this is the right/best way to process a URL but it seems to work.
	basename := path.Base(url)
	fmt.Printf("WARNING: add-url may take some time as the url file has to be downloaded to determine the checksum\n")
	fmt.Printf("- basename: %v\n", basename)

	flavour, err := identifyFlavour(basename)
	if err != nil {
		return UrlDetails{}, err
	}
	fmt.Printf("- flavour: %v\n", flavour)

	os, err := identifyOS(basename)
	if err != nil {
		return UrlDetails{}, err
	}
	fmt.Printf("- os: %v\n", os)

	version, shortVersion, err := identifyVersion(basename)
	if err != nil {
		return UrlDetails{}, err
	}
	fmt.Printf("- version: %v\n", version)
	fmt.Printf("- shortVersion: %v\n", shortVersion)

	arch, err := identifyArchitecture(basename)
	if err != nil {
		return UrlDetails{}, err
	}
	fmt.Printf("- arch: %v\n", arch)

	minimal := identifyMinimal(basename)
	if err != nil {
		return UrlDetails{}, err
	}
	fmt.Printf("- minimal: %v\n", minimal)

	// get the size by querying the url (requires web access)
	size, err := checkRemoteUrl(url)
	if err != nil {
		return UrlDetails{}, err
	}
	fmt.Printf("- size: %v\n", size)

	// getting the checksum requires downloading the file from the URL
	checksum, err := identifyChecksum(url)
	if err != nil {
		return UrlDetails{}, err
	}
	fmt.Printf("- checksum: %v\n", checksum)

	return UrlDetails{
		OS:           os,
		Architecture: arch,
		Checksum:     checksum,
		Minimal:      minimal,
		ShortVersion: shortVersion,
		Version:      version,
		Flavour:      flavour,
		Size:         size,
	}, nil
}

// identify if this is a minimal image
func identifyMinimal(basename string) bool {
	// mysql and Percona-Server have this in the filename
	return strings.Contains(basename, "minimal")
}

// identifyFlavour identifies the flavour based on the basename
func identifyFlavour(basename string) (string, error) {
	patterns := []struct {
		pattern string
		flavour string
	}{
		{"^mysql-5.7", "mysql"},
		{"^mysql-8", "mysql"},
		{"^mysql-cluster-8", "ndb"},
		{"^mysql-shell-", "shell"},
		{"^Percona-Server-", "percona"},
	}

	for _, p := range patterns {
		r := regexp.MustCompile(p.pattern)
		if r.MatchString(basename) {
			return p.flavour, nil
		}
	}

	return "", fmt.Errorf("unable to identify flavor of %v", basename)
}

// identifyArchitecture returns the architecture based on the basename
func identifyArchitecture(basename string) (string, error) {
	pattern := `(aarch64|amd64|arm64|x86-64|x86_64)`

	r := regexp.MustCompile(pattern)
	arch := r.FindString(basename)
	if arch != "" {
		return arch, nil
	}

	return "", fmt.Errorf("unable to identify architecture of %v", basename)
}

// identifyVersion returns the version and shortVersion based on the basename
func identifyVersion(basename string) (version string, shortVersion string, err error) {
	pattern := `(5\.7\.\d+|8\.[01234]\.\d+|10\.\d+\.\d+)`

	r := regexp.MustCompile(pattern)
	version = r.FindString(basename)
	if version != "" {
		index := strings.LastIndex(version, ".")
		shortVersion = version[:index]
		return version, shortVersion, nil
	}

	return "", "", fmt.Errorf("unable to identify version/short version of %v", basename)
}

// identifyOS returns the OS based on the basename
func identifyOS(basename string) (string, error) {
	patterns := []struct {
		pattern string
		os      string
	}{
		{"[Ll]inux", "linux"},
		{"macos", "Darwin"},
	}

	for _, p := range patterns {
		r := regexp.MustCompile(p.pattern)
		os := r.FindString(basename)
		if os != "" {
			return p.os, nil
		}
	}

	return "", fmt.Errorf("unable to identify OS of %v", basename)
}

// identifyChecksum returns the MD5 checksum of the file, by downloading it on the fly
func identifyChecksum(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	hash := md5.New()
	if _, err = io.Copy(hash, resp.Body); err != nil {
		return "", err
	}

	return fmt.Sprintf("MD5:%x", hash.Sum(nil)), nil
}
