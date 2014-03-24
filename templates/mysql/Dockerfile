FROM stackbrew/debian:wheezy

[[ updateApt ]]

# Install mysql
RUN apt-get -y --force-yes install mysql-server

# Edit mysql config (use root as user to be able to write in docker mounted volumes)
RUN sed -i -e "s/^user\s*=\s*mysql/user = root/" /etc/mysql/my.cnf
RUN sed -i -e "s/^bind-address\s*=\s*127.0.0.1/bind-address\t\t= 0.0.0.0/" /etc/mysql/my.cnf
RUN grep -v rootfs /proc/mounts > /etc/mtab

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

# Add custom setup script
[[ beforeAfterScripts ]]

# Add master files if needed
[[ if eq (.Container.GetCustomValue "repl") "master" ]]
	ADD init-master.sh /root/init-master.sh
    RUN chmod +x /root/init-master.sh
[[ end ]]

# Add slave files if needed
[[ if eq (.Container.GetCustomValue "repl") "slave" ]]
	ADD init-master.sh /root/init-slave.sh
    RUN chmod +x /root/init-slave.sh
[[ end ]]

CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] /bin/bash /root/setup.sh \
    && /etc/init.d/mysql start \
    [[ if eq (.Container.GetCustomValue "repl") "master" ]] && /bin/bash /root/init-master.sh \ [[ end ]]
    [[ if eq (.Container.GetCustomValue "repl") "slave" ]] && /bin/bash /root/init-slave.sh \ [[ end ]]
    [[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
    && /bin/bash
