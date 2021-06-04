#! /bin/bash

make build
reflex -sr '\.go$' make serve