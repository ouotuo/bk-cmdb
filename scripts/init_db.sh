
#!/bin/bash
set -e

# get local IP.
localIp=`python ip.py`
curl -X POST -H 'BK_USER:migrate' -H 'HTTP_BLUEKING_SUPPLIER_ID:0' http://${localIp}:admin_port_placeholer/migrate/v3/migrate/community/0

echo ""

mkdir -p /data/bk-cmdb
export GOPATH=/data/bk-cmdb 
cd $GOPATH/src
git clone git@github.com:linclin/bk-cmdb.git  configcenter
git remote add bk  git@github.com:Tencent/bk-cmdb.git
git fetch bk
git merge bk/master
git push