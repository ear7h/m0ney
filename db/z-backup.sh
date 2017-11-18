#!/bin/bash

function main_loop {
    while true; do
        echo "dumping in 3 days"
        sleep 3d
        echo "dumping..."
        mysqldump --databases money > /var/ear7h/m0ney/db/latest.sql
        cat /var/ear7h/m0ney/db/latest.sql > /var/ear7h/m0ney/db/$(date +%Y-%m-%d).sql.gz
        echo "dumped"
    done
}

main_loop&