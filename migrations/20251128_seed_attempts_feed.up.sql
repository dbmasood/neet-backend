INSERT INTO "user" (id, telegram_id, display_name, primary_exam_type, role, created_at, updated_at)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'seed-bot', 'Seed User', 'NEET_PG', 'USER', now(), now())
ON CONFLICT (id) DO NOTHING;

WITH attempts AS (
  SELECT q.id
  FROM question q
  JOIN subject s ON s.id = q.subject_id
  WHERE s.name = 'Pathology'
  LIMIT 5
)
INSERT INTO user_question_attempt (id, user_id, exam_type_id, question_id, is_correct, selected_option, time_taken_ms, source)
SELECT gen_random_uuid(), '11111111-1111-1111-1111-111111111111', (SELECT id FROM exam_type_lookup WHERE code = 'NEET_PG'), q.id, TRUE, 1, 1200, 'practice'
FROM attempts q
ON CONFLICT (id) DO NOTHING;

INSERT INTO feed_post (id, type, title, body, image_url, tags, author, cta, likes, comments, read_time)
VALUES
  ('22222222-2222-2222-2222-222222222222', 'announcement', 'Daily Practice Heroes', 'Top ranks just completed 50 questions!', 'https://example.com/daily.png', 'study,practice', 'Neet PG', '/practice', 12, 4, '2 min read'),
  ('33333333-3333-3333-3333-333333333333', 'article', 'Revision Rituals', 'Learn how to revise high-yield topics in 10 minutes', 'https://example.com/revision.png', 'revision,high-yield', 'Neet PG', '/revision', 30, 7, '3 min read');
