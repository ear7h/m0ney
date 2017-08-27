#!/bin/sh

if [ $EAR7H_ENV = "prod" ]
then
    echo waiting 15 seconds for db init
    sleep 15
fi


echo starting server
./m0ney | tee log.txt
