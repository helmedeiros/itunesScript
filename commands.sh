#!/bin/sh
#
# itunesScript commands gateway.
source $(dirname $0)/builtin/play.sh
source $(dirname $0)/builtin/pause.sh
source $(dirname $0)/builtin/next.sh
source $(dirname $0)/builtin/previous.sh
source $(dirname $0)/builtin/mute.sh
source $(dirname $0)/builtin/unmute.sh
source $(dirname $0)/builtin/vol.sh
source $(dirname $0)/builtin/stop.sh
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

# Next song in iTunes
function cmd_next(){
	echo "Changing to the next song.";
	next;
}

# Previous song in iTunes
function cmd_prev(){
	echo "Changing back to the previous song.";
	prev;
}

function cmd_mute(){
	echo "Putting iTunes on mute.";
	mutting;
}

function cmd_unmute(){
	echo "Taking iTunes from mute.";
	unmutting;
}

function cmd_vol(){
	which_volume; old_volume=$?;
	
	vol $1;	
	
	which_volume; new_volume=$?;
	echo "iTunes volume: $new_volume";
	
	increase_or_decrease $old_volume $new_volume;	
}

# Stop song in iTunes
function cmd_stop(){
	echo "Stopping iTunes.";
	stopping;
}

# Quit iTunes
function cmd_quit(){
	echo "Quiting iTunes.";
	quiting;
}