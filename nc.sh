#!/bin/bash
echo "this has started"
keyword="test"
while true
do
  validator=$(echo $keyword | nc server 12345) # server is the intern name on the network, the number is the port
  if [ "$validator" = "$keyword" ]; then
    echo "server is up"
    exit 0
  fi
  echo "server is down"
  sleep 5
done
