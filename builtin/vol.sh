#!/bin/sh
#
# "itunesScript vol" builtin command.
function vol(){
	currentVolume=`osascript -e 'tell application "iTunes" to sound volume as integer'`;
	if [ $1 = "up" ]; then
		newVolume=$(( currentVolume+10 ));
	else
		newVolume=$(( currentVolume-10 ));
	fi
	osascript -e "tell application \"iTunes\" to set sound volume to $newVolume";
}