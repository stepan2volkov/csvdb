[![Go version](https://img.shields.io/github/go-mod/go-version/stepan2volkov/csvdb.svg)](https://github.com/stepan2volkov/csvdb/blob/main/go.mod)
![CI](https://github.com/stepan2volkov/csvdb/actions/workflows/ci.yaml/badge.svg)
[![codecov](https://codecov.io/gh/stepan2volkov/csvdb/branch/main/graph/badge.svg?token=CP0CR6QKOE)](https://codecov.io/gh/stepan2volkov/csvdb)


# CSV DB

## Описание

__Цель программы__:  использовать sql-like синтаксис для работы с csv-файлами.

__Основные операции:__ `AND`, `OR`, `=`, `<`, `>`.

__Пример запроса__: 
```sql
(salary > 150200.99 AND status='working') OR (age > 18 AND age < 40);
```

__"Особенности":__
1. В конце запроса в обязательном порядке должна стоять `;`
2. Оператор `AND` имеет приоритет над оператором `OR`

## Задачи на разработку

- [x] \(CSVDB-1) Реализовать простейший парсер выражения `where`
- [x] \(CSVDB-2) Подключить линтер и устранить его замечания
- [x] \(CSVDB-3) Настроить pre-commit хук
- [x] \(CSVDB-4) Настроить github actions
- [x] \(CSVDB-5) Описать структуру хранения данных, полученных из csv-файла
- [x] \(CSVDB-6) Описать доступные операции над данными
- [ ] \(CSVDB-7) Реализовать построение плана запроса из разобранного выражения `where`

## Локальная настройка окружения

После скачивания репозитория требуется подключить pre-commit хук, чтобы выявлять проблемы до заливки кода к репозиторий.
Требуется выполнить:
```
chmod +x githooks/pre-commit
cp githooks/pre-commit .git/hooks
```