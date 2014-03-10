FROM stackbrew/debian:wheezy

# Update apt-get
RUN echo "deb http://ftp.fr.debian.org/debian/ wheezy main" >> /etc/apt/sources.list
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim git-core curl build-essential openssl libssl-dev

# User defined apt packets
[[ if (.Container.AptPackets) ]]RUN apt-get -y --force-yes install [[range (.Container.AptPackets)]][[.]] [[end]][[end]]

# Install PHP
RUN apt-get -y -f install php5-cli php5-curl curl

# Install composer
RUN curl -sS https://getcomposer.org/installer | php && mv composer.phar /usr/local/bin/composer

ENTRYPOINT ["composer"]
