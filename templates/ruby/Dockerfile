FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

[[ installRvm ]]

# Install custom gems
[[range (.Container.GetCustomValue "gems")]]
    RUN gem install [[.]]
[[end]]

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if (.Container.HasAfterScript) ]]
    CMD [[.Container.AfterScript]] && /bin/bash
[[ else ]]
    CMD /bin/bash
[[ end]]
