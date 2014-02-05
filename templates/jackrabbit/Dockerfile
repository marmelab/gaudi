FROM stackbrew/debian:wheezy

# Update apt-get
RUN apt-get -y update
RUN apt-get -y --force-yes install aptitude wget python-software-properties gcc g++ make

# Install Java 7
RUN aptitude install -y -f openjdk-7-jre

RUN mkdir -p /opt/jackrabbit/

# Install jackrabbit
RUN wget http://archive.apache.org/dist/jackrabbit/2.6.0/jackrabbit-standalone-2.6.0.jar -P /opt/jackrabbit/

CMD cd /opt/jackrabbit/ && java -jar jackrabbit-standalone-2.6.0.jar --port [[ .Container.GetFirstPort ]] && /bin/bash
