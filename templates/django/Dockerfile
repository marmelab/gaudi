FROM stackbrew/debian:wheezy

RUN echo "deb http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list
RUN echo "deb-src http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

RUN apt-get install -y -f python2.7 python-dev python-setuptools libmysqlclient-dev
RUN easy_install pip
RUN pip install django==1.6
RUN pip install mysql-python
RUN pip install south

[[range (.Container.GetCustomValue "pip_modules")]]
    RUN pip install [[.]]
[[end]]

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
    [[ if eq (.Collection.IsComponentDependingOf .Container "apache") false ]]
    && python manage.py runserver 0.0.0.0:[[ .Container.GetFirstLocalPort "8000" ]] \
    [[end]]
    && /bin/bash
[[ end ]]
