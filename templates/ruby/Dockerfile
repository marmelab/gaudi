FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

[[ $version := (.Container.GetCustomValue "version" "2.1.2") ]]

# Install ruby via RVM
RUN apt-get install -y --reinstall procps
RUN \curl -L https://get.rvm.io | bash -s stable
ENV PATH $PATH:/usr/local/rvm/bin/
RUN rvm install [[ $version ]] && rvm use [[ $version ]]
ENV PATH $PATH:/usr/local/rvm/rubies/ruby-[[ $version ]]/bin/

# Install gems
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
