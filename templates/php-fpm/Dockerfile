FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

# Install PHP FPM
RUN apt-get -y -f install php5-fpm php5-cli php5-mysql php5-mcrypt php5-curl curl

# Install composer
RUN curl -sS https://getcomposer.org/installer | php && mv composer.phar /usr/local/bin/composer

RUN sed -i 's|listen = /var/run/php5-fpm.sock|listen = [[ .Container.GetFirstLocalPort ]]|g' /etc/php5/fpm/pool.d/www.conf
RUN sed -i 's|;cgi.fix_pathinfo=0|cgi.fix_pathinfo=0|g' /etc/php5/fpm/pool.d/www.conf
RUN sed -i 's|;date.timezone =|date.timezone = "Europe/Paris"|g' /etc/php5/fpm/php.ini

ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

# Add custom setup script
[[ beforeAfterScripts ]]

CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] /bin/bash /root/setup.sh \
	&& php5-fpm -R \
	[[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
	&& /bin/bash
