#!/bin/sh
source $(dirname $0)/builtin/play.sh

# Open Itunes Command
function cmd_open(){
	echo "Starting iTunes.";
	open -a iTunes;
}

# Play song in iTunes
function cmd_play(){
	echo "Playing iTunes.";
	playing $1;
}

# Stop song in iTunes
function cmd_stop(){
	echo "Stopping2 iTunes.";
	osascript -e 'tell application "iTunes" to stop';
}