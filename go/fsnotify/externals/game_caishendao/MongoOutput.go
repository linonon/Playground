package caishendao

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //some special mysql type
	p "github.com/mozilla-services/heka/pipeline"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoOutput is main data of this package
type MongoOutput struct {
	*MongoOutputConfig
	Mg       *mgo.Session
	Db       *gorm.DB
	PackChan chan *p.PipelinePack
	wg       *sync.WaitGroup
	or       p.OutputRunner
}

//MongoOutputConfig is config from toml file
type MongoOutputConfig struct {
	ProjectId            string   `toml:"project_id"`
	MongoAddr            []string `toml:"mongo_addr"`
	ConnectTimeout       uint32   `toml:"connect_timeout"`
	Procs                uint32   `toml:"procs"`
	DbName               string   `toml:"db_name"`
	TotalCollection      string   `toml:"total"`
	DailyCollection      string   `toml:"daily"`
	GameCollection       string   `toml:"game"`
	UserRecCollection    string   `toml:"user_record"`
	GameOnlineCollection string   `toml:"game_online_record"`
	Debug                bool     `toml:"debug"`
}

// ConfigStruct defines default configs
func (o *MongoOutput) ConfigStruct() interface{} {
	return &MongoOutputConfig{
		ProjectId:            "komoku",
		MongoAddr:            []string{"localhost:27017"},
		ConnectTimeout:       3,
		Procs:                5,
		DbName:               "game",
		TotalCollection:      "user_stats",
		DailyCollection:      "daily_stats",
		GameCollection:       "caishendao",
		UserRecCollection:    "user_record",
		GameOnlineCollection: "game_online_record",
		Debug:                false,
	}
}

//Init is a `must` method for this plugin, called by heka
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

	insertData := make(map[string]bson.M)
	if len(data.Set) > 0 {
		insertData["$set"] = data.Set
	}
	if len(data.Inc) > 0 {
		insertData["$inc"] = data.Inc
	}
	err = SaveWithRetry(db, o.TotalCollection, data.ID, insertData)
	if err != nil {
		o.or.LogError(err)
	}

	temp = data.Set["updated_at"]
	if updatedAt, ok = temp.(time.Time); !ok {
		return fmt.Errorf("can't convert updated time: %v", data.Set)
	}
	key := fmt.Sprintf("%d-%s", data.ID, GetLocalStr(updatedAt))
	insertData["$set"]["user_id"] = data.ID
	insertData["$set"]["day"] = GetLocalStr(updatedAt)
	err = SaveWithRetry(db, o.DailyCollection, key, insertData)
	if err != nil {
		o.or.LogError(err)
	}

	err = InsertNew(db, o.GameCollection, data.Data)
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

// Process is called by goroutine
func (o *MongoOutput) Process(pack *p.PipelinePack, db *mgo.Database) (err error) {
	var data DecodedData
	payload := pack.Message.GetPayload()
	if err = json.Unmarshal([]byte(payload), &data); err != nil {
		return err
	}
	ConvertTime(&data.Set)
	ConvertTime(&data.Data)
	ConvertTime(&data.UserRecord)
	ConvertTime(&data.GameOnlineRecord)
	if data.Type == "spin" {
		err = o.saveStats(&data, db)
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
	p.RegisterPlugin("CaishendaoOutput", func() interface{} {
		return new(MongoOutput)
	})
}
