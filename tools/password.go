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

package main

import (
	"crypto/sha512"
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go run password.go <password>",
	Short: "Create an sha512 password hash",
	Long:  "This creates an sha512 password hash used for user authentication in the etc/airshipui.json conf file",
	Run:   launch,
}

// take the password argument and turn it into a hash
func launch(cmd *cobra.Command, args []string) {
	if len(args) == 1 {
		// create and disply the sha512 hash for the password
		hash := sha512.New()
		hash.Write([]byte(args[0]))
		fmt.Printf("%x\n", hash.Sum(nil))
	} else {
		fmt.Println("There should be 1 password argument")
	}
}

func main() {
	rootCmd.Execute()
}
