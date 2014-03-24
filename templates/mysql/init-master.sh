#!/bin/sh

mysql -e "GRANT REPLICATION SLAVE ON *.* TO repl@'%' IDENTIFIED BY 'repl'; FLUSH PRIVILEGES;" -uroot

sed -i -e "s/\[mysqld\]/[mysqld]\nlog-bin\nserver-id=1/" /etc/mysql/my.cnf

