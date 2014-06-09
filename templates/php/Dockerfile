FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

# Install PHP 5.4
RUN apt-get -y -f install php5 php5-mysql php5-mcrypt php5-curl curl

# Add custom setup script
[[ beforeAfterScripts ]]

CMD [[ if (.Container.HasAfterScript) ]] /bin/bash /root/after-setup.sh && [[end]] /bin/bash
