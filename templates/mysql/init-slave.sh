#!/bin/sh

sed -i -e "s/\[mysqld\]/[mysqld]\nlog-bin\nserver-id=2\nmaster-host = $[[(.Container.GetCustomValue "master" | ToUpper)]]_PORT_[[ ($.Collection.Get . ).GetFirstPort ]]_TCP_ADDR/" /etc/mysql/my.cnf

mysql -e "CHANGE MASTER TO  MASTER_HOST='$[[(.Container.GetCustomValue "master" | ToUpper)]]_PORT_[[ ($.Collection.Get . ).GetFirstPort ]]_TCP_ADDR', MASTER_LOG_FILE='', MASTER_LOG_POS=4; START SLAVE;" -uroot
