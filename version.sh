#!/bin/sh
#
# itunesScript version definition.
version='1.0'
progname=`basename $0`

#  Version and help.
function version() {
    echo "$progname version $version"
}