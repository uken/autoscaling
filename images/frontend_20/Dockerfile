FROM uken/trusty
MAINTAINER Uken Ops "ops@uken.com"
ENV PATH /usr/local/rvm/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

RUN apt-get -y update && apt-get -y install nodejs
RUN gpg --keyserver hkp://keys.gnupg.net --ignore-time-conflict --recv-keys D39DC0E3
RUN curl -L https://get.rvm.io | bash -s stable
RUN echo 'source /usr/local/rvm/scripts/rvm' >> /etc/bash.bashrc
RUN /usr/local/rvm/bin/rvm-shell -l -c 'rvm requirements'
RUN /usr/local/rvm/bin/rvm-shell -l -c 'rvm install 2.0.0-p598'
RUN /usr/local/rvm/bin/rvm-shell -l -c 'rvm use 2.0.0-p598 --default'
RUN /usr/local/rvm/bin/rvm-shell -l -c 'gem install bundler --no-ri --no-rdoc'

RUN mkdir /root/.ssh
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config

ONBUILD ADD tmp/deploy_key /root/.ssh/id_rsa
ONBUILD RUN chmod 600 /root/.ssh/id_rsa

ONBUILD ADD ./Gemfile /app/Gemfile
ONBUILD ADD ./Gemfile.lock /app/Gemfile.lock
ONBUILD RUN /usr/local/rvm/bin/rvm-shell -l -c 'cd /app && bundle install --jobs 4 --deployment --without development test debug'
ONBUILD ADD ./ /app
