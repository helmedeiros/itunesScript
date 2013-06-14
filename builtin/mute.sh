#!/bin/sh
#
# "itunesScript mute" builtin command.
function mutting(){
	osascript -e 'tell application "iTunes" to set mute to true';
}