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

package ctl

import (
	"fmt"
	"strings"

	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
	"opendev.org/airship/airshipui/pkg/statistics"
)

// This is close to but not exactly a transaction structure from statistics, it's redone here because reasons
type record struct {
	SubComponent configs.WsSubComponentType
	User         *string
	ActionType   *string
	Target       *string
	Success      bool
	Started      int64
	Elapsed      int64
	Stopped      int64
}

// HandleHistoryRequest will flop between requests so we don't have to have them all mapped as function calls
// This will wait for the sub component to complete before responding.  The assumption is this is an async request
func HandleHistoryRequest(user *string, request configs.WsMessage) configs.WsMessage {
	response := configs.WsMessage{
		Type:         configs.CTL,
		Component:    configs.History,
		SubComponent: request.SubComponent,
	}

	var err error

	subComponent := request.SubComponent
	switch subComponent {
	case configs.GetDefaults:
		response.Data, err = getData(nil, nil)
	default:
		err = fmt.Errorf("Subcomponent %s not found", request.SubComponent)
	}

	if err != nil {
		e := err.Error()
		response.Error = &e
	}

	return response
}

// getData will return all rows within a specific date range
func getData(notBefore *int64, notAfter *int64) (map[string][]record, error) {
	var wherePstmt strings.Builder
	where := false
	// because we may want data within a range add range slice where statements
	if notBefore != nil {
		where = true
		wherePstmt.WriteString(fmt.Sprintf(" where started >= %d", notBefore))
	}
	if notAfter != nil {
		if where {
			wherePstmt.WriteString(fmt.Sprintf(" and stopped <= %d", notAfter))
		} else {
			wherePstmt.WriteString(fmt.Sprintf(" where stopped <= %d", notAfter))
		}
	}

	data := map[string][]record{}
	for _, table := range statistics.Tables {
		// create a basic prepared statement to get data
		// Why a prepared statement?  Little Bobby Tables is why:
		// https://xkcd.com/327/
		pstmt := fmt.Sprintf("select * from %s", table)
		if where {
			pstmt += wherePstmt.String()
		}

		// Dark Helmet: Why are we always preparing?  Just go!
		stmt, err := statistics.DB.Prepare(pstmt)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		// Get the rows back from the query
		rows, err := stmt.Query()
		if err != nil {
			log.Error(err)
			return nil, err
		}
		defer rows.Close()

		records := []record{}
		for rows.Next() {
			var r record
			err = rows.Scan(
				&r.SubComponent,
				&r.User,
				&r.ActionType,
				&r.Target,
				&r.Success,
				&r.Started,
				&r.Elapsed,
				&r.Stopped,
			)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			records = append(records, r)
		}

		err = rows.Err()
		if err != nil {
			log.Error(err)
			return nil, err
		}

		if len(records) > 0 {
			data[table] = records
		}
	}

	return data, nil
}
