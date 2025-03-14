package tests

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsm"
	"time"
)

const HumanSchema = `
	{
	  "name": "human",
	  "type": "object",
	  "required": ["id","name"],
	  "description": "人员基本信息",
	  "properties": {
		"id": {
		  "type": ["string"],
		  "title": "ID"
		},
		"age": {
		  "type": ["integer", "null"],
		  "title": "年龄",
		  "order": 100
		},
		"analyse": {
		  "type": ["string", "null"],
		  "title": "分析状态",
		  "order": 101
		},
		"birthday": {
		  "type": ["date", "null"],
		  "title": "出生日期",
		  "order": 102
		},
		"gender": {
		  "type": ["string", "null"],
		  "title": "性别",
		  "order": 103
		},
		"name": {
		  "type": ["string", "null"],
		  "title": "姓名",
		  "order": 104
		},
		"peopleType": {
		  "type": "array",
		  "items": { "type": "string" },
		  "title": "人员类型",
		  "order": 105
		},
		"tags": {
		  "title": "标签",
		  "type": "array",
		  "items": { "type": "string" },
		  "order": 106
		},
		"remark": {
		  "name": "remark",
		  "type": ["string", "null"],
		  "title": "备注"
		},
		"createdTime": {
		  "name": "createdTime",
		  "type": ["datetime"],
		  "title": "创建时间",
		  "order": 9002,
          "dbField":{
		    "updatable":false
		  }
		},
		"creatorId": {
		  "name": "creatorId",
		  "type": ["string","null"],
		  "title": "创建人Id",
		  "notNull": true,
		  "order": 9003,
          "dbField":{
		    "updatable":false
		  }
		},
		"creatorName": {
		  "name": "creatorName",
		  "type": ["string","null"],
		  "title": "创建人",
		  "notNull": true,
		  "order": 9004,
          "dbField":{
		    "updatable":false
		  }
		},
		"updatedTime": {
		  "name": "updatedTime",
		  "type": ["datetime","null"],
		  "title": "修改时间",
		  "notNull": true,
		  "order": 9005
		},
		"updaterId": {
		  "name": "updaterId",
		  "type": ["string","null"],
		  "title": "修改人Id",
		  "notNull": true,
		  "order": 9006
		},
		"updaterName": {
		  "name": "updaterName",
		  "type": ["string","null"],
		  "title": "修改人",
		  "notNull": true,
		  "order": 9007
		}
	  }
	}
	`

const HumanRelSchema = `
	{
	  "name": "humanRel",
	  "type": "object",
	  "required": ["id","name"],
	  "description": "人员基本信息",
	  "properties": {
		"id": {
		  "type": ["string"],
		  "title": "ID"
		},
		"startId": {
		  "type": ["string"],
		  "title": "startId"
		},
		"endId": {
		  "type": ["string"],
		  "title": "endId"
		},
		"relType": {
		  "type": ["string"],
		  "title": "relType"
		},
		"age": {
		  "type": ["integer", "null"],
		  "title": "年龄",
		  "order": 100
		},
		"analyse": {
		  "type": ["string", "null"],
		  "title": "分析状态",
		  "order": 101
		},
		"birthday": {
		  "type": ["date", "null"],
		  "title": "出生日期",
		  "order": 102
		},
		"gender": {
		  "type": ["string", "null"],
		  "title": "性别",
		  "order": 103
		},
		"name": {
		  "type": ["string", "null"],
		  "title": "姓名",
		  "order": 104
		},
		"peopleType": {
		  "type": "array",
		  "items": { "type": "string" },
		  "title": "人员类型",
		  "order": 105
		},
		"tags": {
		  "title": "标签",
		  "type": "array",
		  "items": { "type": "string" },
		  "order": 106
		},
		"remark": {
		  "name": "remark",
		  "type": ["string", "null"],
		  "title": "备注"
		},
		"createdTime": {
		  "name": "createdTime",
		  "type": ["datetime"],
		  "title": "创建时间",
		  "order": 9002,
          "dbField":{
		    "updatable":false
		  }
		},
		"creatorId": {
		  "name": "creatorId",
		  "type": ["string","null"],
		  "title": "创建人Id",
		  "notNull": true,
		  "order": 9003,
          "dbField":{
		    "updatable":false
		  }
		},
		"creatorName": {
		  "name": "creatorName",
		  "type": ["string","null"],
		  "title": "创建人",
		  "notNull": true,
		  "order": 9004,
          "dbField":{
		    "updatable":false
		  }
		},
		"updatedTime": {
		  "name": "updatedTime",
		  "type": ["datetime","null"],
		  "title": "修改时间",
		  "notNull": true,
		  "order": 9005
		},
		"updaterId": {
		  "name": "updaterId",
		  "type": ["string","null"],
		  "title": "修改人Id",
		  "notNull": true,
		  "order": 9006
		},
		"updaterName": {
		  "name": "updaterName",
		  "type": ["string","null"],
		  "title": "修改人",
		  "notNull": true,
		  "order": 9007
		}
	  }
	}
	`

type Human struct {
	Id         string    `gorm:"primaryKey"`
	Name       string    `gorm:"name"`
	Age        int       `gorm:"age"`
	Analyse    string    `gorm:"analyse"`
	Birthday   time.Time `gorm:"birthday"`
	PeopleType []string  `gorm:"people_type;type:text[]"`
}

type EnvConfig struct {
}

func (e EnvConfig) GetAppId() string {
	return "test"
}

func (e EnvConfig) GetAppName() string {
	return "app"
}

func (e EnvConfig) GetAppHttpHost() string {
	return "localhost"
}

func (e EnvConfig) GetAppHttpPort() int {
	return 0
}

func (e EnvConfig) GetDaprHost() string {
	return "localhost"
}

func (e EnvConfig) GetDaprHttpPort() int64 {
	return 0
}

func (e EnvConfig) GetDaprGrpcPort() int64 {
	return 0
}

func (e EnvConfig) GetFsManager() *fsm.Manager {
	return nil
}

func (e EnvConfig) GetHServerSrcPath() string {
	return ""
}

func (e EnvConfig) GetHServerEnable() bool {
	return false
}

func NewEnvConfig() *EnvConfig {
	return &EnvConfig{}
}
