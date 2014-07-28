FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

RUN export DEBIAN_FRONTEND=noninteractive
RUN dpkg --configure -a
RUN apt-get install -y -f debconf-utils
RUN apt-get -q -y install phpmyadmin
RUN echo 'phpmyadmin phpmyadmin/reconfigure-webserver multiselect apache2' | debconf-set-selections
RUN echo 'phpmyadmin phpmyadmin/dbconfig-install boolean false' | debconf-set-selections

RUN sed -i 's|DocumentRoot /var/www|DocumentRoot /usr/share/phpmyadmin|g' /etc/apache2/sites-enabled/000-default
RUN sed -i 's|VirtualHost \*:80|VirtualHost \*:[[ .Container.GetFirstLocalPort ]]|g' /etc/apache2/sites-enabled/000-default

RUN sed -i 's|NameVirtualHost \*:80|NameVirtualHost \*:[[ .Container.GetFirstLocalPort ]]|g' /etc/apache2/ports.conf
RUN sed -i 's|Listen 80|Listen [[ .Container.GetFirstLocalPort ]]|g' /etc/apache2/ports.conf

RUN cp /usr/share/phpmyadmin/config.sample.inc.php /usr/share/phpmyadmin/config.inc.php

[[ $firstLinked := .Container.FirstLinked]]

RUN sed -i "s|if (empty(\$dbserver)) \$dbserver = 'localhost';|if (empty(\$dbserver)) \$dbserver = getenv('[[ $firstLinked.Name | ToUpper ]]_PORT_[[ $firstLinked.GetFirstLocalPort]]_TCP_ADDR');|g" /etc/phpmyadmin/config.inc.php
RUN sed -i "s|//\$cfg\['Servers'\]\[\$i\]\['host'\] = 'localhost';|\$cfg\['Servers'\]\[\$i\]\['host'\] = getenv('[[ $firstLinked.Name | ToUpper ]]_PORT_[[ $firstLinked.GetFirstLocalPort]]_TCP_ADDR');|g" /etc/phpmyadmin/config.inc.php
RUN sed -i "s|// \$cfg\['Servers'\]\[\$i\]\['AllowNoPassword'\] = TRUE;|\$cfg\['Servers'\]\[\$i\]\['AllowNoPassword'\] = TRUE;|g" /etc/phpmyadmin/config.inc.php

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

[[ if .EmptyCmd ]]
CMD /bin/bash
[[ else ]]
CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] /bin/bash /root/setup.sh \
    && /etc/init.d/apache2 start \
    [[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
    && /bin/bash
[[ end ]]
