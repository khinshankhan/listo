-- name: SearchNodes :many
-- Fuzzy title search over entries and groups, one hit per node by its best
-- matching title.
select n.public_id, n.slug, n.type,
       coalesce(e.title, g.name)::text as title,
       hit.matched::text as matched,
       g.kind as group_kind, m.label as medium_label, e.start_date,
       hit.score::float8 as score
from (
  select distinct on (t.node_id) t.node_id, t.title as matched,
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
order by hit.score desc, title
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
-- How many entries each medium has, for the browse filters.
select m.medium, m.label, count(e.id)::int as entries
from catalog.medium m
join catalog.entry e on e.medium = m.medium
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
