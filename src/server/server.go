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

package server

import (
	"config"
	"cryptography"
	"fmt"
	"html/template"
	"log"
	"message"
	"net/http"
)

// Run launches the server
func Run(config *config.Config) error {

	// setup static URL
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// handle decoy on all URLs except the gate
	http.HandleFunc("/", decoyHandler)

	// handle gate
	http.HandleFunc(config.Server.GateEndpoint, gateHandler)
	fmt.Println("[+] Listening on port:", config.Server.Port)

	// Launch web server, either using TLS or not
	if config.Server.UseTLS {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", config.Server.Port),
			config.Server.CrtPath, config.Server.KeyPath, nil))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), nil))
	}

	return nil
}

// decoyHandler will render Decoy page
func decoyHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
}

// gateHandler handles all calls to the gate
func gateHandler(w http.ResponseWriter, r *http.Request) {
	// Extract message from HTTP request
	msg, err := message.ExtractMessage(r)
	if err != nil {
		decoyHandler(w, r)
	}
	// Init AesParams (request scoped)
	var aesParams = new(cryptography.AesParams)
	// Decrypt message
	decryptedMsg, err := cryptography.Decrypt(msg, aesParams)
	if err != nil {
		fmt.Println(err)
		decoyHandler(w, r)
	}
	// Process decrypted message
	response, err := message.ProcessMessage(decryptedMsg)
	if err != nil {
		fmt.Println(err)
		decoyHandler(w, r)
		return
	}
	// Encrypt message
	encryptedMsg, err := cryptography.Encrypt(response, aesParams)
	if err != nil {
		fmt.Println(err)
		decoyHandler(w, r)
		return
	}
	// Server response
	fmt.Fprintf(w, *encryptedMsg)
}
