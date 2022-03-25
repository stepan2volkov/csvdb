[![Go version](https://img.shields.io/github/go-mod/go-version/stepan2volkov/csvdb.svg)](https://github.com/stepan2volkov/csvdb/blob/main/go.mod)
![CI](https://github.com/stepan2volkov/csvdb/actions/workflows/ci.yaml/badge.svg)
[![codecov](https://codecov.io/gh/stepan2volkov/csvdb/branch/main/graph/badge.svg?token=CP0CR6QKOE)](https://codecov.io/gh/stepan2volkov/csvdb)


# CSV DB

## Описание

__Цель программы__:  использовать sql-like синтаксис для работы с csv-файлами.

__Основные операции:__ `AND`, `OR`, `=`, `<`, `>`.

__Пример запроса__: 
```sql
SELECT region, country, item_type, sales_channel, total_cost, total_profit FROM sales WHERE country = 'South Africa' AND item_type = 'Clothes' and sales_channel='Online' AND total_profit > 400000;
```

__"Особенности":__
1. В конце запроса в обязательном порядке должна стоять `;`
2. Оператор `AND` имеет приоритет над оператором `OR`

## Использование

Загрузка csv-файла
```
\load file.csv config.yaml
```

Формат yaml-файла
```yaml
name: tablename         # Наименование таблицы
sep: ','                # Разделитель значений
lazyQuotes: true        # true, если значения заключены в двойные кавычки
fields:                 # Список полей в таблице
- name: lastname        # Наименование поля
  type: string          # Тип поля: string или number
- name: firstname
  type: string
- name: salary
  type: number
```

Список загруженных таблиц
```
\list
```

Удаление таблицы
```
 \drop tablename
 ```
## Задачи на разработку

- [x] \(CSVDB-1) Реализовать простейший парсер выражения `where`
- [x] \(CSVDB-2) Подключить линтер и устранить его замечания
- [x] \(CSVDB-3) Настроить pre-commit хук
- [x] \(CSVDB-4) Настроить github actions
- [x] \(CSVDB-5) Описать структуру хранения данных, полученных из csv-файла
- [x] \(CSVDB-6) Описать доступные операции над данными
- [x] \(CSVDB-7) Реализовать построение плана запроса из разобранного выражения `where`
- [x] \(CSVDB-8) Реализовать загрузку файла в RAM
- [x] \(CSVDB-9) Реализовать конфигурирование посредством yaml-файла
- [x] \(CSVDB-10) Реализовать секцию `from`
- [x] \(CSVDB-11) Реализовать секцию `select`
- [x] \(CSVDB-11) Перенести логирование в файлы `access.log`  и `error.log`

## Локальная настройка окружения

После скачивания репозитория требуется подключить pre-commit хук, чтобы выявлять проблемы до заливки кода к репозиторий.
Требуется выполнить:
```
chmod +x githooks/pre-commit
cp githooks/pre-commit .git/hooks
```