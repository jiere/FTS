#!/bin/sh
# Generate test coverage statistics for Go packages.
#
# Works around the fact that `go test -coverprofile` currently does not work
# with multiple packages, see https://code.google.com/p/go/issues/detail?id=6909
#

set -e

workdir=ut
outfile="gotest.out"
go_mod="fts.local"

generate_test_data() {
    rm -rf "$workdir"
    mkdir "$workdir"

    for pkg in "$@"; do
	#pkg=`echo $pkg`
	if [[ $pkg == "$go_mod/pkg/transaction" || $pkg == "$go_mod/pkg/account" ]];
	then
        f="$workdir/$(echo $pkg | tr / -).cover"
        go test -v -cover >  "$f" "$pkg"
	else continue;
    fi
    done

    cat "$workdir"/*.cover >> "$outfile"
}

generate_test_data $(go list ./...)
