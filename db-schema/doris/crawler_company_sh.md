CREATE TABLE `crawler_company_sh` ( 
  `id` VARCHAR(64) NOT NULL COMMENT '主键ID' ,
  `company_id` VARCHAR(255) NULL COMMENT '公司ID' ,
  `company_name` VARCHAR(255) NULL COMMENT '公司名称' ,
  `in_out` VARCHAR(64) NULL COMMENT '投资方向' ,
  `nodes_map` JSON NULL COMMENT '投资对象' ,
  `tree_list` JSON NULL COMMENT '投资关系' ,
  `created_at` DATETIME NULL COMMENT '创建时间' ,
  `updated_at` DATETIME NULL COMMENT '更新时间' 
)
ENGINE = olap
DUPLICATE KEY(`id`)
DISTRIBUTED BY HASH(`id`) BUCKETS 10
PROPERTIES (
"replication_num" = "1"   
);
