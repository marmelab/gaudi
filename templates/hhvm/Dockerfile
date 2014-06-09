FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

# Update apt to install HHVM
RUN echo deb http://dl.hhvm.com/debian wheezy main | tee /etc/apt/sources.list.d/hhvm.list

# Install HHVM
RUN apt-get -y update
RUN apt-get install -y --force-yes -f hhvm php5-cli curl

[[ $memoryLimit := .Container.GetCustomValue "memoryLimit" "128M" ]]
[[ $maxExecutionTime := .Container.GetCustomValue "maxExecutionTime" "30" ]]
RUN sed -i -e 's|; php options|; php options\nmemory_limit = [[ $memoryLimit ]]\nmax_execution_time = [[ $maxExecutionTime ]]\ndisplay_startup_errors = On\nerror_reporting = E_ALL\ndisplay_errors = On|' /etc/hhvm/php.ini

[[ if .EmptyCmd ]]
CMD /bin/bash
[[ else ]]
CMD (hhvm --mode daemon -vServer.Type=fastcgi -vServer.Port=[[ .Container.GetFirstLocalPort ]] &) && /bin/bash
[[ end ]]
