#!/bin/bash
# Go to the program dir
cd $GOPATH/src/bitbucket.org/ChacaS0/scripts

# get the code from bitbucket #!careful, you need an access key!
git pull origin master

# build the muthafukkah
go install $GOPATH/src/bitbucket.org/ChacaS0/scripts/tempest/tempest.go

