#!/bin/sh
# This is for the event multiple restart scripts are in the checkpoint directory.
# Lists the files in reverse order based on modification time, picks the top one.
# Then, just execute the script.
# Called during "singularity checkpoint restart instance"
latest=$( ls -t1 /.checkpoint/dmtcp_restart_script_*.sh |  head -n 1 )
"$latest"
