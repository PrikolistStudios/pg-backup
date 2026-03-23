# pg-backup

Инструмент командной строки для легкого удаления и резервного копирования баз данных PostgreSQL. Разработан в рамках тестового задания для трудоустройства.

## Установка

Загрузить собранную программу для Linux можно в Releases.

## Использование

```commandline
pg-backup --help
```

```text
This CLI tool is used for easy removal and backup of PostgreSQL databases.
It supports action on multiple databases in one command and globbing

Each argument is considered a glob pattern, and the user must have necessary rights to perform an action on databases which
names match it (read to back up, drop to remove). Be aware that glob gets expanded to commandline arguments before the program is launched, so running:

pg-backup test_db* 

might lead to unexpected results if you have matching files in working directory.

When removing the database, it is not backed up implicitly. Database backups are stored in working directory and made with pg-dump.

Glob pattern syntax is as follows:

pattern:
    { term }

term:
    `*`         matches any sequence of non-separator characters
    `**`        matches any sequence of characters
    `?`         matches any single non-separator character
    `[` [ `!` ] { character-range } `]`
                character class (must be non-empty)
    `{` pattern-list `}`
                pattern alternatives
    c           matches character c (c != `*`, `**`, `?`, `\`, `[`, `{`, `}`)
    `\` c       matches character c

character-range:
    c           matches character c (c != `\\`, `-`, `]`)
    `\` c       matches character c
    lo `-` hi   matches character c for lo <= c <= hi

pattern-list:
    pattern { `,` pattern }
                comma-separated (without spaces) patterns
```

## Пример

```text
pg-backup test_db_* -P 1234 -U my-user -m backup
? Perform the action (backup) on the following databases?
test_db_1
test_db_2
test_db_3
 Yes
Performing backup
```

После чего в рабочей директории должны появиться файлы test_db_[1|2|3].backup.

## Зависимости

Для использования требуется наличие PostgreSQL, и также pg-dump для создания бэкапов.

Для сборки требуется только Go

Для тестирования понадобится Docker, потому что тесты поднимают временные контейнеры.

## Проверка

Для удобства проверки я также опишу API пакетов, чтобы было легче вникнуть в код

### [prog/main.go](prog/main.go)

Запуск приложения [Cobra](https://github.com/spf13/cobra).

### [internal/cmd/root.go](internal/cmd/root.go):

Описания и флагов команды, обработка аргументов командной строки.

### [internal/app](internal/app):

Основная логика приложения. В целом она такая: по списку паттернов получается через `FilterPatterns` список доступных
бд. Потом создается "действие" `DatabaseAction`, которое будет осуществляться `PerformDatabasesAction` над подошедшими
под паттерны бд. Все действия производятся параллельно через горутины.

- [app.go](internal/app/app.go): `PerformDatabasesAction` осуществляет действие над группой баз данных (удаление или
  резервное копирование)
- [remove.go](internal/app/remove.go): `NewRemoveAction` создает `DatabaseAction`, которое удалит названную базу данных.
- [backup.go](internal/app/backup.go): `NewBackupAction` создает `DatabaseAction`, которое создаст резервную копию
  переданной базы данных.
- [connection.go](internal/app/connection.go): `CreateConnection` создает соединение с бд, которое понадобится для
  получения списка всех баз данных, а также для удаления баз данных.
- [filter.go](internal/app/filter.go): `FilterPatterns` по списку паттернов возвращает список подошедших бд.
- [config.go](internal/app/config.go) и [errors.go](internal/app/errors.go) объявляют вспомогательные структуры.
