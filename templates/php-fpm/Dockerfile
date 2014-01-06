FROM stackbrew/debian

# Update apt-get
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim

# Install PHP FPM
RUN apt-get -y -f install php5-fpm php5-mysql php5-mcrypt php5-curl curl

RUN sed -i 's|listen = /var/run/php5-fpm.sock|listen = 9000|g' /etc/php5/fpm/pool.d/www.conf

ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

# Little hack to start a process in a container without supervisor
CMD /bin/bash /root/setup.sh && service php5-fpm start && tail -f /var/log/php5-fpm.log
