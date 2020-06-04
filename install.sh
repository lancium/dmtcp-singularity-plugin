#!/bin/sh -e

newplugin=false
scriptdir="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

# Build and set environment variable for DMTCP location.
function setupDMTCP {
	if [[ -z "${SINGULARITY_DMTCP}" ]]; then
		echo "export SINGULARITY_DMTCP=$scriptdir" | sudo tee -a /etc/profile
	fi
	cd $scriptdir/dmtcp/
	./configure
	make -j
}

# Build plugin for Singularity version <=3.5
function setupSingularityClassic {
	if [ -z "$1" ]; then
		echo "Singularity location must be provided as arg."
		exit
	fi
	echo "Absolute Singularity build path (for plugin build, command line): $1"
	mkdir -p $1/dmtcp-plugin
	cp -r $scriptdir/plugin/* $1/dmtcp-plugin/
	cd $1
	singularity plugin compile ./dmtcp-plugin/
	sudo singularity plugin install ./dmtcp-plugin/dmtcp-plugin.sif
}

# Build plugin for Singularity version > 3.5
function setupSingularityNew {
	rm -rf $scriptdir/plugin-build
	mkdir -p $scriptdir/plugin-build
	singularity plugin create $scriptdir/plugin-build dmtcp
	cp $scriptdir/3.6-plugin/main.go $scriptdir/plugin-build/
	singularity plugin compile $scriptdir/plugin-build/
	sudo singularity plugin install $scriptdir/plugin-build/plugin-build.sif
}

echo "Working directory (will be persistent): $scriptdir"
while getopts ":n" OPT; do
	case $OPT in
		n) echo "Selected build for Singularity >= 3.6"
		newplugin=true
		;;
	esac
done

setupDMTCP
if ! $newplugin ; then
	echo "Will be building in Singularity directory"
	setupSingularityClassic $1
else
	echo "Will be building in $scriptdir/plugin-build"
	setupSingularityNew
fi
