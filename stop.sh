#!/bin/sh
WORK_DIR=`dirname $0`
if [ "$WORK_DIR" = "." ];then
    WORK_DIR=$PWD
fi

PNAME=gofly

pid=`ps -ef | grep $PNAME | grep -v grep | awk '{print $2}'`
if [ ! -z "$pid" ]; then
     kill -9 $pid
     echo "success stop pid:"$pid
fi
