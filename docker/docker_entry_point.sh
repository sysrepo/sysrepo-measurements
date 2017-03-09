#!/bin/bash

#add root password
echo "root:root" | chpasswd

# remove printf from the sysrepo callback in the application_example
cd /opt/dev/sysrepo
sed -i '63d' ./examples/application_example.c
sed -i '61d' ./examples/application_example.c
cd build
make

# run the sysrepo application for the yang model ietf-interfaces
/opt/dev/sysrepo/build/examples/application_example "ietf-interfaces" &
