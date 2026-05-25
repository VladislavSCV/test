# w — Еда.Быстро

Сгенерировано [ekz](https://github.com/VladislavSCV/ekz). Конфигурация: `project.yaml`.

> **Коммиты ДЭ:** минимум 3 фиксации (БД → UI → доработка).

## Предметная область

Доставка еды на дом

## Таблицы БД

| Таблица | Назначение |
|---------|------------|
| `users` | Пользователи (регистрация, админ) |
| `food_orders` | Заказ |
| `reviews` | Отзывы к записям |

### Поля `food_orders`

| Столбец | Тип | Подпись |
|---------|-----|---------|
| id, user_id, created_at | системные | — |
| `address` | string | Адрес доставки |
| `dish` | enum | Блюдо |
| `delivery_time` | date | Время доставки |
| `total` | int | Сумма, ₽ |
| status | string | Статус |

## Запуск

```bash
cd backend && go mod tidy && go run .
cd frontend && npm install && npm run dev
```

Админ: **Admin26** / **Demo20**
