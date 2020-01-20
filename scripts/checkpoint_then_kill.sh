#!/bin/sh
pid=$( ps -C "dmtcp_coordinator" -o pid= | awk ' { print $1 } '  )
/.dmtcp/dmtcp/bin/dmtcp_command -c
while [ "$( /.dmtcp/dmtcp/bin/dmtcp_command -s | grep RUNNING |  awk ' { print $1 } ' )" = "RUNNING=no" ];
do
   sleep 1
done
/.dmtcp/dmtcp/bin/dmtcp_command -k
tail --pid=$pid -f /dev/null
