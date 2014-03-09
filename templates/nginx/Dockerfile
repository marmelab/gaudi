FROM stackbrew/debian:wheezy

# Update apt-get
RUN apt-get -y update
RUN apt-get -y install python-software-properties gcc g++ make vim

RUN apt-get install -y -f nginx

# User defined apt packets
[[ if (.Container.AptPackets) ]]RUN apt-get -y --force-yes install [[range (.Container.AptPackets)]][[.]] [[end]][[end]]


RUN echo "daemon off;" >> /etc/nginx/nginx.conf

ADD default /etc/nginx/sites-enabled/default

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

# Add custom setup script
[[ if (.Container.HasBeforeScript) ]]
	[[ if (.Container.HasBeforeScriptFile) ]]
		ADD [[.Container.BeforeScript]] /root/before-setup.sh
	[[ else ]]
		RUN echo '[[.Container.BeforeScript]]' > /root/before-setup.sh
	[[ end ]]

	RUN chmod +x /root/before-setup.sh
[[ end ]]

[[ if (.Container.HasAfterScript) ]]
	[[ if (.Container.HasAfterScriptFile) ]]
		ADD [[.Container.AfterScript]] /root/after-setup.sh
	[[ else ]]
		RUN echo '[[.Container.AfterScript]]' > /root/after-setup.sh
	[[ end ]]

	RUN chmod +x /root/after-setup.sh
[[ end ]]

CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] /bin/bash /root/setup.sh \
	&& /etc/init.d/nginx start \
	[[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
	&& /bin/bash
