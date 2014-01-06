# Arch-o-matic
Builduing your architectures from a single file.

# Starting an architecture

With a simple `apache-php-mysql.yml` :
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

An architecture can be started with :

```sh
go run src/github.com/marmelab/arch-o-matic/main.go --config="src/github.com/marmelab/arch-o-matic/example/apache-php-mysql.yml"
```
