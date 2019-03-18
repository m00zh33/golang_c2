# Boilerplate C2 written in Go for red teams aka gorfice2k
This is an attempt to create a sample C2 server in Go. This repo includes Go code for the server and python 3 code for the sample of an agent. The idea behind this project is to provide sort of a boilerplate template that red teams can customize for their own beacons. 

:warning: **Disclaimer: This project should be used for authorized testing or educational purposes only.**

![C2 Demo](/docs/c2_demo.gif)
![C2 Agent Demo](/docs/c2_agent_demo.gif)

## Features
- [x] Accepting payload in Cookie (like Emotet). Can be changed to anything in HTTP request.
- [x] Cycling AES key and iv on each request
- [x] Public/Private cryptography for the AES key
- [x] Google Protobufs for the messages
- [x] MySQL server for tasks
- [ ] Admin panel
- [ ] gRPC to communicate with the main C2
- [ ] Improve concurrency with goroutines

## Installation
### C2
1. Clone this repo
`git clone https://github.com/prsecurity/golang_c2.git`
2. Source the environment
`cd golang_c2; source env`
3. Install Go dependencies
`go get -u github.com/golang/protobuf/protoc-gen-go google.golang.org/grpc github.com/go-sql-driver/mysql`
4. Connect to your MySQL server and create database
`CREATE DATABASE C2;`
5. Import provided sample DB
`mysql -u<user> -p c2 < c2_sample.sql`
6. Launch the server
`go run main.go`
7. If this works, have fun and hack the code :beers:

### Agent
1. In a separate tab, go into agent_demo folder
`cd agent_demo`
2. Initialize venv
`python3 -m venv venv`
3. Activate venv
`. venv/bin/activate`
4. Install python3 dependencies
`pip3 install -r requirements`
5. Launch the agent and watch it connect
`python3 agent.py`

## How to navigate code
By default, I broke down the code into 5 modules:
* config       
* cryptography 
* db
* message 
* server
### Config
Config is responsible for parsing JSON config of the server. You can add custom structures for your config based on needs. Right now it contains the basics, like port, TLS, database credentials and cryptography. Config is loaded in main and guides how server will behave.
### Cryptography
Cryptography contains 3 basic functions
#### LoadCrypto
LoadCrypto intializes the ciphers
#### Decrypt
Takes in string and returns bytes. You can change how decrypt works based on your needs, but out of the box it supports Emotet style encryption.
#### Encrypt
Encrypt takes in bytes and produces a string. You can edit how encryption works based on expectation from your beacon.
### DB
Database handles all DB communications. I didn't want to use any ORM plugins so everything is done via SQL driver. My approach is to use `LoadDB` for initializing DB connection and then create files based on the logic. For example, I have Task.go file that contains the structure for `tasks` table's rows.
### Message
Message is the core of the C2. When you generate protobufs, they are being added to message package. My approach for designing messages are as such:
* Each message exchanged between C2 and beacon is a serialized protobuf.
* Each protobuf is wrapped into Envelope protobuf containing `messageId` and `message` body. 
* The processing of messages depends on `messageId`.
* Each message gets its own handler that consumers the request and produces a response.
* Current implementation is very basic, agent send Knock request and C2 responds with Task response. 
* `taskId` guides what the agent will be doing. For example, taskId 16 tells agent to execute task body as a shell command.
You can extend the messaging based on your beacon needs.

MessageID | Protobuf
------------ | ------------- 
15 | Knock 
16 | Task
17 | TaskResponse

Don't mix up `MessageID` that identifies Protobuf with `TaskID` that identifies *what the beacon will be doing*. In the current sample:

TaskID | Task
------------ | ------------- 
15 | Knock 
16 | Shell

### Server
Server module initializes the web server and handler functions. This is a pretty standard Go web server.


:warning: This is a very raw code, I did this project to get up to speed with Go. I may continue improving this, but if you need some specific changes, feel free to DM me on [Twitter](https://twitter.com/prsecurity_) 
