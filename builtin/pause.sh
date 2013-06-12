#!/bin/sh
#
# "itunesScript pause" builtin command.
function pausing(){
	osascript -e 'tell application "iTunes" to pause';
}