.phony: bench
bench:
	docker-compose up -d
	go test -bench=. -count=5 > out.txt
	docker-compose down
	benchstat out.txt
