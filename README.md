# gaudi [![Build Status](https://travis-ci.org/marmelab/gaudi.png?branch=master)](https://travis-ci.org/marmelab/gaudi)


**This project is discontinued. Read about why [here](http://marmelab.com/blog/2014/11/06/retiring-gaudi.html)**.

gaudi is a generator of architecture written in Go and using [Docker](http://www.docker.io).
You can use it to start any type of application, and link them together without knowledge of Docker or system configuration.
Using Go, gaudi can build and start your applications in parallel depending on their dependencies.

Check [gaudi's website](http://gaudi.io) and follows [@GaudiBuilder](https://twitter.com/GaudiBuilder) for more information.

[![gaudi screencast](http://gaudi.io/builder/img/gaudi-video.jpg)](https://vimeo.com/97235816)

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
        ports:
            8080: 8080

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

gaudi will try to find a `.gaudi.yml` file in the current folder, and start each application simultaneously, or sequentially if they depend on each other.

# Installation

gaudi requires [Docker](https://www.docker.io/gettingstarted/) to run.

## OSX / Windows: Using Vagrant

The [Cethy/vagrant-gaudi](https://github.com/Cethy/vagrant-gaudi) repository describes how to install gaudi with Vagrant.

## Debian & Ubuntu

```sh
wget -O - http://gaudi.io/apt/gaudi.gpg.key | sudo apt-key add -
echo "deb http://gaudi.io/apt/ precise main" | sudo tee -a /etc/apt/sources.list

sudo apt-get update
sudo apt-get install gaudi
```

## Other linux systems

On other system you need to install [Go 1.2](http://golang.org/doc/install) to install gaudi.

```sh
go get github.com/marmelab/gaudi
```

Check that your `PATH` includes `$GOPATH/bin`:

```sh
export PATH=$GOPATH/bin:/$PATH
```

## Via Puppet

A [puppet module](https://forge.puppetlabs.com/cethy/gaudi) is available to install gaudi.



The `gaudi` application starts containers with Docker's commands which [requires sudo privileges](http://docs.docker.io/en/latest/use/basics/#dockergroup).
Make sure that the `GOPATH` and `GOROOT` environment variables are correctly set for the `root` user (or other user with root privileges).

# How Does It Work?

gaudi uses [Docker](http://www.docker.io) to start all applications in a specific container.
It builds Docker files and specific configuration files from different templates.
All templates are listed in [the `templates/` folder](https://github.com/marmelab/gaudi/tree/master/templates), one for each application type.

# Examples

You can find an example of [how to start a Symfony application](https://github.com/marmelab/gaudi/wiki/HOW-TO:-Run-a-Symfony-Application) in the wiki.

Another examples can be found in [the `examples` folder](https://github.com/marmelab/gaudi/tree/master/examples).

# Options

See [gaudi options](http://gaudi.io/installation.html#options).

# Configuration

Check [How to configure gaudi to build your environment](http://gaudi.io/configuration.html)

## Types

See [all type of applications supported](http://gaudi.io/components.html).

## Binaries

gaudi can also runs binaries in the current folder.
A binary is not always attached to an application so gaudi allows to configure them in a different field `binaries`.

See [all type of binaries supported](http://gaudi.io/binaries.html).

## Build the debian package

### Create a gpg key

```sh
gpg --gen-key
ls / -R
gpg --armor --export your@email.com --output gaudi.gpg.key
```

### Run makefile

```sh
make apt
```

## Contributing

Your feedback about the usage of gaudi in your specific context is valuable, don't hesitate to [open GitHub Issues](https://github.com/marmelab/gaudi/issues) for any problem or question you may have.

All contributions are welcome. New applications or options should be tested  with go unit test tool.

## License

gaudi is licensed under the [MIT Licence](LICENSE), courtesy of [marmelab](http://marmelab.com).
