FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

[[ installPython ]]

RUN pip install whisper carbon graphite-web

# Configure graphite
RUN mv /opt/graphite/conf/carbon.conf.example /opt/graphite/conf/carbon.conf
RUN mv /opt/graphite/conf/storage-schemas.conf.example /opt/graphite/conf/storage-schemas.conf
RUN mv /opt/graphite/conf/aggregation-rules.conf.example /opt/graphite/conf/aggregation-rules.conf
RUN mv /opt/graphite/conf/dashboard.conf.example /opt/graphite/conf/dashboard.conf
RUN mv /opt/graphite/conf/graphTemplates.conf.example /opt/graphite/conf/graphTemplates.conf
RUN mv /opt/graphite/conf/graphite.wsgi.example /opt/graphite/conf/graphite.wsgi
RUN mv /opt/graphite/webapp/graphite/local_settings.py.example /opt/graphite/webapp/graphite/local_settings.py
RUN cd /opt/graphite/webapp/graphite && echo "yes\n\nroot@example.com\nroot\nroot" | python manage.py syncdb

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if .EmptyCmd]]
CMD /bin/bash
[[ else ]]
CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && \ [[end]]
    python /opt/graphite/webapp/graphite/manage.py runfcgi host=127.0.0.1 port=8080 && \
    [[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh && \[[end]]
    /bin/bash
[[ end ]]
