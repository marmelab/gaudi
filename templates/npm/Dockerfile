FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

[[ $version := .Container.GetCustomValue "fastCgi" "0.10.20"]]

# Install nodejs with nvm
RUN git clone https://github.com/creationix/nvm.git /.nvm
RUN echo "/.nvm/nvm.sh" >> /etc/bash.bashrc
RUN /bin/bash -c '. /.nvm/nvm.sh && nvm install v[[ $version ]] && nvm use v[[ $version ]] && nvm alias default v[[ $version ]] && ln -s /.nvm/v[[ $version ]]/bin/node /usr/bin/node && ln -s /.nvm/v[[ $version ]]/bin/npm /usr/bin/npm'

# Install NPM
RUN curl https://www.npmjs.org/install.sh | clean=no sh

# Install modules
[[range (.Container.GetCustomValue "modules")]]
	RUN npm install -g [[.]]
[[end]]

ENV NODE_PATH /usr/local/lib/node_modules

ENTRYPOINT ["npm"]
