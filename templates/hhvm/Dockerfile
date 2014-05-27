FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

# Update apt to install HHVM
RUN echo deb http://dl.hhvm.com/debian wheezy main | tee /etc/apt/sources.list.d/hhvm.list

# Install HHVM
RUN apt-get -y update
RUN apt-get install -y --force-yes -f hhvm php5-cli curl

[[ if .EmptyCmd]]
CMD /bin/bash
[[ else ]]
CMD (hhvm --mode daemon -vServer.Type=fastcgi -vServer.Port=[[ .Container.GetFirstLocalPort ]] &) && /bin/bash
[[ end ]]
