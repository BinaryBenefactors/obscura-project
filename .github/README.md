# GitHub Actions

В этом проекте используются GitHub Actions для непрерывной интеграции и доставки, с workflow, настроенными на self-hosted runner'ы.

## Workflows

- `ci.yml`: Запускается при каждом пуше в ветки main/develop и pull request'ах в main. Включает:
  - Настройку окружений Node.js и Go
  - Запуск тестов для frontend (Next.js) и backend (Go)
  - Сборку frontend и backend приложений с помощью Docker
  - Запуск Docker Compose для интеграционного тестирования

## Self-Hosted Runner'ы

В проекте используются self-hosted runner'ы для лучшего контроля над окружением выполнения и обработки ресурсоемких задач, таких как сборка Docker-образов.

### Настройка Self-Hosted Runner'ов

См. документацию GitHub о том, как настроить self-hosted runner'ы: https://docs.github.com/en/actions/hosting-your-own-runners/about-self-hosted-runners

## Необходимые секреты

Следующие секреты должны быть настроены в настройках репозитория GitHub:

- `AWS_ACCESS_KEY_ID`: AWS ключ доступа для деплоя
- `AWS_SECRET_ACCESS_KEY`: AWS секретный ключ для деплоя
- `SSH_PRIVATE_KEY`: Приватный ключ для серверов деплоя