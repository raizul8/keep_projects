#!/bin/bash

# cd /home/raz2/razIT/scriptsExtensions/ShellScripts/GasServerProject/versions/gasAplication_0.1

# dart server.dart & dartium http://127.0.0.1:8080/gasApplication/main_gas.html https://myline-eon.ro/login;
# pidof dart server.dart | xargs kill

cd ~/gas_current/src/gas

# ./gas & ~/Downloads/dartium/chrome http://localhost:8080/ https://myline-eon.ro/login;

./gas & /opt/google/chrome/google-chrome http://localhost:8080/ https://myline-eon.ro/login;

pidof gas | xargs kill
