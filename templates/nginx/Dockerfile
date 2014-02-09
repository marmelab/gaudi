FROM stackbrew/debian

# Update apt-get
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim

RUN apt-get install -y -f nginx

# User defined apt packets
[[ if (.Container.AptPackets) ]]RUN apt-get -y --force-yes install [[range (.Container.AptPackets)]][[.]] [[end]][[end]]


RUN echo "daemon off;" >> /etc/nginx/nginx.conf

ADD default /etc/nginx/sites-enabled/default

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

CMD /bin/bash /root/setup.sh && /etc/init.d/nginx start && /bin/bash
