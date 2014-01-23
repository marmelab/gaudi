#!/bin/bash

/usr/bin/mysqld_safe &
sleep 5s
echo "GRANT ALL ON *.* TO root@'%';" | mysql
mysqladmin shutdown
