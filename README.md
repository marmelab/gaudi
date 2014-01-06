# Arch-o-matic
Builduing your architectures from a single file.

# Starting an architecture

Arch-o-matic can build an architecture with Apache, MySQL and PHP from a file like `apache-php-mysql.yml` :
```yml
containers:
    front1:
        type: apache
        links: [app]
        volumes:
            php: /var/www

    app:
        type: php-fpm
        links: [db]
        ports:
            9000: 9000
        volumes:
            php: /var/www

    db:
        type: mysql
        ports:
            3306: 3306

```

The architecture is started with :

```sh
go run src/github.com/marmelab/arch-o-matic/main.go --config="src/github.com/marmelab/arch-o-matic/example/apache-php-mysql.yml"
```

You can now retrieve the IP of the `db`  container :

```sh
docker inspect db
```

Create a database on it :
```sh
mysql -u root -p -h [IP of db container]

CREATE DATABASE project CHARACTER SET utf8 COLLATE utf8_general_ci;

USE users;
CREATE TABLE users (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(100)
);

INSERT INTO users (id, username) VALUES (NULL, 'manu');
```

Retrieve the IP of the `front1` container:
```sh
docker inspect front1
```

Retrieve the list of users :
```sh
wget http://[IP of the front1 container]/list-users.php
```
