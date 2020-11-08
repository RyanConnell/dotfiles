#!/usr/bin/env bash

killall -q polybar
polybar top >> /tmp/polybar_top.log 2>&1 &
polybar bottom >> /tmp/polybar_top.log 2>&1 &
