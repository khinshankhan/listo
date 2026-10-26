-- name: EntryIDBySlug :one
-- An entry's id by slug, for tools that name entries the way seeds do.
select e.id
from catalog.node n
join catalog.entry e on e.id = n.id
where n.slug = $1;

-- name: AddEntryImage :execrows
-- Records an image of an entry, not canon yet. The same image twice is a
-- no-op, zero rows.
insert into catalog.entry_image (entry_id, kind, storage_key, width, height, fit, source, source_url)
values ($1, $2, $3, $4, $5, $6, $7, $8)
on conflict (entry_id, storage_key) do nothing;

-- name: CanonImage :one
-- An entry's canon image of a kind, to weigh a new one against.
select storage_key, width, height
from catalog.entry_image
where entry_id = $1 and kind = $2 and is_canon;

-- name: ClearCanonImage :exec
-- Unseats an entry's canon image of a kind, before another takes its place.
update catalog.entry_image
set is_canon = false, updated_at = now()
where entry_id = $1 and kind = $2 and is_canon;

-- name: SetCanonImage :exec
-- Makes one image its entry's canon; clear the old one first.
update catalog.entry_image
set is_canon = true, updated_at = now()
where entry_id = $1 and storage_key = $2;
