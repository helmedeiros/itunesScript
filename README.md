iTunesScript
============

Shell script project to control itunes from console.

## Installation

### Setup

iTunesScript has some parts, but I've tried to simplify installation - as much as I could!

First, clone down the repository:

    git clone https://github.com/helmedeiros/itunesScript.git
    
Next, you need to make the command executable:

    chmod +x itunesScript.sh
    
Now, create a symbolic link to this script in the path, /usr/bin like below:

    ln -s /Users/<user name>/<folder of cloned project>/itunesScript.sh /usr/bin/itunesScript


### Starting It Up

Open a new terminal and type the following line:

    itunesScript open
    
After that, you can start playing the first music of your iTunes library:

    itunesScript play

or another one from a specific artist:

    itunesScript play Maroon 5
