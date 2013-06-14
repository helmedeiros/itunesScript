#!/bin/sh
#
# "itunesScript vol" builtin command.
function vol(){
	volume;
	my_volume=$?;
	
	if [ $1 = "up" ] || [ $1 = "down" ]; then
		if [ $1 = "up" ]; then
			new_volume=$(( my_volume+10 ));
		fi
		
		if [ $1 = "down" ]; then
			new_volume=$(( my_volume-10 ));
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

function volume(){
	my_volume=`osascript -e 'tell application "iTunes" to sound volume as integer'`;
	return "$my_volume"
}

function increase_or_decrease(){
	echo "iTunes volume: $new_volume";
	
	if [ $old_volume -lt $new_volume ]; then
		echo "Increasing iTunes volume.";
	else
		echo "Decrease iTunes volume.";
	fi
}