#!/bin/sh

cat > backend/src/utils/version.go << EOF
package utils

func GetVersion() string {
        return "$(git describe --tags HEAD)"
}
EOF
