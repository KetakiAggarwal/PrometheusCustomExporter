#!/bin/sh
sleep 10
df -ihP | sed -n '2,$p' | awk '{print "em_result="$1"|"$6"||"$5}' | sed 's/.$//'
