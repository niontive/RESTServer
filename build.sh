#!/bin/sh

# Output directory
OUTDIR="bin/"

# Executable name
OUTFILE="server"

# Files to build
FILES="main.go metadata.go"

# Programs used for building
GO="/bin/go"

build_main() {
    ${GO} build -o ${OUTDIR}${OUTFILE} ${FILES}
}

# Begin script execution
build_main