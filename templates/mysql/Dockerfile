FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

# Install mysql
RUN apt-get -y --force-yes install mysql-server

# Edit mysql config (use root as user to be able to write in docker mounted volumes)
RUN sed -i -e "s/^user\s*=\s*mysql/user = root/" /etc/mysql/my.cnf
RUN sed -i -e "s/^bind-address\s*=\s*127.0.0.1/bind-address\t\t= 0.0.0.0/" /etc/mysql/my.cnf

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if .EmptyCmd ]]
CMD /bin/bash
[[ else ]]
CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] /bin/bash /root/setup.sh \
    && /etc/init.d/mysql start \
    [[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
    && /bin/bash
[[ end ]]
