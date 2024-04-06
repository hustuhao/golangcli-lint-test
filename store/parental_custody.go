package store

import (
	"database/sql/driver"
	"encoding/json"
	"golangcli-lint-test/store/db"
)

// ParentalCustodyConfInfo 是家长监护配置的数据结构
type ParentalCustodyConfInfo struct {
	ID         int64                                  `db:"id"`
	AccountID  int64                                  `db:"account_id"`
	Content    *JsonColumn[ParentalCustodyDetailConf] `db:"content"`
	CreateTime int64                                  `db:"create_time"`
	UpdateTime int64                                  `db:"update_time"`
}

type ParentalCustodyDetailConf struct {
	Open                   bool   `json:"Open"`
	Pwd                    string `json:"Pwd"`
	ChargeLimitEveryTime   int    `json:"ChargeLimitEveryTime"`
	ChargeLimitEveryMonth  int    `json:"ChargeLimitEveryMonth"`
	PlayTimeLimitEveryTime int    `json:"PlayTimeLimitEveryTime"`
	PlayTimeLimitEveryDay  int    `json:"PlayTimeLimitEveryDay"`
	NoWorldLegionChat      bool   `json:"NoWorldLegionChat"`
	TeamFriendOnly         bool   `json:"TeamFriendOnly"`
	NoChat                 bool   `json:"NoChat"`
}

type ParentalCustodyPlayTimeMonitor struct {
	Id                      int64  `db:"id"`
	AccountId               int64  `db:"account_id"`
	Uid                     string `db:"uid"`
	LoginDay                int64  `db:"login_day"`
	CreateTime              int64  `db:"create_time"`
	UpdateTime              int64  `db:"update_time"`
	DayOnlineDuration       int    `db:"day_online_duration"`
	DayMaxOnlineDuration    int    `db:"day_max_online_duration"`
	SingleOnlineDuration    int    `db:"single_online_duration"`
	SingleMaxOnlineDuration int    `db:"single_max_online_duration"`
}

type JsonColumn[T any] struct {
	v *T
}

func (j *JsonColumn[T]) Scan(src any) error {
	if src == nil {
		j.v = nil
		return nil
	}
	j.v = new(T)
	return json.Unmarshal(src.([]byte), &j.v)
}

func (j *JsonColumn[T]) Value() (driver.Value, error) {
	raw, err := json.Marshal(j.v)
	return raw, err
}

func (j *JsonColumn[T]) Get() *T {
	return j.v
}

func (j *JsonColumn[T]) Set(v *T) {
	j.v = v
	return
}

func NewParentalCustodyConfInfo() *ParentalCustodyConfInfo {
	return &ParentalCustodyConfInfo{
		Content: &JsonColumn[ParentalCustodyDetailConf]{
			v: DefaultParentalCustodyDetailConfInfo(),
		},
	}
}

func DefaultParentalCustodyDetailConfInfo() *ParentalCustodyDetailConf {
	c := new(ParentalCustodyDetailConf)
	c.ChargeLimitEveryTime = 100
	c.ChargeLimitEveryMonth = 400
	c.PlayTimeLimitEveryTime = 30
	c.PlayTimeLimitEveryDay = 90
	return c
}

func GetParentalCustodyDetailConfByAccount(accountId int64) (*ParentalCustodyDetailConf, error) {
	conf := DefaultParentalCustodyDetailConfInfo()
	query := "SELECT content FROM parental_custody_conf WHERE account_id = ?"
	rows, err := db.MainDB.Query(query, accountId)
	if err != nil {
		return conf, err // 查询出错
	}
	var str string
	for rows.Next() {
		rows.Scan(&str)
		break
	}
	if str == "" {
		return conf, nil
	}
	err = json.Unmarshal([]byte(str), conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}
