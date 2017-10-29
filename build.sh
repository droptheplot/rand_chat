#!/bin/bash

git pull
go build github.com/droptheplot/rand_chat
sudo service rand-chat restart
