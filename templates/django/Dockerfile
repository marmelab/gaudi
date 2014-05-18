FROM stackbrew/debian:wheezy

RUN echo "deb http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list
RUN echo "deb-src http://ftp.fr.debian.org/debian/ wheezy non-free" >> /etc/apt/sources.list

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR /app

RUN apt-get install -y -f python3 python-dev python-setuptools libmysqlclient-dev
RUN easy_install pip
RUN pip install virtualenv
RUN pip install mysql-python
RUN pip install south

# Install django & configure it
RUN pip install django==1.3
RUN django-admin.py startproject project
RUN python ./project/manage.py startapp myapp

RUN sed -i "1s/^/import os\n/" ./project/settings.py
RUN sed -i -e "s/'ENGINE': 'django.db.backends.'/'ENGINE': 'django.db.backends.mysql'/" ./project/settings.py
RUN sed -i -e "s/'NAME': ''/'NAME': 'django'/" ./project/settings.py
RUN sed -i -e "s/'USER': ''/'User': 'root'/" ./project/settings.py
RUN sed -i -e "s/'HOST': ''/'HOST': os.environ['DB_PORT_3306_TCP_ADDR']/" ./project/settings.py
RUN sed -i -e "s/# 'django.contrib.admin'/'django.contrib.admin'/" ./project/settings.py

RUN sed -i -e "s/# from django.contrib import admin/from django.contrib import admin/" ./project/urls.py
RUN sed -i -e "s/# admin.autodiscover()/admin.autodiscover()/" ./project/urls.py
RUN sed -i -e "s/# url(r'^admin\/', include(admin.site.urls))/url(r'^admin\/', include(admin.site.urls))/" ./project/urls.py

# Add custom setup script
[[ beforeAfterScripts ]]

CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] python ./project/manage.py runserver 0.0.0.0:8000 \
	[[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
	&& /bin/bash
