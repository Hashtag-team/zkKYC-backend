# zkKYC Backend

Серверная часть системы управления zkKYC (Zero-Knowledge Know Your Customer) на Go с использованием:
- 🚀 [Chi Router](https://go-chi.io)
- 🐘 PostgreSQL
- 🔑 Ethereum/DID интеграция
- 📜 Swagger документация

## Особенности

- REST API для управления пользователями KYC

## Установка и запуск

### Требования
- Go 1.19+
- PostgreSQL 14+
- Docker (опционально)

### Быстрый старт
```bash
# Клонировать репозиторий
git clone https://github.com/Hashtag-team/zkKYC-backend.git
cd zkkyc-backend

# Создать файл конфигурации
cp .env.example .env

# Запустить в Docker
docker-compose up --build