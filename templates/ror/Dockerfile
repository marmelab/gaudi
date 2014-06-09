FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

[[ $version := (.Container.GetCustomValue "version" "2.1.2") ]]
[[ $serverType := (.Container.GetCustomValue "serverType" "standalone") ]]

[[ installRvm ]]

# Install nodejs
[[ .Container.SetCustomValue "nodeVersion" "0.10.20"]]
[[ installNodeJS ]]

# Install custom gems
[[range (.Container.GetCustomValue "gems")]]
    RUN gem install [[.]]
[[end]]

[[ if eq $serverType "apache"]]
    RUN apt-get install -y -f apache2
    RUN a2enmod actions alias
    RUN service apache2 reload

    RUN apt-get install -y -f libcurl4-openssl-dev apache2-threaded-dev libapr1-dev libaprutil1-dev
        RUN /bin/bash -l -c 'gem install passenger bundler execjs'
        RUN /bin/bash -l -c 'passenger-install-apache2-module --auto'

        RUN print "LoadModule passenger_module /usr/local/rvm/gems/ruby-[[ $version ]]/gems/passenger-4.0.44/buildout/apache2/mod_passenger.so \n<IfModule mod_passenger.c>\nPassengerRoot /usr/local/rvm/gems/ruby-[[ $version ]]/gems/passenger-4.0.44 \nPassengerDefaultRuby /usr/local/rvm/gems/ruby-[[ $version ]]/wrappers/ruby \n</IfModule>" > /etc/apache2/mods-available/passenger.load
        RUN a2enmod passenger

        ADD 000-default /etc/apache2/sites-enabled/000-default
        ADD ports.conf /etc/apache2/ports.conf
[[ end ]]

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if .EmptyCmd ]]
CMD /bin/bash
[[ else ]]
CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] /bin/bash /root/setup.sh \
    [[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
    && /etc/init.d/apache2 start \
    && /bin/bash
[[ end ]]
