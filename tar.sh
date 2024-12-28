#! /bin/bash
filename=`date +%Y%m%d-%H%M`
tar zcvf plotm_$filename.tar.gz plotm DBConfig.yml ServerConfig.yml DenyIp.txt public/index.html run.sh Measurements.gtpl
