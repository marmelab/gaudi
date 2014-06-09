FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

[[ installPython ]]

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if (.Container.HasAfterScript) ]]
    CMD [[.Container.AfterScript]] && /bin/bash
[[ else ]]
    CMD /bin/bash
[[ end]]
