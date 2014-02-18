FROM stackbrew/debian:wheezy

# Update apt-get
RUN echo "deb http://ftp.fr.debian.org/debian/ wheezy main" >> /etc/apt/sources.list
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim git-core curl build-essential openssl libssl-dev

# Install Nodejs
RUN git clone https://github.com/joyent/node.git /tmp/node
RUN cd /tmp/node && ./configure --openssl-libpath=/usr/lib/ssl && make && make install

# User defined apt packets
[[ if (.Container.AptPackets) ]]RUN apt-get -y --force-yes install [[range (.Container.AptPackets)]][[.]] [[end]][[end]]

# Install NPM
RUN curl https://www.npmjs.org/install.sh | clean=no sh

# Install modules
[[range (.Container.GetCustomValue "modules")]]
	RUN npm install -g [[.]]
[[end]]

ENV NODE_PATH /usr/local/lib/node_modules

ENTRYPOINT ["npm"]
