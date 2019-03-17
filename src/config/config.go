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

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ServerConfig maps config.json
type serverConfig struct {
	UseTLS       bool   `json:"use_tls"`
	Port         int    `json:"port"`
	GateEndpoint string `json:"gate_endpoint"`
	CrtPath      string `json:"crt_path"`
	KeyPath      string `json:"key_path"`
}

// ServerConfig maps config.json
type receiverConfig struct {
	Type    string `json:"type"`
	KeyPath string `json:"key_path"`
	KeySize uint   `json:"key_size"`
}

// ServerConfig maps config.json
type databaseConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// Config defines your config in json
type Config struct {
	Server   serverConfig   `json:"server"`
	Receiver receiverConfig `json:"receiver"`
	Database databaseConfig `json:"database"`
}

// LoadConfig parses json file into Config struct
func LoadConfig() (*Config, error) {
	jsonFile, err := os.Open("config/config.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	cfg := &Config{}
	err = json.Unmarshal(byteValue, cfg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return cfg, nil
}
