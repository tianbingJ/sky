echo "total:"
find . -name "*.go"| xargs cat | wc

echo "test:"
find . -name "*test.go"| xargs cat | wc

