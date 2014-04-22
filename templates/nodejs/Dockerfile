FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

# Install nodejs with nvm
RUN git clone https://github.com/creationix/nvm.git /.nvm
RUN /bin/bash -c '. /.nvm/nvm.sh && nvm install v0.10.20 && nvm use v0.10.20 && nvm alias default v0.10.20 && ln -s /.nvm/v0.10.20/bin/node /usr/bin/node && ln -s /.nvm/v0.10.20/bin/npm /usr/bin/npm'

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
