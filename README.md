# Gaudi

Gaudi is a generator of architecture written in Go and using [Docker](http://www.docker.io).
You can use it to start any type of application, and link them together without knowledge of Docker or system configuration.
Using Go, Gaudi can build and start your applications in parallel depending of they dependencies.

[![Gaudi screencast](http://marmelab.com/gaudi/img/gaudi-video.jpg)](http://showterm.io/83b5d24c67cd39a73de23)

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

```sh
go get github.com/marmelab/gaudi
```

Check that your `PATH` includes `$GOPATH/bin`:

```sh
export PATH=$GOPATH/bin:/$PATH
```

The `gaudi` application starts containers with Docker's commands which [requires sudo privileges](http://docs.docker.io/en/latest/use/basics/#dockergroup).
Make sure that the `GOPATH` and `GOROOT` environment variables are correctly set for the `root` user (or other user with root privileges).

## Optional Build Time Improvement
All containers uses the same base image, to speed up the first build run:

```sh
docker pull stackbrew/debian
```

# Options

- `--config=""` Specify the location of the configuration file
- `--stop` Stop all applications
- `--check` Check if all applications are running

# How Does It Work?

Gaudi uses [Docker](http://www.docker.io) to start all applications in a specific container.
It builds Docker files and specific configuration files from different templates.
All templates are listed in [the `templates/` folder](https://github.com/marmelab/gaudi/tree/master/templates), one for each application type.

You can find an example of [how to start a Symfony application](https://github.com/marmelab/gaudi/wiki/HOW-TO:-Run-a-Symfony-Application) in the wiki.

# Configuration

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

You can add you own files by mounting volumes:

```yaml
applications:
    front1:
        type: apache
        volumes:
            php: /app/php
```

The `php/` folder (absolute or relative to the yml file) will be mounted in the `/app/php` folder in the application.

### Remote Containers

If you want to run an application not yet supported by Gaudi, you can use a prebuilt image, or an image from the [Docker index](https://index.docker.io/):

```yaml
applications:
    server:
        type: nodejs
        links: [redis]
        ports:
            80: 80
        volumes:
            nodejs-redis: /app

    redis:
        type: remote
        image: gary/redis
        path: github.com/manuquentin/docker-redis
        ports:
            6379: 6379
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

## Contributing

Your feedback about the usage of gaudi in your specific context is valuable, don't hesitate to [open GitHub Issues](https://github.com/marmelab/gaudi/issues) for any problem or question you may have.

All contributions are welcome. New applications or options should be tested  with go unit test tool.

## License

Gaudi is licensed under the [MIT Licence](LICENSE), courtesy of [marmelab](http://marmelab.com).
