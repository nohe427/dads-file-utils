// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package duplicatefinder is a struct that most methods for finding
// duplicate files in a users file system
package duplicatefinder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nohe427/dup-file-finder-lib/database"
	"github.com/nohe427/dup-file-finder-lib/hasher"
)

// DuplicateFinder is the main struct for holding data relevant to finding
// duplicate files on the users system
type DuplicateFinder struct {
	appDir       string
	folderToScan string
	appDb        database.Database
	SearchStyle  DuplicateSearchStyle
}

// DuplicateSearchStyle is an enumeration of all the different search options available
// for searching for duplicate files.
type DuplicateSearchStyle int

const (
	// BySize is to find based on file size
	BySize DuplicateSearchStyle = iota
	// ByContents is to find based on file contents (sha256)
	ByContents DuplicateSearchStyle = iota
)

// New generates a new duplicate finder.  The duplicate finder
// manages an internal state and can be used to help find duplicates
// existing in the same drive or directory.  Default search style is BySize.
func New(preferredAppDir string, db database.Database) (*DuplicateFinder, error) {
	dupFinder := &DuplicateFinder{}
	if err := createAppDir(preferredAppDir); err != nil {
		return nil, err
	}
	dupFinder.appDb = db
	dupFinder.appDir = preferredAppDir
	dupFinder.SearchStyle = BySize
	return dupFinder, nil
}

// createAppDir takes the desired path for a App Directory to be located
// and generates that app directory on the file system.  If there is an
// error, it will be returned.
func createAppDir(appDataPath string) error {
	stats, err := os.Stat(appDataPath)
	if !os.IsNotExist(err) {
		if !stats.IsDir() {
			return fmt.Errorf("file exists at %q, expected a directory", appDataPath)
		}
		return nil
	}
	return os.MkdirAll(appDataPath, 0666)
}

func isValidSearchDir(dirToSearch string) (bool, error) {
	stats, err := os.Stat(dirToSearch)
	if !os.IsNotExist(err) {
		if !stats.IsDir() {
			return false, fmt.Errorf("file exists at %q, expected a directory", dirToSearch)
		}
	}
	if os.IsNotExist(err) {
		return false, fmt.Errorf("folder to search %q does not exist", dirToSearch)
	}
	return true, nil
}

// FindDuplicateFiles is the method that starts the search for duplicate files based on the directory
// that is provided to it.
func (df *DuplicateFinder) FindDuplicateFiles(dirToSearch string) error {
	exists, err := isValidSearchDir(dirToSearch)
	if !exists || err != nil {
		return fmt.Errorf("Not a valid search directory : %v", err)
	}
	filepath.Walk(dirToSearch, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if df.SearchStyle == BySize {
				df.appDb.Add(string(info.Size()), path)
			}
			if df.SearchStyle == ByContents {
				hashedString, err := hasher.HashFile(path)
				if err != nil {
					return err
				}
				df.appDb.Add(hashedString, path)
			}
		}
		return nil
	})
	return nil
}
