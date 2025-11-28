
#  API Endpoints для управления задачами

## Сводная таблица методов

| Метод | URL | Описание |
|-------|-----|-----------|
| GET | `/api/todos` | Получить все задачи |
| POST | `/api/todos` | Создать новую задачу |
| PATCH | `/api/todos/:id` | Пометить задачу выполненной |
| DELETE | `/api/todos/:id` | Удалить задачу |

##  Детали API

### GET /api/todos
**Описание:** Получить все задачи

**Response:**
```json
[
  {
    "id": 1,
    "body": "Изучить Go",
    "completed": false
  }
]
```

---

### POST /api/todos
**Описание:** Создать новую задачу

**Request:**
```json
{
  "body": "Новая задача"
}
```

**Response:**
```json
{
  "id": 1,
  "body": "Новая задача",
  "completed": false
}
```

---

### PATCH /api/todos/:id
**Описание:** Пометить задачу выполненной

**Path Parameters:**
- `id` - ID задачи

**Response:**
```json
{
  "id": 1,
  "body": "Новая задача",
  "completed": true
}
```

---

### DELETE /api/todos/:id
**Описание:** Удалить задачу

**Path Parameters:**
- `id` - ID задачи

**Response:**
```json
{
  "success": true
}
```