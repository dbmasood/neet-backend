-- Creates the core NEET PG schema so application data migrations have a foundation.
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TYPE exam_category AS ENUM ('NEET_PG', 'NEET_UG', 'JEE', 'UPSC');
CREATE TYPE question_choice_type AS ENUM ('single', 'multi');

CREATE TABLE exam_type_lookup (
  id SERIAL PRIMARY KEY,
  code exam_category UNIQUE NOT NULL,
  name TEXT NOT NULL,
  description TEXT
);

INSERT INTO exam_type_lookup (code, name, description)
VALUES
  ('NEET_PG', 'NEET PG', 'NEET PG exam content'),
  ('NEET_UG', 'NEET UG', 'NEET UG exam content'),
  ('JEE', 'Joint Entrance Examination', 'JEE exam content'),
  ('UPSC', 'UPSC Civil Services', 'UPSC exam content')
ON CONFLICT (code) DO NOTHING;

CREATE TABLE subject (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  exam_type_id INT NOT NULL REFERENCES exam_type_lookup(id),
  name TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  UNIQUE (exam_type_id, name)
);

CREATE TABLE topic (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  subject_id UUID NOT NULL REFERENCES subject(id),
  name TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  UNIQUE (subject_id, name)
);

CREATE TABLE question (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  exam_type_id INT NOT NULL REFERENCES exam_type_lookup(id),
  subject_id UUID NOT NULL REFERENCES subject(id),
  topic_id UUID NOT NULL REFERENCES topic(id),
  question_text TEXT NOT NULL,
  option_a TEXT NOT NULL,
  option_b TEXT NOT NULL,
  option_c TEXT NOT NULL,
  option_d TEXT NOT NULL,
  correct_option SMALLINT NOT NULL CHECK (correct_option BETWEEN 1 AND 4),
  explanation TEXT,
  difficulty_level SMALLINT NOT NULL DEFAULT 1,
  choice_type question_choice_type NOT NULL DEFAULT 'single',
  is_clinical BOOLEAN NOT NULL DEFAULT FALSE,
  is_image_based BOOLEAN NOT NULL DEFAULT FALSE,
  is_high_yield BOOLEAN NOT NULL DEFAULT FALSE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE practice_session (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  exam_type_id INT NOT NULL REFERENCES exam_type_lookup(id),
  mode TEXT,
  status TEXT,
  total_questions_planned INT,
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at TIMESTAMPTZ
);

CREATE TABLE practice_session_question (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id UUID NOT NULL REFERENCES practice_session(id) ON DELETE CASCADE,
  question_id UUID NOT NULL REFERENCES question(id),
  sequence_index INT NOT NULL,
  selected_option SMALLINT,
  is_correct BOOLEAN,
  time_taken_ms INT,
  answered_at TIMESTAMPTZ
);

CREATE TABLE "user" (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  telegram_id TEXT,
  email TEXT,
  display_name TEXT,
  primary_exam_type exam_category NOT NULL,
  role TEXT NOT NULL DEFAULT 'USER',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_question_attempt (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES "user"(id),
  exam_type_id INT NOT NULL REFERENCES exam_type_lookup(id),
  question_id UUID NOT NULL REFERENCES question(id),
  session_id UUID,
  is_correct BOOLEAN NOT NULL,
  selected_option SMALLINT,
  time_taken_ms INT,
  source TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
