#!/bin/sh

# Open Itunes Command
function cmd_open(){
	echo "Starting iTunes.";
	open -a iTunes;
}

# Play song in iTunes
function cmd_play(){
	echo "Playing iTunes.";
	osascript -e 'tell application "iTunes" to play';
}

# Stop song in iTunes
function cmd_stop(){
	echo "Stopping2 iTunes.";
	osascript -e 'tell application "iTunes" to stop';
}