# Arch-o-matic
Builduing your architectures from a single file.

# Starting an Architecture
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

Here the port 80 will be mapped to the 8080

### Volumes
You can add you own files by mounting volumes :
```yml
containers:
	front1:
		type: apache
		volumes:
			php:/app/php
```

The the php folder (absolute or relative to the yml files) will be mounted in the /app/php folder in the container.

## Types
### Apache
```yml
containers:
    [name]:
        type: apache
```
