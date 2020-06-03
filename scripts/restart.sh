#!/bin/sh
latest=$( ls -t1 ./lancium-checkpoint/ | grep .sh |  head -n 1 )
"./lancium-checkpoint/$latest"