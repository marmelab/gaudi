FROM stackbrew/debian

# Update apt-get
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim

RUN apt-get install -y -f
RUN apt-get install -y -f nginx

RUN echo "daemon off;" >> /etc/nginx/nginx.conf

ADD default /etc/nginx/sites-enabled/default

CMD nginx
