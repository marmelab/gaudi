#!/bin/bash

composer install --prefer-source

mysql -e 'create database IF NOT EXISTS sandbox;' -u root

php app/console doctrine:phpcr:init:dbal -e=test
