#!/bin/sh
#
# "itunesScript status" builtin command.
function status(){
	if [ $1 != "stopped" ]; then
		artist=`osascript -e 'tell application "iTunes" to artist of current track as string'`;
		track=`osascript -e 'tell application "iTunes" to name of current track as string'`;
		album=`osascript -e 'tell application "itunes" to album of current track as string'`;
		echo "The current track is $track at $album by $artist";
	fi
}

