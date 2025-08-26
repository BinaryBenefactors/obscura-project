.PHONY: test test-docker db-manager

# Запуск тестов локально
test:
	go test -v ./internal/domain/repository/...

# Запуск тестов в Docker
test-docker:
	docker-compose up --build backend-test

# Сборка и запуск db-manager
db-manager:
	docker-compose exec backend-app go run cmd/db-manager/main.go $(ARGS)

# Показать всех пользователей
list-users:
	$(MAKE) db-manager ARGS="-cmd list-users"

# Показать информацию о пользователе
get-user:
	$(MAKE) db-manager ARGS="-cmd get-user -id $(id)"

# Показать информацию о пользователе по email
get-user-by-email:
	$(MAKE) db-manager ARGS="-cmd get-user -email $(email)"

# Создать нового пользователя
create-user:
	$(MAKE) db-manager ARGS="-cmd create-user -email $(email) -name '$(name)' -password '$(password)'"

# Удалить пользователя
delete-user:
	$(MAKE) db-manager ARGS="-cmd delete-user -id $(id)"

# Обновить пользователя
update-user:
	$(MAKE) db-manager ARGS="-cmd update-user -id $(id) -name '$(name)' -password '$(password)'"

# Показать активные сессии
list-sessions:
	$(MAKE) db-manager ARGS="-cmd list-sessions" 