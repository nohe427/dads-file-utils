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

package duplicatefinder

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-test/deep"
	"github.com/google/go-cmp/cmp"
	"github.com/nohe427/dup-file-finder-lib/database"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	tmpDir, err := ioutil.TempDir("", "_nohe427_test_dir_")
	if err != nil {
		t.Fatalf("could not create test tmpDir : %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})
	return tmpDir
}

func setupTestDb(t *testing.T) *database.InMemoryDB {
	t.Helper()
	db, err := database.GetInMemoryDbInstance()
	if err != nil {
		t.Fatalf("In memory db failed to be created.  This should never happen. : %v", err)
	}
	return db
}

func setupTestData(t *testing.T) string {
	t.Helper()
	dataDir, err := ioutil.TempDir("", "_nohe427_test_data_dir_")
	os.Mkdir(filepath.Join(dataDir, "dir1"), 0777)
	if err != nil {
		t.Fatalf("could not create test data dir : %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(dataDir)
	})
	writeTestFile(t, filepath.Join(dataDir, "randoFile"), "0000")
	writeTestFile(t, filepath.Join(dataDir, "randoFile2"), "0000")
	writeTestFile(t, filepath.Join(dataDir, "dir1", "randoFile"), "0000")
	writeTestFile(t, filepath.Join(dataDir, "dir1", "randoFileOfAnotherName"), "0000")
	writeTestFile(t, filepath.Join(dataDir, "randoFile1"), "0001")
	writeTestFile(t, filepath.Join(dataDir, "randoFile3"), "0001")
	writeTestFile(t, filepath.Join(dataDir, "dir1", "randoFile1"), "0001")
	writeTestFile(t, filepath.Join(dataDir, "dir1", "randoFileOfAnotherName3"), "0001")
	writeTestFile(t, filepath.Join(dataDir, "randoFile7"), "0002")
	writeTestFile(t, filepath.Join(dataDir, "randoFile8"), "0003")
	writeTestFile(t, filepath.Join(dataDir, "dir1", "randoFile9"), "0004")
	writeTestFile(t, filepath.Join(dataDir, "dir1", "randoFileOfAnotherName0"), "0005")
	return dataDir
}

func writeTestFile(t *testing.T, fileLocation string, contents string) {
	t.Helper()
	if err := ioutil.WriteFile(fileLocation, []byte(contents), 0777); err != nil {
		t.Fatalf("could not write temp file : %v", err)
	}
}

func TestCreateAppDirWithExistingFile(t *testing.T) {
	tmpDir := setupTestDir(t)
	appDir := filepath.Join(tmpDir, "app_dir")
	if err := ioutil.WriteFile(appDir, nil, 0666); err != nil {
		t.Fatalf("could not create test file : %v", err)
	}
	err := createAppDir(appDir)
	wantErr := fmt.Errorf("file exists at %q, expected a directory", appDir)
	if diff := cmp.Diff(err.Error(), wantErr.Error()); diff != "" {
		t.Errorf("difference of errors : %v", diff)
	}
}

func TestCreateAppDirWithExistingDir(t *testing.T) {
	tmpDir := setupTestDir(t)
	appDir := filepath.Join(tmpDir, "app_dir")
	if err := os.MkdirAll(appDir, 0666); err != nil {
		t.Fatalf("creating the test dir failed")
	}
	err := createAppDir(appDir)
	if err != nil {
		t.Errorf("did not expect an error : %v", err)
	}
}

func TestCreateAppDir(t *testing.T) {
	tmpDir := setupTestDir(t)
	appDir := filepath.Join(tmpDir, "app_dir")
	if err := createAppDir(appDir); err != nil {
		t.Errorf("something went wrong creating the app dir : %v", err)
	}
	info, err := os.Stat(appDir)
	if err != nil {
		t.Errorf("something went wrong creating the app dir : %v", err)
	}
	if !info.IsDir() {
		t.Error("created app dir is not showing as a dir")
	}
}

func TestNewDupFinder(t *testing.T) {
	testDir := setupTestDir(t)
	db := setupTestDb(t)
	dupFinder, err := New(testDir, db)
	if err != nil {
		t.Errorf("canont create a new dupfinder : %v", err)
	}
	defer os.RemoveAll(dupFinder.appDir)
	want := &DuplicateFinder{}
	want.appDir = testDir
	if diff := deep.Equal(dupFinder, want); diff != nil {
		t.Errorf("The New() dupfinder is now what we wanted : %v", diff)
	}
}

func TestNewDupFinderOnInvalidDir(t *testing.T) {
	tmpDir := setupTestDir(t)
	db := setupTestDb(t)
	appDir := filepath.Join(tmpDir, "app_dir")
	if err := ioutil.WriteFile(appDir, nil, 0666); err != nil {
		t.Fatalf("could not create test file : %v", err)
	}
	_, err := New(appDir, db)
	wantErr := fmt.Errorf("file exists at %q, expected a directory", appDir)
	if diff := cmp.Diff(err.Error(), wantErr.Error()); diff != "" {
		t.Errorf("difference of errors : %v", diff)
	}
}

func TestNewDupFinderOnExistingDir(t *testing.T) {
	tmpDir := setupTestDir(t)
	db := setupTestDb(t)
	appDir := filepath.Join(tmpDir, "app_dir")
	if err := os.MkdirAll(appDir, 0666); err != nil {
		t.Fatalf("creating the test dir failed")
	}
	_, err := New(appDir, db)
	if err != nil {
		t.Errorf("did not expect an error : %v", err)
	}
}

func TestFindNonExistentFileDir(t *testing.T) {
	testDir := setupTestDir(t)
	db := setupTestDb(t)
	dupFinder, err := New(testDir, db)
	if err != nil {
		t.Errorf("canont create a new dupfinder : %v", err)
	}
	defer os.RemoveAll(dupFinder.appDir)
	testFolder := "/idontexist"
	want := fmt.Errorf("Not a valid search directory : folder to search %q does not exist", testFolder)
	if err := dupFinder.FindDuplicateFiles(testFolder); err.Error() != want.Error() {
		t.Errorf("wanted : %v ||| got : %v", want, err)
	}
}

func TestFindFileAsDir(t *testing.T) {
	testDir := setupTestDir(t)
	db := setupTestDb(t)
	dupFinder, err := New(testDir, db)
	if err != nil {
		t.Errorf("canont create a new dupfinder : %v", err)
	}
	defer os.RemoveAll(dupFinder.appDir)
	tmpFilePath := filepath.Join(os.TempDir(), "file")
	defer os.RemoveAll(tmpFilePath)
	tmpFile, err := os.Create(tmpFilePath)
	if err != nil {
		t.Fatalf("tmpFile could not be created : %v", err)
	}
	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("tmpFile could not be closed : %v", err)
	}
	want := fmt.Errorf("Not a valid search directory : file exists at %q, expected a directory", tmpFilePath)
	if err := dupFinder.FindDuplicateFiles(tmpFilePath); err.Error() != want.Error() {
		t.Errorf("wanted : %v ||| got : %v", want, err)
	}
}

func TestFindDuplicateFilesByContents(t *testing.T) {
	testDir := setupTestDir(t)
	db := setupTestDb(t)
	dupFinder, err := New(testDir, db)
	if err != nil {
		t.Errorf("canont create a new dupfinder : %v", err)
	}
	testDataDir := setupTestData(t)
	err = dupFinder.FindDuplicateFiles(testDataDir)
	if err != nil {
		t.Errorf("should not have failed but did : %v", err)
	}
	// TODO : Write a DupFinder Export function and check the
	// duplicates match what is expected.
}
