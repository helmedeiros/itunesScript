#!/bin/sh
#
# "itunesScript play" builtin command.
function playing(){
	if [ $# -eq 1 ]; then
		artistName=$1
        osascript -e "tell application \"iTunes\" to play (get item 1 of (every track where artist contains \"$artistName\"))";
	else
		osascript -e 'tell application "iTunes" to play';
	fi
}