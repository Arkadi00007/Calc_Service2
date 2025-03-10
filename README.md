# Сервис подсчёта арифметических выражений 2.0

## Описание

Этот проект реализует распределённый калькулятор с оркестратором и агентами.Оркестратор принимает математические выражения, разбивает их на операции и раздаёт агентам для вычисления.

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
- Пользователь может вводить математические выражения через консоль по этому эндпоинту - /api/v1/calculate.
- Пользователь может получить информацию о выражении сделав GET запрос по этому эндпоинту - /api/v1/expressions или получить сведения о конкретном выражении по его id  - /api/v1/expressions/:id
- Результат вычисления или ошибка отображаются в консоли.
---

### Эндпоинты
- /api/v1/calculate - сюда пользователь вводит арифметическое выражение и получает id 
- /api/v1/expressions - сюда пользователь может сделать  GET запрос и получить данные о всех ранее введённых выражениях
- /api/v1/expressions/:id - GET запрос (можно узнать данные только определённого выражения)
- /internal/task - отсюда агент получает простые выражения и сюда же возвращает ответ
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

### 1.Запрос на http://localhost:8080/api/v1/calculate (код 201)
```
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression":"6 + 3 * 5"}'
```

Ответ:
```
{
"id": 1
}
```
---
Запрос (код 422)

```
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression":"6 + r * 5"}'
```

Ответ:
```
Invalid expression
```

---
Запрос (код 500)

```
curl -X GET http://localhost:8080/api/v1/expressions \
-H "Content-Type: application/json" \
-d '{"expression":"6 +  5"}'
```

Ответ:

```
Method not allowed
```

### 2.Запрос на http://localhost:8080/api/v1/expressions (код 200)

До этого, например, запрос
```
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression":"6 + 3 * 5"}'
``` 
потом
```
curl -X GET http://localhost:8080/api/v1/expressions \
-H "Content-Type: application/json"
```

Ответ:
```
{
    "expressions": [{
            "id": 3,
            "result": 21,
            "status": "completed"
        }]
}
```

Запрос (код 405)
````
curl -X POST http://localhost:8080/api/v1/expressions \
-H "Content-Type: application/json"
````

Ответ:
```
 Method not allowed
```



### 3.Запрос на http://localhost:8080/api/v1/expressions/{id} (код 200)
(до этого запроса делаем всё тот же запрос:

```
curl -X POST http://localhost:8080/api/v1/calculate \
-H "Content-Type: application/json" \
-d '{"expression":"6 + 3 * 5"}'
``` 
)

Потом
```
curl -X GET http://localhost:8080/api/v1/expressions/1 \
-H "Content-Type: application/json"
```

Ответ:
```
{
    "id": 3,
    "result": 21,
    "status": "completed"
}
```
