/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package db

import (
	"database/sql"
	"fmt"

	"config"
)

// Global object for SQL connection. SQL connection is meant to be open for a long time.
var Conn *sql.DB

// LoadDB initializes SQL connection based on settings supplied in config
func LoadDB(config *config.Config) error {
	s := fmt.Sprintf("%s:%s@/%s?parseTime=true", config.Database.User, config.Database.Password, config.Database.Database)
	db, err := sql.Open("mysql", s)
	if err != nil {
		return err
	}
	// You have to Ping the db to make sure Open worked
	err = db.Ping()
	if err != nil {
		return err
	}
	Conn = db
	Conn.SetMaxIdleConns(5)
	return nil
}
