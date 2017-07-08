#!/bin/bash
docker pull golang:alpine
docker build -t toomore/gogrs:latest ./
