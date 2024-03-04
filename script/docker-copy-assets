#!/usr/bin/env bash
BASEDIR=$(builtin cd ..; pwd)
CCTL_ASSETS=/home/cctl/cctl/assets
# clear the assets folder
rm -rf  ${BASEDIR}/assets
mkdir ${BASEDIR}/assets
mkdir ${BASEDIR}/assets/net-1
mkdir ${BASEDIR}/assets/chainspec
mkdir ${BASEDIR}/assets/faucet
# copy net-1 users
docker cp cspr-cctl:$CCTL_ASSETS/users/. ${BASEDIR}/assets/net-1
# copy net-1 chainspec
docker cp cspr-cctl:$CCTL_ASSETS/genesis ${BASEDIR}/assets/net-1/chainspec
# copy faucet keys
docker cp cspr-cctl:$CCTL_ASSETS/faucet/ ${BASEDIR}/assets/net-1/faucet
