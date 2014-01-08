
## PHP
```sh
docker build -t arch_o_matic/php src/github.com/marmelab/arch-o-matic/templates/php/

docker run -t -i -link db:db -name app -v=/vagrant/go/example/php:/var/www arch_o_matic/php /bin/bash
```

## Mysql
```sh
docker build -t arch_o_matic/mysql src/github.com/marmelab/arch-o-matic/templates/mysql/

docker run -d -p 3306 -name db arch_o_matic/mysql

CREATE DATABASE users CHARACTER SET utf8 COLLATE utf8_general_ci;

USE users;
CREATE TABLE users (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(100)
);

INSERT INTO users (id, username) VALUES (NULL, 'manu');
```

## Apache
```sh
docker build -t arch_o_matic/apache src/github.com/marmelab/arch-o-matic/templates/apache/

docker run -d -v=/vagrant/go/example/php:/var/www -link app:app -link db:db -name front1 arch_o_matic/apache
OR
docker run -i -t -v=/vagrant/go/src/github.com/marmelab/arch-o-matic/example/php:/var/www -link app:app -name front1 arch_o_matic/apache /bin/bash
```

## Php FPM
```sh
docker build -t arch_o_matic/php-fpm src/github.com/marmelab/arch-o-matic/templates/php-fpm/

docker run -d -p 9000:9000 -v=/vagrant/go/example/php:/var/www -link db:db -name app arch_o_matic/php-fpm
OR
docker run -t -i -p 9000:9000 -v=/vagrant/go/example/php:/var/www -link db:db -name app arch_o_matic/php-fpm /bin/bash
```

## Nodejs
```sh
docker build -t arch_o_matic/nodejs src/github.com/marmelab/arch-o-matic/templates/nodejs/

docker run -t -i arch_o_matic/nodejs /bin/bash
```

## All
List env variables

```sh
printenv
```

Install ps :
apt-get install -y --reinstall procps
