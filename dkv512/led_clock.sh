#!/bin/bash

LEDS="25 26 27 28"

for led in $LEDS; do
    gpio mode $led out
done

for ((i=0; i<7; i++)); do
    for led in $LEDS; do
        gpio write $led 1
        sleep 0.2
        gpio write $led 0
    done
done

