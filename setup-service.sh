#! /bin/bash
go build main.go
sudo rm -R /usr/lib/monitor
sudo mkdir /usr/lib/monitor
sudo cp testbed-monitor.service /usr/lib/monitor/testbed-monitor.service
sudo cp configuration.yaml /usr/lib/monitor/configuration.yaml
sudo cp .env /usr/lib/monitor/.env
sudo cp main /usr/lib/monitor/main
sudo cp testbed-monitor.service /etc/systemd/system/testbed-monitor.service
sudo systemctl enable testbed-monitor