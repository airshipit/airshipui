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
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	serverAddr string = "localhost:8080"
)

func init() {
	go WebServer()
	// wait for the webserver to come up
	time.Sleep(250 * time.Millisecond)
}

func TestRootURI(t *testing.T) {
	resp, err := http.Get("http://" + serverAddr)
	require.NoError(t, err)
	defer resp.Body.Close()
	// this will be not found because of where the webservice starts
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
