#!/bin/bash

# We use root as user to be able to write in docker mounted volumes (mysql can't)
mysql_install_db -user=root -ldata=/usr/lib/mysql

/usr/bin/mysqld_safe &
sleep 5s
echo "GRANT ALL ON *.* TO root@'%';" | mysql

[[ $repl := .Container.GetCustomValue "repl" "" ]]

[[ if eq $repl "master" ]]

sed -i -e "s/\[mysqld\]/[mysqld]\nlog-bin\nserver-id = 1/" /etc/mysql/my.cnf
mysql -e "GRANT REPLICATION SLAVE ON *.* TO repl@'%' IDENTIFIED BY 'repl'; FLUSH PRIVILEGES;" -uroot

[[ end ]]

[[ if eq $repl "slave" ]]
[[ $masterName := (.Container.GetCustomValue "master") ]]
[[ $master := $.Collection.Get $masterName ]]

# Connect to master & retrieve current log file & position
MASTER_STATUS=$(mysql -u root -h $[[ $masterName | ToUpper ]]_PORT_[[ $master.GetFirstPort ]]_TCP_ADDR -e "show master status\G")
MASTER_LOG_FILE=$(echo "$MASTER_STATUS" | grep File | sed 's/File://' | sed 's/^ *//;s/ *$//')
MASTER_LOG_POS=$(echo "$MASTER_STATUS" | grep Position | sed 's/Position://' | sed 's/^ *//;s/ *$//')

sed -i -e "s/\[mysqld\]/[mysqld]\nlog-bin\nserver-id = 2/" /etc/mysql/my.cnf
mysql -e "CHANGE MASTER TO MASTER_HOST='$[[$masterName | ToUpper ]]_PORT_[[ $master.GetFirstPort ]]_TCP_ADDR', MASTER_LOG_POS=245, MASTER_USER='repl', MASTER_PASSWORD='repl', MASTER_LOG_FILE='$MASTER_LOG_FILE', MASTER_LOG_POS=$MASTER_LOG_POS; START SLAVE;" -uroot

[[ end ]]

mysqladmin shutdown
