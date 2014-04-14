FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

# Install jekyll
RUN apt-get -y --force-yes install ruby1.9.1-dev
RUN gem install jekyll

# User defined apt packets
[[ if (.Container.AptPackets) ]]RUN apt-get -y --force-yes install [[range (.Container.AptPackets)]][[.]] [[end]][[end]]

ENTRYPOINT ["jekyll"]
