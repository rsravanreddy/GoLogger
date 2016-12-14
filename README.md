# GoLogger

Go logger is a high performance logger which logs the messages to kafka and comsumers can be implemented to process logs from kafka

## Setup
Before you start the service you need to run  kafka cluster and get valid kafka broker()

### step 1:
Start a kafka cluster,default broker address is ```localhost:9092``` if you used the out of box kafka server properties

### step 2:
`svchost: :8080`  <br />
`kafkaaddress: localhost:9092 /* address of kafka broker obtained in step 1 */`  <br />
`dbuser: root /*optioal*/ `  <br /> 
`dbpassword:  /*optioal*/`  <br />
`dbhost: localhost/*optioal*/ ` <br />
`dbname: Todo /*optioal*/` <br />

### step 3:
`make build` <br />
`./cmd/server/server --config config.yaml server` <br />

You should see if everything went ok

`[GIN-debug] POST   /log                      --> logger/service.(*LogResource).CreateLog-fm (4 handlers)` <br />
`[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default` <br />
[GIN-debug] Listening and serving HTTP on :8080` <br />
