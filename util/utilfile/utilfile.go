/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package utilfile

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Exists returns if a file or directory exists.
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// FilePath returns the absolute path for a file name.
func FilePath(dirPath string, fileName string) (string, error) {
	dir, dirPathErr := DirPath(dirPath)
	if dirPathErr != nil {
		return "", errors.WithStack(dirPathErr)
	}
	return filepath.Join(dir, fileName), nil
}

// DirPath returns the absolute path for a directory name.
func DirPath(dirPath string) (string, error) {
	pwd, getwdErr := os.Getwd()
	if getwdErr != nil {
		return "", errors.WithStack(getwdErr)
	}

	return filepath.Join(pwd, dirPath), nil
}
