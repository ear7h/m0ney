package sessions

import (
	"m0ney/data"
	"time"
	"fmt"
)


//this represents a practice/training session
type Session struct {
	SessStart   time.Time
	CurrentTime time.Time     `json:"currentTime"`
	EndTime     time.Time     `json:"endTime"`
	Scale    time.Duration `json:"interval"`
	Ticker      string        `json:"ticker"`
	Table       data.Table    `json:"table"`
}


//returns moment or historical struct
func (s *Session) Next() interface{} {

	bottomStr := s.CurrentTime.Format(data.SQL_TIME)
	topStr := s.CurrentTime.Add(s.Scale).Format(data.SQL_TIME)

	rows, err := data.DB.Query(
		"SELECT * FROM `?` WHERE updated_at > ? AND updated_at <= ? ORDER BY updated_at DESC LIMIT 1;", s.Table ,topStr, bottomStr)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()


	var ret interface{}

	rows.Next()
	if s.Table == "moment" {
		ret := data.Moment{}
		err = rows.Scan(
			&ret.AskPrice,
			&ret.AskSize,
			&ret.BidPrice,
			&ret.BidSize,
			&ret.LastTradePrice,
			&ret.Symbol,
			&ret.UpdatedAt,
		)
	} else if s.Table == "historical" {
		//todo
	}

	if err != nil {
		return nil
	}

	return ret
}
