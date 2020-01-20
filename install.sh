#!/bin/sh
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi
if [ -z "$1" ]
  then
    echo "Singularity location must be provided as arg."
	exit
fi
echo "Current directory (target in git repo): $( pwd )"
echo "Absolute Singularity build path (for plugin build, command line): $1"
ln -s $( pwd )/plugin $1/plugins/dmtcp-singularity-plugin
echo "export SINGULARITY_DMTCP=$( pwd )" >> /etc/profile
cd ./dmtcp/
./configure
make -j
cd $1
singularity plugin compile ./plugins/dmtcp-singularity-plugin/
singularity plugin install ./plugins/dmtcp-singularity-plugin/dmtcp-singularity-plugin.sif
