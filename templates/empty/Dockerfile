# Inspired from svendowideit/ambassador
FROM stackbrew/debian:wheezy

[[ updateApt ]]

[[ if (.Container.HasAfterScript) ]]
    CMD [[.Container.AfterScript]] && /bin/bash
[[ else ]]
    CMD /bin/bash
[[ end]]
