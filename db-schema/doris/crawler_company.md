CREATE TABLE `crawler_company` (
  `id` VARCHAR(64) NOT NULL COMMENT '主键ID' ,
  `name` VARCHAR(255) NULL COMMENT '名称' ,
  `reg_no` VARCHAR(255) NULL COMMENT '工商注册号' ,
  `oper_status` VARCHAR(64) NULL COMMENT '经营状态' ,
  `credit_code` VARCHAR(64) NULL COMMENT '统一社会信用代码' ,
  `iden` VARCHAR(64) NULL COMMENT '纳税人识别号' ,
  `appr_date` VARCHAR(64) NULL COMMENT '核准日期' ,
  `create_date` VARCHAR(64) NULL COMMENT '成立时间' ,
  `qual` VARCHAR(255) NULL COMMENT '纳税人资质' ,
  `ent_type` VARCHAR(255) NULL COMMENT '企业类型' ,
  `reg_authority` VARCHAR(255) NULL COMMENT '登记机关' ,
  `eng_name` VARCHAR(255) NULL COMMENT '英文名称' ,
  `addr` VARCHAR(500) NULL COMMENT '注册地址' ,
  `specific_addr` VARCHAR(500) NULL COMMENT '详细地址' ,
  `legal_person` VARCHAR(255) NULL COMMENT '法人' ,
  `bus_info` JSON NULL COMMENT '工商信息' ,
  `sh_info` JSON NULL COMMENT '股东信息' ,
  `key_person` JSON NULL COMMENT '主要人员' ,
  `is_transform` TINYINT NULL DEFAULT 0  COMMENT '是否转换' ,
  `source_url` VARCHAR(500) NULL COMMENT '来源URL' ,
  `source_id` VARCHAR(64) NULL COMMENT '来源ID(业务唯一键)' ,
  `source_system` VARCHAR(64) NULL COMMENT '来源系统' ,
  `remark` STRING NULL COMMENT '备注' ,
  `created_at` DATETIME NULL COMMENT '创建时间' ,
  `updated_at` DATETIME NULL COMMENT '更新时间' ,
  `prev_updated_at` DATETIME NULL COMMENT '上次更新时间'
)
ENGINE = olap
DUPLICATE KEY(`id`)
DISTRIBUTED BY HASH(`id`) BUCKETS 10
PROPERTIES (
"replication_num" = "1"   
);