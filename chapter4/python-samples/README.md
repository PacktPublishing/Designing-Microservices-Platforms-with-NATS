## How to run the sample
1) Start the nats server with -V flag

`./nats-server -V`

2) Run the sample

`python3 sample.py`

3) You should see the a log entry similar to the below one in the command execution window confirming that message is sent to the NATS server

`Received a message on 'patient.profile.json ': Hello
Received a message on 'patient.treatments.json ': World`


4) You should see a set of log entries in the nats-server console similar to below

`[1420936] 2021/05/23 16:38:00.404866 [TRC] 127.0.0.1:41836 - cid:3 - <<- [CONNECT {"echo": true, "lang": "python3", "pedantic": false, "protocol": 1, "verbose": false, "version": "0.11.4"}]
[1420936] 2021/05/23 16:38:00.405026 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - <<- [PING]
[1420936] 2021/05/23 16:38:00.405037 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - ->> [PONG]
[1420936] 2021/05/23 16:38:00.405609 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - <<- [SUB patient.*.json  1]
[1420936] 2021/05/23 16:38:00.405660 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - <<- [PUB patient.profile.json  5]
[1420936] 2021/05/23 16:38:00.405678 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - <<- MSG_PAYLOAD: ["Hello"]
[1420936] 2021/05/23 16:38:00.405789 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - ->> [MSG patient.profile.json 1 5]
[1420936] 2021/05/23 16:38:00.405851 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - <<- [PUB patient.treatments.json  5]
[1420936] 2021/05/23 16:38:00.405863 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - <<- MSG_PAYLOAD: ["World"]
[1420936] 2021/05/23 16:38:00.405887 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - ->> [MSG patient.treatments.json 1 5]
[1420936] 2021/05/23 16:38:00.405964 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - <<- [PING]
[1420936] 2021/05/23 16:38:00.405990 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - ->> [PONG]
[1420936] 2021/05/23 16:38:00.406028 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - <<- [UNSUB 1 1]
[1420936] 2021/05/23 16:38:00.406519 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - <<- [PING]
[1420936] 2021/05/23 16:38:00.406533 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - ->> [PONG]
[1420936] 2021/05/23 16:38:00.406927 [TRC] 127.0.0.1:41836 - cid:3 - "v0.11.4:python3" - <-> [DELSUB 1]`
