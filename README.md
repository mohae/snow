# snow

Snow is a [security now podcast](https://twit.tv/shows/security-now) downloader.

By default, only the latest episode is downloaded, unless it already exists in the destination directory.

Episodes can also be downloaded by:
* last n episodes
* specific episode number
* a range of episode numbers
* all episodes

Currently, only downloading of the MP3 files is supported, both the regular version and the "low-quality" version.

Snow should work on any platform that Go supports.

## Usage
### Compile
Assuming you have [Go](https://golang.org) installed:

    go install github.com/mohae/snow

A `snow` executable will be in your `$GOPATH/bin` directory.

### Run

This will download the latest episode:

    $ snow

This will download low-quality versions of the last 10 episodes:

    $ snow -lastn 10 -lq

This will download episode 500:

    $ snow -start 500

This will download episode 11-42, inclusive:

    $ snow -start 11 -stop 42

This will download all episodes:

    $ snow -lastn 1

or

    $ snow -start 1

### Flags

Flag | Type | Default | Description  
|:--|:--|:--|:--  
help, h|false|bool|help output  
low|false|bool|download the low quality version: 16Kbps mp3  
overwrite|false|bool|overwrite existing file, if one exists  
verbose|false|bool|verbose output
concurrency|1|int|number of episodes to concurrently download  
lastn|1|int|download the last n episodes; 0 means all  
start|0|int|episode number from which to start downloading  
stop|0|int|episode number at which to stop downloading  
savedir|$HOME/Downloads/security-now|string|save directory  

## License
Apache License, Version 2.0
