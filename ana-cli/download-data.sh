#!/bin/bash

# get the directory of the script
cd $(dirname $0)

# download the reddit jokes dataset
curl -L -o ../data/text_emotion.csv https://query.data.world/s/gvmyfoxhetec4dbecnvdwsmohfj5cw
