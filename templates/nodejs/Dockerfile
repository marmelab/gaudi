FROM stackbrew/debian

# Update apt-get
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim git-core curl build-essential openssl libssl-dev

# Install Nodejs
RUN git clone https://github.com/joyent/node.git /tmp/node
RUN cd /tmp/node && ./configure --openssl-libpath=/usr/lib/ssl && make && make install

# User defined apt packets
[[ if (.Container.AptPackets) ]]RUN apt-get -y --force-yes install [[range (.Container.AptPackets)]][[.]] [[end]][[end]]

# Install NPM
RUN git clone git://github.com/isaacs/npm.git /tmp/npm
RUN cd /tmp/npm && make install

# Install modules
[[range (.Container.GetCustomValue "modules")]]
	RUN npm install -g [[.]]
[[end]]

ENV NODE_PATH /usr/local/lib/node_modules

[[ if (.Container.HasAfterScript) ]]
	CMD [[.Container.AfterScript]] && /bin/bash
[[ else ]]
	CMD ["/usr/local/bin/node"] && /bin/bash
[[ end]]
