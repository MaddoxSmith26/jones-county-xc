-- name: GetAllAthletes :many
SELECT id, name, grade, personal_record, events, created_at
FROM athletes
ORDER BY name;

-- name: GetAthleteByID :one
SELECT id, name, grade, personal_record, events, created_at
FROM athletes
WHERE id = ?;

-- name: GetAllMeets :many
SELECT id, name, date, location, description, created_at
FROM meets
ORDER BY date;

-- name: GetResultsByMeetID :many
SELECT r.id, r.athlete_id, r.meet_id, r.time, r.place, r.created_at,
       a.name AS athlete_name, a.grade AS athlete_grade
FROM results r
JOIN athletes a ON r.athlete_id = a.id
WHERE r.meet_id = ?
ORDER BY r.place;

-- name: CreateResult :execresult
INSERT INTO results (athlete_id, meet_id, time, place)
VALUES (?, ?, ?, ?);

-- name: CreateAthlete :execresult
INSERT INTO athletes (name, grade, personal_record, events)
VALUES (?, ?, ?, ?);

-- name: CreateMeet :execresult
INSERT INTO meets (name, date, location, description)
VALUES (?, ?, ?, ?);

-- name: DeleteAthlete :exec
DELETE FROM athletes WHERE id = ?;

-- name: DeleteMeet :exec
DELETE FROM meets WHERE id = ?;

-- name: DeleteResult :exec
DELETE FROM results WHERE id = ?;

-- name: GetTopTimes :many
SELECT r.id, r.athlete_id, r.meet_id, r.time, r.place,
       a.name AS athlete_name, a.grade AS athlete_grade,
       m.name AS meet_name, m.date AS meet_date
FROM results r
JOIN athletes a ON r.athlete_id = a.id
JOIN meets m ON r.meet_id = m.id
ORDER BY r.time ASC
LIMIT 10;
