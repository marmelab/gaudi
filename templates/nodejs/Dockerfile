FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

[[ installNodeJS ]]

# Install modules
[[range (.Container.GetCustomValue "modules")]]
    RUN npm install -g [[.]]
[[end]]

ENV NODE_PATH /usr/local/lib/node_modules

[[ if .EmptyCmd ]]
CMD /bin/bash
[[ else ]]
    [[ if (.Container.HasAfterScript) ]]
        CMD [[.Container.AfterScript]] && /bin/bash
    [[ else ]]
        CMD ["/usr/local/bin/node"] && /bin/bash
    [[ end]]
[[ end ]]
