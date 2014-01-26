FROM stackbrew/debian

# Update apt-get
RUN apt-get -y update
RUN apt-get -y --force-yes install python-software-properties gcc g++ make

RUN apt-get -y --force-yes install mysql-server

# Edit mysql config (use root as user to be able to write in docker mounted volumes)
RUN sed -i -e "s/^user\s*=\s*mysql/user = root/" /etc/mysql/my.cnf
RUN sed -i -e "s/^bind-address\s*=\s*127.0.0.1/bind-address\t\t= 0.0.0.0/" /etc/mysql/my.cnf

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

# Add custom setup script
[[ if (.Container.HasBeforeScript) ]]
	[[ if (.Container.HasBeforeScriptFile) ]]
		ADD [[.Container.BeforeScript]] /root/custom-setup.sh
	[[ else ]]
		RUN echo [[.Container.BeforeScript]] > /root/custom-setup.sh
	[[ end ]]

	RUN chmod +x /root/custom-setup.sh
[[ end ]]

CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/custom-setup.sh && [[end]] /bin/bash /root/setup.sh && /etc/init.d/mysql start && /bin/bash
