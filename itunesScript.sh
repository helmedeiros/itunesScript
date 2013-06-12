#!/bin/sh
#
#Command-line controller for Apple's iTunes.
#
#Main script to manage input options.

source $(dirname $0)/help.sh
source $(dirname $0)/commands.sh

version='1.0'
progname=`basename $0`
itunes_script_usage_string="$progname [--version] [--help] <command> [<args>]";

#  Version and help.
function version() {
    echo "$progname version $version"
}

function handle_options(){
	while [ $# -gt 0 ]; do
	    arg=$1;
		case $arg in
			"open"		) cmd_open;
	            break;;
			"play"		) cmd_play $2
				break ;;
			"pause"		) cmd_pause
				break ;;
			"stop"		) cmd_stop
					break ;;
			"next"		) cmd_next
					break ;;
			"prev"		) cmd_prev
					break ;;
			"quit"		) cmd_quit
					break ;;
			"--version"	) version;
	            break;;
			"--help" | *) list_common_cmds_help;
				break;;
		esac
	done
}

function list_commands(){
	printf "usage: %s\n\n" "$itunes_script_usage_string";
	
	list_common_cmds_help
}

########################~Main~############################
# The user didn't specify a command; give them help
if [ $# = 0 ]; then
	list_commands;
fi
	
handle_options $1 $2;
