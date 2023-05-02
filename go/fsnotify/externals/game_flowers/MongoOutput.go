package flower

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //some special mysql type
	p "github.com/mozilla-services/heka/pipeline"

	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoOutput is main data of this package
type MongoOutput struct {
	*MongoOutputConfig
	Mg       *mgo.Session
	Rd       *redis.Pool
	Db       *gorm.DB
	PackChan chan *p.PipelinePack
	wg       *sync.WaitGroup
	or       p.OutputRunner
}

// MongoOutputConfig is config from toml file
type MongoOutputConfig struct {
	ProjectId            string   `toml:"project_id"`
	MongoAddr            []string `toml:"mongo_addr"`
	ConnectTimeout       uint32   `toml:"connect_timeout"`
	Procs                uint32   `toml:"procs"`
	DbName               string   `toml:"db_name"`
	TotalCollection      string   `toml:"total"`
	DailyCollection      string   `toml:"daily"`
	GameCollection       string   `toml:"game"`
	LogCollection        string   `toml:"game_log"`
	StrategyCollection   string   `toml:"strategy_record"`
	UserRecCollection    string   `toml:"user_record"`
	GameOnlineCollection string   `toml:"game_online_record"`
	Debug                bool     `toml:"debug"`
}

// ConfigStruct defines default configs
func (o *MongoOutput) ConfigStruct() interface{} {
	return &MongoOutputConfig{
		ProjectId:            "dark3",
		MongoAddr:            []string{"localhost:27017"},
		ConnectTimeout:       3,
		Procs:                5,
		DbName:               "game",
		TotalCollection:      "user_stats",
		DailyCollection:      "daily_stats",
		GameCollection:       "flower",
		LogCollection:        "flower_log",
		StrategyCollection:   "strategy_record",
		UserRecCollection:    "user_record",
		GameOnlineCollection: "game_online_record",
		Debug:                false,
	}
}

// Init is a `must` method for this plugin, called by heka
func (o *MongoOutput) Init(config interface{}) (err error) {
	o.MongoOutputConfig = config.(*MongoOutputConfig)
	o.PackChan = make(chan *p.PipelinePack, int(o.Procs))
	// init mongo connect
	dialInfo := mgo.DialInfo{Addrs: o.MongoAddr,
		Timeout: time.Duration(o.ConnectTimeout) * time.Second}
	if o.Mg, err = mgo.DialWithInfo(&dialInfo); err != nil {
		o.or.LogError(err)
		return err
	}
	return nil
}

// Run is called after Init, this is the old style before 0.10.0
func (o *MongoOutput) Run(or p.OutputRunner, h p.PluginHelper) (err error) {
	o.or = or
	o.wg = new(sync.WaitGroup)

	for i := 0; i < int(o.Procs); i++ {
		o.wg.Add(1)
		go save(o)
	}
	for pack := range or.InChan() {
		o.PackChan <- pack
	}
	if o.Mg != nil {
		o.Mg.Close()
	}
	close(o.PackChan)
	o.wg.Wait()
	return nil
}

// goroutine to dispatch data
func save(o *MongoOutput) {
	var (
		err error
		db  *mgo.Database
	)
	defer o.Db.Close()
	for pack := range o.PackChan {
		session := o.Mg.Copy()
		db = session.DB(o.DbName)
		if err = o.Process(pack, db); err != nil {
			o.or.LogError(err)
			pack.Recycle(p.NewRetryMessageError(err.Error()))
		} else {
			o.or.UpdateCursor(pack.QueueCursor)
			pack.Recycle(nil)
		}
		session.Close()
	}
	o.wg.Done()
}

func (o *MongoOutput) saveStats(data *DecodedData, db *mgo.Database) (err error) {
	var (
		ok        bool
		temp      interface{}
		updatedAt time.Time
	)

	result := make(map[string]interface{})
	if err = db.C(o.TotalCollection).Find(bson.M{"_id": data.ID}).One(&result); err != nil {
		return
	}

	insertData := make(map[string]bson.M)
	if len(data.UserSet) > 0 {
		insertData["$set"] = data.UserSet
	}
	if len(data.UserInc) > 0 {
		insertData["$inc"] = data.UserInc
	}
	err = SaveWithRetry(db, o.TotalCollection, data.ID, insertData)
	if err != nil {
		o.or.LogError(err)
	}
	temp = data.UserSet["updated_at"]
	if updatedAt, ok = temp.(time.Time); !ok {
		return fmt.Errorf("can't convert updated time: %v", data.UserSet)
	}
	key := fmt.Sprintf("%d-%s", data.ID, GetLocalStr(updatedAt))
	insertData["$set"]["user_id"] = data.ID
	insertData["$set"]["day"] = GetLocalStr(updatedAt)
	err = SaveWithRetry(db, o.DailyCollection, key, insertData)
	if err != nil {
		o.or.LogError(err)
	}
	return
}

