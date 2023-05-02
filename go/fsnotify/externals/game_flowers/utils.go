package flower

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DecodedData is got from luckydecoder
type DecodedData struct {
	ID          uint64
	Type        string
	UserSet     bson.M
	UserInc     bson.M
	GameSet     bson.M
	GameInc     bson.M
	GamePush    bson.M
	LogPush     bson.M
	LogSet      bson.M
	TaxSet      bson.M
	StrategySet bson.M
	Data        bson.M
	UserRecord bson.M
	GameOnlineRecord bson.M
}

// ConvertTime convert all field endswith `_at` to utc time
func ConvertTime(data *bson.M) {
	var (
		err error
	)
	for k, v := range *data {
		if strings.HasSuffix(k, "_at") {
			timestamp := v.(string)
			(*data)[k], err = time.Parse("2006-01-02T15:04:05.999999", timestamp)
			if err != nil {
				log.Printf("time %s format wrong", timestamp)
				(*data)[k] = time.Now().UTC()
			}
		}
	}
}

// GetLocalStr change utc time to local date str
func GetLocalStr(base time.Time) string {
	h, _ := time.ParseDuration("1h")
	return base.Add(8 * h).Format("2006-01-02")
}

func GetNanoTimestamp(base time.Time) int64 {
	return base.UnixNano() / 1e6
}

func GetTimestamp(base time.Time) int64 {
	return base.Unix()
}

// IsSameDay check if two time is same day locally
func IsSameDay(l time.Time, r time.Time) bool {
	return GetLocalStr(l) == GetLocalStr(r)
}

// SaveWithRetry call upsert twice for mgo bug
func SaveWithRetry(db *mgo.Database, c string, id interface{}, data map[string]bson.M) error {
	_, err := db.C(c).UpsertId(id, data)
	if err != nil {
		if mgo.IsDup(err) {
			_, err = db.C(c).UpsertId(id, data)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// InsertNew insert new record without specified id
func InsertNew(db *mgo.Database, c string, data bson.M) error {
	err := db.C(c).Insert(data)
	return err
}

//ConvertToInt guess Num format and convert to Int
func ConvertToInt(temp interface{}) (int, error) {
	switch t := temp.(type) {
	case int:
		return int(t), nil
	case float64, float32:
		return int(reflect.ValueOf(t).Float()), nil
	case int64, int32:
		return int(reflect.ValueOf(t).Int()), nil
	default:
		return 0, fmt.Errorf("can't convert to int:%v", temp)
	}
}

var floatType = reflect.TypeOf(float64(0))

//ConvertToFloat64 guess Num format and convert to Float64
func ConvertToFloat64(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}

const (
	OffSuit  = 1 // 散牌
	Pair     = 2 // 对子
	Straight = 3 // 顺子
	Flush    = 4 // 同花
	Sflush   = 5 // 同花顺
	Set      = 6 // 豹子
)

const (
	ROOM_LOW    = 0 // 初级房
	ROOM_MID    = 1 // 中级房
	ROOM_HIGH   = 2 // 高级房
	ROOM_MASTER = 3 // 大师房
)

const (
	BET_BASE    = 1
	BET_CALL    = 2
	BET_RAISE   = 3
	BET_DOUBLE  = 4
	BET_COMPARE = 5
	BET_PK_ALL  = 6
)

