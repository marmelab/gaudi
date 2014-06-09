# Inspired from svendowideit/ambassador
FROM stackbrew/debian:wheezy

[[ updateApt ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

RUN apt-get install -y -f socat

CMD env | grep _TCP= | sed 's/.*_PORT_\([0-9]*\)_TCP=tcp:\/\/\(.*\):\(.*\)/socat TCP4-LISTEN:\1,fork,reuseaddr TCP4:\2:\3 \&/' | sh \
    && /bin/bash
