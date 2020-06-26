# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  # use an iamge of debian, in this case buster 64
  config.vm.box = "debian/buster64"

  # forward ports 80 and 443 to the host system
  # this will allow accessing the webserver from the real system. 
  config.vm.network "forwarded_port", guest: 80, host: 80
  config.vm.network "forwarded_port", guest: 443, host: 443

  # share the factory folder in /factory/
  config.vm.synced_folder "distillery/", "/distillery/"

  # for performance, we setup 4GB of memort and 2 CPUs
  config.vm.provider "virtualbox" do |vb|
    vb.memory = 4096
    vb.cpus = 2
  end

  # tell the user where things are. 
  config.vm.post_up_message = "Ready to distil and make WissKIs. Scripts can be found in /distillery/. "
end
