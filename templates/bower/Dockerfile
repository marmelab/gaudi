FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

[[ installNodeJS ]]

# Install NPM
RUN curl https://www.npmjs.org/install.sh | clean=no sh

# Install bower
RUN npm install -g bower

ENV NODE_PATH /usr/local/lib/node_modules

ENTRYPOINT ["bower", "--allow-root"]
