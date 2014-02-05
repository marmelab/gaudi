FROM stackbrew/debian:wheezy

# Update apt-get
RUN apt-get -y update
RUN apt-get -y --force-yes install python-software-properties gcc g++ make

RUN apt-get -y --force-yes install mysql-server

# Edit mysql config (use root as user to be able to write in docker mounted volumes)
RUN sed -i -e "s/^user\s*=\s*mysql/user = root/" /etc/mysql/my.cnf
RUN sed -i -e "s/^bind-address\s*=\s*127.0.0.1/bind-address\t\t= 0.0.0.0/" /etc/mysql/my.cnf
RUN grep -v rootfs /proc/mounts > /etc/mtab

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
	&& /etc/init.d/mysql start \
	[[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
	&& /bin/bash
