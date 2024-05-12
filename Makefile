test:
	go test ./...

bench:
	go test -benchmem -run=^$$ -bench . github.com/Gandalf-Le-Dev/abyssnet.protocol.gows

cover:
	go test -coverprofile=./bin/cover.out --cover ./...
