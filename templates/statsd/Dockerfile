# @see: mikaelhg/docker-rabbitmq
FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

# Install nodejs
[[ .Container.SetCustomValue "nodeVersion" "0.10.20"]]
[[ installNodeJS ]]

RUN git clone https://github.com/etsy/statsd /.statsd
WORKDIR /.statsd

RUN cp /.statsd/exampleConfig.js /.statsd/local.js

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if .EmptyCmd ]]
	CMD /bin/bash
[[ else ]]
	CMD (node stats.js local.js &) && /bin/bash
[[ end ]]
