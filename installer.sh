#!/bin/bash

echo "|------------------------------------------------|"
echo "|           Conference Camera installer          |"
echo "| https://github.com/smart-cab/conference-camera |"
echo "|                                                |"
echo "|                                    by ALEGOR   |"
echo "|------------------------------------------------|"

rm .env # remove .env file if exists

apt install golang
apt install nginx

read -p "Enter the ip for the service to work: " ip
read -p "Enter the port for the service to work: " port
read -p "Enable auto generation qr code (1 - yes / 0 - no): " autoqr

# generation .env file
echo "IP=$ip" >> .env
echo "PORT=$port" >> .env
echo "AUTO_QR_CODE=$autoqr" >> .env
echo "DEBUG=0" >> .env # disable debug on production

go build main.go
