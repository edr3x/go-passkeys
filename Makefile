.PHONY: watch migrate
	
include .env

watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

migrate:
	@psql "$(POSTGRESQL_URL)" -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
	POSTGRESQL_URL="$(POSTGRESQL_URL)" tern migrate --migrations migrations/
	 
rollback:
	POSTGRESQL_URL="$(POSTGRESQL_URL)" tern migrate --migrations migrations/ --destination ${dest}
