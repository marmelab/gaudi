FROM stackbrew/debian:wheezy

# Update apt-get
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim

# Install PHP 5.4
RUN apt-get -y -f install php5 php5-mysql php5-mcrypt php5-curl curl

CMD ["/bin/bash"]
