#! /bin/bash
filename=`date +%Y%m%d-%H%M`
tar zcvf plotm_$filename.tar.gz \
	plotm \
	Measurements.gtpl \
	DBConfig.yml \
	ServerConfig.yml \
	DenyIp.txt \
	public/index.html \
	public/count.html \
	tmp/0000000000000000000000000000000000000000000000000000000000000000.yml \
	YmlFiles/Default_000.yml \
	YmlFiles/SingleshotToPeriodic.yml \
	run.sh \
	tar.sh
