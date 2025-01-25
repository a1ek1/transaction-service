# Сервис по работе с транзакциями

HTTP сервер, реализующий REST API, выполненный в рамках тестового задания на стажировку для компании _"Infotecs"_.


## Подготовка к запуску сервера

---

### Необходимо, чтобы на вашем устройстве были установлены

- Docker
- Docker compose
- goose (для применения миграций)
- Go 1.23

## Запуск сервера

1. Склонируйте проект

```bash
    git clone git@github.com:a1ek1/transaction-service.git
    cd transaction-service
```

2. Запустите контейнер с базой данных

```bash
    # Команду выполняем из директории transaction-service
    docker-compose up
```

После выполнения данной команды у вас должен запуститься контейнер с PostgreSQL, достуный на порту **5434**

3. Примените миграции

Если у вас не установлен goose, то нужно в командной строке выполнить команду
```bash
    go install github.com/pressly/goose/v3/cmd/goose@latest
```

После установки проверьте корректность работы
```bash
    goose --version
```
Дальше нужно _перейти в директорию transaction-service/migrations_ и выполнить следующую команду

```bash
    # Указаны значения по умолчанию. Если поменяете в конфиге и в docker-compose.yml, то здесь будут другие данные
    goose postgres "host=localhost user=postgres port=5434 password=postgres database=transaction_service sslmode=disable" up
```

Вам должно вывестить сообщение об успешном применении миграций

4. Запустите файл main.go

```bash
    # Из директории transaction-service
    go run cmd/main.go
```
После выполнения этой команды сервер будет доступен на **_localhost:8080_**
## Тестирование работы

1. Перевод средств с одного счета на другой

```bash
    curl -X POST http://localhost:8080/api/send \
          -H "Content-Type: application/json" \
          -d '{
            "from": "{номер_кошелька_отправителя}",
            "to": "{номер_кошелька_получателя}",
            "amount": {сумма_перевода}
          }'
```
2. Просмотр баланса кошелька

```bash
    # Введите номер нужного кошелька
    curl -X GET http://localhost:8080/api/wallet/{номер_кошелька}/balance
```

3. Получение N последних транзаций (в данной реализации учтены только выполненные транзакции)

```bash
    # Введите число нужных транзакций
    curl -X GET "http://localhost:8080/api/transactions?count=N"
```

4. Получение всех кошельков

```bash
    curl -X GET http://localhost:8080/api/wallets
```