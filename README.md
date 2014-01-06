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

```sh
go run src/github.com/marmelab/arch-o-matic/main.go --config="src/github.com/marmelab/arch-o-matic/example/apache-php-mysql.yml"
```

# Killing all the containers

Arch-o-matic does not kill previous running container yet. So you have to run :

```sh
docker ps | grep 'ago' | awk '{print $1}' | xargs docker kill || docker ps -a | grep 'ago' | awk '{print $1}' | xargs docker rm
```

# TODO
- [x] Mysql
- [x] PHP
- [x] PHP-FPM
- [x] Apache
- [ ] Nodejs
- [ ] Nginx
- [ ] Bower
- [ ] Jackrabbit
- [ ] Composer
- [ ] Grunt
- [ ] CouchDB
- [ ] Mongodb
- [ ] Elastic search
- [ ] Supervisor
- [ ] Golang
- [ ] Haproxy
- [ ] HHVM
- [ ] Memcached
- [ ] Postgresql
- [ ] RabbitMQ
- [ ] Redis
- [ ] Sentry
- [ ] Tomcat
- [ ] Varnish
- [ ] Graphite
- [ ] Nagios
- [ ] Jenkins
- [ ] SQLite
- [ ] statsd
- [ ] Hadoop
- [ ] Django
- [ ] Kibana
- [ ] Cassandra
- [ ] Collectd
- [ ] LDAP
