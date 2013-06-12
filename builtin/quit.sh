#!/bin/sh
#
# "itunesScript quit" builtin command.
function quiting(){
	 osascript -e 'tell application "iTunes" to quit';
}