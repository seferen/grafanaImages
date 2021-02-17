# grafanaImages
## Оглавление
* [Подготовка](#Подготовка)
* [Настройка](#Настройка)
## Подготовка
Данная утилита предназанчена для работы с Grafana и не работает без плагина Grafana renderImage, обязательно проверьте наличие данного плагина перед использованием утилиты
Для информации по утановке плагина на Grafana см. https://grafana.com/grafana/plugins/grafana-image-renderer
## Настройка
Для работы приложения требуется конфигурационный файл config.json, по умолчанию данный файл необходимо разместить в корневом каталоге приложения
```json
{
    "url": "http://grafanahost:3000",
    "token": "TOKEN",
    "test": {
        "timeStart": "2021-01-01 00:00:00",
        "timeEnd": "2021-01-23 23:59:59"
    },
    "dashboards": {
        "nameOfDashbord": [
            {
                "nameOfParametr1": "volume1",
                "nameOfParametr1": "volume2"
            }
        ]
    }
}

// Если параметры dashbord должны использоваться по умолчанию

{
    "url": "http://grafanahost:3000",
    "token": "TOKEN",
    "test": {
        "timeStart": "2021-01-01 00:00:00",
        "timeEnd": "2021-01-23 23:59:59"
    },
    "dashboards": {
        "nameOfDashbord": []
    }
}
```
* url - адрес расположения grafana
* token - TOKEN для доступа к grafana
* test - начало и завершения теста, timeStart - начало теста, timeEnd - завершение теста, время указывается в формате "YYYY-MM-DD hh:mm:ss"
* dashboards - хранит наименования dashboards 
Для использования друго конфигурационного файла используется флаг -fileConfig
```Shell
-fileConfib anotherConfigFile.json
```

