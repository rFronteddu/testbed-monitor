# testbed-monitor
## Description
This is the Testbed Monitor app, created by rFronteddu and maintained by cb218. It is a simple component designed to receive and log status reports from a host. The program generates daily and weekly aggregate reports and emails them to the address specified.

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
DESTINATION=
```
The **configuration.yaml** file should declare the following variables:
```
# The port the receiver runs on
RECEIVE_PORT: 
# How often the system checks if it is time to send a report
AGGREGATE_PERIOD: 
# The hour the report is sent
AGGREGATE_HOUR: 
# How long should the receiver wait for a host report before trying to ping
EXPECTED_REPORT_PERIOD: 
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