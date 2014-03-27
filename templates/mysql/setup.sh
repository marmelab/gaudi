#!/bin/bash

# We use root as user to be able to write in docker mounted volumes (mysql can't)
mysql_install_db -user=root -ldata=/usr/lib/mysql

/usr/bin/mysqld_safe &
sleep 5s
echo "GRANT ALL ON *.* TO root@'%';" | mysql


[[ if eq (.Container.GetCustomValue "repl") "master" ]]

sed -i -e "s/\[mysqld\]/[mysqld]\nlog-bin\nserver-id = 1/" /etc/mysql/my.cnf
mysql -e "GRANT REPLICATION SLAVE ON *.* TO repl@'%' IDENTIFIED BY 'repl'; FLUSH PRIVILEGES;" -uroot

[[ end ]]

[[ if eq (.Container.GetCustomValue "repl") "slave" ]]

sed -i -e "s/\[mysqld\]/[mysqld]\nlog-bin\nserver-id = 2/" /etc/mysql/my.cnf
mysql -e "CHANGE MASTER TO MASTER_HOST='$[[(.Container.GetCustomValue "master") | ToUpper ]]_PORT_[[ ($.Collection.Get (.Container.GetCustomValue "master") ).GetFirstPort ]]_TCP_ADDR', MASTER_LOG_POS=245, MASTER_USER='repl', MASTER_PASSWORD='repl', MASTER_LOG_FILE='/var/lib/mysql/master/mysqld-bin.000005', MASTER_LOG_POS=4; START SLAVE;" -uroot

[[ end ]]

mysqladmin shutdown
