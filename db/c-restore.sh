#!/bin/bash

if [ -f /var/ear7h/m0ney/db/latest.sql ]; then
    echo "restoring from backup"
    echo "/var/ear7h/m0ney/db/latest.sql"
    mysql < /var/ear7h/m0ney/db/latest.sql
    else
        echo "/var/ear7h/m0ney/db/latest.sql not found"
fi
