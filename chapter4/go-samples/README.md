## How to run the sample
1) Start the nats server with -V flag

`./nats-server -V`

2) Run the sample

`go run go-sample.go`

3) You should see the a log entry similar to the below one in the command execution window confirming that message is sent to the NATS server

`2021/05/23 16:39:17 Published [hello] : 'world'`

4) You should see a set of log entries in the nats-server console similar to below

`[2988049] 2021/05/28 18:01:37.469037 [TRC] 127.0.0.1:43476 - cid:3 - <<- [CONNECT {"verbose":false,"pedantic":false,"tls_required":false,"name":"","lang":"go","version":"1.11.0","protocol":1,"echo":true,"headers":true,"no_responders":true}]
[2988049] 2021/05/28 18:01:37.469313 [TRC] 127.0.0.1:43476 - cid:3 - "v1.11.0:go" - <<- [PING]
[2988049] 2021/05/28 18:01:37.469378 [TRC] 127.0.0.1:43476 - cid:3 - "v1.11.0:go" - ->> [PONG]
[2988049] 2021/05/28 18:01:37.475668 [TRC] 127.0.0.1:43476 - cid:3 - "v1.11.0:go" - <<- [PUB patient.profile 18]
[2988049] 2021/05/28 18:01:37.475742 [TRC] 127.0.0.1:43476 - cid:3 - "v1.11.0:go" - <<- MSG_PAYLOAD: ["{\"name\":\"parakum\"}"]
[2988049] 2021/05/28 18:01:37.475772 [TRC] 127.0.0.1:43476 - cid:3 - "v1.11.0:go" - <<- [PING]
[2988049] 2021/05/28 18:01:37.475785 [TRC] 127.0.0.1:43476 - cid:3 - "v1.11.0:go" - ->> [PONG]`
