#!/bin/bash
# if the second argument is "battery.charge" return 100
# otherwise return "OL"
if [ "$2" = "battery.charge" ]; then
    echo 100
else
    echo "OL"
fi