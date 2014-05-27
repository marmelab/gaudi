FROM stackbrew/debian:wheezy

[[ updateApt ]]
[[ addUserFiles ]]

# Install jekyll
RUN apt-get -y --force-yes install ruby1.9.1-dev
RUN gem install jekyll execjs therubyracer

ENTRYPOINT ["jekyll"]
