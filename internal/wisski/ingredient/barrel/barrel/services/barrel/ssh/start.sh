#!/bin/bash

# create the sshd directory
if [ ! -d /run/sshd ]; then
   mkdir /run/sshd
   chmod 0755 /run/sshd
fi

# regenerate key files if they do not yet exist
[[ -f "/ssh/hostkeys/ssh_host_rsa_key" ]] || ssh-keygen -q -N "" -t dsa -f /ssh/hostkeys/ssh_host_rsa_key
[[ -f "/ssh/hostkeys/ssh_host_ecdsa_key" ]] || ssh-keygen -q -N "" -t ecdsa -f /ssh/hostkeys/ssh_host_ecdsa_key
[[ -f "/ssh/hostkeys/ssh_host_ed25519_key" ]] || ssh-keygen -q -N "" -t ed25519 -f /ssh/hostkeys/ssh_host_ed25519_key

/usr/sbin/sshd -e -D -f /ssh/sshd_config