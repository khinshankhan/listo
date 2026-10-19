-- name: EntryByPublicID :one
-- An entry by public id, with its medium and format labels. The caller
-- checks the slug.
select n.id, n.public_id, n.slug, e.medium, m.label as medium_label, e.format,
       mf.label as format_label, e.country, e.title, e.start_date, e.start_precision,
       e.unit_kind, e.unit_total
from catalog.node n
join catalog.entry e on e.id = n.id
join catalog.medium m on m.medium = e.medium
join catalog.medium_format mf on mf.medium = e.medium and mf.format = e.format
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
-- A node's relations, labelled and anchored from its own side.
-- NOTE: spoilers are left out until there's UI to hide them.
select (case when r.from_node = sqlc.arg(node_id)::uuid then rk.label
             else coalesce(rk.inverse_label, rk.label) end)::text as reading,
       rk.family, o.public_id, o.slug, o.type,
       coalesce(e.title, g.name)::text as title,
       m.label as medium_label, e.start_date, r.note,
       mine.label as my_unit, theirs.label as their_unit
from catalog.relation r
join catalog.relation_kind rk on rk.kind = r.kind
join catalog.node o on o.id = case when r.from_node = sqlc.arg(node_id)::uuid then r.to_node else r.from_node end
left join catalog.entry e on e.id = o.id
left join catalog.grp g on g.id = o.id
left join catalog.medium m on m.medium = e.medium
left join catalog.unit mine on mine.id = case when r.from_node = sqlc.arg(node_id)::uuid then r.from_unit else r.to_unit end
left join catalog.unit theirs on theirs.id = case when r.from_node = sqlc.arg(node_id)::uuid then r.to_unit else r.from_unit end
where sqlc.arg(node_id)::uuid in (r.from_node, r.to_node)
  and not r.spoiler
order by rk.family, reading, e.start_date nulls last, title;

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
