#!/bin/bash
#===============  - ALPHA -  ===============#
ALPHA_VERSION=$GOPATH/src/bitbucket.org/ChacaS0/scripts/tempest
#! @Deprecated
# # Go to the program dir
# cd $GOPATH/src/bitbucket.org/ChacaS0/scripts

# # get the code from bitbucket #!careful, you need an access key!
# git pull origin master

# # build the muthafukkah
# go install $GOPATH/src/bitbucket.org/ChacaS0/scripts/tempest/tempest.go

#===============  - BETA -  ===============#
BETA_VERSION=$GOPATH/src/github.com/ChacaS0/tempest

#* Check if Alpha version exists
if [ -d "$ALPHA_VERSION" ]; then 
  if [ -L "$ALPHA_VERSION" ]; then
    # It is a symlink!
    # Don't give a fck about this shit if exists. Why the fck would that exist
    rm -r $ALPHA_VERSION
	 echo ":: Removed what seems to be an alpha version ... wth?"
  else
	# It's a directory!
	
	# The bin exists too ? GTFO!!
	if [ -f "$GOBIN/tempest" ]; then
		cd $ALPHA_VERSION
		go clean -i tempest.go
	else
		# all good then
	fi

	# Remove folder then
	rm -r $ALPHA_VERSION #? ``rm -rf`` needed?
	echo ":: Removed the alpha version"
  fi

fi


#* Regular installation for Beta version
#! Doesn't work ``go get -u``!! Doesn't install
# Go to the program dir
# cd $GOPATH/src/github.com/ChacaS0/tempest
# Check if Beta version already exists
if [ -d "$BETA_VERSION" ]; then
	# Then it's all good
else
	# # Doesn't exist, so create it
	mkdir -p $GOPATH/src/github.com/ChacaS0
	cd $GOPATH/src/github.com/ChacaS0
	git clone https://github.com/ChacaS0/tempest.git
	echo ":: End of the migration towards Beta"
fi
# #* [Update]
echo ":: Updating"
cd $BETA_VERSION
git pull origin master 
go install $BETA_VERSION/tempest.go
echo "[Update]:: Installed!"