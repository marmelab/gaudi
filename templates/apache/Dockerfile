FROM stackbrew/debian:wheezy

RUN echo "deb http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list
RUN echo "deb-src http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list

[[ updateApt ]]
[[ addUserFiles ]]

RUN apt-get install -y -f
RUN apt-get install -y -f apache2-mpm-worker libapache2-mod-fastcgi php5-fpm
RUN sed -i 's/;pm.start_servers/pm.start_servers/g' /etc/php5/fpm/pool.d/www.conf
RUN a2enmod actions fastcgi alias

[[range (.Container.GetCustomValue "modules")]]
	RUN a2enmod [[.]]
[[end]]
RUN service apache2 reload

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

ADD fastcgi.conf /etc/apache2/mods-enabled/fastcgi.conf
ADD 000-default /etc/apache2/sites-enabled/000-default
ADD ports.conf /etc/apache2/ports.conf

# Add custom setup script
[[ beforeAfterScripts ]]

CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] /bin/bash /root/setup.sh \
	&& /etc/init.d/apache2 start \
	[[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
	&& /bin/bash
