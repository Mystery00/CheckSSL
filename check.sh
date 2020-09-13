#!/usr/bin/env bash

mkdir './tmp' -p

curl "https://$1" -k -v -s -o /dev/null 2>"./tmp/$1.result"

subject=$(grep 'subject:' "./tmp/$1.result" | sed 's|\(.*\): \(.*\)|\2|g')
issuer=$(grep 'issuer:' "./tmp/$1.result" | sed 's|\(.*\): \(.*\)|\2|g')
start=$(grep 'start date:' "./tmp/$1.result" | sed 's|\(.*\): \(.*\)|\2|g')
start=$(date -d "$start" '+%s%3N')
expire=$(grep 'expire date:' "./tmp/$1.result" | sed 's|\(.*\): \(.*\)|\2|g')
expire=$(date -d "$expire" '+%s%3N')
message=$(grep 'SSL certificate verify' "./tmp/$1.result" | sed 's|\* *\(.*\)|\1|g')

echo "{\"subject\":\"$subject\",\"issuer\":\"$issuer\",\"start\":\"$start\",\"expire\":\"$expire\",\"message\":\"$message\"}"
