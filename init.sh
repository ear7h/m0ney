#!/bin/sh

echo waiting 15 seconds for db init
sleep 15

#echo hello world
nohup m0ney > server.out
nohup daemon/daemon > daemon.out