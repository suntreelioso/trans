#!/bin/bash

version=$(git describe --tags 2>/dev/null)

mkdir -p release

cd build
for dir in $(ls)
do
    if [ -d $dir ]; then
        zip_name="${dir}-${version}.zip"
        zip -r ../release/$zip_name $dir
    fi
done
