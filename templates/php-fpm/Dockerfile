FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

# Install PHP FPM
RUN apt-get -y -f install php5-fpm php5-cli php5-mysql php5-mcrypt php5-curl curl

# Install composer
RUN curl -sS https://getcomposer.org/installer | php && mv composer.phar /usr/local/bin/composer

[[ $memoryLimit := .Container.GetCustomValue "memoryLimit" "128M" ]]
[[ $maxExecutionTime := .Container.GetCustomValue "maxExecutionTime" "30" ]]
[[ $maxInputTime := .Container.GetCustomValue "maxInputTime" "60" ]]
[[ $locale := .Container.GetCustomValue "locale" "Europe/Paris" ]]

RUN sed -i 's|listen = /var/run/php5-fpm.sock|listen = [[ .Container.GetFirstLocalPort ]]|g' /etc/php5/fpm/pool.d/www.conf
RUN sed -i 's|;cgi.fix_pathinfo=0|cgi.fix_pathinfo=0|g' /etc/php5/fpm/pool.d/www.conf
RUN sed -i 's|;date.timezone =|date.timezone = "[[ $locale ]]"|g' /etc/php5/fpm/php.ini
RUN sed -i 's|memory_limit = 128M|memory_limit = [[ $memoryLimit ]]|g' /etc/php5/fpm/php.ini
RUN sed -i 's|max_execution_time = 30|max_execution_time = [[ $maxExecutionTime ]]|g' /etc/php5/fpm/php.ini
RUN sed -i 's|max_input_time = 60|max_input_time = [[ $maxInputTime ]]|g' /etc/php5/fpm/php.ini

ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if .EmptyCmd ]]
CMD /bin/bash
[[ else ]]
CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] /bin/bash /root/setup.sh \
    && php5-fpm -R \
    [[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
    && /bin/bash
[[ end ]]
