#!/bin/bash
find . -name "*_test.go" -exec sed -i 's/"device"/testDeviceParameter/g' {} +
find . -name "*_test.go" -exec sed -i 's/"createInfo"/testCreateInfoParameter/g' {} +
find . -name "*_test.go" -exec sed -i 's/"memory"/testMemoryParameter/g' {} +
find . -name "*_test.go" -exec sed -i 's/"nil createInfo"/testNilCreateInfo/g' {} +
find . -name "*_test.go" -exec sed -i 's/"ValidationError"/testValidationErrorType/g' {} +
find . -name "*_test.go" -exec sed -i 's/"nil device"/testNilDevice/g' {} +
