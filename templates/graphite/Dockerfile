FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

RUN apt-get install -y -f libcairo2 libcairo2-dev python-cairo libbz2-dev libsqlite3-dev nginx-light

[[ .Container.SetCustomValue "pythonVersion" "2.7.6"]]
[[ installPython ]]

[[ .Container.SetCustomValue "djangoVersion" "1.5"]]
[[ installDjango ]]

RUN pip install whisper carbon graphite-web django-tagging pysqlite flup daemonize gunicorn 'Twisted<12.0'

# Configure graphite
RUN mv /opt/graphite/conf/carbon.conf.example /opt/graphite/conf/carbon.conf
RUN mv /opt/graphite/conf/storage-schemas.conf.example /opt/graphite/conf/storage-schemas.conf
RUN mv /opt/graphite/conf/aggregation-rules.conf.example /opt/graphite/conf/aggregation-rules.conf
RUN mv /opt/graphite/conf/dashboard.conf.example /opt/graphite/conf/dashboard.conf
RUN mv /opt/graphite/conf/graphTemplates.conf.example /opt/graphite/conf/graphTemplates.conf
RUN mv /opt/graphite/conf/graphite.wsgi.example /opt/graphite/conf/graphite.wsgi
RUN mv /opt/graphite/webapp/graphite/local_settings.py.example /opt/graphite/webapp/graphite/local_settings.py

# Create locations for pid/log files
RUN mkdir -p /var/run/graphite && chown www-data /var/run/graphite
RUN mkdir -p /var/log/carbon && chown www-data /var/log/carbon

# Initialize the webapp
RUN cd /opt/graphite/webapp/graphite && python manage.py syncdb --noinput 

# Add custom setup script
[[ beforeAfterScripts ]]

ADD ./nginx.conf /etc/nginx/nginx.conf

[[ if .EmptyCmd]]
CMD /bin/bash
[[ else ]]
CMD /opt/graphite/bin/carbon-cache.py --debug start & \
	gunicorn_django -b127.0.0.1:8000 -w2 /opt/graphite/webapp/graphite/settings.py & \
	/usr/sbin/nginx & \
	/bin/bash
[[ end ]]
