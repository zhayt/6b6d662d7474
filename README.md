# Описание проекта

Данный проект представляет собой web-сервис на Golang & gorilla/mux, который позволяет получать и сохранять данные из публичного API национального банка и хранить их в локальной базе данных MS SQL Server. Таже доставать данные из базы данных по дате и по куду валюты.

# Прошу обратить внимание
### У меня не получилось локально развернуть MS SQL Server, потому что у меня стоит Kali Linux на машине, а MS SQL Server доступен только для Windows и Linux (Ubuntu, RedHat и SUSE) [Подробнее можно прочитать по этой ссылке](https://superuser.com/questions/1655788/how-to-install-mssql-server-on-kali-linux). Поэтому я запустил MS SQL Server в контейнере. Также мне пришлось изменить тип данных поля Title и Code с VARCHAR на NVARCHAR, так как данные, содержащие кириллические символы, заменялись на вопросительные знаки (?) при добавлении в базу данных.

# Как запустить
## Предварительные требования
- Установленный Docker
- Установленный Docker Compose

## Запуск приложения

1. Клонируйте репозиторий проекта на свой компьютер:

``` shell
git clone https://github.com/zhayt/6b6d662d7474
```

2. Перейдите в директорию проекта:

```shell
cd 6b6d662d7474
```

3. Запустите приложение с помощью Docker Compose:

```shell
docker-compose up --build
```

4. Чтобы остановить приложение, выполните команду:

``` shell
docker-compose down
```

Swagger документация доступна по адресу http://localhost:8080/swagger/
Также есть .http файл для тестовых запросов на эндпойнты