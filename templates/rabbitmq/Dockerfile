# @see: mikaelhg/docker-rabbitmq
FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

# Install RabbitMQ
RUN echo "deb http://www.rabbitmq.com/debian/ testing main" > /etc/apt/sources.list.d/rabbitmq.list
RUN wget -O - http://www.rabbitmq.com/rabbitmq-signing-key-public.asc | apt-key add -
RUN apt-get -y update
RUN apt-get -y -f install rabbitmq-server 
RUN /usr/sbin/rabbitmq-plugins enable rabbitmq_management
RUN echo "[{rabbit, [{loopback_users, []}]}]." > /etc/rabbitmq/rabbitmq.config

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if .EmptyCmd ]]
	CMD /bin/bash
[[ else ]]
	CMD (/usr/sbin/rabbitmq-server &) && /bin/bash
[[ end ]]
