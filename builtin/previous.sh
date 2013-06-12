#!/bin/sh
#
# "itunesScript prev" builtin command.
function prev(){
	osascript -e 'tell application "iTunes" to previous track';
}