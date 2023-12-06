## Tools

#### pg_stat_statements;
```sql
CREATE EXTENSION pg_stat_statements;

SELECT query, total_exec_time, calls, rows
FROM pg_stat_statements
ORDER BY total_exec_time DESC;

select  from pg_stat_statements;
select query, total_exec_time, mean_exec_time, calls, rows from pg_stat_statements where query ilike '%from users%';

-- reset
SELECT pg_stat_statements_reset();
```

#### https://pganalyze.com/