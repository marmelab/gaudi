FROM stackbrew/debian

# Update apt-get
RUN apt-get -y update
RUN apt-get -y --force-yes install python-software-properties gcc g++ make

RUN apt-get -y --force-yes install mysql-server

# Edit mysql config
RUN sed -i -e"s/^bind-address\s*=\s*127.0.0.1/bind-address = 0.0.0.0/" /etc/mysql/my.cnf

RUN echo "mysqld_safe &" > /root/setup.sh
RUN echo "sleep 5s" >> /root/setup.sh
RUN echo "mysql -u root -e \"GRANT ALL ON *.* TO root@'%';\"" >> /root/setup.sh
RUN echo "mysqladmin shutdown" >> /root/setup.sh
RUN chmod +x /root/setup.sh

CMD ./root/setup.sh && /usr/sbin/mysqld
