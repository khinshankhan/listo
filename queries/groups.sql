-- name: NodeByPublicID :one
-- A node's type, to redirect an id used under the wrong path.
select public_id, slug, type from catalog.node where public_id = $1;

-- name: GroupByPublicID :one
select n.id, n.public_id, n.slug, g.kind, g.name, g.standing
from catalog.node n
join catalog.grp g on g.id = n.id
where n.public_id = $1;

-- name: GroupOrderings :many
-- A group's orderings, recommended first, one row per step.
select o.id as ordering_id, o.kind, o.name, o.rationale, o.origin,
       oi.position, oi.branch, oi.step, oi.section,
       n.public_id, n.slug, e.title, e.medium, e.format, mf.label as format_label,
       e.start_date, e.start_precision, e.unit_kind, e.unit_total,
       uf.label as unit_from, ut.label as unit_to, uf.ordinal as from_ordinal, ut.ordinal as to_ordinal,
       oi.necessity, oi.note
from catalog.ordering o
join catalog.ordering_item oi on oi.ordering_id = o.id
join catalog.entry e on e.id = oi.entry_id
join catalog.node n on n.id = e.id
join catalog.medium_format mf on mf.medium = e.medium and mf.format = e.format
left join catalog.unit uf on uf.id = oi.unit_from
left join catalog.unit ut on ut.id = oi.unit_to
where o.group_id = $1
order by case o.kind when 'recommended' then 0 when 'release' then 1 else 2 end,
         o.name, oi.position, oi.branch, oi.step;

-- name: GroupEntries :many
-- Every entry under a group at any depth, with the branch it sits under (the
-- group's direct member it was reached through), and its role and position in
-- its own group. An entry under two branches comes back once for each.
with recursive down (branch, group_id, member_id) as (
  select ms.member_id, ms.group_id, ms.member_id
  from catalog.membership ms
  where ms.group_id = $1
  union all
  select down.branch, ms.group_id, ms.member_id
  from down
  join catalog.membership ms on ms.group_id = down.member_id
) cycle member_id set looped using path
select n.public_id, n.slug, e.title, e.medium, m.label as medium_label, e.format,
       mf.label as format_label, e.start_date, e.start_precision, e.unit_kind, e.unit_total,
       ms.role, ms.position_label,
       b.public_id as branch_public_id, b.slug as branch_slug, b.type as branch_type,
       bg.name as branch_name, bg.kind as branch_kind, bg.standing as branch_standing
from down
join catalog.entry e on e.id = down.member_id
join catalog.node n on n.id = e.id
join catalog.medium m on m.medium = e.medium
join catalog.medium_format mf on mf.medium = e.medium and mf.format = e.format
join catalog.membership ms on ms.group_id = down.group_id and ms.member_id = down.member_id
join catalog.node b on b.id = down.branch
left join catalog.grp bg on bg.id = b.id
where not down.looped
order by e.start_date nulls last, e.title;

-- name: GroupSiblings :many
-- The other groups that share a parent with this one, with every entry
-- beneath each counted.
with recursive down (root, id) as (
  select s.member_id, s.member_id
  from catalog.membership p
  join catalog.membership s on s.group_id = p.group_id
  where p.member_id = $1 and s.member_id <> $1
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
