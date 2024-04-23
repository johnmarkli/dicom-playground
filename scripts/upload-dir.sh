#!/usr/bin/env bash
# 
# Upload directory of DICOMs to server through API
# Usage
# ./upload-dir.sh [-d] [-u]

dir=$(PWD)
url="localhost:8080"

while getopts "d:u:h" arg; do
  case $arg in
    d) dir=$OPTARG;;
    u) url=$OPTARG;;
    h) echo "Usage $(basename $0) [-d] [-u]";;
  esac
done

echo "DIR $dir"
echo "URL $url"
echo "NAME $(basename $0)"

for f in $dir/*
do
  if [[ -f $f ]]; then
    echo "Uploading $f"
    curl --form file="@$f" "$url/dicoms"
  fi
done

echo "Upload finished"

