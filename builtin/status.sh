#!/bin/sh
#
# "itunesScript status" builtin command.
function status(){
	if [ $1 != "stopped" ]; then
		artist=`osascript -e 'tell application "iTunes" to artist of current track as string'`;
		track=`osascript -e 'tell application "iTunes" to name of current track as string'`;
		album=`osascript -e 'tell application "itunes" to album of current track as string'`;
		status_message="The current track is $track at $album by $artist";
		echo $status_message;
	
		read_status $1 "$status_message";
	fi
}

function read_status(){
	if [ $1 != "playing" ]; then
		osascript -e "say \"[$2]\" using \"Alex\""
	fi
}

