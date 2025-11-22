## `SCHEMA.md`

````md
# SCHEMA.md
Database Schema – Multi-Exam Learning Platform

Primary DB: **PostgreSQL**

This document describes the logical schema used by the backend.  
Migrations should be created using `golang-migrate` according to this spec.

---

## 1. Enums

### 1.1 `exam_category`

```sql
CREATE TYPE exam_category AS ENUM ('NEET_PG', 'NEET_UG', 'JEE', 'UPSC');
````

Represents the four main exam categories supported by the platform.

### 1.2 `question_choice_type`

```sql
CREATE TYPE question_choice_type AS ENUM ('single', 'multi');
```

* `single` – single correct option
* `multi` – multiple correct options (future use)

### 1.3 `practice_mode`

```sql
CREATE TYPE practice_mode AS ENUM ('smart', 'custom', 'revision', 'exam');
```

* `smart` – AI-guided based on weakness
* `custom` – user-selected subject/topic/difficulty
* `revision` – spaced-repetition flow
* `exam` – practice using exam config

### 1.4 `practice_session_status`

```sql
CREATE TYPE practice_session_status AS ENUM ('in_progress', 'completed', 'abandoned');
```

### 1.5 `exam_type`

```sql
CREATE TYPE exam_type AS ENUM ('MOCK', 'SUBJECT_TEST', 'REWARD_EVENT', 'DAILY_TEST');
```

### 1.6 `exam_status`

```sql
CREATE TYPE exam_status AS ENUM ('DRAFT', 'SCHEDULED', 'ONGOING', 'COMPLETED');
```

### 1.7 `wallet_tx_type`

```sql
CREATE TYPE wallet_tx_type AS ENUM (
  'REWARD',
  'EXAM_ENTRY',
  'COUPON',
  'ADJUSTMENT',
  'REFERRAL',
  'SPIN',
  'BONUS'
);
```

### 1.8 `referral_status`

```sql
CREATE TYPE referral_status AS ENUM ('INVITED', 'JOINED', 'ACTIVATED');
```

### 1.9 `revision_result`

```sql
CREATE TYPE revision_result AS ENUM ('CORRECT', 'INCORRECT');
```

### 1.10 `user_role`

```sql
CREATE TYPE user_role AS ENUM ('USER', 'ADMIN', 'SUPER_ADMIN');
```

---

## 2. Lookup Tables

### 2.1 `exam_type_lookup`

Maps enum `exam_category` to DB IDs and metadata.

```sql
CREATE TABLE exam_type_lookup (
  id           SERIAL PRIMARY KEY,
  code         exam_category UNIQUE NOT NULL,
  name         TEXT NOT NULL,
  description  TEXT
);
```

---

## 3. Subjects & Topics

### 3.1 `subject`

```sql
CREATE TABLE subject (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  exam_type_id    INT NOT NULL REFERENCES exam_type_lookup(id),
  name            TEXT NOT NULL,
  is_active       BOOLEAN NOT NULL DEFAULT TRUE,
  UNIQUE (exam_type_id, name)
);
```

### 3.2 `topic`

```sql
CREATE TABLE topic (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  subject_id      UUID NOT NULL REFERENCES subject(id),
  name            TEXT NOT NULL,
  is_active       BOOLEAN NOT NULL DEFAULT TRUE,
  UNIQUE (subject_id, name)
);
```

---

## 4. Users & Profiles

### 4.1 `user`

```sql
CREATE TABLE "user" (
  id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  telegram_id        TEXT,
  email              TEXT,
  display_name       TEXT,
  primary_exam_type  exam_category NOT NULL,
  role               user_role NOT NULL DEFAULT 'USER',
  is_blocked         BOOLEAN NOT NULL DEFAULT FALSE,
  created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX user_telegram_id_idx ON "user"(telegram_id) WHERE telegram_id IS NOT NULL;
CREATE UNIQUE INDEX user_email_idx ON "user"(email) WHERE email IS NOT NULL;
```

### 4.2 `user_exam_profile`

Per-user, per-exam aggregated stats.

```sql
CREATE TABLE user_exam_profile (
  id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id             UUID NOT NULL REFERENCES "user"(id),
  exam_type_id        INT NOT NULL REFERENCES exam_type_lookup(id),
  total_questions     INT NOT NULL DEFAULT 0,
  total_correct       INT NOT NULL DEFAULT 0,
  total_time_seconds  INT NOT NULL DEFAULT 0,
  overall_level       INT NOT NULL DEFAULT 1, -- 1-5
  current_streak_days INT NOT NULL DEFAULT 0,
  longest_streak_days INT NOT NULL DEFAULT 0,
  last_login_at       TIMESTAMPTZ,
  CONSTRAINT user_exam_profile_uniq UNIQUE (user_id, exam_type_id)
);
```

---

## 5. Questions & Practice

### 5.1 `question`

Core MCQ bank.

```sql
CREATE TABLE question (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  exam_type_id     INT NOT NULL REFERENCES exam_type_lookup(id),
  subject_id       UUID NOT NULL REFERENCES subject(id),
  topic_id         UUID NOT NULL REFERENCES topic(id),
  question_text    TEXT NOT NULL,
  option_a         TEXT NOT NULL,
  option_b         TEXT NOT NULL,
  option_c         TEXT NOT NULL,
  option_d         TEXT NOT NULL,
  correct_option   SMALLINT NOT NULL CHECK (correct_option BETWEEN 1 AND 4),
  explanation      TEXT,
  difficulty_level SMALLINT NOT NULL DEFAULT 1, -- 1-5
  choice_type      question_choice_type NOT NULL DEFAULT 'single',
  is_clinical      BOOLEAN NOT NULL DEFAULT FALSE,
  is_image_based   BOOLEAN NOT NULL DEFAULT FALSE,
  is_high_yield    BOOLEAN NOT NULL DEFAULT FALSE,
  is_active        BOOLEAN NOT NULL DEFAULT TRUE,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 5.2 `practice_session`

Represents a practice session for a user (smart/custom/revision/exam).

```sql
CREATE TABLE practice_session (
  id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id                 UUID NOT NULL REFERENCES "user"(id),
  exam_type_id            INT NOT NULL REFERENCES exam_type_lookup(id),
  mode                    practice_mode NOT NULL,
  status                  practice_session_status NOT NULL DEFAULT 'in_progress',
  total_questions_planned INT,
  started_at              TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at            TIMESTAMPTZ
);
```

### 5.3 `practice_session_question`

Ordered list of questions inside a practice session.

```sql
CREATE TABLE practice_session_question (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id      UUID NOT NULL REFERENCES practice_session(id) ON DELETE CASCADE,
  question_id     UUID NOT NULL REFERENCES question(id),
  sequence_index  INT NOT NULL,
  selected_option SMALLINT,
  is_correct      BOOLEAN,
  time_taken_ms   INT,
  answered_at     TIMESTAMPTZ,
  UNIQUE (session_id, sequence_index)
);
```

### 5.4 `user_question_attempt`

Flat log of all attempts (useful for analytics, SRS, weakness detection).

```sql
CREATE TABLE user_question_attempt (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         UUID NOT NULL REFERENCES "user"(id),
  exam_type_id    INT NOT NULL REFERENCES exam_type_lookup(id),
  question_id     UUID NOT NULL REFERENCES question(id),
  session_id      UUID REFERENCES practice_session(id),
  is_correct      BOOLEAN NOT NULL,
  selected_option SMALLINT,
  time_taken_ms   INT,
  source          TEXT NOT NULL, -- 'practice', 'exam', 'revision'
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

## 6. Exams & Events

### 6.1 `exam_config`

Configuration for exams / events (mock test, subject test, reward event, etc.).

```sql
CREATE TABLE exam_config (
  id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  exam_type_id       INT NOT NULL REFERENCES exam_type_lookup(id),
  name               TEXT NOT NULL,
  type               exam_type NOT NULL,
  description        TEXT,
  num_questions      INT NOT NULL,
  time_limit_minutes INT NOT NULL,
  marks_per_correct  NUMERIC(5,2) NOT NULL DEFAULT 4.0,
  negative_per_wrong NUMERIC(5,2) NOT NULL DEFAULT -1.0,
  entry_fee_cents    INT NOT NULL DEFAULT 0,
  schedule_start_at  TIMESTAMPTZ,
  schedule_end_at    TIMESTAMPTZ,
  status             exam_status NOT NULL DEFAULT 'DRAFT',
  created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 6.2 `exam_question`

Mapping of questions to exam_config with ordering.

```sql
CREATE TABLE exam_question (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  exam_config_id UUID NOT NULL REFERENCES exam_config(id) ON DELETE CASCADE,
  question_id    UUID NOT NULL REFERENCES question(id),
  sequence_index INT NOT NULL,
  UNIQUE (exam_config_id, sequence_index)
);
```

### 6.3 `exam_registration`

User registrations for exams (if needed).

```sql
CREATE TABLE exam_registration (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  exam_config_id UUID NOT NULL REFERENCES exam_config(id),
  user_id        UUID NOT NULL REFERENCES "user"(id),
  registered_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (exam_config_id, user_id)
);
```

### 6.4 `exam_attempt`

Stores actual attempts & scores.

```sql
CREATE TABLE exam_attempt (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  exam_config_id UUID NOT NULL REFERENCES exam_config(id),
  user_id        UUID NOT NULL REFERENCES "user"(id),
  started_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at   TIMESTAMPTZ,
  score          NUMERIC(8,2),
  percentage     NUMERIC(5,2),
  rank           INT,
  reward_cents   INT DEFAULT 0,
  status         TEXT NOT NULL DEFAULT 'completed'
);
```

---

## 7. Podcasts

### 7.1 `podcast_episode`

```sql
CREATE TABLE podcast_episode (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  exam_type_id     INT NOT NULL REFERENCES exam_type_lookup(id),
  subject_id       UUID REFERENCES subject(id),
  topic_id         UUID REFERENCES topic(id),
  title            TEXT NOT NULL,
  description      TEXT,
  audio_url        TEXT NOT NULL,
  duration_seconds INT,
  tags             JSONB,
  is_active        BOOLEAN NOT NULL DEFAULT TRUE,
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

## 8. Wallet, Coupons & Referrals

### 8.1 `wallet_transaction`

No dedicated `wallet` table – balance is derived.

```sql
CREATE TABLE wallet_transaction (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         UUID NOT NULL REFERENCES "user"(id),
  exam_type_id    INT REFERENCES exam_type_lookup(id),
  amount_cents    INT NOT NULL, -- positive or negative
  tx_type         wallet_tx_type NOT NULL,
  description     TEXT,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 8.2 `coupon`

```sql
CREATE TABLE coupon (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code              TEXT NOT NULL UNIQUE,
  description       TEXT,
  type              TEXT NOT NULL, -- 'FIXED', 'ENTRY_PASS', 'BONUS'
  amount_cents      INT,
  max_uses_total    INT,
  max_uses_per_user INT,
  expires_at        TIMESTAMPTZ,
  is_active         BOOLEAN NOT NULL DEFAULT TRUE,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 8.3 `coupon_redemption`

```sql
CREATE TABLE coupon_redemption (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  coupon_id    UUID NOT NULL REFERENCES coupon(id),
  user_id      UUID NOT NULL REFERENCES "user"(id),
  redeemed_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  amount_cents INT,
  UNIQUE (coupon_id, user_id)
);
```

### 8.4 `referral`

```sql
CREATE TABLE referral (
  id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  referrer_user_id   UUID NOT NULL REFERENCES "user"(id),
  referred_user_id   UUID REFERENCES "user"(id),
  referral_code_used TEXT,
  status             referral_status NOT NULL DEFAULT 'INVITED',
  bonus_cents        INT DEFAULT 0,
  created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  activated_at       TIMESTAMPTZ
);
```

---

## 9. Revision / SRS

### 9.1 `revision_item`

```sql
CREATE TABLE revision_item (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id         UUID NOT NULL REFERENCES "user"(id),
  exam_type_id    INT NOT NULL REFERENCES exam_type_lookup(id),
  question_id     UUID NOT NULL REFERENCES question(id),
  next_review_at  TIMESTAMPTZ NOT NULL,
  interval_index  INT NOT NULL DEFAULT 0,
  times_reviewed  INT NOT NULL DEFAULT 0,
  last_result     revision_result,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, question_id)
);
```

---

## 10. Notes

* All `UUID` fields assume `pgcrypto` or `gen_random_uuid()` is available.
* Monetary values are stored as **cents/paise** (INT) for accuracy.
* Many aggregate views (streaks, levels, dashboards) will be computed from these base tables.

````
