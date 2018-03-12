# <a href="https://github.com/ChacaS0/tempest" ><img src="https://image.ibb.co/c829L7/logo_v5_small.gif" width="25" height="25" /></a> TEMPest
*TEMPest is a tool to manage easily temporary folders/files*

[![Build Status](https://travis-ci.org/ChacaS0/tempest.svg?branch=master)](https://travis-ci.org/ChacaS0/tempest) [![codecov](https://codecov.io/gh/ChacaS0/tempest/branch/master/graph/badge.svg)](https://codecov.io/gh/ChacaS0/tempest) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/25460f686b864a68a44f6871ec3b9ad2)](https://www.codacy.com/app/ChacaS0/tempest?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ChacaS0/tempest&amp;utm_campaign=Badge_Grade) [![Go Report Card](https://goreportcard.com/badge/github.com/ChacaS0/tempest)](https://goreportcard.com/report/github.com/ChacaS0/tempest) [![GoDoc](https://godoc.org/github.com/ChacaS0/tempest?status.svg)](https://godoc.org/github.com/ChacaS0/tempest) 
[![apm](https://img.shields.io/apm/l/vim-mode.svg)](https://github.com/ChacaS0/tempest/blob/master/LICENSE)
 [![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=TQYU8BPPV5YWA) [![beta](https://img.shields.io/badge/edition-beta-ff69b4.svg)](https://github.com/ChacaS0/tempest)

## Index
- [Installation](#installation)
	- [From source](#from-source-github)
		- [Requirements](#requirements)
		- [Command](#command)
		- [Temporarily ``GOBIN``](#gobin-permanently-in-path)
- [Initialization](#initialization)
	- [Command line ``init``](#command-line-init)
	- [Parameters](#parameters)
- [Add a new temp folder to the list](#add-a-new-temp-folder-to-the-list)
	- [To add the current directory](#to-add-the-current-directory)
		- [Through command-line](#through-command-line)
		- [Through text editor](#through-text-editor)
- [List the current directories added to **TEMPest**](#list-the-current-directories-added-to-tempest)
	- [Using **TEMPest**](#using-tempest)
	- [Viewing the file ``~/.tempestcf``](#viewing-the-file-tempestcf)
- [Runing a global purge](#runing-a-global-purge)
	- [Test mode](#test-mode)
	- [Real one](#real-one)
- [Purging one directory](#purging-one-directory)
	- [Test mode](#test-mode-1)
	- [Real one](#real-one-1)
- [Access the documentation](#access-the-documentation-usage)

## Installation
### From source (Github)
#### Requirements
* [git]("https://git-scm.com/downloads")
* [go (golang)]("https://golang.org/")
#### Command
First choose if you want to add permanantly ``GOBIN`` to your ``PATH``
##### ``GOBIN`` temporarily in ``PATH``
```bash
export GOBIN=/bin
```
or
```bash
export GOBIN=/usr/bin
```

Then
```bash
go get -v -u github.com/ChacaS0/tempest
```
##### ``GOBIN`` permanently in ``PATH``
Add this line to ``/etc/environment``:
```bash
GOBIN=<PATH_OF_YOUR_CHOICE>
PATH=$PATH:$GOBIN
```

Then
```bash
go get -v -u github.com/ChacaS0/tempest
```

## Initialization
### Command line ``init``
It is very easy to use.  
First, to initialize it the first time, run:

```bash
> tempest init
```

This will generate a ``~/.tempestcf`` file. 
It will hold the list of all the temp directories.

> If there is an issue and the file can't be created somehow, you can still create it at its default location: ``$HOME/.tempestcf`` and leave it empty for now.
  
This will also generate a ``~/.tempest.yaml`` file.
It will hold the configuration of **TEMPest**.  

> If there is an issue and the file can't be created somehow, you can still crate it at its default location: ``$HOME/.tempest.yaml`` with default content:
```yaml
duration: 5
auto-mode: false
```
#### Parameters:

* **duration** : This is the maximum age the content of the temp directories, choose it carefully!

> *You have to choose a duration greater than 1, it is for your own safety!!*

## Add a new temp folder to the list
### To add the current directory
Positionate yourself to the deried directory. 
For example, if you want to add /tmp, use this instructions:  

```bash
$ cd /temp
$ tempest add
```

### Through command-line
Just run:

```bash
tempest add ...<PATH>
```

* **``<PATH>``** being the path to the directory to be added to the list of temp directories  
* **``...``** meaning that many arguments can be passed  

> **By convention we will give the name ``temp`` to the directories to be added to ``tempest``**

### Through text editor
Just open ``~/.tempestcf`` and add a new line with the <u>absolute path</u> of the temp directory to be added.

## List the current directories added to TEMPest
### Using TEMPest

```bash
$ tempest list
```

### Viewing the file ``~/.tempestcf``
``$ cat ~/.tempestcf``  
Or  
``$ vi ~/.tempestcf``  

## Updating TEMPest
```bash
$ tempest update
```

## Runing a global purge
*The age of the files deleted will be the one older than the number of days set as "duration" in ~/.tempest.yaml*
### Test mode
In this mode, it will display the file it <u>would</u> delete plus the size.  
**Nothing gets deleted**. To do so, try:

```bash
$ tempest start -t
```

### Real one
**Runing this will actually delete files/directories, make sure everything inside ``~/.tempestcf`` is meant to be there with ``tempest list`` first.**

```bash
$ tempest start
```

## Purging one directory
>*Still based on the config file*  

It is possible to purge a directory even if it is not added to tempest. There is also a test mode for this one.
### Test mode

```bash
$ tempest purge -p <PATH> -t
```

* **``<PATH>``** is the path you want to purge
* **``-t``** declare the test mode

### Real one

```bash
$ tempest purge -p <PATH>
```

* **``<PATH>``** is the path you want to purge

## Access the documentation (usage)
It is recommanded to have **Showdown** installed. If you don't but are interested, checkout this [link]("https://github.com/craigbarnes/showdown").  

```bash
$ tempest doc
```  

There is also a "man like" view of the documentation.  

```bash
$ tempest doc -m
```


-------------------
*Thanks to*
> <a href="https://github.com/golang/go" target="_blank"><img src="https://upload.wikimedia.org/wikipedia/commons/2/23/Golang.png" width=33%></a> 

> <a href="https://github.com/spf13/cobra" target="_blank"><img src="https://cloud.githubusercontent.com/assets/173412/10886352/ad566232-814f-11e5-9cd0-aa101788c117.png" width=33%></a> 

> <a href="https://github.com/spf13/viper" target="_blank"><img src="https://cloud.githubusercontent.com/assets/173412/10886745/998df88a-8151-11e5-9448-4736db51020d.png" width=33%></a> 