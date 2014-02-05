
## PHP
```sh
docker build -t gaudi/php src/github.com/marmelab/gaudi/templates/php/

docker run -t -i -link db:db -name app -v=/vagrant/go/example/php:/var/www gaudi/php /bin/bash
```

## Mysql
```sh
docker build -t gaudi/mysql src/github.com/marmelab/gaudi/templates/mysql/

docker run -d -p 3306 -name db gaudi/mysql

CREATE DATABASE symfony CHARACTER SET utf8 COLLATE utf8_general_ci;

INSERT INTO symfony.user (username) VALUES ("Manu");
```

## Apache
```sh
docker build -t gaudi/apache src/github.com/marmelab/gaudi/templates/apache/

docker run -d -v=/vagrant/go/example/php:/var/www -link app:app -link db:db -name front1 gaudi/apache
OR
docker run -i -t -v=/vagrant/go/src/github.com/marmelab/gaudi/example/php:/var/www -link app:app -name front1 gaudi/apache /bin/bash
```

## Php FPM
```sh
docker build -t gaudi/php-fpm src/github.com/marmelab/gaudi/templates/php-fpm/

docker run -d -p 9000:9000 -v=/vagrant/go/example/php:/var/www -link db:db -name app gaudi/php-fpm
OR
docker run -t -i -p 9000:9000 -v=/vagrant/go/example/php:/var/www -link db:db -name app gaudi/php-fpm /bin/bash
```

## Nodejs
```sh
docker build -t gaudi/nodejs src/github.com/marmelab/gaudi/templates/nodejs/

docker run -t -i gaudi/nodejs /bin/bash
```

## Varnish
```sh
docker build -t gaudi/varnish src/github.com/marmelab/gaudi/templates/varnish

docker run -t -i -p 80:80 -name varnish -link db:db gaudi/varnish /bin/bash
```

## Jackrabbit
```sh
docker build -t gaudi/jackrabbit src/github.com/marmelab/gaudi/templates/jackrabbit

docker run -t -i -p 80:80 -name jackrabbit gaudi/jackrabbit /bin/bash
```

## All
List env variables

```sh
printenv
```

Install ps :
apt-get install -y --reinstall procps
