FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

# Install PHP
RUN apt-get -y -f install php5-cli php5-curl curl

# Install composer
RUN curl -sS https://getcomposer.org/installer | php && mv composer.phar /usr/local/bin/composer

ENTRYPOINT ["composer"]
