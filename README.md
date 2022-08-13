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

## Connections
**UDP** - The testbed monitor receives the host report on port 8758.<br>
**ICMP** - If the testbed monitor has not received a report from a host in 60 minutes, it pings that host via ICMP.

## Installation
### Prerequisites
* go > 1.6
### Install
```
    go get 
    go build
```
