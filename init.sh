#!/bin/sh

sleep 1000000

#echo hello world
nohup m0ney > server.out
nohup daemon/daemon > daemon.out