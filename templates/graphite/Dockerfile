# @see hopsoft/graphite-statsd
FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

RUN apt-get -y --force-yes install vim  nginx  python-flup expect git memcached sqlite3 libcairo2 libcairo2-dev python-cairo pkg-config build-essential python-dev libsqlite3-dev

RUN wget -P /opt http://python-distribute.org/distribute_setup.py
RUN python /opt/distribute_setup.py
RUN easy_install pip

RUN pip install django==1.3 whisper==0.9.12 carbon graphite-web django-tagging==0.3.1 pysqlite flup daemonize gunicorn twisted==11.1.0 python-memcached==1.53 txAMQP==0.6.2

# Configure graphite
RUN mv /opt/graphite/conf/carbon.conf.example /opt/graphite/conf/carbon.conf
RUN mv /opt/graphite/conf/storage-schemas.conf.example /opt/graphite/conf/storage-schemas.conf
RUN mv /opt/graphite/conf/aggregation-rules.conf.example /opt/graphite/conf/aggregation-rules.conf
RUN mv /opt/graphite/conf/dashboard.conf.example /opt/graphite/conf/dashboard.conf
RUN mv /opt/graphite/conf/graphTemplates.conf.example /opt/graphite/conf/graphTemplates.conf
RUN mv /opt/graphite/conf/graphite.wsgi.example /opt/graphite/conf/graphite.wsgi
RUN mv /opt/graphite/webapp/graphite/local_settings.py.example /opt/graphite/webapp/graphite/local_settings.py
RUN echo "\n\n[stats]\npattern = ^stats.*\nretentions = 10s:6h,1min:6d,10min:1800d" >> /opt/graphite/conf/storage-schemas.conf

# Create locations for pid/log files
RUN mkdir -p /var/run/graphite && chown www-data /var/run/graphite
RUN mkdir -p /var/log/carbon && chown www-data /var/log/carbon

# Initialize the webapp
sed -i -e "s|#SECRET_KEY = 'UNSAFE_DEFAULT'|SECRET_KEY = 'OJNOKdsqds!d987§8'|" /opt/graphite/webapp/graphite/local_settings.py
ADD ./graphite_syncdb /tmp/graphite_syncdb
RUN chmod 775 /tmp/graphite_syncdb
RUN /tmp/graphite_syncdb

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
