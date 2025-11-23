INSERT INTO practice_session (id, user_id, exam_type_id, mode, status, total_questions_planned, started_at, completed_at)
VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid, '11111111-1111-1111-1111-111111111111'::uuid, (SELECT id FROM exam_type_lookup WHERE code = 'NEET_PG'), 'smart', 'completed', 5, now() - interval '2 hours', now()),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'::uuid, '11111111-1111-1111-1111-111111111111'::uuid, (SELECT id FROM exam_type_lookup WHERE code = 'NEET_PG'), 'custom', 'in_progress', 3, now() - interval '1 hours', NULL);

WITH selected_questions AS (
  SELECT q.id, row_number() OVER (ORDER BY q.id) AS idx
  FROM question q
  JOIN subject s ON s.id = q.subject_id
  WHERE s.name = 'Pathology'
  LIMIT 5
)
INSERT INTO practice_session_question (session_id, question_id, sequence_index, selected_option, is_correct, time_taken_ms, answered_at)
SELECT 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid, id, idx, 1, TRUE, 1500, now()
FROM selected_questions
UNION ALL
SELECT 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb'::uuid, id, idx, NULL, NULL, NULL, NULL
FROM selected_questions;
