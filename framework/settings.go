package bot

import (
	"database/sql"
	"encoding/json"
)

type Json map[string]interface{}

type Settings struct {
	Db *sql.DB
}

func (s *Settings) Set(id, key string, value Json) error {

	marshaled, err := json.Marshal(value)

	if err != nil {
		return err
	}

	_, err = s.Db.Query(`insert into settings(id, key, value) values($1, $2, $3)
		on conflict (key) do update set key = $2, value = $3`, id, key, marshaled)

	return err
}

func (s *Settings) Get(id, key string) (Json, error) {

	row := s.Db.QueryRow(`select value from settings where id = $1 and key = $2`, id, key)

	if err := row.Err(); err == sql.ErrNoRows {
		return nil, err
	}

	var raw []byte
	err := row.Scan(&raw)

	if err != nil {
		return nil, err
	}

	var output Json
	err = json.Unmarshal(raw, &output)

	return output, err
}

func (s *Settings) Has(id, key string) bool {
	row := s.Db.QueryRow(`select exists(select 1 from settings where id = $1 and key = $2)`, id, key)

	if row.Err() != nil {
		return false
	}

	var has bool
	row.Scan(&has)

	return has
}

func (s *Settings) Delete(id, key string) error {
	_, err := s.Db.Query(`delete from settings where id = $1 and key = $2`, id, key)
	return err
}
