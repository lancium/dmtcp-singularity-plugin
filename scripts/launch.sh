#!/bin/sh
echo "Launching DMTCP..."
mkdir -p ./lancium-checkpoint/
/.dmtcp/dmtcp/bin/dmtcp_launch --ckptdir ./lancium-checkpoint/ --ckpt-open-files $@