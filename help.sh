#!/bin/sh
#
# itunesScript helps messages.
itunes_script_usage_string="$progname [--version] [--help] <command> [<args>]";

function list_commands(){
	printf "usage: %s\n\n" "$itunes_script_usage_string";
	
	list_common_cmds_help
}

function list_common_cmds_help(){
	echo "The most commonly used `basename $0` commands are";
	echo "    open		Start iTunes.";
	echo "    status		Show iTunes status.";
	echo "    play		Start playing iTunes musics.";
	echo "    pause		Pausing iTunes musics.";
	echo "    next		Send to the next iTunes musics.";
	echo "    prev		Back to the previous iTunes musics.";
	echo "    mute		Mute iTunes.";
    echo "    unmute		Unmute iTunes.";
	echo "    vol up		Increase iTunes vol by 10%."
	echo "    vol down		Decrease iTunes vol by 10%."
	echo "    vol #		Change iTunes vol to # [0-100%]."
	echo "    stop		Stop iTunes.";
	echo "    quit		Quit iTunes.";
}