#!/bin/sh

go build full.go -o ./full
export FROMENV='value from env'
./full -FROMCMDLINE='value from command line'

echo test generalitem
export GENERALITEM='generalitem value from env'
./full -GENERALITEM='generalitem value from command line'

./full

export GENERALITEM=''
./full