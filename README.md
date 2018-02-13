# TEMPest
*TEMPest is a tool to manage easily temp folders/files*

> *Not up to date yet!*

## Installation
### From AUR (ArchLinux)
*coming soon*
### From source (Github)
*coming soon*
## Initialization
### Command line ``init``
It is very easy to use.  
First, to initialize it the first time, run:

```bash
> tempest init
```

This will generate a ``~/.tempestcf`` file. 
It will hold the list of all the temp directories.

### Config file for tempest ``~/.tempest.yaml``
#### First check if it exists

```bash
$ cat $HOME/.tempest.yaml
```

#### If it doesn't, create it
For now, you **have to** create this file manually. 
Idealy you want to place it in the home directory. You can achieve that by runing:

```bash
$ touch $HOME/.tempest.yaml
```

Then add the following to the file and save it:

```yaml
{
	"duration": 5
}
```

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
*Still buggy*
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