#!/bin/sh
pid = 0 #get pid of dmtcp coordinator
/.dmtcp/dmtcp/bin/dmtcp_command -c
tail --pid=$pid -f /dev/null