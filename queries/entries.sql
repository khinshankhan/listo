-- name: EntryByPublicID :one
-- An entry by public id, with its medium and format labels and its cover:
-- the canon one, else the widest. The caller checks the slug.
select n.id, n.public_id, n.slug, e.medium, m.label as medium_label, e.format,
       mf.label as format_label, e.country, e.title, e.description, e.status,
       e.start_date, e.start_precision, e.end_date, e.end_precision,
       e.unit_kind, e.unit_total, e.unit_available, e.unit_minutes, e.length_minutes, e.source_type,
       cover.storage_key as cover_key, cover.width as cover_width, cover.height as cover_height, cover.fit as cover_fit
from catalog.node n
join catalog.entry e on e.id = n.id
join catalog.medium m on m.medium = e.medium
join catalog.medium_format mf on mf.medium = e.medium and mf.format = e.format
left join catalog.entry_image cover on cover.id = (
  select i.id
  from catalog.entry_image i
  where i.entry_id = e.id and i.kind = 'cover'
  order by i.is_canon desc, i.width desc, i.id
  limit 1
)
where n.public_id = $1;

-- name: NodeTitles :many
-- Every title a node goes by, the likeliest to be searched for first.
select kind, language, title
from catalog.title
where node_id = $1
order by array_position(array['english', 'native', 'romanized', 'localized', 'synonym'], kind), title;

-- name: NodeLinks :many
-- Jump-out links: only ids are stored, URLs come from the id_source vocabulary.
select s.site, x.source, x.external_id,
       replace(s.url_template, '{id}', x.external_id)::text as url
from catalog.external_id x
join catalog.id_source s on s.source = x.source
where x.node_id = $1
order by s.site, x.source;

-- name: NodeRelations :many
-- A node's relations, labelled and anchored from its own side, an entry with
-- its cover.
-- NOTE: spoilers are left out until there's UI to hide them.
select (case when r.from_node = sqlc.arg(node_id)::uuid then rk.label
             else coalesce(rk.inverse_label, rk.label) end)::text as reading,
       rk.family, o.public_id, o.slug, o.type,
       coalesce(e.title, g.name)::text as title,
       m.label as medium_label, e.start_date, r.note,
       mine.label as my_unit, theirs.label as their_unit,
       cover.storage_key as cover_key, cover.width as cover_width, cover.height as cover_height, cover.fit as cover_fit
from catalog.relation r
join catalog.relation_kind rk on rk.kind = r.kind
join catalog.node o on o.id = case when r.from_node = sqlc.arg(node_id)::uuid then r.to_node else r.from_node end
left join catalog.entry e on e.id = o.id
left join catalog.grp g on g.id = o.id
left join catalog.medium m on m.medium = e.medium
left join catalog.unit mine on mine.id = case when r.from_node = sqlc.arg(node_id)::uuid then r.from_unit else r.to_unit end
left join catalog.unit theirs on theirs.id = case when r.from_node = sqlc.arg(node_id)::uuid then r.to_unit else r.from_unit end
left join catalog.entry_image cover on cover.id = (
  select i.id
  from catalog.entry_image i
  where i.entry_id = e.id and i.kind = 'cover'
  order by i.is_canon desc, i.width desc, i.id
  limit 1
)
where sqlc.arg(node_id)::uuid in (r.from_node, r.to_node)
  and not r.spoiler
order by rk.family, reading, e.start_date nulls last, coalesce(e.title, g.name);

-- name: NodeGroups :many
-- Every group a node is in, nearest first, with a cycle guard. Only direct
-- groups have a position.
with recursive up (id, depth) as (
  select m.group_id, 1
  from catalog.membership m
  where m.member_id = $1
  union all
  select m.group_id, up.depth + 1
  from up
  join catalog.membership m on m.member_id = up.id
) cycle id set looped using path
select n.public_id, n.slug, g.kind, g.name, ms.position_label, min(up.depth)::int as depth
from up
join catalog.grp g on g.id = up.id
join catalog.node n on n.id = g.id
left join catalog.membership ms on ms.group_id = g.id and ms.member_id = $1
where not up.looped
group by n.public_id, n.slug, g.kind, g.name, ms.position_label
order by depth, g.name;

-- name: NodeTags :many
select t.kind, t.name
from catalog.node_tag nt
join catalog.tag t on t.id = nt.tag_id
where nt.node_id = $1
order by t.kind, t.name;

-- name: EntryCompanies :many
-- Who made an entry, by role then name.
select c.name, ec.role
from catalog.entry_company ec
join catalog.company c on c.id = ec.company_id
where ec.entry_id = $1
order by array_position(array['studio', 'developer', 'publisher', 'producer', 'licensor'], ec.role), c.name;

-- name: EntryUnits :many
-- An entry's units in order, for the episode list and what airs next.
select ordinal, kind, number, label, title, air_date, minutes
from catalog.unit
where entry_id = $1
order by ordinal;
