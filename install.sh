#!/bin/sh -e
if [ -z "$1" ]
  then
    echo "Singularity location must be provided as arg."
	exit
fi
echo "Current directory (target in git repo): $( pwd )"
echo "Absolute Singularity build path (for plugin build, command line): $1"
mkdir -p $1/plugins
ln -s $( pwd )/plugin $1/plugins/dmtcp-singularity-plugin
if [[ -z "${SINGULARITY_DMTCP}" ]]; then
    echo "export SINGULARITY_DMTCP=$( pwd )" | tee -a /etc/profile
fi
cd ./dmtcp/
./configure
make -j
cd $1
singularity plugin compile ./plugins/dmtcp-singularity-plugin/
sudo singularity plugin install ./plugins/dmtcp-singularity-plugin/dmtcp-singularity-plugin.sif
