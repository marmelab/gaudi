FROM stackbrew/debian:wheezy

# Update apt-get
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim

# Install PHP FPM
RUN apt-get -y -f install php5-fpm php5-cli php5-mysql php5-mcrypt php5-curl curl

# Install composer
RUN curl -sS https://getcomposer.org/installer | php && mv composer.phar /usr/local/bin/composer

RUN sed -i 's|listen = /var/run/php5-fpm.sock|listen = 9000|g' /etc/php5/fpm/pool.d/www.conf
RUN sed -i 's|;cgi.fix_pathinfo=0|cgi.fix_pathinfo=0|g' /etc/php5/fpm/pool.d/www.conf

ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

CMD /bin/bash /root/setup.sh && service php5-fpm start && /bin/bash
