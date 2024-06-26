#!/bin/bash

echo "|------------------------------------------------|"
echo "|           Conference Camera installer          |"
echo "| https://github.com/smart-cab/conference-camera |"
echo "|                                                |"
echo "|                                    by ALEGOR   |"
echo "|------------------------------------------------|"

rm .env # remove .env file if exists

apt install nodejs
apt install golang
apt install nginx
apt install uvcdynctrl

read -p "Enter the ip for the service to work: " ip
read -p "Enter the port for the service to work: " port
read -p "Enable auto generation qr code (1 - yes / 0 - no): " autoqr
read -p "Enter school: " school

# generation .env file
echo "IP=$ip" >> .env
echo "AUTO_QR_CODE=$autoqr" >> .env
echo "DEBUG=0" >> .env # disable debug on production
echo "TOKEN_LENGTH=8" >> .env
echo "REACT_APP_SCHOOL=$school" >> frontend/.env.production

cd frontend
npm i
cd ../
go build .
