-- name: SearchNodes :many
-- Fuzzy title search over entries and groups, one hit per node by its best
-- matching title, an entry with its cover. Ties go a group before its
-- entries, then oldest first.
select n.public_id, n.slug, n.type,
       coalesce(e.title, g.name)::text as title,
       hit.matched::text as matched, hit.matched_kind::text as matched_kind,
       g.kind as group_kind, e.medium, m.label as medium_label, e.format, mf.label as format_label,
       e.start_date, e.start_precision, e.unit_kind, e.unit_total,
       cover.storage_key as cover_key, cover.width as cover_width, cover.height as cover_height, cover.fit as cover_fit,
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
left join catalog.entry_image cover on cover.id = (
  select i.id
  from catalog.entry_image i
  where i.entry_id = e.id and i.kind = 'cover'
  order by i.is_canon desc, i.width desc, i.id
  limit 1
)
-- a medium narrows before the limit, so its hits aren't crowded out
where sqlc.narg(medium)::text is null or e.medium = sqlc.narg(medium)::text
order by hit.score desc, title, e.start_date nulls first, n.public_id
limit sqlc.arg(max_results)::int;

-- name: SearchMedia :many
-- How many nodes a search finds in each medium, however many that is; groups
-- have no medium.
select e.medium, m.label as medium_label, count(*)::int as hits
from (
  select distinct t.node_id
  from catalog.title t
  where t.title_norm % catalog.norm(sqlc.arg(q)::text)
     or catalog.norm(sqlc.arg(q)::text) <% t.title_norm
) hit
left join catalog.entry e on e.id = hit.node_id
left join catalog.medium m on m.medium = e.medium
group by e.medium, m.label
order by hits desc, m.label nulls last;

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
-- A page of the entries the filters let through, newest first, each row
-- carrying its cover and how many match in all. A format is only ever asked
-- for with its medium. A season is 0 for winter to 3 for fall, and only
-- matches entries dated to the month.
select n.public_id, n.slug, e.title, e.medium, m.label as medium_label, e.format, mf.label as format_label,
       e.start_date, e.start_precision, e.unit_kind, e.unit_total,
       cover.storage_key as cover_key, cover.width as cover_width, cover.height as cover_height, cover.fit as cover_fit,
       count(*) over ()::int as total
from catalog.entry e
join catalog.node n on n.id = e.id
join catalog.medium m on m.medium = e.medium
join catalog.medium_format mf on mf.medium = e.medium and mf.format = e.format
left join catalog.entry_image cover on cover.id = (
  select i.id
  from catalog.entry_image i
  where i.entry_id = e.id and i.kind = 'cover'
  order by i.is_canon desc, i.width desc, i.id
  limit 1
)
where (sqlc.narg(medium)::text is null or e.medium = sqlc.narg(medium)::text)
  and (sqlc.narg(format)::text is null or e.format = sqlc.narg(format)::text)
  and (sqlc.narg(genre)::text is null or exists (
        select 1
        from catalog.node_tag nt
        join catalog.tag t on t.id = nt.tag_id
        where nt.node_id = e.id and t.kind = 'genre' and t.name = sqlc.narg(genre)::text))
  and (sqlc.narg(year)::int is null or extract(year from e.start_date)::int = sqlc.narg(year)::int)
  and (sqlc.narg(season)::int is null
       or (e.start_precision <> 'year' and (extract(month from e.start_date)::int - 1) / 3 = sqlc.narg(season)::int))
order by e.start_date desc nulls last, e.title, n.public_id
limit sqlc.arg(page_size)::int offset sqlc.arg(skip)::int;

-- name: BrowseShelves :many
-- Each medium's newest entries, at most per_shelf apiece, with a cover each
-- when there is one.
select shelf.public_id, shelf.slug, shelf.title, shelf.medium, shelf.medium_label, shelf.format, shelf.format_label,
       shelf.start_date, shelf.start_precision, shelf.unit_kind, shelf.unit_total,
       cover.storage_key as cover_key, cover.width as cover_width, cover.height as cover_height, cover.fit as cover_fit
from (
  select e.id, n.public_id, n.slug, e.title, e.medium, m.label as medium_label, e.format, mf.label as format_label,
         e.start_date, e.start_precision, e.unit_kind, e.unit_total,
         row_number() over (partition by e.medium order by e.start_date desc nulls last, e.title, n.public_id) as rank
  from catalog.entry e
  join catalog.node n on n.id = e.id
  join catalog.medium m on m.medium = e.medium
  join catalog.medium_format mf on mf.medium = e.medium and mf.format = e.format
) shelf
left join catalog.entry_image cover on cover.id = (
  select i.id
  from catalog.entry_image i
  where i.entry_id = shelf.id and i.kind = 'cover'
  order by i.is_canon desc, i.width desc, i.id
  limit 1
)
where shelf.rank <= sqlc.arg(per_shelf)::int
order by shelf.medium, shelf.rank;

-- name: BrowseGenres :many
-- Genres some entry has, for the filter.
select t.name
from catalog.tag t
where t.kind = 'genre'
  and exists (select 1 from catalog.node_tag nt join catalog.entry e on e.id = nt.node_id where nt.tag_id = t.id)
order by t.name;

-- name: BrowseFormats :many
-- Every medium's formats some entry has, for the filter, so it can follow
-- the type without asking again.
select mf.medium, mf.format, mf.label
from catalog.medium_format mf
where exists (select 1 from catalog.entry e where e.medium = mf.medium and e.format = mf.format)
order by mf.medium, mf.label;

-- name: BrowseYears :many
-- Years some entry started in, newest first, for the filter.
select distinct extract(year from start_date)::int as year
from catalog.entry
where start_date is not null
order by year desc;
