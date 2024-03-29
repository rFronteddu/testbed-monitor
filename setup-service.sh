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
sudo cp traps.json /usr/lib/monitor/traps.json
sudo cp report_template.html /usr/lib/monitor/report_template.html
sudo cp notification_template.html /usr/lib/monitor/notification_template.html
sudo cp endpoints_mapping.yaml /usr/lib/monitor/endpoints_mapping.yaml
sudo cp main /usr/lib/monitor/main

# copy the service to the
sudo rm /etc/systemd/system/testbed-monitor.service
sudo cp testbed-monitor.service /etc/systemd/system/testbed-monitor.service
sudo systemctl enable testbed-monitor

# kill previously running programs and start
sudo pkill main
sudo systemctl start testbed-monitor

# view daemon status at anytime with >sudo systemctl status testbed-monitor