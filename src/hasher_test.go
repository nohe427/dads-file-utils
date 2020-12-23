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

package hasher

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func createTestDir(t *testing.T) string {
	t.Helper()
	tmpDir, err := ioutil.TempDir("", "_nohe427_test_dir_")
	if err != nil {
		t.Fatalf("could not create test temp dir : %v", err)
	}
	t.Cleanup(func() {
		defer os.RemoveAll(tmpDir)
	})
	return tmpDir
}

func createTestFile(t *testing.T) string {
	t.Helper()
	tmpDir := createTestDir(t)
	testFilePath := filepath.Join(tmpDir, "testfile")
	if err := ioutil.WriteFile(testFilePath, []byte("test"), 0666); err != nil {
		t.Fatalf("could not write a test temp file : %v", err)
	}
	return testFilePath
}

func TestHashFile(t *testing.T) {
	testFilePath := createTestFile(t)
	hash, err := HashFile(testFilePath)
	if err != nil {
		t.Errorf("Generated a hashing error : %v", err)
	}
	want := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"
	if !cmp.Equal(hash, want) {
		t.Errorf("HashFile('test') got : %q want : %q", hash, want)
	}
}
