FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

[[ $version := (.Container.GetCustomValue "version" "3.3.5") ]]

# Install python via pyenv
RUN apt-get install -y -f git
RUN git clone git://github.com/yyuu/pyenv.git /.pyenv
ENV PYENV_ROOT $HOME/.pyenv
ENV PATH $PYENV_ROOT/shims:$PYENV_ROOT/bin:$PATH

RUN pyenv install [[ $version ]]
ENV PATH /.pyenv/versions/[[ $version ]]/bin:$PATH

# Install pip modules
[[range (.Container.GetCustomValue "pip_modules")]]
	RUN pip install [[.]]
[[end]]

# Add custom setup script
[[ beforeAfterScripts ]]

[[ if (.Container.HasAfterScript) ]]
	CMD [[.Container.AfterScript]] && /bin/bash
[[ else ]]
	CMD /bin/bash
[[ end]]
