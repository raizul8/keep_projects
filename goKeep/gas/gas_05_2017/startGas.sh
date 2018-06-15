#!/bin/bash

# cd /home/raz/IdeaProjects1/goKeep/gas/gas_05_2017/src/gasHttp
cd /home/raz/CodeProjects/keep_projects/goKeep/gas/gas_05_2017/src/gasHttp

./gasHttp false & /opt/google/chrome/google-chrome http://localhost:8000/static/gasWeb/gas.html https://myline-eon.ro/login;

pidof gasHttp | xargs kill




