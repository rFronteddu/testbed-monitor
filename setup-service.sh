#! /bin/bash
# This script will set up (and re-set-up) the daemon for the testbed-monitor

# build the main program
go build main.go

# copy all the necessary program files for the daemon
sudo rm -R /usr/lib/monitor
sudo mkdir /usr/lib/monitor
sudo cp testbed-monitor.service /usr/lib/monitor/testbed-monitor.service
sudo cp configuration.yaml /usr/lib/monitor/configuration.yaml
sudo cp .env /usr/lib/monitor/.env
sudo cp template.html /usr/lib/monitor/template.html
sudo cp main /usr/lib/monitor/main

# copy the service to the
sudo rm /etc/systemd/system/testbed-monitor.service
sudo cp testbed-monitor.service /etc/systemd/system/testbed-monitor.service
sudo systemctl enable testbed-monitor

# kill previously running programs and start
sudo pkill main
sudo systemctl start testbed-monitor

# view daemon status at anytime with >sudo systemctl status testbed-monitor