#!/bin/bash

# Replace php-fpm user and group
sed -i 's/user = www-data/user = 0/g' /etc/php5/fpm/pool.d/www.conf
sed -i 's/group = www-data/group = 0/g' /etc/php5/fpm/pool.d/www.conf

# Add envvars to PHP-FPM configuration files
envs=`printenv`

for env in $envs
do
    IFS== read name value <<< "$env"

	echo "env[$name] = $value" >> /etc/php5/fpm/php-fpm.conf
done
