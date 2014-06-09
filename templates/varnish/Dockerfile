FROM stackbrew/debian:wheezy

RUN echo "deb http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list
RUN echo "deb-src http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

RUN curl http://repo.varnish-cache.org/debian/GPG-key.txt | apt-key add -
RUN echo "deb http://repo.varnish-cache.org/debian/ wheezy varnish-3.0" >> /etc/apt/sources.list
RUN apt-get update
RUN apt-get -y install varnish

ADD varnish.conf /etc/default/varnish
ADD default.vcl /etc/varnish/default.vcl

# Add setup script
ADD setup.sh /root/setup.sh
RUN chmod +x /root/setup.sh

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if .EmptyCmd ]]
CMD /bin/bash
[[ else ]]
CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] /bin/bash /root/setup.sh \
    && varnishd -f /etc/varnish/default.vcl -s malloc,100M -a 0.0.0.0:[[ .Container.GetFirstLocalPort ]] \
    [[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
    && /bin/bash
[[ end ]]
