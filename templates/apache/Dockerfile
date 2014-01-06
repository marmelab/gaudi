FROM stackbrew/debian

RUN echo "deb http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list
RUN echo "deb-src http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list

# Update apt-get
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim

RUN apt-get install -y -f
RUN apt-get install -y -f apache2-mpm-worker libapache2-mod-fastcgi php5-fpm
RUN sed -i 's/;pm.start_servers/pm.start_servers/g' /etc/php5/fpm/pool.d/www.conf
RUN a2enmod actions fastcgi alias
RUN service apache2 reload

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

RUN echo '<IfModule mod_fastcgi.c>\n\tAddHandler php5-fcgi .php\n\tAction php5-fcgi /cgi-bin/php5.external\n\n\t<Location "/cgi-bin/php5.external">\n\t\tOrder Deny,Allow\n\t\tDeny from All\n\t\tAllow from env=REDIRECT_STATUS\n\t</Location>\n</IfModule>' > /etc/apache2/mods-enabled/fastcgi.conf
RUN sed -i 's|</VirtualHost>|\tFastCgiExternalServer /var/www/cgi-bin/php5.external -host ${APP_PORT_9000_TCP_ADDR}:${APP_PORT_9000_TCP_PORT}\n\tAlias /cgi-bin/ /var/www/cgi-bin/\n</VirtualHost>|g' /etc/apache2/sites-enabled/000-default
RUN sed -i 's|ScriptAlias /cgi-bin/ /usr/lib/cgi-bin/||g' /etc/apache2/sites-enabled/000-default

CMD /bin/bash /root/setup.sh && /usr/sbin/apache2ctl -D FOREGROUND
