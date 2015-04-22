# -*- mode: ruby -*-
# vi: set ft=ruby :
ROOT = File.dirname(File.expand_path(__FILE__))

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "precise64"
  config.ssh.forward_agent = true

  if Dir.glob("#{File.dirname(__FILE__)}/.vagrant/machines/default/*/id").empty?
    # Add lxc-docker package
    pkg_cmd = "wget -q -O - https://get.docker.io/gpg | apt-key add -;" \
      "echo deb http://get.docker.io/ubuntu docker main > /etc/apt/sources.list.d/docker.list;" \
      "apt-get update -qq; apt-get install -q -y --force-yes lxc-docker; "
    # Add vagrant user to the docker group
    pkg_cmd << "usermod -a -G docker vagrant; "

    # Add consul
    pkg_cmd << "apt-get install -q -y --force-yes unzip curl build-essential;" \
    "cd /srv && wget -q https://dl.bintray.com/mitchellh/consul/0.5.0_linux_amd64.zip;" \
    "cd /srv && unzip 0.5.0_linux_amd64.zip; "

    # Add consul template
    pkg_cmd << "cd /srv && wget -q https://github.com/hashicorp/consul-template/releases/download/v0.7.0/consul-template_0.7.0_linux_amd64.tar.gz;" \
      "cd /srv && tar xvf consul-template_0.7.0_linux_amd64.tar.gz; "

    config.vm.provision :shell, :inline => pkg_cmd
  end
end