func (o *MongoOutput) trackApply(data *DecodedData, db *mgo.Database) (err error) {
	err = o.saveStats(data, db)
	err = o.trackGame(data, db)
	return
}

func (o *MongoOutput) trackBet(data *DecodedData, db *mgo.Database) (err error) {
	result := make(map[string]interface{})
	if err = db.C(o.TotalCollection).Find(bson.M{"_id": data.ID}).One(&result); err != nil {
		//机器人
		data.GameInc["virtual_bet_amount"] = data.UserInc["flower.bet_amount"]
		if bet_amount, ok := data.UserInc["flower.bet_amount"].(float64); !ok {
			o.or.LogError(fmt.Errorf("flower.bet_amount %s is not a float64", data.UserInc["flower.bet_amount"]))
		} else {
			data.GameInc["virtual_gain_amount"] = 0 - bet_amount
		}
	} else {
		//真实玩家
		data.GameInc["real_bet_amount"] = data.UserInc["flower.bet_amount"]
		if bet_amount, ok := data.UserInc["flower.bet_amount"].(float64); !ok {
			o.or.LogError(fmt.Errorf("flower.bet_amount %s is not a float64", data.UserInc["flower.bet_amount"]))
		} else {
			data.GameInc["real_gain_amount"] = 0 - bet_amount
		}
		err = o.saveStats(data, db)
	}
	err = o.trackGame(data, db)
	return
}

func (o *MongoOutput) trackAward(data *DecodedData, db *mgo.Database) (err error) {
	isVirtual := 0
	result := make(map[string]interface{})
	if err = db.C(o.TotalCollection).Find(bson.M{"_id": data.ID}).One(&result); err != nil {
		//机器人
		isVirtual = -1
		data.GameInc["virtual_win_amount"] = data.UserInc["flower.win_amount"]
		data.GameInc["virtual_gain_amount"] = data.UserInc["flower.win_amount"]
	} else {
		//真实玩家
		data.GameInc["real_win_amount"] = data.UserInc["flower.win_amount"]
		data.GameInc["real_gain_amount"] = data.UserInc["flower.win_amount"]
	}
	data.GameSet["status"] = "done"
	err = o.trackGame(data, db)
	if isVirtual == 0 {
		data.UserSet["flower.last_res"] = 1
		err = o.saveStats(data, db)
	}
	// track输家记录
	gameStats := make(map[string]interface{})
	err = db.C(o.GameCollection).Find(bson.M{"_id": data.Data["game_id"]}).One(&gameStats)
	if err != nil {
		o.or.LogError(fmt.Errorf("trackAward: fail to find flower:%s %s", data.Data["game_id"], err.Error()))
		return err
	}
	if applyers, ok := gameStats["applyers"]; ok {
		for _, uid := range applyers.([]interface{}) {
			if uid.(float64) != float64(data.ID) {
				stat := DecodedData{
					ID:   uint64(uid.(float64)),
					Type: "lose",
					UserSet: bson.M{
						"flower.last_res": -1,
						"updated_at":      time.Now().UTC(),
					},
				}
				o.saveStats(&stat, db)
			}
		}
	}
	return
}

func (o *MongoOutput) trackGame(data *DecodedData, db *mgo.Database) (err error) {
	insertData := make(map[string]bson.M)
	insertData["$set"] = data.GameSet
	if len(data.GameInc) > 0 {
		insertData["$inc"] = data.GameInc
	}
	if len(data.GamePush) > 0 {
		insertData["$addToSet"] = data.GamePush
	}
	err = SaveWithRetry(db, o.GameCollection, data.GameSet["game_id"], insertData)
	if err != nil {
		o.or.LogError(err)
	}
	return
}

func (o *MongoOutput) trackStrategy(data *DecodedData, db *mgo.Database) (err error) {
	err = InsertNew(db, o.StrategyCollection, data.StrategySet)
	if err != nil {
		o.or.LogError(err)
	}
	return
}

