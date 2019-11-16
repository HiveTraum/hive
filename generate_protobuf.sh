#!/bin/bash

cd inout || exit
protoc --go_out=. inout.proto