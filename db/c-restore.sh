#!/bin/bash

if [ -f /tmp/foo.txt ]; then
    echo "restoring from backup"
    echo "/var/ear7h/m0ney/db/latest.sql"
    mysql < /var/ear7h/m0ney/db/latest.sql
fi