func (o *MongoOutput) trackLog(data *DecodedData, db *mgo.Database) (err error) {
	insertData := make(map[string]bson.M)
	if len(data.LogPush) > 0 {
		insertData["$addToSet"] = data.LogPush
	}
	if len(data.LogSet) > 0 {
		insertData["$set"] = data.LogSet
	}
	err = SaveWithRetry(db, o.LogCollection, data.GameSet["game_id"], insertData)
	if err != nil {
		o.or.LogError(err)
	}
	return
}

func (o *MongoOutput) trackLogInfo(data *DecodedData, db *mgo.Database) (err error) {
	// track user login or out
	var (
		ok      bool
		userIp  string
		userAid string
	)
	result := make(map[string]interface{})
	err = db.C(o.TotalCollection).Find(bson.M{"_id": data.ID}).One(&result)
	if err == nil {
		if userIp, ok = result["active_ip"].(string); ok {
			data.UserRecord["active_ip"] = userIp
		}
		if userAid, ok = result["active_aid"].(string); ok {
			data.UserRecord["active_aid"] = userAid
		}
		if userAid, ok = result["active_chn"].(string); ok {
			data.UserRecord["active_chn"] = userAid
		}
		err = InsertNew(db, o.UserRecCollection, data.UserRecord)
		if err != nil {
			o.or.LogError(err)
		}
	}
	return
}

func (o *MongoOutput) trackGameOnlineInfo(data *DecodedData, db *mgo.Database) (err error) {
	var (
		time_stamp int
		game_id    int
	)

	if time_stamp, err = ConvertToInt(data.GameOnlineRecord["time"]); err != nil {
		o.or.LogError(err)
		return
	}

	if game_id, err = ConvertToInt(data.GameOnlineRecord["game_id"]); err != nil {
		o.or.LogError(err)
		return
	}

	key := fmt.Sprintf("%s-%s-%d-%d", data.GameOnlineRecord["mer_code"],
		data.GameOnlineRecord["brand_code"],
		game_id,
		time_stamp)

	insertData := make(map[string]bson.M)
	insertData["$set"] = data.GameOnlineRecord

	err = SaveWithRetry(db, o.GameOnlineCollection, key, insertData)
	if err != nil {
		o.or.LogError(err)
	}

	return
}

// CleanupForRestart is called before restart plugin
func (o *MongoOutput) CleanupForRestart() {
	return
}

// Process is called by goroutine
func (o *MongoOutput) Process(pack *p.PipelinePack, db *mgo.Database) (err error) {
	var (
		data DecodedData
	)
	payload := pack.Message.GetPayload()
	if err = json.Unmarshal([]byte(payload), &data); err != nil {
		return err
	}
	ConvertTime(&data.UserSet)
	ConvertTime(&data.GameSet)
	ConvertTime(&data.TaxSet)
	ConvertTime(&data.StrategySet)
	ConvertTime(&data.UserRecord)
	ConvertTime(&data.GameOnlineRecord)
	if data.Type != "announce" {
		err = o.trackLog(&data, db)
	}
	if data.Type == "start" {
		err = o.trackGame(&data, db)
	} else if data.Type == "apply" {
		err = o.trackApply(&data, db)
	} else if data.Type == "bet" {
		err = o.trackBet(&data, db)
	} else if data.Type == "compute" {
		err = o.trackAward(&data, db)
	} else if data.Type == "fold" {
		err = o.trackGame(&data, db)
	} else if data.Type == "compare" {
		err = o.trackGame(&data, db)
	} else if data.Type == "pk_all" {
		err = o.trackGame(&data, db)
	} else if data.Type == "check" {
		err = o.trackGame(&data, db)
	} else if data.Type == "changed" {
		err = o.trackGame(&data, db)
	} else if data.Type == "strategy_sys" || data.Type == "strategy_user" {
		err = o.trackStrategy(&data, db)
	} else if data.Type == "announce" {
		err = o.trackGame(&data, db)
	} else if data.Type == "login" || data.Type == "logout" || data.Type == "disconnect" || data.Type == "kick_out" || data.Type == "reconnect" {
		err = o.trackLogInfo(&data, db)
	} else if data.Type == "online" {
		err = o.trackGameOnlineInfo(&data, db)
	} else {
		o.or.LogError(fmt.Errorf("unknow msg type:%s", data.Type))
	}
	return
}

func init() {
	p.RegisterPlugin("FlowerOutput", func() interface{} {
		return new(MongoOutput)
	})
}
