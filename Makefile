up-db:
	docker run -d --name mssql-22 -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=kursPswdsuper1' -p 1433:1433  --rm mcr.microsoft.com/mssql/server:2022-latest

stop-db:
	docker stop mssql-22


init-db: up-db
	sleep 10 && docker exec -it kmftw_mssql_1 /opt/mssql-tools/bin/sqlcmd -U sa -P kursPswdsuper1 -Q "CREATE DATABASE TEST"

run:
	docker-compose up

stop:
	docker-compose down