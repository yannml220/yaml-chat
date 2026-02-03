build:
	go build -o api ./cmd/api

start :
	cd build/ && docker-compose up

down :
	cd build/ && docker-compose down

start-fiber-server :
	go run ./cmd/api/main.go


start-postgres :
	cd build/ && docker-compose up -d postgres

stop-postgres :
	cd build/ && docker-compose stop postgres

run-postgres : start-postgres
	sleep 5 && cd build/ &&  docker-compose exec yann psql -U yann -d chat_agent_db

refresh-group-permission : 
	newgrp docker 

goose-init :
	goose -dir "./internal/app/db/migrations"  create init sql && goose -dir "./internal/app/db/migrations" up


migration-create : 
	goose -dir "./internal/app/db/migrations" create $(NAME) sql 


migration-apply :
	GOOSE_DBSTRING=$(DSN) goose -dir "./internal/app/db/migrations" postgres  up 

migration-apply-easy :
	GOOSE_DBSTRING=postgres://yann:yann@localhost:5432/chat_agent_db?sslmode=disable goose -dir "./internal/app/db/migrations" postgres  up 

start-python-service :
	cd services/python/app/ && uvicorn main:app --reload --host 127.0.0.1 --port 8000


migration-rollback :
	GOOSE_DBSTRING=$(DSN)  goose -dir  "./internal/app/db/migrations" postgres down-to $(VERSION) 


migration-reset :
	GOOSE_DBSTRING=$(DSN) goose -dir "./internal/app/db/migrations" postgres reset 


migration-rollback-last :
	GOOSE_DBSTRING=$(DSN)  goose -dir  "./internal/app/db/migrations" postgres down
