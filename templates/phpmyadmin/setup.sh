#!/bin/bash

# Add envvars to apache configuration files
envs=`printenv`

for env in $envs
do
    IFS== read name value <<< "$env"

    echo "export $name='$value'" >> /tmp/dockerenv
done

echo ". /tmp/dockerenv" >> /etc/apache2/envvars
mkdir -p /var/www/cgi-bin
