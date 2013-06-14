#!/bin/sh
#
# "itunesScript stop" builtin command.
function stopping(){
	osascript -e 'tell application "iTunes" to stop';
}