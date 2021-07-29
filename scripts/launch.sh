#!/bin/sh
# Quick script that launches DMTCP
# Any options desired for DMTCP should be set here.

#echo "Launching DMTCP..."
/.dmtcp/dmtcp/bin/dmtcp_launch --ckptdir /.checkpoint --no-gzip --ckpt-open-files $@
