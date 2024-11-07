#!/bin/bash
docker build -t compile-server .
docker run -p 1235:1235 -d -it --rm --name golang-server compile-server
