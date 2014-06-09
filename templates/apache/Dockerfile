FROM stackbrew/debian:wheezy

RUN echo "deb http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list
RUN echo "deb-src http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

RUN apt-get install -y -f apache2-mpm-worker libapache2-mod-fastcgi
RUN a2enmod actions fastcgi alias

[[range (.Container.GetCustomValue "modules")]]
    RUN a2enmod [[.]]
[[end]]
RUN service apache2 reload

[[ if(.Container.DependsOf "django" )]]
    RUN apt-get install -y -f python2.7 python-dev python-setuptools libmysqlclient-dev
    RUN easy_install pip
    RUN pip install django==1.6
    RUN pip install mysql-python

    RUN apt-get install -y -f libapache2-mod-wsgi

    RUN echo "WSGIPythonPath /app/[[ (.Collection.GetType "django").GetCustomValue "project_name" "project" ]]:/usr/local/lib/python2.7/site-packages" >> /etc/apache2/httpd.conf
[[ end ]]

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

ADD fastcgi.conf /etc/apache2/mods-enabled/fastcgi.conf
ADD 000-default /etc/apache2/sites-enabled/000-default
ADD ports.conf /etc/apache2/ports.conf

[[ $fastCgi := .Collection.Get (.Container.GetCustomValueAsString "fastCgi") (.Collection.GetType "php-fpm")  ]]
[[ if and $fastCgi (or (eq $fastCgi.Type "php-fpm") (eq $fastCgi.Type "hhvm"))]]
    [[ $memoryLimit := $fastCgi.GetCustomValue "memoryLimit" "128M" ]]
    [[ $maxExecutionTime := $fastCgi.GetCustomValue "maxExecutionTime" "30" ]]
    [[ $maxInputTime := $fastCgi.GetCustomValue "maxInputTime" "60" ]]
    [[ $locale := $fastCgi.GetCustomValue "locale" "Europe/Paris" ]]

    RUN apt-get install -y -f php5-fpm
    RUN sed -i 's|;date.timezone =|date.timezone = "[[ $locale ]]"|g' /etc/php5/fpm/php.ini
    RUN sed -i 's|memory_limit = 128M|memory_limit = [[ $memoryLimit ]]|g' /etc/php5/fpm/php.ini
    RUN sed -i 's|max_execution_time = 30|max_execution_time = [[ $maxExecutionTime ]]|g' /etc/php5/fpm/php.ini
    RUN sed -i 's|max_input_time = 60|max_input_time = [[ $maxInputTime ]]|g' /etc/php5/fpm/php.ini
    RUN sed -i 's|;pm.start_servers|pm.start_servers|g' /etc/php5/fpm/pool.d/www.conf
[[ end ]]
# Add custom setup script
[[ beforeAfterScripts ]]

[[ if .EmptyCmd]]
CMD /bin/bash
[[ else ]]
CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] /bin/bash /root/setup.sh \
    && /etc/init.d/apache2 start \
    [[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
    && /bin/bash
[[ end ]]
