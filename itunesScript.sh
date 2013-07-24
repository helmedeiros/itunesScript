#!/bin/sh
#
#Command-line controller for Apple's iTunes.
#
#Main script to manage input options.
source $(dirname $0)/version.sh
source $(dirname $0)/help.sh
source $(dirname $0)/commands.sh

function handle_options(){
	# The user didn't specify a command; give them help
	if [ $# = 0 ]; then
		list_commands;
	else 
		while [ $# -gt 0 ]; do
			arg=$1;
			case $arg in
				"open"		) cmd_open;
				break;;
				"play"		) cmd_play $2;
				break ;;
				"pause"		) cmd_pause;
				break ;;
				"next"		) cmd_next;
				break ;;
				"prev"		) cmd_prev;
				break ;;
				"mute"		) cmd_mute;
				break ;;
				"unmute"	) cmd_unmute;
				break ;;
				"vol"		) cmd_vol $2;
				break;;
				"status"	) cmd_status;
				break;;
				"shuffle"	) cmd_shuffle;
				break;;
				"stop"		) cmd_stop;
				break ;;
				"quit"		) cmd_quit;
				break ;;
				"--version"	) version;
				break;;
				"--help" | *) list_common_cmds_help;
				break;;
			esac
		done
	fi
}
	
handle_options $@;
