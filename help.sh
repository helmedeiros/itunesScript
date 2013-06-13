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
	echo "    play		Start playing iTunes musics.";
	echo "    pause		Pausing iTunes musics.";
	echo "    stop		Stop iTunes.";
	echo "    quit		Quit iTunes.";
}