.PHONY: pg-up pg-down ch-up ch-down

tool-install:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	
pg-up:
	goose -dir ./migration/pg postgres "postgres://admin:admin@host.docker.internal:5432/test" up 
pg-down:
	goose -dir ./migration/pg postgres "postgres://admin:admin@host.docker.internal:5432/test" down

ch-up:
	goose -dir ./migration/ch clickhouse "clickhouse://admin:admin@host.docker.internal:9000/test?secure=false&skip_verify=true" up 
ch-down:
	goose -dir ./migration/ch clickhouse "clickhouse://admin:admin@host.docker.internal:9000/test?secure=false&skip_verify=true" down