# testbed-monitor
## Description
This is the Testbed Monitor app, created by rFronteddu and maintained by cb218. It is a simple component designed to receive and log status reports from a host. The program generates daily and weekly aggregate reports and emails them to the address specified.

The app is designed to run in conjunction with the **Host Monitor** (https://github.com/rFronteddu/host-monitor) app, another utility created by rFronteddu.

## Inputs
**main.go** - Declare the hosts to be monitored in TowerIPs[]<br>
**mailer.go** - Declare the email to send reports to in the Mail() function

## Connections
**UDP** - The testbed monitor receives the host report on port 8758.<br>
**ICMP** - If the testbed monitor has not received a report from a host in 60 minutes, it pings that host via ICMP.

## Outputs
The reports and returned pings are saved in .txt files named from the date and time *(YYYY-MM-DD_HHMMSS)*.

## Installation
### Prerequisites
* go > 1.6
### Install
```
    go get 
    go build
```
