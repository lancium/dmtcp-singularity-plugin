#!/bin/sh
echo "Launching DMTCP..."
#/.dmtcp/dmtcp/bin/dmtcp_launch --join-coordinator --ckptdir ./lancium-checkpoint/ --ckpt-open-files $@
/.dmtcp/dmtcp/bin/dmtcp_launch --ckptdir /.checkpoint/ --ckpt-open-files $@
