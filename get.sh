#!/usr/bin/env bash

if [ -z $GOPATH ]; then
    echo "error: env: GOPATH not exists!!"
    exit 1
fi

ORG_DIR="$GOPATH/src/github.com/argcv"
PROJ_DIR="$ORG_DIR/manul"

if [ ! -d $PROJ_DIR ]; then
    mkdir -p $PROJ_DIR
    git clone git@github.com:argcv/manul.git $PROJ_DIR
    echo "Cloned to $PROJ_DIR"
fi

pushd $PROJ_DIR > /dev/null 2>&1
if [[ $(git status --porcelain) ]]; then
    echo "Found uncommitted change"
    exit 2
else
    echo "Update..."
    git pull origin
fi
popd > /dev/null 2>&1

