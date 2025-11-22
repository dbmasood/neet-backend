# PLANS.md

**Master Development Plan – Multi-Exam Learning Platform (NEET PG, NEET UG, JEE, UPSC)**
**Backend: Go (go-clean-template) · Frontend Admin: shadcn/ui · App: Telegram Mini App + Web App**

---

## 1. Overview

This project is a full-scale multi-exam learning platform designed for NEET PG, NEET UG, JEE, and UPSC aspirants.
It includes:

* A Telegram Mini App & Web App for students
* An Admin Dashboard for content management
* A backend written in **Go** following **Clean Architecture**
* PostgreSQL + Redis + Vercel Blob Storage
* AI-powered features: Weakness Detection + Revision Engine (SRS)
* Gamification: wallet, rewards, coupons, spin wheel, streaks, levels

The system must be modular, scalable, and AI-friendly, supporting continuous extension.

---

# 2. Project Stages

## **Stage 1 — Architecture & Foundations (Backend + Admin UI)**

**Goal:** Get the system skeleton ready.

### Backend (Go Clean Template)

* [x] Finalize `AGENTS.md`
* [ ] Add `SCHEMA.md`
* [ ] Add `ROUTES.md`
* [x] Add `openapi.yaml` (v1 passed)
* [ ] Generate entities (Go structs)
* [ ] Create module folders (question, exam, user, wallet, podcast…)
* [ ] Wire DI in `internal/app/app.go`
* [ ] Add initial HTTP routing in /v1
* [ ] Implement JWT middlewares (UserAuth, AdminAuth)

### Admin Dashboard (shadcn Admin Template)

* [ ] Dashboard overview widgets
* [ ] Question CRUD pages
* [ ] Subject / Topic CRUD
* [ ] Exam & Event creation screens (multi-tab form)
* [ ] Podcast management
* [ ] Coupons & reward rules
* [ ] Users analytics pages
* [ ] AI Settings page
* [ ] Referral stats & leaderboard
* [ ] Menu & layout polished (Binance theme)

---

# 3. Stage 2 — Core Learning Experience (User-facing)

## **Practice System**

* Smart Practice (AI guided)
* Custom practice by subject/topic/difficulty
* Immediate feedback with explanations
* Level progression (L1 → L5)
* Weak topic identification
* Time tracking per question

## **Revision Engine (SRS)**

* Implement revision queue generator
* Track intervals per question
* Daily review UI

## **Daily Tests & Mock Exams**

* Exams with scoring, ranking, rewards
* Registration system
* Multi-exam support
* Negative marking rules
* Time-bound sessions

---

# 4. Stage 3 — Content & Engagement

## **Podcast System**

* List podcasts by exam > subject > topic
* AI-generated episodes uploaded via admin
* Audio streaming via Vercel Blob or Cloud storage
* Listening tracking (time spent)

## **Microlearning Feed**

* Instagram-like infinite scroll
* Cards: flashcards, questions, facts, memes, tips
* TikTok/Reel-style content consumption

## **Gamification**

* Spin wheel (daily reward)
* Wallet points as INR (in-app only)
* Earn by:

  * right answers
  * streaks
  * referrals
  * completing tests
  * watching podcasts
* Redeem via coupons or events

---

# 5. Stage 4 — Economics & Monetization (Future)

## **Referral System**

* Referral codes
* Invite → Join → Activate → Earn cycle
* Payout rules configured from Admin

## **Coupons**

* Fixed amount coupons
* Entry-pass coupons
* Bonus/reward boost coupons

## **Marketplace (Future)**

* Merch items (t-shirts, books, posters)
* Users redeem wallet points

---

# 6. Stage 5 — AI Features

## Weakness Detection

* Identify user's low-accuracy topics
* Monitor accuracy over time
* Auto-adjust difficulty

## Smart Practice Engine

* Monitor:

  * depth
  * accuracy
  * speed
  * time spent
* Generate best next set of questions

## Revision Engine (SRS)

* Uses spaced repetition intervals
* Customize intervals via admin

---

# 7. Stage 6 — Analytics & Observability

### Admin Analytics

* User growth
* DAU / MAU
* Questions answered
* Accuracy trends
* Subject performance
* Topic weakness map
* Exam participation
* Wallet transactions
* Referral performance
* Revenue dashboard (future)

### Technical Observability

* Prometheus metrics
* CPU / mem / p95 latency
* Request logs via zerolog
* RabbitMQ/NATS health

---

# 8. Backend Modules (Go Clean Architecture)

Each module includes:

* `entity/`
* `usecase/`
* `repo/persistent/`
* `controller/http/v1/`

### Modules:

1. **User**
2. **Auth**
3. **Subject**
4. **Topic**
5. **Question**
6. **Practice**
7. **Revision**
8. **Exam**
9. **Podcast**
10. **Wallet**
11. **Coupon**
12. **Referral**
13. **AI Settings**
14. **Analytics**

---

# 9. Database (Postgres + Redis)

* PostgreSQL for relational data
* Redis for:

  * caching subject/topic lists
  * caching user stats
  * queueing spin wheel jobs
  * session throttling

Migration tooling: golang-migrate.

---

# 10. Deployment Plan

**App (NextJS / Vercel)**

* Hosted on Vercel
* Uses Vercel Blob for file storage

**Backend**

* Deployed to Railway / Fly.io / AWS ECS
* Requires:

  * Postgres
  * Redis
  * RabbitMQ (optional)
  * S3 or Vercel Blob for audio uploads

**Admin Dashboard**

* Static build deployed on Vercel or S3+CloudFront

---

# 11. Milestones & Deliverables

## Milestone 1 — Backend Skeleton & Admin MVP

* Entities
* Usecases (empty but structured)
* Admin UI built
* Subject/Topic/Question basic CRUD
* Swagger for all endpoints

## Milestone 2 — Practice & Revision System

* Full flow working
* Smart practice logic
* Revision engine

## Milestone 3 — Exams & Events

* Admin can schedule exams
* User registration & attempts
* Scoring + ranking

## Milestone 4 — Wallet, Rewards & Gamification

* Wallet tracking
* Coupons
* Referral activation
* Spin wheel
* Streaks

## Milestone 5 — Podcasts & Feed

* Audio library
* Microlearning stream

## Milestone 6 — Analytics Layer

* Admin analytics (basic → advanced)

---

# 12. Rules for All Contributors & AI Agents

* Follow Clean Architecture strictly
* No business logic inside controllers
* Usecases must depend on interfaces
* Repositories must NOT import Fiber
* Entities must NOT import DB packages
* Every route must have Swagger annotations
* Everything must match OpenAPI spec
* Keep code modular
* Keep folder structure consistent

---

# 13. Future Scope

* Web3-style tokenization (optional future)
* Adaptive testing engine (CAT)
* Teacher dashboard
* Coaching classes management
* Social learning features
* Marketplace

---

# 14. Current Focus

🔵 **Focus for now (Phase 1):**

* Implement bare-minimum skeleton for **users**, **subjects**, **topics**, **questions**, and **admin CRUD**
* Connect openapi.yaml → controllers → usecases → repos
* Build admin dashboard UI for these modules
* Build Telegram Mini App initial screens

Once Phase 1 is solid, we move to practice + revision + exams.
