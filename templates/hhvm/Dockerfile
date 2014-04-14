FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

# Update apt to install HHVM
RUN wget -O - http://dl.hhvm.com/conf/hhvm.gpg.key | apt-key add -
RUN apt-key adv --recv-keys --keyserver keyserver.ubuntu.com 5A16E7281BE7A449
RUN apt-key adv --recv-keys --keyserver keyserver.ubuntu.com 40976EAF437D05B5

RUN echo deb http://dl.hhvm.com/debian wheezy main | tee /etc/apt/sources.list.d/hhvm.list
RUN echo deb http://archive.ubuntu.com/ubuntu precise main universe >> /etc/apt/sources.list

# Install HHVM
RUN apt-get -y update
RUN apt-get -y -f install hhvm-fastcgi php5-cli curl

# Install composer
RUN curl -sS https://getcomposer.org/installer | php && mv composer.phar /usr/local/bin/composer

CMD (hhvm --mode daemon -vServer.Type=fastcgi -vServer.Port=[[ .Container.GetFirstLocalPort ]] &) && /bin/bash
