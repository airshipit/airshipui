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

package webservice

import (
	"github.com/pkg/errors"

	"net/http"

	"opendev.org/airship/airshipui/util/utilfile"
	"opendev.org/airship/airshipui/util/utilhttp"
)

const (
	clientPath = "client/dist/airshipui-ui"
)

func serveFile(w http.ResponseWriter, r *http.Request) {
	filePath, filePathErr := utilfile.FilePath(clientPath, r.URL.Path)
	if filePathErr != nil {
		utilhttp.HandleErr(w, errors.WithStack(filePathErr), http.StatusInternalServerError)
		return
	}
	fileExists, fileExistsErr := utilfile.Exists(filePath)
	if fileExistsErr != nil {
		utilhttp.HandleErr(w, errors.WithStack(fileExistsErr), http.StatusInternalServerError)
		return
	}
	if fileExists {
		http.ServeFile(w, r, filePath)
	}
}
