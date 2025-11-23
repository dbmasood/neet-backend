package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var subjectList = []string{
	"Anesthesia",
	"Anatomy",
	"Biochemistry",
	"Dental",
	"ENT",
	"Forensic Medicine (FM)",
	"Obstetrics and Gynecology (O&G)",
	"Medicine",
	"Microbiology",
	"Ophthalmology",
	"Orthopedics",
	"Pathology",
	"Pediatrics",
	"Pharmacology",
	"Physiology",
	"Psychiatry",
	"Radiology",
	"Skin",
	"Preventive & Social Medicine (PSM)",
	"Surgery",
}

var subjectAliases = func() map[string]string {
	m := make(map[string]string, len(subjectList))
	for _, subj := range subjectList {
		m[strings.ToLower(strings.TrimSpace(subj))] = subj
	}
	return m
}()

type questionRecord struct {
	Question      string  `json:"question"`
	Explanation   *string `json:"exp"`
	CorrectOption int     `json:"cop"`
	OptionA       *string `json:"opa"`
	OptionB       *string `json:"opb"`
	OptionC       *string `json:"opc"`
	OptionD       *string `json:"opd"`
	SubjectName   string  `json:"subject_name"`
	TopicName     *string `json:"topic_name"`
	ID            string  `json:"id"`
	ChoiceType    *string `json:"choice_type"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <medmcqa-dev-ndjson> <migration-prefix>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	inputPath := os.Args[1]
	prefix := os.Args[2]

	subjectTopics, err := gatherTopics(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gather topics: %v\n", err)
		os.Exit(1)
	}

	if err := writeUp(prefix+".up.sql", inputPath, subjectTopics); err != nil {
		fmt.Fprintf(os.Stderr, "write up migration: %v\n", err)
		os.Exit(1)
	}

	if err := writeDown(prefix + ".down.sql"); err != nil {
		fmt.Fprintf(os.Stderr, "write down migration: %v\n", err)
		os.Exit(1)
	}
}

func gatherTopics(path string) (map[string]map[string]struct{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	subjects := make(map[string]map[string]struct{}, len(subjectList))
	for _, subj := range subjectList {
		subjects[subj] = map[string]struct{}{
			"General": {},
		}
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var rec questionRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			return nil, fmt.Errorf("unmarshal: %w", err)
		}

		subject := canonicalSubject(rec.SubjectName)
		if subject == "" {
			continue
		}

		topic := canonicalTopic(rec.TopicName)
		if topic == "" {
			topic = "General"
		}

		subjects[subject][topic] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}

func canonicalSubject(value string) string {
	return subjectAliases[strings.ToLower(strings.TrimSpace(value))]
}

func canonicalTopic(value *string) string {
	if value == nil {
		return "General"
	}
	clean := strings.TrimSpace(*value)
	if clean == "" {
		return "General"
	}
	return clean
}

func writeUp(path, dataPath string, subjects map[string]map[string]struct{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "-- Generated NEET PG MedMCQA migration loading subjects/topics/questions from %s\n", dataPath)
	fmt.Fprintln(f, "-- Adds the dataset into the existing `subject`, `topic`, and `question` tables in the `NEET_PG` exam.")

	writeSubjectInserts(f)
	writeTopicInserts(f, subjects)
	return writeQuestionInserts(f, dataPath)
}

func writeSubjectInserts(f *os.File) {
	fmt.Fprintln(f, `
WITH exam AS (
    SELECT id FROM exam_type_lookup WHERE code = 'NEET_PG'
)
INSERT INTO subject (exam_type_id, name, is_active)
SELECT exam.id, data.name, TRUE
FROM exam
CROSS JOIN (VALUES`)

	for i, subj := range subjectList {
		fmt.Fprintf(f, "    ('%s')%s\n", escape(subj), comma(i, len(subjectList)))
	}

	fmt.Fprintln(f, `) data(name)
ON CONFLICT (exam_type_id, name) DO UPDATE SET is_active = EXCLUDED.is_active;`)
}

func writeTopicInserts(f *os.File, subjects map[string]map[string]struct{}) {
	type pair struct {
		subject string
		topic   string
	}

	var topicPairs []pair
	for _, subj := range subjectList {
		for topic := range subjects[subj] {
			topicPairs = append(topicPairs, pair{subject: subj, topic: topic})
		}
	}

	sort.Slice(topicPairs, func(i, j int) bool {
		if topicPairs[i].subject == topicPairs[j].subject {
			return topicPairs[i].topic < topicPairs[j].topic
		}
		return topicPairs[i].subject < topicPairs[j].subject
	})

	for _, pair := range topicPairs {
		fmt.Fprintf(f, `
INSERT INTO topic (subject_id, name, is_active)
SELECT sub.id, '%s', TRUE
FROM subject sub
JOIN exam_type_lookup exam ON exam.id = sub.exam_type_id
WHERE exam.code = 'NEET_PG' AND sub.name = '%s'
ON CONFLICT (subject_id, name) DO UPDATE SET is_active = EXCLUDED.is_active;
`, escape(pair.topic), escape(pair.subject))
	}
}

func writeQuestionInserts(f *os.File, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var rec questionRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			return fmt.Errorf("unmarshal questions: %w", err)
		}

		subject := canonicalSubject(rec.SubjectName)
		if subject == "" {
			continue
		}

		topic := canonicalTopic(rec.TopicName)
		if topic == "" {
			topic = "General"
		}

		choiceType := "single"
		if rec.ChoiceType != nil && strings.ToLower(strings.TrimSpace(*rec.ChoiceType)) == "multi" {
			choiceType = "multi"
		}

		correct := rec.CorrectOption
		if correct < 1 || correct > 4 {
			correct = 1
		}

		fmt.Fprintf(f, `INSERT INTO question (id, exam_type_id, subject_id, topic_id, question_text, option_a, option_b, option_c, option_d, correct_option, explanation, choice_type, difficulty_level, is_clinical, is_image_based, is_high_yield, is_active)
SELECT '%s',
       exam.id,
       sub.id,
       top.id,
       '%s',
       '%s',
       '%s',
       '%s',
       '%s',
       %d,
       %s,
       '%s',
       1,
       FALSE,
       FALSE,
       FALSE,
       TRUE
FROM exam_type_lookup exam
JOIN subject sub ON sub.exam_type_id = exam.id AND sub.name = '%s'
JOIN topic top ON top.subject_id = sub.id AND top.name = '%s'
WHERE exam.code = 'NEET_PG'
ON CONFLICT (id) DO NOTHING;

`, escape(rec.ID), escape(rec.Question), escapePtr(rec.OptionA), escapePtr(rec.OptionB), escapePtr(rec.OptionC), escapePtr(rec.OptionD), correct, nullOrQuote(rec.Explanation), choiceType, escape(subject), escape(topic))

		if count++; count%1000 == 0 {
			fmt.Fprintf(f, "-- inserted %d questions\n", count)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func writeDown(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "-- remove NEET PG questions inserted by the migration")
	fmt.Fprintln(f, `DELETE FROM question
WHERE exam_type_id = (SELECT id FROM exam_type_lookup WHERE code = 'NEET_PG');`)

	return nil
}

func escape(value string) string {
	clean := strings.ReplaceAll(value, "'", "''")
	clean = strings.ReplaceAll(clean, "\n", " ")
	return clean
}

func escapePtr(value *string) string {
	if value == nil {
		return ""
	}
	return escape(*value)
}

func nullOrQuote(value *string) string {
	if value == nil {
		return "NULL"
	}
	return "'" + escape(*value) + "'"
}

func comma(idx, total int) string {
	if idx == total-1 {
		return ""
	}
	return ","
}
