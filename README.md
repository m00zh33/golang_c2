# Boilerplate C2 written in Go for red teams
This is an attempt to create a sample C2 server in Go. This repo includes Go code for the server and python 3 code for the sample of an agent.

Features:
- [x] Accepting payload in Cookie (like Emotet)
- [x] Cycling AES key and iv on each request
- [x] Public/Private cryptography for the AES key
- [x] Google Protobufs for the messages
- [x] MySQL server for tasks
- [ ] Admin panel
- [ ] gRPC to communicate with the main C2

