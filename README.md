# testbed-monitor
## Description
This is the Testbed Monitor app, created by rFronteddu and  cb218. It is a simple component designed to receive and process status reports from a host. The program generates daily and weekly aggregate reports and emails them to the address specified and posts them to the API specified. Notifications related to field value limits can also be configured.

The app is designed to run in conjunction with the **Host Monitor** (https://github.com/rFronteddu/host-monitor) app, another utility created by rFronteddu.

## Inputs
The **.env** file should declare the following variables:
```
# The email to send the report from
EMAIL=
PASSWORD=
# The SMTP host
HOST=
# The SMTP port
MAIL_PORT=
# The email to send the report to
# Multiple emails can be separated with a comma, ex: DESTINATION=email1,email2,email3
DESTINATION=
```
The **configuration.yaml** file should declare the following variables:
```
# If the program should monitor hosts, set MONITOR_HOSTS to true
MONITOR_HOSTS: false
# The port the receiver runs on
RECEIVE_PORT: "8758"
# How often the system checks if it is time to send a report
AGGREGATE_PERIOD: "2"
# The hour the report is sent
AGGREGATE_HOUR: "23"
# How long should the receiver wait for a host report before trying to ping
EXPECTED_REPORT_PERIOD: "1"
# The port to access GQL playground on server
GQL_PORT: "8081"
# The IP address and port to post/get data to the API
API_IP: ""
API_PORT: "4100"
# MQTT subscription information
MQTTBroker: ""
MQTTTopic: ""

# If the program should monitor a testbed via ping, set MONITOR_TESTBED to true
MONITOR_TESTBED: true
# Provide the IP address of the testbed to be monitored
TESTBED_IP: "127.0.0.1"
# How often the testbed should be pinged
PING_PERIOD: "30"
```
The **traps.json** file should declare the following variables, if notifications at field value limits are desired:
```
[{
  "id": "",
  "description": "",
  "operator": "",
    - Currently supported operators are <, >, and =
  "field": "",
    - Currently supported fields are Reboots, UsedRAM, UsedDisk, Cpu, and Temperature
  "trigger": "",
  "period": ""
    - Unit: hours
}]
```


## Installation
### Prerequisites
* go > 1.6
### Install
If using a linux system, the **testbed-monitor.service** daemon can be used to build and run the program.
```
    ./setup-service.sh
```
On other systems, the program can be started manually after building.
```
    go build main.go
    ./main
```