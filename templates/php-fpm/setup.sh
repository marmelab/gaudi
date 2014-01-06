#!/bin/bash

# Add envvars to PHP-FPM configuration files
envs=`printenv`

for env in $envs
do
    IFS== read name value <<< "$env"

	echo "env[$name] = $value" >> /etc/php5/fpm/php-fpm.conf
done
