FROM stackbrew/debian

RUN echo "deb http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list
RUN echo "deb-src http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list

# Update apt-get
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim curl

RUN curl http://repo.varnish-cache.org/debian/GPG-key.txt | apt-key add -
RUN echo "deb http://repo.varnish-cache.org/debian/ wheezy varnish-3.0" >> /etc/apt/sources.list
RUN apt-get update
RUN apt-get -y install varnish

ADD varnish.conf /etc/default/varnish
ADD default.vcl /etc/varnish/default.vcl

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

CMD /bin/bash /root/setup.sh && /etc/init.d/varnish start && varnishlog
