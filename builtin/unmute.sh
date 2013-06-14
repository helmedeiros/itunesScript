#!/bin/sh
#
# "itunesScript unmute" builtin command.
function unmutting(){
	osascript -e 'tell application "iTunes" to set mute to false';
}