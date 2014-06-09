# @see: mikaelhg/docker-rabbitmq
FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

# Install nodejs
[[ .Container.SetCustomValue "nodeVersion" "0.10.20"]]
[[ installNodeJS ]]

RUN git clone https://github.com/etsy/statsd /.statsd
WORKDIR /.statsd

[[ $firstLinked := .Container.FirstLinked]]

RUN echo "{\n  graphiteHost: process.env.[[ $firstLinked.Name | ToUpper ]]_PORT_2003_TCP_ADDR, \n  graphitePort: 2003, \n  port: [[ .Container.GetFirstLocalPort ]], \n  flushInterval: 10000 \n }" > /.statsd/local.js

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if .EmptyCmd ]]
	CMD /bin/bash
[[ else ]]
	CMD (node stats.js local.js &) && /bin/bash
[[ end ]]
