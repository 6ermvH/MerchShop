### TODO

## Баги
- [FIXED] Не получается купить предмет `/api/buy/{item}` возврат `501 | db error`
- Пользователь может создать бесконечное кол-во корректных `JWT`
- [FIXED] `/api/info` не выводит поля, которые `nil` (мб не баг)
- [FIXED] `api` не отдаёт ошибки
- [FIXED] `/api/info` даёт неверную информацию о транзакциях

## Доделать

```bash
commit 44c258f2e154907c0d43bcfb5cd38a650521e109
Author: German Feskov <g.feskov@yandex.ru>
Date:   Mon Nov 10 12:03:33 2025 +0300

    Add: middleware auth test
```

- Переделать `package repo`, т.к тяжело генерировать моки [+]
- добавить Debug и Info логгирование []
- переписать точку входа []
- Реализовать unit-tests [+]
- Реализовать E2E-tests []

```bash
commit e25d22a6a0dfc4b52d24b9e26c503bb45d4c7670 (HEAD -> main, github/main, github/HEAD)
Author: German Feskov <g.feskov@yandex.ru>
Date:   Mon Nov 10 20:11:01 2025 +0300

    Rm: coverage files
```

- Переписать Makefile []
- Отрефакторить тесты на линтерах []
- Поправить линтеры []  