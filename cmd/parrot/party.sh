#!/bin/bash
cd "$(dirname "$0")"
if [ ! -d ../../bin/cmd/parrot/parrots ]
then
  make party
else
  echo '[parrot] Skip party!'
fi
