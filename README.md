# Boilerplate C2 written in Go for red teams
This is an attempt to create a sample C2 server in Go. This repo includes Go code for the server and python 3 code for the sample of an agent. The idea behind this project is to provide sort of a boilerplate template that red teams can customize for their own beacons. 

**Disclaimer: This project should be used for authorized testing or educational purposes only.**

![C2 Demo](/docs/c2_demo.gif)
![C2 Agent Demo](/docs/c2_agent_demo.gif)
Features:
- [x] Accepting payload in Cookie (like Emotet)
- [x] Cycling AES key and iv on each request
- [x] Public/Private cryptography for the AES key
- [x] Google Protobufs for the messages
- [x] MySQL server for tasks
- [ ] Admin panel
- [ ] gRPC to communicate with the main C2

## Installation
1. Clone this repo
`git clone https://github.com/prsecurity/golang_c2.git`
2. Source the environment
`cd golang_c2; source env`
