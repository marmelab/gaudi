#!/bin/bash

# We use root as user to be able to write in docker mounted volumes (mysql can't)
mysql_install_db -user=root -ldata=/usr/lib/mysql

/usr/bin/mysqld_safe &
sleep 5s
echo "GRANT ALL ON *.* TO root@'%';" | mysql
mysqladmin shutdown
