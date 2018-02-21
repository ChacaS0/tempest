#!/bin/bash
#===============  - BETA -  ===============#
# BETA_VERSION=$GOPATH/src/github.com/ChacaS0/tempest

# Go to the program dir
cd $GOPATH/src/github.com/ChacaS0/tempest

# get the code from bitbucket
git pull origin master

# build the muthafukkah
go install -i $GOPATH/src/github.com/ChacaS0/tempest/tempest.go