FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

WORKDIR [[ .Container.GetFirstMountedDir ]]

[[ $version := (.Container.GetCustomValue "version" "1.2") ]]

# Install go
RUN apt-get -y -f install make git curl mercurial bison bzr
RUN wget https://go.googlecode.com/files/go[[ $version ]].linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go[[ $version ]].linux-amd64.tar.gz && \
    rm go[[ $version ]].linux-amd64.tar.gz

# Set GOPATH and GOROOT environment variables
ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

# Install deps
[[range (.Container.GetCustomValue "modules")]]
    RUN go get [[.]]
[[end]]

[[ if .EmptyCmd ]]
CMD /bin/bash
[[ else ]]
    [[ if (.Container.HasAfterScript) ]]
        CMD [[.Container.AfterScript]] && /bin/bash
    [[ else ]]
        CMD /bin/bash
    [[ end]]
[[ end ]]
