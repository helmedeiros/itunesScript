#!/bin/sh
source $(dirname $0)/builtin/play.sh
source $(dirname $0)/builtin/pause.sh
source $(dirname $0)/builtin/next.sh
source $(dirname $0)/builtin/quit.sh

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

# Pause song in iTunes
function cmd_pause(){
	echo "Pausing iTunes.";
	pausing;
}

function cmd_next(){
	echo "Changing to the next song";
	next;
}

# Quit iTunes
function cmd_quit(){
	echo "Quiting iTunes.";
	quiting;
}

# Stop song in iTunes
function cmd_stop(){
	echo "Stopping2 iTunes.";
	osascript -e 'tell application "iTunes" to stop';
}