## How to run the sample
1) Start the nats server with -V flag

`./nats-server -V`

2) Run the sample

`go run go-sample.go`

3) You should see the a log entry similar to the below one in the command execution window confirming that message is sent to the NATS server

`2021/05/23 16:39:17 Published [hello] : 'world'`

4) You should see a set of log entries in the nats-server console similar to below

`[1420936] 2021/05/23 16:39:17.273691 [TRC] 127.0.0.1:41878 - cid:4 - <<- [CONNECT {"verbose":false,"pedantic":false,"tls_required":false,"name":"","lang":"go","version":"1.11.0","protocol":1,"echo":true,"headers":true,"no_responders":true}]
[1420936] 2021/05/23 16:39:17.274182 [TRC] 127.0.0.1:41878 - cid:4 - "v1.11.0:go" - <<- [PING]
[1420936] 2021/05/23 16:39:17.274390 [TRC] 127.0.0.1:41878 - cid:4 - "v1.11.0:go" - ->> [PONG]
[1420936] 2021/05/23 16:39:17.275034 [TRC] 127.0.0.1:41878 - cid:4 - "v1.11.0:go" - <<- [PUB hello 5]
[1420936] 2021/05/23 16:39:17.275107 [TRC] 127.0.0.1:41878 - cid:4 - "v1.11.0:go" - <<- MSG_PAYLOAD: ["world"]
[1420936] 2021/05/23 16:39:17.275138 [TRC] 127.0.0.1:41878 - cid:4 - "v1.11.0:go" - <<- [PING]
[1420936] 2021/05/23 16:39:17.275158 [TRC] 127.0.0.1:41878 - cid:4 - "v1.11.0:go" - ->> [PONG]`
