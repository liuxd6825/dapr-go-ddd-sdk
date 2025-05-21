package xtest

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/jsonschema/v6"
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
		  "title": "ID",
          "order": 1
		},
		"tenantId": {
		  "type": ["string"],
		  "title": "租户ID",
			"order": 2
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
		  "title": "备注",
           "order": 107
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

const CompanySchema = `
{
  "name": "company",
  "type": "object",
  "description": "公司基本信息",
  "meta": {
    "ddd": {
      "aggField": "id",
      "aggType": "master.company",
      "isPubEvent": false
    },
    "dbTable": {
      "name": "company",
      "dbKey": "mysql",
      "properties": {
        "neo4j": {
          "type": "node"
        }
      }
    },
    "form": {},
    "attributes": {
      "neo4j":{
        "labels": ["company"],
        "type": ["node"]
      }
    }
  },
  "properties": {
    "id": {
      "name": "id",
      "type": [
        "string"
      ],
      "title": "ID",
      "order": 1001,
      "meta": {
        "dbField": {
          "primaryKey": true,
          "creatable": false,
          "updatable": false,
          "readable": true,
          "relEndId": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "tenantId": {
      "name": "tenantId",
      "title": "租户Id",
      "type": [
        "string",
        "null"
      ],
      "order": 1002,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true,
          "nodeLabelFormat": "tenant_%s",
          "nodeLabel": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "remark": {
      "name": "remark",
      "type": [
        "string",
        "null"
      ],
      "title": "备注",
      "order": 9001,
      "meta": {
        "form": {
          "ctlType": "ui5-input"
        },
        "dbField": {
          "size": 200
        }
      }
    },
    "createdTime": {
      "name": "createdTime",
      "type": [
        "datetime",
        "null"
      ],
      "title": "创建时间",
      "order": 9002,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "creatorId": {
      "name": "creatorId",
      "type": [
        "string",
        "null"
      ],
      "title": "创建人Id",
      "order": 9003,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "creatorName": {
      "name": "creatorName",
      "type": [
        "string",
        "null"
      ],
      "title": "创建人",
      "order": 9004,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "updatedTime": {
      "name": "updatedTime",
      "type": [
        "datetime",
        "null"
      ],
      "title": "修改时间",
      "order": 9005,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "updaterId": {
      "name": "updaterId",
      "type": [
        "string",
        "null"
      ],
      "title": "修改人Id",
      "order": 9006,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "updaterName": {
      "name": "updaterName",
      "type": [
        "string",
        "null"
      ],
      "title": "修改人",
      "order": 9007,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "deletedTime": {
      "name": "deletedTime",
      "type": [
        "datetime",
        "null"
      ],
      "title": "删除时间",
      "order": 9008,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "deleterId": {
      "name": "deleterId",
      "type": [
        "string",
        "null"
      ],
      "title": "删除人Id",
      "order": 9009,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "deleterName": {
      "name": "deleterName",
      "type": [
        "string",
        "null"
      ],
      "title": "删除人",
      "order": 9010,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "isDeleted": {
      "name": "isDeleted",
      "type": [
        "boolean",
        "null"
      ],
      "title": "是否删除",
      "order": 9011,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "caseId": {
      "name": "caseId",
      "title": "项目Id",
      "type": "string",
      "order": 100,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true,
          "nodeLabelFormat": "case_%s",
          "nodeLabel": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "name": {
      "name": "name",
      "title": "公司名称",
      "type": [
        "string",
        "null"
      ],
      "order": 1000,
      "meta": {
        "dbField": {
          "creatable": true,
          "updatable": true,
          "readable": true
        },
        "form": {
          "required": true
        },
        "query": {
          "allowQuery": true
        }
      }
    },
    "regNo": {
      "name": "regNo",
      "title": "工商注册号",
      "type": [
        "string",
        "null"
      ],
      "order": 1001
    },
    "operatingStatus": {
      "name": "operatingStatus",
      "title": "经营状态",
      "type": [
        "string",
        "null"
      ],
      "order": 1004
    },
    "creditCode": {
      "name": "creditCode",
      "title": "统一社会信用代码",
      "type": [
        "string",
        "null"
      ],
      "order": 1005
    },
    "identification": {
      "name": "identification",
      "title": "纳税人识别号",
      "type": [
        "string",
        "null"
      ],
      "order": 1006
    },
    "orgCode": {
      "title": "组织机构代码",
      "type": [
        "string",
        "null"
      ],
      "order": 1007
    },
    "businessTerm": {
      "title": "营业期限",
      "type": [
        "string",
        "null"
      ],
      "order": 1008
    },
    "qualification": {
      "title": "纳税人资质",
      "type": [
        "string",
        "null"
      ],
      "order": 1009
    },
    "approvalDate": {
      "title": "核准日期",
      "type": [
        "date",
        "null"
      ],
      "order": 1010
    },
    "createDate": {
      "title": "成立时间",
      "type": [
        "date",
        "null"
      ],
      "order": 1011
    },
    "entType": {
      "title": "企业类型",
      "type": [
        "string",
        "null"
      ],
      "order": 1012
    },
    "industry": {
      "title": "行业",
      "type": [
        "string",
        "null"
      ],
      "order": 1013
    },
    "personnelSize": {
      "title": "人员规模",
      "type": [
        "integer",
        "null"
      ],
      "order": 1014
    },
    "insuredPersons": {
      "title": "参保人数",
      "type": [
        "string",
        "null"
      ],
      "order": 1015
    },
    "regAuthority": {
      "title": "登记机关",
      "type": [
        "string",
        "null"
      ],
      "order": 1016
    },
    "nameUsed": {
      "title": "曾用名",
      "type": [
        "string",
        "null"
      ],
      "order": 1017
    },
    "engName": {
      "title": "英文名称",
      "type": [
        "string",
        "null"
      ],
      "order": 1018
    },
    "addr": {
      "name": "addr",
      "title": "注册地址",
      "type": [
        "string",
        "null"
      ],
      "order": 1019
    },
    "addrType": {
      "title": "地址类型",
      "type": [
        "string",
        "null"
      ],
      "order": 1021
    },
    "regCapitalType": {
      "name": "regCapitalType",
      "title": "注册资本类型",
      "type": [
        "string",
        "null"
      ],
      "order": 1022
    },
    "regCapital": {
      "name": "regCapital",
      "title": "注册资本",
      "type": [
        "string",
        "null"
      ],
      "order": 1023
    },
    "paidInCapital": {
      "name": "paidInCapital",
      "title": "实缴资本",
      "type": [
        "string",
        "null"
      ],
      "order": 1024
    },
    "remark": {
      "name": "remark",
      "type": [
        "string",
        "null"
      ],
      "title": "备注",
      "order": 9001,
      "meta": {
        "form": {
          "ctlType": "ui5-input",
          "colSpan": "S1 M2 L2 XL2"
        },
        "dbField": {
          "size": 200
        }
      }
    },
    "natureOfBusiness": {
      "title": "经营范围",
      "type": [
        "string",
        "null"
      ],
      "order": 9100,
      "meta": {
        "dbField": {
          "size": 2000
        },
        "form": {
          "colSpan": "-1",
          "ctlType": "ui5-textarea",
          "maxlength": "100",
          "rows": "8"
        },
        "column": {
          "width": "300"
        }
      }
    }
  },
  "order": 1001
}
`

const CompanyCompanySchema = `
{
  "name": "company_company",
  "type": "object",
  "description": "合同公司关系",
  "meta": {
    "ddd": {
      "aggField": "companyId",
      "aggType": "master.company",
      "isPubEvent": false
    },
    "dbTable": {
      "name": "company_company",
      "dbKey": "mysql"
    },
    "attributes": {
      "neo4j":{
        "labels": ["company"],
        "type": ["rel","node"]
      }
    }
  },
  "properties": {
    "id": {
      "name": "id",
      "type": [
        "string"
      ],
      "title": "ID",
      "order": 1001,
      "meta": {
        "dbField": {
          "primaryKey": true,
          "creatable": false,
          "updatable": false,
          "readable": true,
          "relEndId": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "tenantId": {
      "name": "tenantId",
      "title": "租户Id",
      "type": [
        "string",
        "null"
      ],
      "order": 1002,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true,
          "nodeLabelFormat": "tenant_%s",
          "nodeLabel": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "remark": {
      "name": "remark",
      "type": [
        "string",
        "null"
      ],
      "title": "备注",
      "order": 9001,
      "meta": {
        "form": {
          "ctlType": "ui5-input"
        },
        "dbField": {
          "size": 200
        }
      }
    },
    "createdTime": {
      "name": "createdTime",
      "type": [
        "datetime",
        "null"
      ],
      "title": "创建时间",
      "order": 9002,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "creatorId": {
      "name": "creatorId",
      "type": [
        "string",
        "null"
      ],
      "title": "创建人Id",
      "order": 9003,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "creatorName": {
      "name": "creatorName",
      "type": [
        "string",
        "null"
      ],
      "title": "创建人",
      "order": 9004,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "updatedTime": {
      "name": "updatedTime",
      "type": [
        "datetime",
        "null"
      ],
      "title": "修改时间",
      "order": 9005,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "updaterId": {
      "name": "updaterId",
      "type": [
        "string",
        "null"
      ],
      "title": "修改人Id",
      "order": 9006,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "updaterName": {
      "name": "updaterName",
      "type": [
        "string",
        "null"
      ],
      "title": "修改人",
      "order": 9007,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "deletedTime": {
      "name": "deletedTime",
      "type": [
        "datetime",
        "null"
      ],
      "title": "删除时间",
      "order": 9008,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "deleterId": {
      "name": "deleterId",
      "type": [
        "string",
        "null"
      ],
      "title": "删除人Id",
      "order": 9009,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "deleterName": {
      "name": "deleterName",
      "type": [
        "string",
        "null"
      ],
      "title": "删除人",
      "order": 9010,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "isDeleted": {
      "name": "isDeleted",
      "type": [
        "boolean",
        "null"
      ],
      "title": "是否删除",
      "order": 9011,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true
        },
        "column": {
          "hide": true
        },
        "form": {
          "hide": true
        },
        "query": {
          "allowQuery": false
        }
      }
    },
    "caseId": {
      "name": "caseId",
      "title": "项目Id",
      "type": "string",
      "order": 2001,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true,
          "nodeLabelFormat": "case_%s",
          "nodeLabel": true
        },
        "column": {
          "hide": true
        }
      }
    },
    "companyId": {
      "name": "companyId",
      "title": "公司ID",
      "type": "string",
      "order": 2001,
      "meta": {
        "dbField": {
          "creatable": false,
          "updatable": false,
          "readable": true,
          "relStartId": true
        },
        "column": {
          "hide": true
        }
      }
    },
    "relationType": {
      "name": "relationType",
      "order": 108,
      "title": "关系类型",
      "type": ["string"],
      "meta": {
        "column": {
          "type": "dropdown",
          "source": "[\"控股\", \"子公司\", \"分公司\"]"
        },
        "dbField": {
          "relType": true
        }
      }
    },
    "name": {
      "name": "name",
      "type": ["string", "null"],
      "title": "公司名称",
      "order": 2003
    },
    "legalPerson": {
      "name": "legalPerson",
      "type": ["string", "null"],
      "title": "法人",
      "order": 2007
    },
    "phone": {
      "name": "phone",
      "type": ["string", "null"],
      "title": "预留电话",
      "order": 2005
    },
    "identNum": {
      "name": "identNum",
      "type": ["string", "null"],
      "title": "纳税人识别号",
      "order": 2006
    },
    "addr": {
      "name": "addr",
      "type": ["string", "null"],
      "title": "注册地址",
      "order": 2009
    }
  }
}

`

func GetHumanSchema() *jsonschema.Schema {
	return schema.NewJsonSchemaWithJson("human.json", HumanSchema)
}

func GetCompanySchema() *jsonschema.Schema {
	return schema.NewJsonSchemaWithJson("company.json", CompanySchema)
}

func GetCompanyCompanySchema() *jsonschema.Schema {
	return schema.NewJsonSchemaWithJson("company_company.json", CompanyCompanySchema)
}
