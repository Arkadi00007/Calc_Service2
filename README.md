# Сервис подсчёта арифметических выражений 2.0

## Описание

Этот проект реализует распределённый калькулятор с оркестратором и агентами. Оркестратор принимает математические выражения, разбивает их на операции и раздаёт агентам для вычисления. Сервис поддерживает JWT-аутентификацию и хранение данных о выражениях в базе данных.
Перед использованием пользователь должен пройти регистрацию и авторизацию.
---

## Как запустить
### Установка зависимостей
````
go mod tidy
````

(запуск агента и оркестратора)
````bash
cd Calc_Service2
````
````bash
go run cmd/main.go
````
По умолчанию порт - ":8080"

---


###  Как пользоваться:
(все адреса начинаются с "localhost:8080")
1. Регистрация пользователя.
- Отправьте POST-запрос на эндпоинт /api/v1/register для создания нового аккаунта.
2. Авторизация
- Получите JWT-токен, выполнив POST-запрос на /api/v1/login с вашими учётными данными 
3. Ввод выражения для вычисления 
- Отправьте математическое выражение на обработку с помощью POST-запроса на /api/v1/calculate
4. Пользователь может получить информацию о выражении сделав GET запрос по этому эндпоинту - /api/v1/expressions или получить сведения о конкретном выражении по его id  - /api/v1/expressions/:id
- Результат вычисления или ошибка отображаются в консоли.
---

### Оркестратор и агент взаимодействуют по gRPC

### Эндпоинты
- /api/v1/register - Регистрация нового пользователя({ "login": , "password": })
- POST /api/v1/login - Авторизация пользователя.При успешном входе возвращает JWT-токен, который используется для доступа к защищённым маршрутам API.Токен необходимо указывать в заголовке Authorization: Bearer <ваш_токен> при последующих запросах
- /api/v1/calculate - сюда пользователь вводит(POST запрос) арифметическое выражение и получает id(в теле запроса должен быть JWT токен) 
- /api/v1/expressions - по этому эндпоинту пользователь может сделать  GET запрос и получить данные о всех ранее введённых выражениях(в теле запроса должен быть JWT токен)
- /api/v1/expressions/:id - GET запрос (можно узнать данные только определённого выражения)(в теле запроса должен быть JWT токен)
---



#### Пример тела запроса (/api/v1/calculate):
```
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
"expression": <строка с выражение>
}'
```
```json
{
  "expression": "2+2"
}
```
---


# Примеры запросов:

### Регистрация (http://localhost:8080/api/v1/register)

- Запрос:
```
curl -X POST http://localhost:8080/api/v1/register \
-H "Content-Type: application/json" \
-d '{"username": "testuser", "password": "securepassword"}'
```

- Ответ при успешной регистрации (HTTP статус 201)
```
{
  "message": "Registered "
}
```
- Ответ, если такой пользователь уже существует(HTTP статус 409)
```
{
  "error": "Username already exists"
}
```
Ответ при ошибке валидации (HTTP статус 400):
```
{
  "error": "Invalid input"
}

```

---

### 2. Авторизация пользователя
- Запрос:
```
curl -X POST http://localhost:8080/api/v1/login \
-H "Content-Type: application/json" \
-d '{"username": user", "password": "secret"}'
```
- Ответ при успешной авторизации (HTTP статус 200) — возвращается JWT токен:
```
{
  "token": "your.jwt.token.here"
}
```
- Ответ при ошибке (неверные данные — неправильный логин или пароль)  (HTTP статус 401):
```
{
  "error": "Invalid credentials"
}
```


---
### 3. Запрос на http://localhost:8080/api/v1/calculate (код 201)
- Запрос(в теле запроса укажите ваш JWT токен)
```
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-H "Authorization: Bearer ваш_jwt_токен" \
-d '{"expression":"6 + 3 * 5"}'
```

- Ответ:
```
{
"id": 1
}
```
---
- Запрос (код 422)

```
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-H "Authorization: Bearer ваш_jwt_токен" \
-d '{"expression":"6 + r * 5"}'

```

- Ответ:
```
Invalid expression
```

---
- Запрос (код 500)

```
curl -X GET http://localhost:8080/api/v1/expressions \
-H "Content-Type: application/json" \
-H "Authorization: Bearer ваш_jwt_токен" \
-d '{"expression":"6 + 5"}'
```

- Ответ:
```
Method not allowed
```
---
### 4. Запрос на http://localhost:8080/api/v1/expressions (код 200)
### Объязательно укажите ваш JWT токен в теле запроса
Например, делаем запрос
```
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-H "Authorization: Bearer ваш_jwt_токен" \
-d '{"expression": "6 + 3 * 5"}'
``` 
потом получаем выражения:
```
curl -X GET http://localhost:8080/api/v1/expressions \
-H "Authorization: Bearer ваш_jwt_токен"
```

- Ответ:
```
{
    "expressions": [{
            "id": 3,
            "result": 21,
            "status": "completed"
        }]
}
```
---
- Запрос (код 405)
````
curl -X POST http://localhost:8080/api/v1/expressions \
-H "Content-Type: application/json"
````

- Ответ:
```
 Method not allowed
```

---

### 3.Запрос на http://localhost:8080/api/v1/expressions/{id} (код 200)
(делаем всё тот же запрос:

```
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-H "Authorization: Bearer ваш_jwt_токен" \
-d '{"expression": "6 + 3 * 5"}'

``` 
)

- Получаем выражение по id:
```
curl -X GET http://localhost:8080/api/v1/expressions/1 \
-H "Authorization: Bearer ваш_jwt_токен"

```

- Ответ:
```
{
    "id": 3,
    "result": 21,
    "status": "completed"
}
```
