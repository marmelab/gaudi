FROM stackbrew/debian:wheezy

# Update apt-get
RUN echo "deb http://ftp.fr.debian.org/debian/ wheezy main" >> /etc/apt/sources.list
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim git-core curl build-essential openssl libssl-dev wget

# User defined apt packets
[[ if (.Container.AptPackets) ]]RUN apt-get -y --force-yes install [[range (.Container.AptPackets)]][[.]] [[end]][[end]]

# Install nodejs with nvm
RUN git clone https://github.com/creationix/nvm.git /.nvm
RUN echo "/.nvm/nvm.sh" >> /etc/bash.bashrc
RUN /bin/bash -c '. /.nvm/nvm.sh && nvm install v0.10.20 && nvm use v0.10.20 && nvm alias default v0.10.20 && ln -s /.nvm/v0.10.20/bin/node /usr/bin/node && ln -s /.nvm/v0.10.20/bin/npm /usr/bin/npm'

# Install NPM
RUN curl https://www.npmjs.org/install.sh | clean=no sh

# Install bower
RUN npm install -g bower

ENV NODE_PATH /usr/local/lib/node_modules

ENTRYPOINT ["bower", "--allow-root"]
