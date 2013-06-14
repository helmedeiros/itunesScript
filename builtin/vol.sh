#!/bin/sh
#
# "itunesScript vol" builtin command.
function vol(){
	current_volume=`osascript -e 'tell application "iTunes" to sound volume as integer'`;
	if [ $1 = "up" ] || [ $1 = "down" ]; then
		if [ $1 = "up" ]; then
			new_volume=$(( current_volume+10 ));
		fi
		
		if [ $1 = "down" ]; then
			new_volume=$(( current_volume-10 ));
		fi
	else
		if [ $1 -gt 0 ]; then
			new_volume=$1;
		fi
	fi
	change_volume $new_volume;
}

function change_volume(){
	new_volume=$1;
	osascript -e "tell application \"iTunes\" to set sound volume to $new_volume";
}