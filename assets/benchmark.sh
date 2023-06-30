#!/bin/bash
echo Testing without proxy
fstart=`date +%s.%N`
file="./domain.txt"
while IFS= read -r line
do
    # display $line or do somthing with $line
    nslookup "$line" > /dev/null
done <"$file"
fend=`date +%s.%N`

fruntime=$( echo "$fend - $fstart" | bc -l )
echo exec time without proxy $fruntime

echo Testing with proxy
start=`date +%s.%N`
file="./domain.txt"
while IFS= read -r line
do
    # display $line or do somthing with $line
    nslookup -port=53  "$line" '127.0.0.1'  > /dev/null
done <"$file"
end=`date +%s.%N`
runtime=$( echo "$end - $start" | bc -l )
echo exec time with proxy $runtime
