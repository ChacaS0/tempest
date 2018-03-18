---
layout: default
title: How to install
---

# [](#header-2)Install TEMPest
## [](#header-2-1)Requirements
- <a href="https://git-scm.com/book/en/v1/Getting-Started-Installing-Git" target="_blank">git</a>
- <a href="https://golang.org/doc/install" target="_blank">go</a>
- **TEMPest** :sweat_smile:

## [](#header-2-2)Let's begin
So now you should have golang installed, you have access to ``GOPATH``. If you haven't set up ``GOBIN``, you should, because this is where the **TEMPest** binary will go once installed.  
  
### [](#header-2-2-1)Linux
#### [](#header-2-2-1-1)Permanent configuration of ``GOPATH``
Add this line to ``~/.bash_profile``:
```bash
GOBIN=<PATH_OF_YOUR_CHOICE>
```

Save and exit the text editor and source the ``.bash_profile``:
```bash
source ~/.bash_profile
```

Then add ``GOBIN`` to your ``PATH``

#### [](#header-2-2-1-2)Temporary configuration of ``GOBIN``
Simply run:
```bash
export GOBIN=/usr/bin/
```

### [](#header-2-2-2)Windows
> *Coming soon :stuck_out_tongue_winking_eye:*

### [](#header-2-2-3)Install **TEMPest**
So easy:
```bash
go get -v -u -t github.com/ChacaS0/tempest
```
  
Great now we just need to generate the files **TEMPest** needs to work by initialization.

## [](#header-2-3)Initialization
### [](#header-2-3-1)Command line
**TEMPest** should be installed by now, and running ``tempest -v`` should work!  
  
To initialize, run the following command:
```bash
tempest init
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
 **duration:** *You have to choose a duration greater than 1, it is for your own safety!!*

#### [](#header-2-3-1-1)Parameters

| Variable name | Description                                                                                          |
|:--------------|:-----------------------------------------------------------------------------------------------------|
| ``duration``  | The maximum age of the files to be deleted (editing an element will actualize the age of the target) |
| ``auto-mode`` | Whether **TEMPest** should run automatically (not implemented yet)                                   |

## [](#header-2-4)What's next?
Once this is done, your **TEMPest** is ready to rock!  
  
:point_right: You can see basic usage at the <a href="{{site.url}}/usage">Usage section</a>.
