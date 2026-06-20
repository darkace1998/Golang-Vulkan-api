#!/bin/bash
git checkout synchronization_test.go
sed -i -e '150,152d' synchronization_test.go
sed -i -e '3a\	"errors"' synchronization_test.go
