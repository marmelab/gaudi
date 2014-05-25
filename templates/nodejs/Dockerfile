FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

[[ $version := .Container.GetCustomValue "version" "0.10.20"]]

# Install nodejs with nvm
RUN git clone https://github.com/creationix/nvm.git /.nvm
RUN /bin/bash -c '. /.nvm/nvm.sh && nvm install v[[ $version ]] && nvm use v[[ $version ]] && nvm alias default v[[ $version ]] && ln -s /.nvm/v[[ $version ]]/bin/node /usr/bin/node && ln -s /.nvm/v[[ $version ]]/bin/npm /usr/bin/npm'

# Install modules
[[range (.Container.GetCustomValue "modules")]]
	RUN npm install -g [[.]]
[[end]]

ENV NODE_PATH /usr/local/lib/node_modules

[[ if .EmptyCmd]]
CMD /bin/bash
[[ else ]]
	[[ if (.Container.HasAfterScript) ]]
		CMD [[.Container.AfterScript]] && /bin/bash
	[[ else ]]
		CMD ["/usr/local/bin/node"] && /bin/bash
	[[ end]]
[[ end ]]
