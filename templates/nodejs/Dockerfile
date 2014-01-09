FROM stackbrew/debian

# Update apt-get
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim git-core curl build-essential openssl libssl-dev

# Install Nodejs
RUN git clone https://github.com/joyent/node.git /tmp/node
RUN cd /tmp/node && ./configure --openssl-libpath=/usr/lib/ssl && make && make install

# Install NPM
RUN git clone git://github.com/isaacs/npm.git /tmp/npm
RUN cd /tmp/node && make install

ENTRYPOINT ["/usr/local/bin/node"]
