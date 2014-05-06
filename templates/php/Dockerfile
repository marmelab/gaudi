FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

# Install PHP 5.4
RUN apt-get -y -f install php5 php5-mysql php5-mcrypt php5-curl curl

[[ if (.Container.HasAfterScript) ]]
	CMD [[.Container.AfterScript]] && /bin/bash
[[ else ]]
	CMD /bin/bash
[[ end]]

