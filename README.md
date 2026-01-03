# URL Shortener

## Запуск

```bash
# Запуск проекта
docker compose up -d

# Остановка
docker compose down

# Логи
docker compose logs -f
```

## Миграции

Миграции выполняются автоматически при `docker compose up`. Для ручного управления:

```bash
# Накатить миграции
docker compose run --rm migrate up

# Откатить последнюю миграцию
docker compose run --rm migrate down

# Откатить все миграции
docker compose run --rm migrate reset

# Статус миграций
docker compose run --rm migrate status
```

## TODO

- [x] Base functionality
- [ ] Tests
- [ ] Metrics
- [ ] Redis
- [ ] CI/CD
- [ ] Docs
