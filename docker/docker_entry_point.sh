#!/bin/bash

#add root password
echo "root:root" | chpasswd

/opt/dev/sysrepo/build/examples/application_example "ietf-interfaces" &
