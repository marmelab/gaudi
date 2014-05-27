define DISTRIBUTION
Origin: gaudi.io
Label: apt repository
Codename: precise
Version: $(GAUDI_VERSION)
Architectures: amd64 i386 source
Components: main non-free contrib
Description: gaudi debian package repository
SignWith: yes
Pull: precise
endef

GAUDI_VERSION = $(shell go run main.go -v)

export DISTRIBUTION

apt:
	goxc -tasks-=go-test -pv=$(GAUDI_VERSION)

	cd templates/ && tar -cvf $(GOPATH)/bin/gaudi-xc/$(GAUDI_VERSION)/templates.tar . && cd -

	mkdir -p $(GOPATH)/bin/gaudi-xc/$(GAUDI_VERSION)/conf/
	echo "$$DISTRIBUTION" > $(GOPATH)/bin/gaudi-xc/$(GAUDI_VERSION)/conf/distributions

	cp /usr/gaudi/gaudi.gpg.key $(GOPATH)/bin/gaudi-xc/$(GAUDI_VERSION)/gaudi.gpg.key

	cd $(GOPATH)/bin/gaudi-xc/$(GAUDI_VERSION) && \
        dpkg-sig --sign builder gaudi_$(GAUDI_VERSION)_amd64.deb && \
        dpkg-sig --sign builder gaudi_$(GAUDI_VERSION)_i386.deb && \
        reprepro --ask-passphrase remove precise gaudi && \
        reprepro --ask-passphrase -Vb . -S main includedeb precise gaudi_$(GAUDI_VERSION)_amd64.deb && \
        reprepro --ask-passphrase -Vb . -S main includedeb precise gaudi_$(GAUDI_VERSION)_i386.deb
