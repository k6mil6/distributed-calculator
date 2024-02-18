# По всем вопросам tg - @kamil_66

Для проверки задания потребуется установленные 

docker (https://docs.docker.com/get-docker/) и git (https://git-scm.com/downloads)
postman - для удобства отправки запросов, по желанию (https://www.postman.com/downloads/)

Для развертывания приложения необходимо:

1. клонировать репозиторий в удобную папку (git clone https://github.com/k6mil6/distributed-calculator.git .)
2. перейти в папку 
3. создать файл config.hcl (файл конфигурации), ниже перечислил что можно указать (можно просто скопировать)
```
goroutine_number=5 #изменяет кол-во горутин(воркеров) на агента(количество агентов можно изменить в файле docker-compose.yml)
heartbeat_timeout=10s #время, раз в которое отправляется хартбит для проверки состояния воркера
worker_timeout=30s #время раз в которое воркеры делают запрос на свободные задания
fetcher_timeout=10s #время, раз в которое выражения делятся на подвыражения
```
4. прописать docker-compose up
5. ниже представлены примеры для проверки работоспособности

пример вычисления выражения с передачей таймаутов (mac/linux)
```
curl --location 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "id": "3422b448-2460-4fd2-9183-8000de6f8348",
    "expression": "2+2+3",
    "timeouts": {
        "+": 10,
        "-": 20,
        "/": 10,
        "*": 5
    }
}'
```
windows (powershell)
```
$body = @{
    id = "3422b448-2460-4fd2-9183-8000de6f8348"
    expression = "2+2+3"
    timeouts = @{
        "+" = 10
        "-" = 20
        "/" = 10
        "*" = 5
    }
} | ConvertTo-Json -Compress

$response = Invoke-WebRequest -Uri 'http://localhost:8080/calculate' -Method Post -ContentType 'application/json' -Body $body
Write-Output $response.Content
```

timeouts - параметр, который можно не указывать, будет использовано последнее добавленное значение
также важно, чтобы id был формата uuid, для проверки нескольких выражений можно изменять последние цифры самого id

пример вычисления выражения без передачи таймаутов(mac/linux)
```
curl --location 'http://localhost:8080/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "id": "3422b448-2460-4fd2-9183-8000de6f8348",
    "expression": "2+2+3"
}'
```
windows(powershell)
```
$body = @{
    id = "3422b448-2460-4fd2-9183-8000de6f8348"
    expression = "2+2+3"
} | ConvertTo-Json -Compress

$response = Invoke-WebRequest -Uri 'http://localhost:8080/calculate' -Method Post -ContentType 'application/json' -Body $body
Write-Output $response.Content
```


пример получения выражения по id (macos/linux)

```
curl --location 'http://localhost:8080/expression/3422b448-2460-4fd2-9183-8000de6f8346'
```

windows(powershell)
```
$response = Invoke-WebRequest -Uri 'http://localhost:8080/expression/3422b448-2460-4fd2-9183-8000de6f8346' -Method Get
$response.Content
```

пример получения всех выражений (macos/linux)

```
curl --location 'http://localhost:8080/all_expressions'
```

windows(powershell)
```
$response = Invoke-WebRequest -Uri 'http://localhost:8080/all_expressions' -Method Get
$response.Content
```



