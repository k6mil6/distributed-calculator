# По всем вопросам tg - @kamil_66

Диаграмма, как тут всё работает
https://excalidraw.com/#json=gWoBa-VziAJJz358b8HIG,b2P-xXy_aJ3oYZj25ndl1Q

Для проверки задания потребуется установленные 

docker (https://docs.docker.com/get-docker/) и git (https://git-scm.com/downloads)
postman - для удобства отправки запросов, по желанию (https://www.postman.com/downloads/)

Для развертывания приложения необходимо:

1. клонировать репозиторий в удобную папку (git clone https://github.com/k6mil6/distributed-calculator.git .)
2. перейти в папку
3. в консоле прописать docker-compose up (если оркестратор не запустился, необходимо перезапустить его, либо через интерфейс docker desktop, либо нажать ctrl+c и написать docker-compose еще раз)
4. ниже представлены примеры для проверки работоспособности

пример вычисления выражения с передачей таймаутов (mac/linux)
```
curl --location 'http://localhost:5441/calculate' \
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

$response = Invoke-WebRequest -Uri 'http://localhost:5441/calculate' -Method Post -ContentType 'application/json' -Body $body
Write-Output $response.Content
```

timeouts - параметр, который можно не указывать, будет использовано последнее добавленное значение
также важно, чтобы id был формата uuid, для проверки нескольких выражений можно изменять последние цифры самого id

пример вычисления выражения без передачи таймаутов(mac/linux)
```
curl --location 'http://localhost:5441/calculate' \
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

$response = Invoke-WebRequest -Uri 'http://localhost:5441/calculate' -Method Post -ContentType 'application/json' -Body $body
Write-Output $response.Content
```


пример получения выражения по id (macos/linux)

```
curl --location 'http://localhost:5441/expression/3422b448-2460-4fd2-9183-8000de6f8346'
```

windows(powershell)
```
$response = Invoke-WebRequest -Uri 'http://localhost:5441/expression/3422b448-2460-4fd2-9183-8000de6f8346' -Method Get
$response.Content
```

пример получения всех выражений (macos/linux)

```
curl --location 'http://localhost:5441/all_expressions'
```

windows(powershell)
```
$response = Invoke-WebRequest -Uri 'http://localhost:5441/all_expressions' -Method Get
$response.Content
```

пример установки таймаутов для операций (macos/linux)

```
curl --location 'http://localhost:5441/set_timeouts' \
--header 'Content-Type: application/json' \
--data '{
"timeouts": {
        "+": 10, 
        "-": 10,
        "*": 10,
        "/": 10
    }
}'
```

windows(powershell)
```
$body = @{
    timeouts = @{
        "+" = 10
        "-" = 20
        "/" = 10
        "*" = 5
    }
} | ConvertTo-Json -Compress

$response = Invoke-WebRequest -Uri 'http://localhost:5441/set_timeouts' -Method Post -ContentType 'application/json' -Body $body
Write-Output $response.Content
```

пример получения актуальных таймаутов (macos/linux)

```
curl --location 'http://localhost:5441/actual_timeouts'
```

windows(powershell)
```
$response = Invoke-WebRequest -Uri 'http://localhost:5441/actual_timeouts' -Method Get
$response.Content
```




