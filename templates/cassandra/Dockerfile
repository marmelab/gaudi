FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

RUN echo "deb http://www.apache.org/dist/cassandra/debian 11x main" > /etc/apt/sources.list.d/cassandra.list
RUN echo "deb-src http://www.apache.org/dist/cassandra/debian 11x main" >> /etc/apt/sources.list.d/cassandra.list

RUN apt-get update
RUN apt-get install -y -f --force-yes cassandra

RUN sed -i -e "s/listen_address:\slocalhost/listen_address: 0.0.0.0/" /etc/cassandra/cassandra.yaml
RUN sed -i -e "s/rpc_address:\slocalhost/rpc_address: 0.0.0.0/" /etc/cassandra/cassandra.yaml

RUN sed -i -e "s/#MAX_HEAP_SIZE=\"4G\"/MAX_HEAP_SIZE=\"[[ (.Container.GetCustomValue "maxHeapSize" "512M") ]]\"/" /etc/cassandra/cassandra-env.sh
RUN sed -i -e "s/#HEAP_NEWSIZE=\"800M\"/HEAP_NEWSIZE=\"[[ (.Container.GetCustomValue "heapNewSize" "128M") ]]\"/" /etc/cassandra/cassandra-env.sh

[[ if .EmptyCmd ]]
CMD /bin/bash
[[ else ]]
CMD /etc/init.d/cassandra start && /bin/bash
[[ end ]]
