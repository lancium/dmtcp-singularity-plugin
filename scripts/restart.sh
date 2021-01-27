#!/bin/sh
latest=$( ls -t1 /.checkpoint/ | grep .sh |  head -n 1 )
"/.checkpoint/$latest"
