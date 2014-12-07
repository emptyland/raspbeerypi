#!/bin/bash

LEDS="25 26 27 28"

CMD="/usr/local/bin/gpio"

for led in $LEDS; do
    $CMD mode $led out
done

for ((i=0; i<7; i++)); do
    for led in $LEDS; do
        $CMD write $led 1
        sleep 0.2
        $CMD write $led 0
    done
done

