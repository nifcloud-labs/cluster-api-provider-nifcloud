#!/bin/bash

for i in $(find .. -name '*.go')
do
  if ! grep -q Copyright "$i"
  then
    cat <(cat boilerplate.go.txt | sed -e "s/_YEAR_/$(date '+%Y')/g") "$i" > "$i".new && mv "$i".new "$i"
  fi
done
