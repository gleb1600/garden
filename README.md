# Garden Simulator: — консольное приложение для симуляции ухода за растениями. Цифровой сад на PostgreSQL

## Особенности проекта
- Динамическая модель роста растений с учетом времени
- Работа с PostgreSQL через pgx
- CRUD операции
- Автоматическое ухудшение характеристик растений (жажда, голод)
- Интерактивное управление через командную строку
- Визуализацию состояния в реальном времени

## Работа

1. **Установите PostgreSQL локально**:

https://www.postgresql.org/download/

2. **Откройте "SQL Shell (psql) и введите пароль от postgres"**
```psql
Server [localhost]: 
Database [postgres]: 
Port [5432]: 
Username [postgres]: 
Пароль: ваш пароль от postgres
```
3. **В psql создайте пользователя 'garden' с паролем 'secret' и БД 'gardendb'**:
```sql
CREATE USER garden WITH PASSWORD 'secret';
CREATE DATABASE garden OWNER garden;
\q
```
4. **Откройте PowerShell и выполните миграцию и инициализацию данных**:
```bash
& 'C:\Program Files\PostgreSQL\17\bin\psql.exe' -U garden -d gardendb -f migrations\init.sql

& 'C:\Program Files\PostgreSQL\17\bin\psql.exe' -U garden -d gardendb -f scripts\seed.sql
```
5. **Запуск**:
```bash
go run main.go
```
6. **Управляйте садом**:
```bash
Доступные команды:
  list                 - Показать все растения
  add [species] [name] - Посадить новое растение
  remove [id]          - Удалить растение
  water [id]           - Полить растение
  fertilize [id]       - Удобрить растение
  help                 - Показать помощь
  exit                 - Выйти из программы
```

## Структура проекта
```bash
garden
    ├── migrations         # SQL-миграции
    │      └── init.sql    # Миграция
    ├── plants             # Логика работы с структурами
    │      └── plantsfunc  # Работа с растением
    ├── scripts            # Скрипты инициализации
    │      └── seed.sql    # Инициализация
    ├── storage            # Логика хранения данных
    │      └── dbfunc      # Работа с БД
    ├── main.go            # Точка входа
    └── go.mod             # Зависимости
```
