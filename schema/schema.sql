-- listo schema. No migrations until there's data to keep: make db-reset
-- rebuilds from this and seed/. data/constraints_test.go tries one bad write
-- per rule.
--
-- catalog is shared, tracking is per user; catalog never references tracking.

create extension pg_trgm;
create extension unaccent;

create schema catalog;
create schema tracking;

-- A title's search form: lowercase, unaccented, punctuation collapsed.
-- NOTE: pinning the dictionary makes unaccent IMMUTABLE, so it can feed a
-- generated column.
create function catalog.norm(t text) returns text
language sql immutable parallel safe strict
return btrim(regexp_replace(lower(public.unaccent('public.unaccent'::regdictionary, t)), '[[:punct:][:space:]]+', ' ', 'g'));

-- Vocabularies are rows, not CHECK lists, so adding one is an insert.

create table catalog.medium (
  medium         text primary key,
  label          text not null,            -- "visual novel"
  verb_current   text not null,            -- "watching": status current, rendered per medium
  verb_planning  text not null,            -- "plan to watch"
  verb_repeating text not null,            -- "rewatching"
  part_label     text                      -- what unit.part means here: "season", "volume"
);

create table catalog.medium_format (
  medium text not null references catalog.medium,
  format text not null,
  label  text not null,                    -- what the UI prints: anime + movie reads "film"
  primary key (medium, format)
);

create table catalog.relation_kind (
  kind          text primary key check (kind ~ '^[a-z_]+$'),
  label         text not null,             -- reads from -> to: "sequel of"
  inverse_label text,                      -- reads to -> from: "prequel of"; null when symmetric
  is_symmetric  boolean not null,
  family        text not null check (family in ('sequence', 'derivative', 'version', 'setting')),
  unique (kind, is_symmetric),             -- target of relation's (kind, is_symmetric) FK
  check (is_symmetric = (inverse_label is null))
);

-- An id space, not a site: MAL and TMDB number each medium separately.
create table catalog.id_source (
  source       text primary key,           -- mal_anime, tmdb_tv, wikidata
  site         text not null,              -- "MyAnimeList"
  url_template text not null               -- https://myanimelist.net/anime/{id}
);

-- Every catalog thing is a node, so titles, ids, relations, memberships and
-- tags attach to entries and groups alike.

create table catalog.node (
  id         uuid primary key default uuidv7(),
  -- for URLs, /entries/{public_id}/{slug}; the slug is just decoration
  public_id  bigint not null generated always as identity unique,
  type       text not null check (type in ('entry', 'group')),
  slug       text not null unique check (slug ~ '^[a-z0-9]+(-[a-z0-9]+)*$'),
  source     text not null default 'staff',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (id, type)                        -- target of the subtype FKs
);

create table catalog.entry (
  id              uuid primary key,
  type            text not null default 'entry' check (type = 'entry'),
  medium          text not null,
  format          text not null,
  country         char(2),
  title           text not null,           -- display default; real titles live in catalog.title
  start_date      date,
  start_precision text check (start_precision in ('day', 'month', 'year')),
  unit_kind       text not null default 'none'
                  check (unit_kind in ('episode', 'chapter', 'volume', 'route', 'none')),
  unit_total      int,
  foreign key (id, type) references catalog.node (id, type) on delete cascade,
  constraint entry_format_fits_medium foreign key (medium, format)
    references catalog.medium_format (medium, format),
  check ((start_date is null) = (start_precision is null))
);

-- "group" is a reserved word, hence grp.
create table catalog.grp (
  id       uuid primary key,
  type     text not null default 'group' check (type = 'group'),
  kind     text not null check (kind in ('series', 'franchise', 'universe', 'continuity')),
  name     text not null,
  standing text check (standing in ('primary', 'active', 'concluded', 'abandoned', 'disputed')),
  foreign key (id, type) references catalog.node (id, type) on delete cascade,
  constraint grp_standing_only_on_continuity check ((kind = 'continuity') = (standing is not null))
);

create table catalog.title (
  node_id    uuid not null references catalog.node on delete cascade,
  kind       text not null check (kind in ('native', 'romanized', 'english', 'localized', 'synonym')),
  language   text,
  title      text not null,
  title_norm text generated always as (catalog.norm(title)) stored,
  primary key (node_id, kind, title)
);
create index title_norm_trgm on catalog.title using gin (title_norm gin_trgm_ops);

create table catalog.external_id (
  source      text not null references catalog.id_source,
  external_id text not null,
  node_id     uuid not null references catalog.node on delete cascade,
  constraint external_id_dedupe primary key (source, external_id)
);
create index on catalog.external_id (node_id);

