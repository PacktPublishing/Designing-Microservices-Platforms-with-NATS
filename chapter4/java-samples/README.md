## How to run the sample
1) Start the nats server with -V flag

`./nats-server -V`

2) Run the sample

`./gradlew run`

3) You should see the a log entry similar to the below one in the command execution window confirming that message is sent to the NATS server

`> Task :app:run
Hello World!
`

4) You should see a set of log entries in the nats-server console similar to below

`[1420936] 2021/05/23 16:40:00.915548 [TRC] 127.0.0.1:41896 - cid:5 - <<- [CONNECT {"lang":"java","version":"2.11.2","protocol":1,"verbose":false,"pedantic":false,"tls_required":false,"echo":true,"headers":true,"no_responders":true}]
[1420936] 2021/05/23 16:40:00.917094 [TRC] 127.0.0.1:41896 - cid:5 - "v2.11.2:java" - <<- [PING]
[1420936] 2021/05/23 16:40:00.920313 [TRC] 127.0.0.1:41896 - cid:5 - "v2.11.2:java" - ->> [PONG]
[1420936] 2021/05/23 16:40:00.936999 [TRC] 127.0.0.1:41896 - cid:5 - "v2.11.2:java" - <<- [PUB patient.treatments 39]
[1420936] 2021/05/23 16:40:00.937086 [TRC] 127.0.0.1:41896 - cid:5 - "v2.11.2:java" - <<- MSG_PAYLOAD: ["{“tablets”:[panadol, asithromizin]}"]
[1420936] 2021/05/23 16:40:00.960413 [TRC] 127.0.0.1:41896 - cid:5 - "v2.11.2:java" - <<- [PING]
[1420936] 2021/05/23 16:40:00.960484 [TRC] 127.0.0.1:41896 - cid:5 - "v2.11.2:java" - ->> [PONG]`


