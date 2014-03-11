FROM stackbrew/debian:wheezy

[[ updateApt ]]

# Install PHP 5.4
RUN apt-get -y -f install php5 php5-mysql php5-mcrypt php5-curl curl

# User defined apt packets
[[ if (.Container.AptPackets) ]]RUN apt-get -y --force-yes install [[range (.Container.AptPackets)]][[.]] [[end]][[end]]

CMD ["/bin/bash"]
