#!/bin/sh

cat > ./src/models/version.go << EOF
package models

func getVersion() string {
    return "$(git describe --tags HEAD)"
}

EOF
