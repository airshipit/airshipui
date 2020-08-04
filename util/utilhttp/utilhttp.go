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

package utilhttp

import (
	"fmt"
	"net/http"
)

// Header stuct for a http header.
type Header struct {
	Key   string
	Value string
}

// Oauth stuct for oauth parameters.
type Oauth struct {
	ClientID     string
	ClientSecret string
}

// API struct for apis.
type API struct {
	URL     string
	Headers []Header
	Oauth   Oauth
}

// HandleErr handles a http error.
func HandleErr(w http.ResponseWriter, err error, code int) {
	fmt.Printf("[Error] %+v\n\n", err)
	http.Error(w, err.Error(), code)
}
