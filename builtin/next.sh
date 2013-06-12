#!/bin/sh
#
# "itunesScript next" builtin command.
function next(){
	osascript -e 'tell application "iTunes" to next track';
}