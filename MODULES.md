
---

## `MODULES.md`

```md
# MODULES.md
Domain Modules & Responsibilities

This document defines each logical module in the backend and what it owns.  
Modules follow Clean Architecture:

- Entities in `internal/entity`
- Usecases in `internal/usecase/<module>`
- Repositories in `internal/repo/persistent`
- HTTP handlers in `internal/controller/http/v1`

---

## 1. auth

**Package:** `internal/usecase/auth`  
**Controller:** `internal/controller/http/v1/auth_routes.go`  
**Entities:** uses `entity.User` and tokens (JWT).

### Responsibilities

- Login / signup via Telegram (`POST /v1/auth/telegram`)
- Issue UserAuth / AdminAuth JWTs
- Validate roles & permissions (used by middleware, not directly here)

### Repositories / dependencies

- `UserRepository` (create/find user by telegramId/email)
- Token generator (JWT util, in `pkg`)

---

## 2. user

**Package:** `internal/usecase/user`  
**Controller:** `internal/controller/http/v1/user_routes.go`  
**Entities:** `User`, `ExamProfile`.

### Responsibilities

- Return current user profile + exam stats (`GET /v1/me`)
- Manage user-level stats & exam profile aggregates
- Update last login / streaks (triggered after login)

### Repos

- `UserRepository`
- `UserExamProfileRepository`

---

## 3. subject

**Package:** `internal/usecase/subject`  
**Controllers:**
- App: `user_routes.go` (subjects/topics listing)
- Admin: `admin_subjects_topics_routes.go`

### Responsibilities

- List subjects by exam (app)
- List topics by subject (app)
- Admin CRUD for subjects & topics

### Repos

- `SubjectRepository`
- `TopicRepository`

---

## 4. question

**Package:** `internal/usecase/question`  
**Controllers:**
- Admin: `admin_questions_routes.go`
- App: used indirectly via practice & revision.

### Responsibilities

- Admin CRUD for questions
- Query questions by filters (exam, subject, topic, difficulty, flags)
- Provide question selection utilities for:
  - practice sessions
  - exams
  - revision engine

### Repos

- `QuestionRepository`

---

## 5. practice

**Package:** `internal/usecase/practice`  
**Controller:** `internal/controller/http/v1/practice_routes.go`  
**Entities:** `PracticeSession`, `PracticeSessionQuestion`.

### Responsibilities

- Create practice sessions (smart/custom/revision/exam)
- Attach selected questions to session
- Track answers & correctness
- Update `user_question_attempt`
- Notify other modules (analytics, revision, weakness detection) via internal events

### Repos

- `PracticeSessionRepository`
- `PracticeSessionQuestionRepository`
- `UserQuestionAttemptRepository`
- Read access to `QuestionRepository`
- Optional: `UserExamProfileRepository` (for stats updates)

---

## 6. revision (SRS)

**Package:** `internal/usecase/revision`  
**Controller:** `revision_routes.go` (`GET /v1/revision/queue`)  

### Responsibilities

- Maintain `revision_item` table for each user & question
- Determine which items are “due” for revision
- Update intervals and `next_review_at` based on answers (events from `practice`)

### Repos

- `RevisionRepository`
- Read from `UserQuestionAttemptRepository`
- Possibly config from `AISettingsRepository`

---

## 7. exam

**Package:** `internal/usecase/exam`  
**Controllers:**
- App: `events_routes.go` (`/v1/events`)
- Admin: `admin_exams_routes.go`

### Responsibilities

- Admin CRUD for `exam_config`
- Manage `exam_question` mapping
- User side:
  - list available events
  - (future) registrations & attempts
- Integrate with wallet for entry fees & rewards

### Repos

- `ExamConfigRepository`
- `ExamQuestionRepository`
- `ExamRegistrationRepository` (future)
- `ExamAttemptRepository`

---

## 8. podcast

**Package:** `internal/usecase/podcast`  
**Controllers:**
- App: `podcast_routes.go`
- Admin: `admin_podcasts_routes.go`

### Responsibilities

- Admin CRUD for podcast episodes
- List/filter podcast episodes by exam/subject/topic
- Provide playback metadata (title, description, duration, audio URL)

### Repos

- `PodcastRepository`

---

## 9. wallet

**Package:** `internal/usecase/wallet`  
**Controllers:**
- App: `wallet_routes.go`
- Admin: (future: manual adjustments, user view)

### Responsibilities

- Compute wallet balance from `wallet_transaction`
- Return summary & transaction history
- Record transactions for:
  - exams
  - rewards
  - coupons
  - referrals
  - spins
  - bonuses

### Repos

- `WalletTransactionRepository`

---

## 10. coupon

**Package:** `internal/usecase/coupon`  
**Controllers:**
- App: `/v1/coupons/redeem`
- Admin: `admin_coupons_routes.go`

### Responsibilities

- Validate coupon code
- Check usage limits, expiry, active status
- Apply rewards via wallet transactions
- Admin CRUD for coupons

### Repos

- `CouponRepository`
- `CouponRedemptionRepository`
- `WalletTransactionRepository`

---

## 11. referral

**Package:** `internal/usecase/referral`  
**Controllers:**
- App: `/v1/referral`
- (Future Admin views can reuse this module)

### Responsibilities

- Track referral relationships
- Compute referral metrics (invited, joined, activated, earned)
- Trigger rewards when referral activates (exam activity, threshold reached)

### Repos

- `ReferralRepository`
- `WalletTransactionRepository`
- `UserRepository`

---

## 12. ai (AI Settings & Logic)

**Package:** `internal/usecase/ai`  
**Controller:** `admin_ai_settings_routes.go`

### Responsibilities

- Persist AI configuration:
  - weaknessMinAttempts
  - weaknessThresholdPercent
  - strongThresholdPercent
  - revisionIntervalsDays
  - includeGuessedCorrect
  - revisionEnabled
- Provide configuration to:
  - `revision` module
  - `practice` / `question` (weakness detection)

### Repos

- `AISettingsRepository`

---

## 13. analytics

**Package:** `internal/usecase/analytics`  
**Controller:** `admin_analytics_routes.go`

### Responsibilities

- Provide aggregates for the admin dashboard:
  - Total users
  - Active users (DAU)
  - Questions answered
  - Average accuracy
  - Study time
  - Rewards paid
- Possibly read from:
  - `user_question_attempt`
  - `user_exam_profile`
  - `exam_attempt`
  - `wallet_transaction`
- May use separate analytics tables or views in the future.

### Repos

- `AnalyticsRepository` (could be read-only views)
- Or uses several existing repositories

---

## 14. common / shared

### 14.1 internal events (future)

We may implement an internal event bus under `pkg/events`:

- `QuestionAnswered`
- `PracticeSessionCompleted`
- `ExamCompleted`
- `WalletTransactionCreated`

Usecases (`practice`, `exam`, `coupon`, `referral`, `analytics`) will **publish** events; other modules may **subscribe**.

### 14.2 Util / infra

- `pkg/logger` – zerolog setup
- `pkg/db` – DB connection & migrations
- `pkg/jwt` – token helpers
- `pkg/validator` – shared request validation

---

## Summary

Each module:

- Owns its **entities** (in `internal/entity`)  
- Exposes its **usecases** (in `internal/usecase/<module>`)  
- Depends only on **interfaces**, with concrete repos in `internal/repo/persistent`  
- Has HTTP entrypoints in `internal/controller/http/v1`  

When adding new features, always:

1. Extend entities if needed.  
2. Add/extend usecase methods.  
3. Implement repository changes.  
4. Add routes & handlers.  
5. Wire everything in `internal/app/app.go`.  
````

