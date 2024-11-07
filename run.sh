#!/bin/bash
docker build -t compile-server
docker run -p 1234:1234 -d -it --rm --name golang-server compile-server
