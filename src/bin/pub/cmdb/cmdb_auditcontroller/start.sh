#!/bin/bash

set -e

# get local IP.
localIp=`python ip.py`

# make log-dir if not exists.
if [[ ! -d "./logs" ]];then
    mkdir ./logs
fi

# set execute 
chmod +x cmdb_auditcontroller
./cmdb_auditcontroller --addrport=${localIp}:50005 --logtostderr=false --log-dir=./logs --v=0 --regdiscv=rd_server_placeholer > ./logs/std.log 2>&1 &
