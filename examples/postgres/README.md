# Postgres

## Use Case
- CRUD
- Soft Delete
- Exist
```postgresql
select exists (select 1 from checkout_hours where checkout_id = $1
and treatment_id = $2 and hours_id = $3)
```
- Round
```postgresql
    round(extract(epoch from cs.end_time - cs.start_time) / 60) as duration,
```
- Min Over Partition
```postgresql
    min(cs.start_time) over (partition by s.id, c.book_date) as min_start_time,
```
- Max Over Partition
```postgresql
    max(cs.end_time) over (partition by s.id, c.book_date) as max_end_time,
```
- Greatest
```postgresql
    greatest(
        round(extract(epoch from cs.end_time - now()::time) / 60), -1
    ) as remaining_minute,
```
- JSON single or array
- Group By