#!/bin/bash

FILE=bundle.go
BACKFILE=bundle.go.bak

if test -f "$FILE"; then
    if test -f "$BACKFILE"; then
        rm -f $BACKFILE
    fi
    mv $FILE $BACKFILE
fi

fyne bundle --pkg main --name Lander_Png -o $FILE lander.png
fyne bundle -a --name GoLogo_Png -o $FILE go-logo-blue.png


