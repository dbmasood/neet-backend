-- remove NEET PG questions inserted by the migration

DELETE FROM question
WHERE exam_type_id = (SELECT id FROM exam_type_lookup WHERE code = 'NEET_PG');

