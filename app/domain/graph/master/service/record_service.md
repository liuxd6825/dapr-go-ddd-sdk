```neo4j
CALL apoc.periodic.iterate(

	"CALL apoc.load.parquet('s3://minioadmin:minioadmin@192.168.120.224:9000/import-data/record/record_hUe2yxR81YNB89AY3sWRIM4U.parquet') YIELD value AS row RETURN row",
	"WITH row,
	      row.tenant_id  AS tenantId,
	      row.case_id    AS caseId,
	      row.acct       AS acct,
	      row.opp_acct   AS oppAcct,
	      date(row.date) AS txDate,
	      row.name       AS name,
	      row.opp_name   AS oppName,
	      coalesce(toFloat(row.amount), 0.0) AS amount,
	      coalesce(toFloat(row.payout), 0.0) AS payout,
	      coalesce(toFloat(row.income),  0.0) AS income,
	      row.ccy AS ccy
	 WHERE tenantId IS NOT NULL AND caseId IS NOT NULL
	   AND acct IS NOT NULL AND oppAcct IS NOT NULL AND txDate IS NOT NULL
	 WITH *, ':tenant_' + tenantId + ':case_' + caseId + ':record' AS labels
	 MERGE (n1:human{id: coalesce(name, acct)})
	   ON CREATE SET n1.name = name, n1.tenant_id = tenantId, n1.case_id = caseId
	 MERGE (a1:account{id: acct})
	   ON CREATE SET a1.name = acct, a1.tenant_id = tenantId, a1.case_id = caseId
	 MERGE (n1)-[:owner]->(a1)
	 MERGE (n2:human{id: coalesce(oppName, oppAcct)})
	   ON CREATE SET n2.name = oppName, n2.tenant_id = tenantId, n2.case_id = caseId
	 MERGE (a2:account{id: oppAcct})
	   ON CREATE SET a2.name = oppAcct, a2.tenant_id = tenantId, a2.case_id = caseId
	 MERGE (n2)-[:owner]->(a2)
	 WITH a1, a2, tenantId, caseId, acct, oppAcct, txDate, amount, payout, income, ccy
	 MERGE (a1)-[r:record{
	      id: 'rec:' + tenantId + ':' + caseId + ':' + acct + ':' + oppAcct + ':' + toString(txDate)
	   }]->(a2)
	 ON CREATE SET
	   r.created_time = timestamp(),
	   r.tenant_id = tenantId, r.case_id = caseId,
	   r.acct = acct, r.opp_acct = oppAcct,
	   r.date = txDate,
	   r.amount = amount, r.payout = payout, r.income = income,
	   r.txn_count = 1, r.ccy = ccy
	 ON MATCH SET
	   r.amount = coalesce(r.amount, 0.0) + amount,
	   r.payout = coalesce(r.payout, 0.0) + payout,
	   r.income  = coalesce(r.income,  0.0) + income,
	   r.txn_count = coalesce(r.txn_count, 0) + 1,
	   r.updated_time = timestamp()",
	{batchSize: 500, parallel: false, iterateList: true,
	 concurrency: 1, retries: 1}

) YIELD batches, total, committedOperations, failedOperations, timeTaken, errorMessages
RETURN batches, total, committedOperations, failedOperations, timeTaken, errorMessages;
```