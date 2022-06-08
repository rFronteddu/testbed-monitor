# testbed-monitor
## Description
This is the Testbed Monitor app, created by rFronteddu and maintained by cbarry. It is a simple component designed to receive and log status reports from a host. When the monitor does not receive a report after a certain amount of time, the host is pinged to ensure it is still running. If the ping is unsuccessful, the host engineer will be notified.
The app is designed to run in conjunction with the **Host Monitor** *(https://github.com/rFronteddu/host-monitor)* app, another utility created by rFronteddu.

## GRPC
The testbed monitor pings the host on port 8090.

## Installation
### Prerequisites
* go > 1.6
### Install
```
    go get 
    go build
```
