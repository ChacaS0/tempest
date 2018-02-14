#!/bin/bash
# # Go to the program dir
# cd $GOPATH/src/bitbucket.org/ChacaS0/scripts

# # get the code from bitbucket #!careful, you need an access key!
# git pull origin master

# # build the muthafukkah
# go install $GOPATH/src/bitbucket.org/ChacaS0/scripts/tempest/tempest.go

#===============

#* Check if Alpha version exists
if [ -d "$ALPHA_VERSION" ]; then 
  if [ -L "$ALPHA_VERSION" ]; then
    # It is a symlink!
    # Don't give a fck about this shit if exists. Why the fck would that exist
    rm $ALPHA_VERSION
	 echo ":: Removed the alpha version anyway"
  else
    # It's a directory!
    rm -r $ALPHA_VERSION
	 echo ":: Removed the alpha version"
  fi
fi

#* Regular installation for Beta version
#? Does ``go get -u`` works just fine?? Does it install too ?
# Go to the program dir
# cd $GOPATH/src/github.com/ChacaS0/tempest
# Check if Beta version already exists
if [ -d "$BETA_VERSION" ]; then
	if [ -L "$BETA_VERSION" ]; then
		# Symlink: don't care
	else
		# # It exists already, then we update only
		# cd $GOPATH/src/github.com/ChacaS0/tempest
		# git pull origin master
		# # Update Success ?
	fi
else
	# # Doesn't exist, so create it
	# mkdir $GOPATH/src/github.com/ChacaS0
	# git clone https://github.com/ChacaS0/tempest.git
	# echo ":: End of the migration towards Beta"
fi