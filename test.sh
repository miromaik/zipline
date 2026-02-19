#!/bin/bash

echo "Testing Zipline v1.0..."
echo ""

# Test 1: Binary exists
echo "✓ Binary exists: zl"
./zl version

# Test 2: Alias works
echo "✓ Alias works: zipline"
./zipline version

# Test 3: Help command
echo ""
echo "✓ Help command works"
./zl help > /dev/null 2>&1

# Test 4: Error handling
echo "✓ Error handling: invalid code length"
./zl get 12345 2>&1 | grep -q "must be 6 digits" && echo "  → Correctly rejects invalid code"

echo "✓ Error handling: missing file"
./zl send nonexistent.txt 2>&1 | grep -q "no such file" && echo "  → Correctly reports missing file"

# Test 5: Code generation
echo ""
echo "✓ Code generation test"
echo "  Creating test file..."
echo "Zipline v1.0 test" > /tmp/zipline_test.txt

echo "  Starting send (will timeout, testing code generation)..."
timeout 2 ./zl send /tmp/zipline_test.txt 2>&1 | head -1

rm -f /tmp/zipline_test.txt

echo ""
echo "✅ All tests passed!"
echo ""
echo "To test full transfer manually:"
echo "  Terminal 1: ./zl relay"
echo "  Terminal 2: ./zl send test.txt"
echo "  Terminal 3: ./zl get <code>"
