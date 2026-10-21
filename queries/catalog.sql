-- name: SearchNodes :many
-- Fuzzy title search over entries and groups, one hit per node by its best
-- matching title.
select n.public_id, n.slug, n.type,
       coalesce(e.title, g.name)::text as title,
       hit.matched::text as matched, hit.matched_kind::text as matched_kind,
       g.kind as group_kind, e.medium, m.label as medium_label, e.format, mf.label as format_label,
       e.start_date, e.start_precision, e.unit_kind, e.unit_total,
       hit.score::float8 as score
from (
  select distinct on (t.node_id) t.node_id, t.title as matched, t.kind as matched_kind,
         greatest(similarity(t.title_norm, catalog.norm(sqlc.arg(q)::text)),
                  word_similarity(catalog.norm(sqlc.arg(q)::text), t.title_norm)) as score
  from catalog.title t
  where t.title_norm % catalog.norm(sqlc.arg(q)::text)
     or catalog.norm(sqlc.arg(q)::text) <% t.title_norm
  order by t.node_id, score desc
) hit
join catalog.node n on n.id = hit.node_id
left join catalog.entry e on e.id = n.id
left join catalog.grp g on g.id = n.id
left join catalog.medium m on m.medium = e.medium
left join catalog.medium_format mf on mf.medium = e.medium and mf.format = e.format
order by hit.score desc, title
limit sqlc.arg(max_results)::int;

-- name: TitleSuggestions :many
-- Titles near a query that found nothing, closest first, scored both ways
-- round: the query within a title, and a title within the query. It scans
-- every title, which is fine while the catalog is small and it only runs on
-- a miss.
select near.title
from (
  select t.title,
         greatest(similarity(t.title_norm, catalog.norm(sqlc.arg(q)::text)),
                  word_similarity(catalog.norm(sqlc.arg(q)::text), t.title_norm),
                  word_similarity(t.title_norm, catalog.norm(sqlc.arg(q)::text))) as score
  from catalog.title t
) near
where near.score >= sqlc.arg(at_least)::float8
group by near.title
order by max(near.score) desc, near.title
limit sqlc.arg(max_results)::int;

-- name: TopGroups :many
-- Top-level groups with every entry beneath them counted.
-- NOTE: union, not union all, so a membership cycle ends.
with recursive down (root, id) as (
  select g.id, g.id
  from catalog.grp g
  where not exists (select 1 from catalog.membership ms where ms.member_id = g.id)
  union
  select down.root, ms.member_id
  from down
  join catalog.membership ms on ms.group_id = down.id
)
select n.public_id, n.slug, g.kind, g.name, count(e.id)::int as entries
from down
join catalog.grp g on g.id = down.root
join catalog.node n on n.id = g.id
left join catalog.entry e on e.id = down.id
group by n.public_id, n.slug, g.kind, g.name
order by g.name;

-- name: MediumCounts :many
-- How many entries each medium has, none included.
select m.medium, m.label, count(e.id)::int as entries
from catalog.medium m
left join catalog.entry e on e.medium = m.medium
group by m.medium, m.label
order by entries desc, m.label;

-- name: BrowseEntries :many
-- Every entry, or one medium's, alphabetically.
select n.public_id, n.slug, e.title, m.label as medium_label, e.start_date
from catalog.entry e
join catalog.node n on n.id = e.id
join catalog.medium m on m.medium = e.medium
where sqlc.narg(medium)::text is null or e.medium = sqlc.narg(medium)::text
order by e.title, e.start_date nulls last;
