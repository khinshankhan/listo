-- name: EntryByPublicID :one
-- An entry by public id, with its format label. The caller checks the slug.
select n.public_id, n.slug, e.medium, e.format, mf.label as format_label, e.country,
       e.title, e.start_date, e.start_precision, e.unit_kind, e.unit_total
from catalog.node n
join catalog.entry e on e.id = n.id
join catalog.medium_format mf on mf.medium = e.medium and mf.format = e.format
where n.public_id = $1;
