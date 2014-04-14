FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

# Install nodejs with nvm
RUN git clone https://github.com/creationix/nvm.git /.nvm
RUN echo "/.nvm/nvm.sh" >> /etc/bash.bashrc
RUN /bin/bash -c '. /.nvm/nvm.sh && nvm install v0.10.20 && nvm use v0.10.20 && nvm alias default v0.10.20 && ln -s /.nvm/v0.10.20/bin/node /usr/bin/node && ln -s /.nvm/v0.10.20/bin/npm /usr/bin/npm'

# Install NPM
RUN curl https://www.npmjs.org/install.sh | clean=no sh

# Install modules
[[range (.Container.GetCustomValue "modules")]]
	RUN npm install -g [[.]]
[[end]]

ENV NODE_PATH /usr/local/lib/node_modules

ENTRYPOINT ["npm"]
