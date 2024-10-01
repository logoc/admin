#!/bin/sh
WORK_DIR=`dirname $0`
if [ "$WORK_DIR" = "." ];then
    WORK_DIR=$PWD
fi

PNAME=gofly

run_daemon() {
    [ -d $WORK_DIR/log ] || mkdir $WORK_DIR/log
    cd $WORK_DIR
    nohup ./bin/${PNAME} >> $WORK_DIR/log/stdout.log 2>>$WORK_DIR/log/stderr.log & 
}

run_daemon