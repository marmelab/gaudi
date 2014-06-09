FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

# Install Java 7
RUN aptitude install -y -f openjdk-7-jre

RUN mkdir -p /opt/jackrabbit/

# Install jackrabbit
RUN wget http://archive.apache.org/dist/jackrabbit/2.6.5/jackrabbit-standalone-2.6.5.jar -P /opt/jackrabbit/

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if .EmptyCmd ]]
CMD /bin/bash
[[ else ]]
CMD [[ if (.Container.HasBeforeScript) ]] /bin/bash /root/before-setup.sh && [[end]] cd /opt/jackrabbit/ \
    && (java -jar jackrabbit-standalone-2.6.5.jar --port [[ .Container.GetFirstLocalPort ]] &) \
    [[ if (.Container.HasAfterScript) ]] && /bin/bash /root/after-setup.sh \[[end]]
    && /bin/bash
[[end]]
