# Boilerplate C2 written in Go for red teams
This is an attempt to create a sample C2 server in Go. This repo includes Go code for the server and python 3 code for the sample of an agent. The idea behind this project is to provide sort of a boilerplate template that red teams can customize for their own beacons. 

**Disclaimer: This project should be used for authorized testing or educational purposes only.**

![C2 Demo](/docs/c2_demo.gif)
![C2 Agent Demo](/docs/c2_agent_demo.gif)

## Features
- [x] Accepting payload in Cookie (like Emotet)
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
7. If this works, have fun and hack the code

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
