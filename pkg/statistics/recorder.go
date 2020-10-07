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

package statistics

import (
	"database/sql"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // this is required for the sqlite driver
	"opendev.org/airship/airshipui/pkg/configs"
	"opendev.org/airship/airshipui/pkg/log"
)

// Transaction will record the details of the CTL transaction and record them to the DB
type Transaction struct {
	Table        configs.WsComponentType
	SubComponent configs.WsSubComponentType
	User         *string
	Target       *string
	Started      int64
}

var (
	writeMutex sync.Mutex
	db         *sql.DB
	tables     = []string{"baremetal", "cluster", "config", "document", "image", "phase", "secret"}
)

const (
	// the table structure used for the records
	tableCreate = `CREATE TABLE IF NOT EXISTS table (
		subcomponent varchar(64) null,
		user varchar(64),
		target varchar(64) null,
		success tinyint(1) default 0,
		started timestamp,
		elapsed bigint,
		stopped timestamp,
		primary key (subcomponent, user, started, stopped))`
	// the prepared statement used for inserts
	// TODO (aschiefe): determine if we need to batch inserts
	insert = "INSERT INTO table(subcomponent, user, target, success, started, elapsed, stopped) values(?,?,?,?,?,?,?)"
)

// Init will create the database if it doesn't exist or open the existing database
func Init() {
	intitTables := false
	// TODO (aschiefe): pull the db location out to the confing
	if _, err := os.Stat("./sqlite/statistics.db"); os.IsNotExist(err) {
		intitTables = true
	}
	// need to define error so that the program will set the global db variable
	var err error
	// TODO (aschiefe): encrypt & password protect the database
	// TODO (aschiefe): pull the db location out to the confing
	db, err = sql.Open("sqlite3", "./sqlite/statistics.db")
	if err != nil {
		log.Fatal(err)
	}
	if intitTables {
		err = createTables()
		if err != nil {
			log.Fatal(err)
		}
	}
}

// createTables is only used when there is no database to write the correct structure for the records
func createTables() error {
	for _, table := range tables {
		stmt, err := db.Prepare(strings.ReplaceAll(tableCreate, "table", table))

		if err != nil {
			return err
		}

		_, err = stmt.Exec()
		if err != nil {
			return err
		}
		log.Tracef("%s table created.", table)
	}
	return nil
}

// NewTransaction establishes the transaction which will record
func NewTransaction(request configs.WsMessage, user *string) *Transaction {
	return &Transaction{
		Table:        request.Component,
		SubComponent: request.SubComponent,
		Target:       request.Target,
		Started:      time.Now().UnixNano() / 1000000,
		User:         user,
	}
}

// Complete will put an entry into the statistics database for the transaction
func (transaction *Transaction) Complete(errorMessagePresent bool) {
	if transaction.User != nil && transaction.isRecordable() {
		stmt, err := db.Prepare(strings.ReplaceAll(insert, "table", string(transaction.Table)))
		if err != nil {
			log.Error(err)
			return
		}

		started := transaction.Started
		stopped := time.Now().UnixNano() / 1000000

		success := 0
		if errorMessagePresent {
			success = 1
		}

		writeMutex.Lock()
		defer writeMutex.Unlock()
		result, err := stmt.Exec(transaction.SubComponent,
			transaction.User,
			transaction.Target,
			success,
			started,
			(stopped - started),
			stopped)

		if err != nil {
			log.Error(err)
			return
		}

		rows, err := result.RowsAffected()
		if err != nil {
			log.Error(err)
			return
		}

		log.Tracef("%d rows inserted into %s.", rows, transaction.Table)
	}
}

// isRecordable will shuffle through the transaction and determine if we should write it to the database
func (transaction *Transaction) isRecordable() bool {
	recordable := true
	if transaction.Table == configs.Auth {
		recordable = false
	}
	switch transaction.SubComponent {
	case configs.GetTarget:
		recordable = false
	case configs.GetPhaseTree:
		recordable = false
	}
	return recordable
}