create table catalog.unit (
  id       uuid primary key default uuidv7(),
  entry_id uuid not null references catalog.entry on delete cascade,
  ordinal  int not null check (ordinal > 0),   -- dense 1..N, for math
  kind     text not null check (kind in ('episode', 'chapter', 'volume', 'route', 'special')),
  number   numeric,                            -- the official number, for display and matching
  part     int,                                -- season (tv) or volume (manga chapters): medium.part_label
  label    text not null,                      -- "Ep. 7", "23β"
  unique (entry_id, ordinal),
  unique (id, entry_id)                        -- lets anchors prove which entry they belong to
);

-- Graph

create table catalog.relation (
  id             uuid primary key default uuidv7(),
  from_node      uuid not null references catalog.node on delete cascade,
  kind           text not null,
  to_node        uuid not null references catalog.node on delete cascade,
  is_symmetric   boolean not null default false,  -- must match relation_kind (FK below)
  from_unit      uuid,                            -- anchors: where on each side the link attaches
  to_unit        uuid,
  note           text,
  spoiler        boolean not null default false,
  established_on date,
  established_by uuid references catalog.entry,
  source         text not null default 'staff',
  constraint relation_kind_symmetry foreign key (kind, is_symmetric)
    references catalog.relation_kind (kind, is_symmetric),
  constraint relation_from_unit_on_from_node foreign key (from_unit, from_node)
    references catalog.unit (id, entry_id),
  constraint relation_to_unit_on_to_node foreign key (to_unit, to_node)
    references catalog.unit (id, entry_id),
  unique (from_node, kind, to_node),
  constraint relation_not_to_self check (from_node <> to_node),
  constraint relation_symmetric_stored_once check (not is_symmetric or from_node < to_node)
);
create index on catalog.relation (to_node);

create table catalog.membership (
  group_id       uuid not null references catalog.grp on delete cascade,
  member_id      uuid not null references catalog.node on delete cascade,
  role           text,                     -- "source material", "main entry"
  position_label text,                     -- "#1.5", "Vol. 3"
  primary key (group_id, member_id),
  constraint membership_not_in_self check (group_id <> member_id)
);
create index on catalog.membership (member_id);

create table catalog.ordering (
  id        uuid primary key default uuidv7(),
  group_id  uuid not null references catalog.grp on delete cascade,
  kind      text not null check (kind in ('release', 'chronological', 'recommended')),
  name      text not null,
  rationale text,
  origin    text not null default 'staff' check (origin in ('staff', 'community', 'auto')),
  unique (group_id, name)
);

-- More than one branch at a position is a pick-one choice; step orders a
-- branch.
create table catalog.ordering_item (
  ordering_id uuid not null references catalog.ordering on delete cascade,
  position    int not null check (position > 0),
  branch      int not null default 1 check (branch > 0),
  step        int not null default 1 check (step > 0),
  entry_id    uuid not null references catalog.entry,
  unit_from   uuid,
  unit_to     uuid,
  necessity   text not null default 'required'
              check (necessity in ('required', 'recommended', 'optional', 'skip_first_time')),
  section     text,
  note        text,
  primary key (ordering_id, position, branch, step),
  foreign key (unit_from, entry_id) references catalog.unit (id, entry_id),
  foreign key (unit_to, entry_id) references catalog.unit (id, entry_id)
);

create table catalog.tag (
  id   uuid primary key default uuidv7(),
  kind text not null check (kind in ('genre', 'theme', 'warning')),
  name text not null,
  unique (kind, name)
);

create table catalog.node_tag (
  node_id uuid not null references catalog.node on delete cascade,
  tag_id  uuid not null references catalog.tag on delete cascade,
  source  text not null default 'staff',
  primary key (node_id, tag_id)
);

-- Tracking

create table tracking.account (
  id                  uuid primary key default uuidv7(),
  username            text not null unique,
  priority_aging_days numeric not null default 7    -- days per +1 of aging; 0 turns aging off
                      check (priority_aging_days >= 0)
);

create table tracking.list_entry (
  account_id uuid not null references tracking.account on delete cascade,
  entry_id   uuid not null references catalog.entry,  -- no cascade: catalog deletes go through a merge
  status     text not null
             check (status in ('planning', 'current', 'paused', 'dropped', 'completed', 'repeating')),
  progress   numeric not null default 0,
  priority   int not null default 0,                   -- the nice number, +/- by hand
  added_at   timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  primary key (account_id, entry_id)
);

-- Plan-to-watch order: the nice number plus one point per priority_aging_days
-- on the list, so old entries rise and nothing starves.
create view tracking.queue as
select le.account_id, le.entry_id, le.priority, le.added_at,
       le.priority + coalesce(extract(epoch from now() - le.added_at) / 86400 / nullif(a.priority_aging_days, 0), 0)
         as effective
from tracking.list_entry le
join tracking.account a on a.id = le.account_id
where le.status = 'planning';
