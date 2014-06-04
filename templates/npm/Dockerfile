FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

[[ installNodeJS ]]

# Install NPM
RUN curl https://www.npmjs.org/install.sh | clean=no sh

# Install modules
[[range (.Container.GetCustomValue "modules")]]
    RUN npm install -g [[.]]
[[end]]

ENV NODE_PATH /usr/local/lib/node_modules

ENTRYPOINT ["npm"]
