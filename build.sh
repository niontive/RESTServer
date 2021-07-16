#!/bin/sh

# Output directory
OUTDIR="bin/"

# Executable name
OUTFILE="server"

# Files to build
MAIN="main.go"

# Programs used for building
GO="/bin/go"

build_main() {
    ${GO} build -o ${OUTDIR}${OUTFILE} ${MAIN}
}

# Begin script execution
build_main