#!/bin/sh
#
# "itunesScript play" builtin command.
function playing(){
	if [ $# -eq 1 ]; then
		artistName=$1
		osascript -e "tell application \"iTunes\"
			tell source \"Library\"
				tell library playlist 1
					set albumList to album of (every track where artist contains \"$artistName\") --list of album names for ARTIST tracks
					set response to (first item of albumList) as string --present a list of unique names and get the user to select one
				end tell
		
				if not ((name of playlists) contains \"Current Album\") then --check if playlist current Album already exists
					set newPlaylist to make new playlist with properties {name:\"Current Album\"} --No? Then make it
				else
					set newPlaylist to playlist \"Current Album\" --Yes? 
					delete every track of newPlaylist --Then delete current references
				end if
			
				tell library playlist 1 to duplicate (every track whose album is response) to newPlaylist --find all tracks that have selected album name and copy them to the playlist
					play newPlaylist --play
				end tell
			end tell";
	else
		osascript -e 'tell application "iTunes" to play';
	fi
}