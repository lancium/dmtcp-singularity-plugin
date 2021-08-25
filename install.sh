#!/bin/sh -e

#newplugin=false
scriptdir="$( cd "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

# Build and set environment variable for DMTCP location.
function setupDMTCP {
	#if [[ -z "${SINGULARITY_DMTCP}" ]]; then
	#	echo "export SINGULARITY_DMTCP=$scriptdir" | sudo tee -a /etc/profile
	#fi
	cd $scriptdir/dmtcp/
	./configure --enable-static-libstdcxx
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
	cp $scriptdir/plugin/main.go $scriptdir/plugin-build/
	singularity plugin compile $scriptdir/plugin-build/
	sudo singularity plugin install $scriptdir/plugin-build/plugin-build.sif
}

echo "Working directory (will be persistent): $scriptdir"
#while getopts ":n" OPT; do
#	case $OPT in
#		n) echo "Selected build for Singularity >= 3.6"
#		newplugin=true
#		;;
#	esac
#done

setupDMTCP
echo "Will be building in $scriptdir/plugin-build"
setupSingularityNew
echo "The following environment variables may need to be set if using the DMTCP version just built:"
echo "	export SINGULARITY_DMTCP_BIN=$scriptdir"
echo "	export SINGULARITY_DMTCP_LIB=$scriptdir/dmtcp/lib"
echo "You may also invoke Singularity like the following to start a container with checkpoint capabilities:"
echo "SINGULARITY_DMTCP_BIN=$scriptdir SINGULARITY_DMTCP_LIB=$scriptdir/dmtcp/lib/dmtcp/ singularity checkpoint start <container> <instance name>"