#!/bin/bash -ex

deps=$1
if [ ! -e "$deps" ]; then
	echo "Usage: $0 <dependencies.tsv>"
	exit 1
fi

while read -r project repotype hash ts; do
	if [ "$repotype" = "git" ]; then
		checkout_project=$(echo $project | sed 's/golang.org\/x/github.com\/golang/')
		tmprepo=$(mktemp -d)
		git clone https://${checkout_project}.git $tmprepo
		(cd $tmprepo; git checkout -b pull-godeps $hash)
		if [ -d vendor/src/$project ]; then
			git subtree pull --squash -P vendor/src/$project $tmprepo pull-godeps
		else
			git subtree add --squash -P vendor/src/$project $tmprepo pull-godeps
		fi
		rm -rf $tmprepo
	fi
done < $deps

