#!/bin/bash

echo "BUILD >>>>> core:cmd "
cd "$(dirname "$0")"
make
r=$?
if [ $r -eq 0 ]
then
    echo " BUILD OK <<<<< core "
else
    echo " BUILD FAIL !!! core "
    exit 1
fi

for d in */
    do
    echo " BUILD >>>>> $d "
    cd $d
    make
    r=$?
    if [ $r -eq 0 ]
    then
        echo " BUILD OK <<<<< $d "
    else
        echo " BUILD FAIL !!! $d "
        exit 1
    fi
    cd ..
    echo
done
