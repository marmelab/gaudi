# Gaudi
Gaudi is a generator of architecture written in Go and using [Docker](http://www.docker.io).
You can use it to start any types of applications and link them together without knowledge of Docker and system configuration.
Using Go, Gaudi can build and start your applications in parallel depending of they dependencies.

# Basic usage
The architecture can be described with a single file (called `.gaudi.yml`) :
```yml
containers:
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

This environment can be started with :

```sh
gaudi
```

Gaudi will try to find a `.gaudi.yml` file in the current folder and start each application simultaneously according to their dependencies.

# Installation
```sh
go get github.com/marmelab/gaudi
```

Check that yout PATH includes `$GOPATH/bin`
```sh
export PATH=$GOPATH/bin:/$PATH
```

# Options
- `--config=""` Specify the location of the configuration file
- `--rebuild` Rebuild all applications (with this option, data not stored in volumes will be lost)
- `--stop` Stop all applications
- `--check` Check if all applications are running

# How does it work

Gaudi uses [Docker](http://www.docker.io) to start all applications in a specific container.
It builds a Docker files and specific configuration files from different templates.
Each templates are listed in the `templates` folder, one for each application type.

# Configuration

## Common Configuration

The YML file describing the architecture should have a section called `containers`.

### Type
You can specify what king a application you want to run :
```yml
containers:
	[Application name]:
		type: [one of the listed type below]
```

Application types are listed below.

### Links
When an applications depends on another, you can link them :
```yml
containers:
	app1:
		type: varnish
		links: [front1, front2]
	front1:
		type: apache
	front2:
		type: apache
```

Here the `app1` application will receive environment variables for each link like :
```
FRONT1_NAME=/front1/app1
FRONT1_PORT=tcp://172.17.0.215:80
FRONT1_PORT_3306_TCP_PORT=80
FRONT1_PORT_3306_TCP_PROTO=tcp
FRONT1_PORT_3306_TCP_ADDR=172.17.0.215
FRONT1_PORT_3306_TCP=tcp://172.17.0.215:80
```

### Ports
To open some ports on an applications :
```yml
containers:
	front1:
		type: apache
		ports:
			80:8080
```

The port 80 inb the host machine will be mapped to the 8080 in the container.

### Volumes
You can add you own files by mounting volumes :
```yml
containers:
	front1:
		type: apache
		volumes:
			php:/app/php
```

The php folder (absolute or relative to the yml files) will be mounted in the /app/php folder in the application.

## Example : Running a Symfony Application

To run a Symfony app, we need an Apache server, PHP-FPM & MySQL. So we create a `.gaudi.yml` in a new folder :
```yml
containers:
    front:
        type: apache
        links: [app]
        ports:
            80: 80
        volumes:
            .: /var/www
        custom:
            fastCgi: app
            documentRoot: /var/www/web
            modules: [rewrite]

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
        after_script: mysql -e "CREATE DATABASE symfony CHARACTER SET utf8 COLLATE utf8_general_ci;" -uroot
        volumes:
            .gaudi/mysql: /var/lib/mysql
```

Here the current folder in mounter in `/var/www` (the default DocumentRoot of Apache), we enable the m`od_rewrite` module for Apache, link it to PHP-FPM which is linked to a MySQL app.
We also mount a `.gaudi/mysql`folder to allow MySQL storage persistence outside of the container.

To start the Symfony app, run in the current folder :
```sh
mkdir -p .gaudi/mysql
gaudi

Cleaning front ...
Cleaning app ...
Cleaning db ...
Building gaudi/front ...
Building gaudi/app ...
Building gaudi/db ...
Application db started (172.17.0.79:3306)
Application app started (172.17.0.80:9000)
Application front started (172.17.0.81:80)
```

The building process can take a few minutes the first time. The next time Docker will find the difference between the last build and run only we it need.

## Types

Each application uses a `custom` section in the configuration to defines them own aspect.

### Varnish
```yml
containers:
    [name]:
        type: varnish
        links: [front1, front2]
    custom:
        backends: [front1, front2]
```

`backends` custom param is used to defines which containers are load balanced by Varnish. Theses containers have to be linked with `links`.

### Nginx

#### As a webserver:
```yml
containers:
    [name]:
        type: nginx
        links: [app]
    custom:
        fastCgi: app
```

#### As a load balancer:
```yml
containers:
    [name]:
        type: nginx
        links: [front1, front2]
    custom:
        backends: [front1, front2]
```

`backends` custom param is used to defines which containers are load balanced by Nginx. Theses containers have to be linked with `links`.


### Apache
```yml
containers:
    [name]:
        type: apache
    custom:
        fastCgi: app
```

`fastCgi` custom param is used to point out an application to forward Fast-CGI scripts.

## Contributing

Your feedback about the usage of gaudi in your specific context is valuable, don't hesitate to [open GitHub Issues](https://github.com/marmelab/gaudi/issues) for any problem or question you may have.

All contributions are welcome. New applications or options should be tested  with go unit test tool.

## License

Gaudi is licensed under the [MIT Licence](LICENSE), courtesy of [marmelab](http://marmelab.com).
