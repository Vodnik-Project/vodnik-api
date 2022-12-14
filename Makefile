include app.env

createdb:
	createdb --username $(DB_USER) $(DB_NAME)
createdbtest:
	createdb --username $(DB_USER) $(DB_NAME_TEST)
dbschema:
	psql $(DB_SOURCE) -f db/schema/schema.sql
dbschematest:
	psql $(DB_SOURCE_TEST) -f db/schema/schema.sql
dropdb:
	dropdb $(DB_NAME)
dropdbtest:
	dropdb $(DB_NAME_TEST)
test:
	go test -v --cover ./...

.PHONY: env createdb createdbtest dbschema dbschematest dropdb dropdbtest test