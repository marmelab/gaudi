FROM stackbrew/debian:wheezy

# Update apt-get
RUN apt-get -y update
RUN apt-get -y --force-yes install aptitude wget python-software-properties gcc g++ make

# User defined apt packets
[[ if (.Container.AptPackets) ]]RUN apt-get -y --force-yes install [[range (.Container.AptPackets)]][[.]] [[end]][[end]]

# Install Java 7
RUN aptitude install -y -f openjdk-7-jre

RUN mkdir -p /opt/jackrabbit/

# Install jackrabbit
RUN wget http://archive.apache.org/dist/jackrabbit/2.6.5/jackrabbit-standalone-2.6.5.jar -P /opt/jackrabbit/

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

CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] cd /opt/jackrabbit/ \
	&& (java -jar jackrabbit-standalone-2.6.5.jar --port [[ .Container.GetFirstLocalPort ]] &) \
	[[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
	&& /bin/bash
