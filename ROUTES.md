
## `ROUTES.md`

```md
# ROUTES.md
HTTP API Routes – v1

Base URL: `/v1`

Security schemes:

- **UserAuth** – Bearer JWT for end users (students).
- **AdminAuth** – Bearer JWT for admins/super-admins.

This file lists the primary routes the backend exposes.  
Details (schemas, examples) live in `openapi.yaml`.

---

## 1. Auth

### 1.1 User Auth

```http
POST /v1/auth/telegram
````

* **Description:** Login or signup via Telegram.
* **Auth:** public (returns UserAuth token)
* **Body:** `{ telegramId, displayName, exam }`
* **Response:** `{ accessToken, user }`

---

## 2. App: User

### 2.1 Get current user & exam profile

```http
GET /v1/me
```

* **Auth:** UserAuth
* **Response:** user info + exam profile (streak, level, totals)

### 2.2 List subjects

```http
GET /v1/subjects?exam={examCategory?}
```

* **Auth:** UserAuth
* **Description:** Subjects for given or primary exam.

### 2.3 List topics

```http
GET /v1/topics?subjectId={uuid}
```

* **Auth:** UserAuth
* **Description:** Topics under a subject.

---

## 3. App: Practice

### 3.1 Create practice session

```http
POST /v1/practice/sessions
```

* **Auth:** UserAuth
* **Body:** mode, exam, subjectIds, topicIds, difficultyLevels, numQuestions, timeLimitMinutes

### 3.2 List user practice sessions

```http
GET /v1/practice/sessions
```

* **Auth:** UserAuth

### 3.3 Get session detail (with questions)

```http
GET /v1/practice/sessions/{id}
```

* **Auth:** UserAuth

### 3.4 Submit answer

```http
POST /v1/practice/sessions/{id}/answers
```

* **Auth:** UserAuth
* **Body:** `{ sessionQuestionId, selectedOption, timeTakenMs }`
* **Response:** Updated question state (isCorrect, etc.)

---

## 4. App: Revision (SRS)

### 4.1 Get revision queue

```http
GET /v1/revision/queue
```

* **Auth:** UserAuth
* **Description:** Items due for revision.

> (Later we can add `/v1/revision/sessions` if needed.)

---

## 5. App: Exams & Events

### 5.1 List available exams/events

```http
GET /v1/events
```

* **Auth:** UserAuth
* **Description:** Exams this user can see (mock tests, reward events, etc.)

> Future:
>
> * `POST /v1/events/{id}/register`
> * `POST /v1/events/{id}/attempts` (start/submit)

---

## 6. App: Podcasts

### 6.1 List podcast episodes

```http
GET /v1/podcasts?subjectId={uuid?}&topicId={uuid?}
```

* **Auth:** UserAuth
* **Description:** Podcasts filtered by exam/subject/topic.

### 6.2 Get one podcast episode

```http
GET /v1/podcasts/{id}
```

* **Auth:** UserAuth

---

## 7. App: Wallet, Coupons, Referral

### 7.1 Wallet summary

```http
GET /v1/wallet
```

* **Auth:** UserAuth
* **Description:** Balance, lifetime earned & spent.

### 7.2 Wallet transactions

```http
GET /v1/wallet/transactions
```

* **Auth:** UserAuth

### 7.3 Redeem coupon

```http
POST /v1/coupons/redeem
```

* **Auth:** UserAuth
* **Body:** `{ code }`
* **Response:** Updated wallet summary.

### 7.4 Referral summary

```http
GET /v1/referral
```

* **Auth:** UserAuth
* **Description:** Referral stats & earnings.

---

## 8. Admin: Questions

Base path: `/v1/admin/questions`

### 8.1 List questions

```http
GET /v1/admin/questions?exam={exam}&subjectId={uuid?}&topicId={uuid?}
```

* **Auth:** AdminAuth

### 8.2 Create question

```http
POST /v1/admin/questions
```

* **Auth:** AdminAuth
* **Body:** QuestionCreateRequest

### 8.3 Get question

```http
GET /v1/admin/questions/{id}
```

* **Auth:** AdminAuth

### 8.4 Update question

```http
PATCH /v1/admin/questions/{id}
```

* **Auth:** AdminAuth

### 8.5 Delete question

```http
DELETE /v1/admin/questions/{id}
```

* **Auth:** AdminAuth

---

## 9. Admin: Subjects & Topics

### 9.1 Subjects

```http
GET  /v1/admin/subjects
POST /v1/admin/subjects
```

* **Auth:** AdminAuth
* **GET:** list all subjects.
* **POST:** create subject `{ exam, name }`.

### 9.2 Topics

```http
GET  /v1/admin/topics?subjectId={uuid?}
POST /v1/admin/topics
```

* **Auth:** AdminAuth
* **POST:** `{ subjectId, name }`.

---

## 10. Admin: Exams & Events

Base path: `/v1/admin/exams`

```http
GET    /v1/admin/exams
POST   /v1/admin/exams
GET    /v1/admin/exams/{id}
PATCH  /v1/admin/exams/{id}
DELETE /v1/admin/exams/{id}
```

* **Auth:** AdminAuth
* Manage `exam_config` records.

---

## 11. Admin: Podcasts

```http
GET    /v1/admin/podcasts
POST   /v1/admin/podcasts
GET    /v1/admin/podcasts/{id}
PATCH  /v1/admin/podcasts/{id}
DELETE /v1/admin/podcasts/{id}
```

* **Auth:** AdminAuth

---

## 12. Admin: Coupons

```http
GET    /v1/admin/coupons
POST   /v1/admin/coupons
GET    /v1/admin/coupons/{id}
PATCH  /v1/admin/coupons/{id}
DELETE /v1/admin/coupons/{id}
```

* **Auth:** AdminAuth

---

## 13. Admin: AI Settings

```http
GET  /v1/admin/ai-settings
PUT  /v1/admin/ai-settings
```

* **Auth:** AdminAuth
* Manages:

  * weaknessMinAttempts
  * weaknessThresholdPercent
  * strongThresholdPercent
  * revisionIntervalsDays
  * includeGuessedCorrect
  * revisionEnabled

---

## 14. Admin: Analytics (Overview)

```http
GET /v1/admin/analytics/overview?exam={exam?}&range={today|7d|30d}
```

* **Auth:** AdminAuth
* Returns high-level metrics for the dashboard (users, activity, accuracy, rewards, etc.)

````
