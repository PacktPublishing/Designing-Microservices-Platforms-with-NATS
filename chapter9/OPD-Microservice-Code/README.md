# Instructions to test the Microservices based OPD application that uses NATS

## Starting the NATS Server
`nats-server`

## Starting the microservices

### Starting the Registration Service
`cd registration-service`

`go run cmd/main.go -dbName opd_data -dbUser root -dbPassword Root@1985`

2021/07/11 13:56:12 Starting NATS Microservices OPD Sample - Registration Service version 0.1.0
2021/07/11 13:56:12 Listening for HTTP requests on 0.0.0.0:9090


### Starting the Inspection Service
`cd inspection-service`

`go run cmd/main.go -dbName opd_data -dbUser root -dbPassword Root@1985`

2021/07/11 13:56:18 Starting NATS Microservices OPD Sample - Inspection Service version 0.1.0
2021/07/11 13:56:18 Listening for HTTP requests on 0.0.0.0:9091


### Starting the Treatment Service
`cd treatment-service`

`go run cmd/main.go -dbName opd_data -dbUser root -dbPassword Root@1985`

2021/07/11 13:56:26 Starting NATS Microservices OPD Sample - Treatment Service version 0.1.0
2021/07/11 13:56:26 Listening for HTTP requests on 0.0.0.0:9092


### Starting the Release Service
`cd release-service`

`go run cmd/main.go -dbName opd_data -dbUser root -dbPassword Root@1985`

2021/07/11 13:54:25 Starting NATS Microservices OPD Sample - Release Service version 0.1.0
2021/07/11 13:54:25 Listening for HTTP requests on 0.0.0.0:9093

## Trying out the use case

### Patient registration

#### Register a Patient
- Request

`curl "http://localhost:9090/opd/patient/register" -X POST -d '{"full_name":"chanaka fernando","address":"44, sw19, london","id":200, "sex":"male", "phone":222222222}'`
- Response

{"id":200,"token":1}

#### View the Patient
- Request

`curl "http://localhost:9090/opd/patient/view/200"`
- Response

{"full_name":"chanaka fernando","address":"44, sw19, london","id":200, "sex":"male", "phone":222222222}

#### Update a Patient
- Request

`curl "http://localhost:9090/opd/patient/update" -X PUT -d '{"full_name":"chanaka fernando","address":"667/280/6, liyanagemulla, seeduwa","id":200, "sex":"male", "phone":222222222}'`
- Response

"Record for Patient updated sucessfully"

#### Generate a Token
- Request

`curl "http://localhost:9090/opd/patient/token/300"`
- Response

{"id":300,"token":2}

### Patient Inspection

#### View the pending inspections
- Request

`curl "http://localhost:9091/opd/inspection/pending"`
- Response

{"id":200,"token":1}

#### Update an inspection report
- Request

`curl "http://localhost:9091/opd/inspection/record" -X POST -d '{"id":200, "time":"2021/07/12 10:21 AM", "observations":"ABC Syncrome", "medication":"XYZ x 3", "tests":"FBT, UC", "notes":"possible Covid-19"}'`
- Response

{"id":200,"medication":"XYZ x 3","tests":"FBT, UC","notes":"possible Covid-19"}

#### View inspection history of a patient
- Request

`curl "http://localhost:9091/opd/inspection/history/200"`
- Response

[{"id":200,"time":"2021/07/12 10:21 AM","observations":"ABC Syncrome","medication":"XYZ x 3","tests":"FBT, UC","notes":"possible Covid-19"}]

### Patient treatment

#### View pending treatments
- Request

`curl "http://localhost:9092/opd/treatment/pending"`
- Response

[{"id":200,"medication":"XYZ x 3","tests":"FBT, UC","notes":"possible Covid-19"}]

#### Add a medication record
- Request

`curl "http://localhost:9092/opd/treatment/medication" -X POST -d '{"id":200,"time":"2021 07 12 4:35 PM","dose":"xyz x 1, abc x 2","notes":"low fever"}'`
- Response

"Record updated successfully"

#### Add a Test record
- Request

`curl "http://localhost:9092/opd/treatment/tests" -X POST -d '{"id":200,"time":"2021 07 12 4:35 PM","test_name":"FBC","status":"sample collected", "notes":"possible covid-19"}'`
- Response

"Test recorded successfully"

#### View medication history
- Request

`curl "http://localhost:9092/opd/treatment/history/200"`
- Response

[{"id":200,"time":"2021 07 12 4:35 PM","notes":"low fever"},{"id":200,"time":"2021 07 12 4:35 PM","notes":"low fever"}]

#### View test history
- Request

`curl "http://localhost:9092/opd/treatment/tests/200"`
- Response

[{"id":200,"time":"2021 07 12 4:35 PM","test_name":"FBC","status":"sample collected","notes":"possible covid-19"}]

#### Initiate a patient release
- Request

`curl "http://localhost:9092/opd/treatment/release" -X POST -d '{"id":200,"time":"2021 07 12 8:35 PM","next_state":"discharge","post_medication":"NM x 10 days"}'`
- Response

"Release event published"

### Patient release

#### View pending releases
- Request

`curl "http://localhost:9093/opd/release/pending"`
- Response

[{"id":200,"time":"2021 07 12 8:35 PM","next_state":"discharge","post_medication":"NM x 10 days"}]

#### Add a patient release record
- Request

`curl "http://localhost:9093/opd/release/discharge" -X POST -d '{"id":200,"time":"2021 07 12 9:35 PM","state":"discharge","post_medication":"NM x 10 days","next_visit":"2021 08 12 09:00AM"}'`
- Response

"Patient discharge recorded successfully"


### Port Numbers for microservices

- Registration Service - 9090
- Inspection Service - 9091
- Treatment Service - 9092
- Release Service - 9093
