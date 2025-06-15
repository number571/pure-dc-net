.PHONY: build
build:
	go build -o bin/dc ./cmd/dc
	echo test/node1/dc_iter.txt test/node2/dc_iter.txt test/node3/dc_iter.txt | xargs -n 1 cp test/_init_/dc_iter.txt
	echo test/node1/dc test/node2/dc test/node3/dc | xargs -n 1 cp bin/dc 
