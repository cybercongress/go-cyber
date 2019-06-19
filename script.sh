#!/bin/bash
for file in ~/build/docs/*
do
if [ -f "$file" ]
then
touch temp.md
echo "---
project: cyberd
---" >> temp.md
cat $file >> temp.md
cat temp.md > $file
rm -rf temp.md
fi
done
