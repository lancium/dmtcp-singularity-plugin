#!/bin/sh
latest=$( ls -t1 /.dmtcp/checkpoint/ | grep .sh |  head -n 1 )
"/.dmtcp/checkpoint/$latest"