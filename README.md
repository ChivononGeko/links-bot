# Links Bot 🚀

Это проект, представляющий собой `Telegram-бота`, который генерирует одноразовые ссылки для регистрации. Бот управляется администраторами, генерирует уникальные токены, проверяет их статус и позволяет пользователям пройти регистрацию через веб-страницу.

В процессе работы создаются токены, которые шифруются перед отправкой пользователю. Бот также поддерживает запросы для получения списка использованных и неиспользованных токенов. Для хранения данных используется `SQLite`. Помимо работы бота, параллельно с ним запускается сервер, который обрабатывает регистрацию через веб-интерфейс и интегрирован с системой `Poster`.

## 📌 Функции

- Генерация уникальных токенов: Создание одноразовых ссылок для регистрации. 🔑

- Проверка токенов: Проверка введенного токена и получение информации о нем. ✅

- Список токенов: Получение списка использованных и неиспользованных токенов. 📜

- Регистрация пользователей: Пользователи проходят регистрацию через веб-форму по предоставленной ссылке. 📝

- Интеграция с Poster: Регистрация данных пользователей интегрирована с внешней системой Poster для дальнейшей обработки данных. 🔗

## 🚀 Установка

Клонируйте репозиторий:

```bash
git clone https://github.com/ChivononGeko/links-bot.git
cd links-bot
```

Создайте файл .env и настройте переменные:

```bash
BOT_TOKEN=your_telegram_bot_token
DB_PATH=path_to_your_database.db
POSTER_TOKEN=your_poster_api_token
BASE_URL=your_base_url
SERVER_PORT=8080
ADMINS=admin_telegram_ids
```

Установите зависимости:

```bash
go mod tidy
```

Соберите и запустите проект:

```bash
go run cmd/main.go
```

Проект автоматически запустит как бота, так и веб-сервер.

## 📝 Команды бота

- `/register` — Генерирует уникальную одноразовую ссылку для регистрации. 🔑

- `/check_token` — Проверяет статус токена (пользователь должен ввести токен после этой команды). 🔍

- `/used_tokens` — Получить список использованных токенов. 📜

- `/unused_tokens` — Получить список неиспользованных токенов. 📋

## 🌐 Веб-сервер

Параллельно с ботом запускается HTTP-сервер, который обрабатывает регистрацию через веб-страницу.

- `/register?token=your_token` — Страница регистрации, куда пользователь переходит по ссылке. 🌍

- `/submit` — Страница для отправки данных (имя, телефон, дата рождения) и регистрации. ✍️

## 🏗️ Архитектура

Проект использует гексагональную архитектуру:

- `Domain`: Логика работы с токенами и регистрацией. 🧠

- `Ports`: Интерфейсы для взаимодействия с сервисами (например, Poster API). 🔌

- `Adapters`: Реализация взаимодействия с внешними сервисами, такими как база данных и Poster API. ⚙️

- `Services`: Логика бизнес-правил для генерации ссылок и управления токенами. 💼

## 🐳 Docker

Для удобства развертывания проекта, имеется Dockerfile для создания контейнера с приложением. Чтобы собрать Docker-образ, используйте следующую команду:

```bash
docker build -t registration-token-bot .
```

Запустить контейнер:

```bash
docker run -p 8080:8080 registration-token-bot
```

Docker-контейнер будет работать как веб-сервер и Telegram-бот одновременно.

## 🔧 Используемые технологии

- `Go` — Язык программирования.

- `SQLite` — Система управления базами данных.

- `Telegram Bot API` — Для создания бота.

- `Poster API` — Для интеграции с внешней системой (Poster).

- `HTML, CSS, JS` — Для создания веб-страниц.

## ✍🏼 Автор

`Damir Usetov`
