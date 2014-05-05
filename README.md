# gaudi [![Build Status](https://travis-ci.org/marmelab/gaudi.png?branch=master)](https://travis-ci.org/marmelab/gaudi)

Check [gaudi's website](http://gaudi.io) for more informations.


gaudi is a generator of architecture written in Go and using [Docker](http://www.docker.io).
You can use it to start any type of application, and link them together without knowledge of Docker or system configuration.
Using Go, gaudi can build and start your applications in parallel depending of they dependencies.

[![gaudi screencast](http://gaudi.io/builder/img/gaudi-video.jpg)](http://showterm.io/83b5d24c67cd39a73de23)

# Basic Usage

Describe any architecture with a simple YAML file (called `.gaudi.yml`). For instance, for a PHP+MySQL combo:

```yaml
applications:
    front1:
        type: apache
        links: [app]
        volumes:
            .: /var/www
        custom:
            fastCgi: app

    app:
        type: php-fpm
        links: [db]
        ports:
            9000: 9000
        volumes:
            .: /var/www

    db:
        type: mysql
        ports:
            3306: 3306
```

Start this environment (with sudo privileges):

```sh
gaudi
```

Gaudi will try to find a `.gaudi.yml` file in the current folder, and start each application simultaneously, or sequentially if they depend on each other.

# Installation

## Requirements

- [Go 1.2](http://golang.org/doc/install)
- [Docker](https://www.docker.io/gettingstarted/)

## Install Gaudi

```sh
go get github.com/marmelab/gaudi
```

Check that your `PATH` includes `$GOPATH/bin`:

```sh
export PATH=$GOPATH/bin:/$PATH
```

The `gaudi` application starts containers with Docker's commands which [requires sudo privileges](http://docs.docker.io/en/latest/use/basics/#dockergroup).
Make sure that the `GOPATH` and `GOROOT` environment variables are correctly set for the `root` user (or other user with root privileges).

# Options

- `--config=""` Specify the location of the configuration file
- `--debug` Display some useful information
- `--no-cache` Do not use docker's cache when building containers (builds will be slower)
- `--quiet` Do not display build & pull output
- `rebuild` Force all containers to rebuild (useful when you change your configuration)
- `stop` Stop all applications
- `clean` Stop & remove all applications
- `check` Check if all applications are running
- `run binaryName [arguments]` Run a specific binary

# How Does It Work?

Gaudi uses [Docker](http://www.docker.io) to start all applications in a specific container.
It builds Docker files and specific configuration files from different templates.
All templates are listed in [the `templates/` folder](https://github.com/marmelab/gaudi/tree/master/templates), one for each application type.

# Examples

You can find an example of [how to start a Symfony application](https://github.com/marmelab/gaudi/wiki/HOW-TO:-Run-a-Symfony-Application) in the wiki.

Another examples can be found in [the `example` folder](https://github.com/marmelab/gaudi/tree/master/example).

# Configuration

After changing your configuration, you should restart gaudi with the `rebuild` argument:

`gaudi rebuild`

## Common Configuration

The YML file describing the architecture should have a section called `applications:`.

### Type

You can specify what kind of application you want to run:

```yaml
applications:
    [Application name]:
        type: [one of the listed type below]
```

Application types are listed below.

### Links

When an application depends on another, you can link them together:

```yaml
applications:
    app1:
        type: varnish
        links: [front1, front2]
    front1:
        type: apache
    front2:
        type: apache
```

Here the `app1` application receives environment variables for each link, as follows:

```
FRONT1_NAME=/front1/app1
FRONT1_PORT=tcp://172.17.0.215:80
FRONT1_PORT_3306_TCP_PORT=80
FRONT1_PORT_3306_TCP_PROTO=tcp
FRONT1_PORT_3306_TCP_ADDR=172.17.0.215
FRONT1_PORT_3306_TCP=tcp://172.17.0.215:80
```

### Ports

To open some ports on an application:

```yaml
applications:
    front1:
        type: apache
        ports:
            80: 8080
```

The port 80 in the host machine will be mapped to the 8080 in the container.

### Volumes

You can add you own folders by mounting volumes:

```yaml
applications:
    front1:
        type: apache
        volumes:
            php: /app/php
```

The `php/` folder (absolute or relative to the yml file) will be mounted in the `/app/php` folder in the application.

### Environment variables

Environment variables can be injected with the `environments` field:

```yaml
applications:
    db:
        type: index
        image: paintedfox/postgresql
        environments:
            USER: docker
            PASS: docker
            DB: gaudi
```

Environment variable will be set thanks to the docker's `-e` argument.

### Apt packets

If you want to install other apt packets, use the `apt_get` parameter:

```yaml
applications:
    app:
        type: apache
        apt_get: [php5-gd, php5-intl]
```

### Add files into container

Files can be added directly into the container with the `add` parameter.
Contrary to `volumes`, `add` can only be used for files.

```yaml
applications:
    app:
        type: apache
        add:
	        conf.ini: /root/conf.ini
	        dir/conf2.ini: /root/conf2.ini
```

`add` takes a map of `relative file on the host machine`: `destination path on the container`.

### Remote Containers

If you want to run an application not yet supported by Gaudi, you can use a prebuilt image (on Github), or an image from the [Docker index](https://index.docker.io/):

Github images can be pulled via the `github` type:
```yaml
applications:
    redis:
        type: github
        image: gary/redis
        path: github.com/manuquentin/docker-redis
        ports:
            6379: 6379
```

Images from Docker index uses the `index` type:

```yaml
applications:
    mongodb:
        type: index
        image: dockerfile/mongodb
        ports:
            27017: 27017
            28017: 28017
```

### Before and after scripts

Each type of application listed bellow can be configured with `before_script` and `after_script` field.
Theses fields can represents a file path or a command. They are executed before or after the main executable of the application.

Example: Start a node server
```yaml
applications:
	front:
		type: nodejs
		volumes:
            .: /app
        after_script: node /app/server.js
```

Before scripts can added to run custom script before the application boots:
```yaml
applications:
	db:
		type: mysql
		volumes:
            .: /app
		before_script: /app/bin/init.sh
```

## Types

Each application uses a `custom` section to define its own custom configuration settings.

### Varnish

```yaml
applications:
    [name]:
        type: varnish
        links: [front1, front2]
    custom:
        backends: [front1, front2]
```

The `backends` custom parameter defines which applications are load balanced by Varnish. Theses applications have to be linked together using `links`.

### Nginx

As a webserver:

```yaml
applications:
    [name]:
        type: nginx
        links: [app]
    custom:
        fastCgi: app
```

As a load balancer:

```yaml
applications:
    [name]:
        type: nginx
        links: [front1, front2]
    custom:
        backends: [front1, front2]
```

The `backends` custom parameter defines which applications are load balanced by Nginx. Theses applications have to be linked together using `links`.


### Apache

```yaml
applications:
    [name]:
        type: apache
    custom:
        fastCgi: app
```

The `fastCgi` custom parameter points out an application where to forward Fast-CGI scripts.

### MySQL

```yaml
applications:
    [name]:
        type: mysql
    custom:
        repl: master # Or "slave"
        master: master # When using "repl: slave": indicate the name of the master
```

The `repl` custom value indication if the MySQL instance is declared as `master` or `slave`.
When a MySQL is defined as slave, you should set it's the master application name in the `master` params.

### PHP

This application is a simple php5 service, if you want to use it with `Apache` or `Nginx`, use the `PHP-FPM` one.

```yaml
applications:
    [name]:
        type: php
```

### PHP-FPM

```yaml
applications:
    [name]:
        type: php-fpm
```

### Nodejs

To start a Node.js application, use the `after_script` parameter. If the `after_script` is not set, Node.js will run without arguments.

```yaml
applications:
    [name]:
        type: nodejs
        after_script: node /app/server.js
```

### Cassandra

```yaml
applications:
    [name]:
        type: cassandra
        ports:
            9160: 9160
            7000: 7000
        custom:
            maxHeapSize: 512M # Optionnal
            heapNewSize: 256M # Optionnal
```

### Jackrabbit

```yaml
applications:
    [name]:
        type: jackrabbit
```

### PhpMyAdmin

```yaml
applications:
    pma:
        type: phpmyadmin
        ports:
            80: 80
        links: [db]

    db:
        type: mysql
        ports:
            3306: 3306
```

## Binaries

Gaudi can also runs binaries in the current folder.
A binary is not always attached to an application so Gaudi allows to configure them in a different field `binaries`:

```yaml
binaries:
    [name]:
        type: [type]

```

To run a binary with gaudi, simple use (for composer for instance) :

```sh
gaudi run [name] [arguments]
```

For `npm`:

```sh
gaudi run npm install
```

Or `bower`:

```sh
gaudi run bower install angularjs --save
```

### Npm

```yaml
binaries:
    npm:
        type: npm
```

### Composer

```yaml
binaries:
    composer:
        type: composer
```

### Bower

```yaml
binaries:
    bower:
        type: bower
```

### Jekyll

```yaml
binaries:
    jekyll:
        type: jekyll
```

## Contributing

Your feedback about the usage of gaudi in your specific context is valuable, don't hesitate to [open GitHub Issues](https://github.com/marmelab/gaudi/issues) for any problem or question you may have.

All contributions are welcome. New applications or options should be tested  with go unit test tool.

## License

Gaudi is licensed under the [MIT Licence](LICENSE), courtesy of [marmelab](http://marmelab.com).
