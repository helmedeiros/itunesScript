iTunesScript
============

Shell script project to control itunes from console.

## Installation

### OS X

iTunesScript is built for the Mac.

### Setup

iTunesScript has some parts, but I've tried to simplify installation - as much as I could!

First, clone down the repository:

    git clone https://github.com/helmedeiros/itunesScript.git
    
Next, you need to make the command executable:

    chmod +x itunesScript
    
To make sure my shell knows where to find iTunesScript you will need to add the addres from where you've cloned the project to your .bashrc file's PATH variable. Here's how mine looks

    export PATH=${PATH}:/Users/helmed/Projects/workspaceShell/itunesScript/
    
Make sure you reload your shell with

    source ~/.bashrc


### Starting It Up

Open a new terminal and type the following line:

    itunesScript open
    
After that, you can start playing the first music of your iTunes library:

    itunesScript play

or another one from a specific artist:

    itunesScript play Maroon 5
