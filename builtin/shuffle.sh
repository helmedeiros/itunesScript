#!/bin/sh
#
# "itunesScript shuffle" builtin command.
function shuffle(){
	echo "Switching shuffle on."
	osascript -e "tell application \"iTunes\"
		if current playlist's shuffle is false then
			set current playlist's shuffle to true
		else
			set current playlist's shuffle to false
		end if
	end tell";
}