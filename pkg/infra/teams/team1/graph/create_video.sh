#!/bin/bash

DOTFILES=$(find pics/ "*.gv")

for FILE in ${DOTFILES}; do 
    dot -Tpng -Kneato -s1 -O $FILE
done

ffmpeg -framerate 2 -i pics/graph%d.gv.png -c:v libx264 -strict -2 -preset slow -pix_fmt yuv420p -vf "scale=trunc(iw/2)*2:trunc(ih/2)*2" -f mp4 out.mp4