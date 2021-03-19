#!/bin/bash
docker run --rm -it --privileged -v $(pwd)/data:/opt/Mashiron-go/bin/data -v $(pwd)/mashiron.ini:/opt/Mashiron-go/bin/mashiron.ini
