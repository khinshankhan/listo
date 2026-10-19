-- name: NodeByPublicID :one
-- A node's type, to redirect an id used under the wrong path.
select public_id, slug, type from catalog.node where public_id = $1;

-- name: GroupByPublicID :one
select n.id, n.public_id, n.slug, g.kind, g.name, g.standing
from catalog.node n
join catalog.grp g on g.id = n.id
where n.public_id = $1;

-- name: GroupMembers :many
-- A group's direct members: subgroups first, then entries in release order.
select n.public_id, n.slug, n.type,
       coalesce(e.title, g.name)::text as title,
       g.kind as group_kind, m.label as medium_label, e.start_date,
       ms.role, ms.position_label
from catalog.membership ms
join catalog.node n on n.id = ms.member_id
left join catalog.entry e on e.id = n.id
left join catalog.grp g on g.id = n.id
left join catalog.medium m on m.medium = e.medium
where ms.group_id = $1
order by n.type desc, e.start_date nulls last, title;

-- name: GroupOrderings :many
-- A group's orderings, one row per step.
select o.id as ordering_id, o.kind, o.name, o.rationale,
       oi.position, oi.branch, oi.step,
       n.public_id, n.slug, e.title,
       uf.label as unit_from, ut.label as unit_to,
       oi.necessity, oi.note
from catalog.ordering o
join catalog.ordering_item oi on oi.ordering_id = o.id
join catalog.entry e on e.id = oi.entry_id
join catalog.node n on n.id = e.id
left join catalog.unit uf on uf.id = oi.unit_from
left join catalog.unit ut on ut.id = oi.unit_to
where o.group_id = $1
order by o.name, oi.position, oi.branch, oi.step;
