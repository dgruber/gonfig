#!/bin/sh

# example for setting up pipeline in concourse

fly -t ci set-pipeline --config pipeline.yml --load-vars-from params.xml -p gonfig
fly -t ci up -p gonfig
