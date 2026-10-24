-- Seed helpers, loaded by make db-reset only. Lookups are strict: a mistyped
-- slug is an error, not a new node.

create function catalog.id(p_slug text) returns uuid
language plpgsql stable as $$
declare v uuid;
begin
  select id into strict v from catalog.node where slug = p_slug;
  return v;
exception when no_data_found then
  raise exception 'no node with slug %', p_slug;
end $$;

-- p_ref is an ordinal ('23') or a label ('23β').
create function catalog.unit_id(p_slug text, p_ref text) returns uuid
language plpgsql stable as $$
declare v uuid;
begin
  select id into strict v from catalog.unit
  where entry_id = catalog.id(p_slug)
    and case when p_ref ~ '^[0-9]+$' then ordinal = p_ref::int else label = p_ref end;
  return v;
exception when no_data_found then
  raise exception 'no unit % on %', p_ref, p_slug;
end $$;

-- Seeds pin public ids so URLs don't move as the seed grows.
-- NOTE: public_id is generated always, hence overriding system value.
create function catalog.add_node(p_type text, p_slug text, p_public_id bigint default null)
returns uuid language plpgsql as $$
declare v uuid;
begin
  if p_public_id is null then
    insert into catalog.node (type, slug) values (p_type, p_slug) returning id into v;
  else
    insert into catalog.node (public_id, type, slug) overriding system value
    values (p_public_id, p_type, p_slug) returning id into v;
  end if;
  return v;
end $$;

-- p_start is '2009', '2011-04' or '2011-04-06'; its length sets the precision.
create function catalog.add_entry(p_slug text, p_medium text, p_format text, p_title text,
    p_start text default null, p_country text default null,
    p_unit_kind text default 'none', p_total int default null,
    p_title_kind text default 'english', p_public_id bigint default null)
returns uuid language plpgsql as $$
declare
  v uuid;
  v_precision text := case length(p_start) when 4 then 'year' when 7 then 'month' when 10 then 'day' end;
begin
  if p_start is not null and v_precision is null then
    raise exception 'bad start date %', p_start;
  end if;
  v := catalog.add_node('entry', p_slug, p_public_id);
  insert into catalog.entry (id, medium, format, country, title, start_date, start_precision, unit_kind, unit_total)
  values (v, p_medium, p_format, p_country, p_title,
          case v_precision
            when 'year' then (p_start || '-01-01')::date
            when 'month' then (p_start || '-01')::date
            when 'day' then p_start::date
          end,
          v_precision, p_unit_kind, p_total);
  insert into catalog.title (node_id, kind, title) values (v, p_title_kind, p_title);
  return v;
end $$;

create function catalog.add_units(p_slug text, p_count int) returns void
language sql as $$
  insert into catalog.unit (entry_id, ordinal, kind, number, label)
  select e.id, n, e.unit_kind, n,
         case e.unit_kind when 'episode' then 'Ep. ' when 'chapter' then 'Ch. ' when 'volume' then 'Vol. ' else '' end || n
  from catalog.entry e, generate_series(1, p_count) n
  where e.id = catalog.id(p_slug)
$$;

create function catalog.add_special(p_slug text, p_label text) returns void
language sql as $$
  insert into catalog.unit (entry_id, ordinal, kind, label)
  select e.id, coalesce(max(u.ordinal), 0) + 1, 'special', p_label
  from catalog.entry e
  left join catalog.unit u on u.entry_id = e.id
  where e.id = catalog.id(p_slug)
  group by e.id
$$;

create function catalog.add_group(p_slug text, p_kind text, p_name text, p_standing text default null,
    p_public_id bigint default null)
returns uuid language plpgsql as $$
declare v uuid;
begin
  v := catalog.add_node('group', p_slug, p_public_id);
  insert into catalog.grp (id, kind, name, standing) values (v, p_kind, p_name, p_standing);
  insert into catalog.title (node_id, kind, title) values (v, 'english', p_name);
  return v;
end $$;

create function catalog.alias(p_slug text, p_kind text, p_title text) returns void
language sql as $$
  insert into catalog.title (node_id, kind, title) values (catalog.id(p_slug), p_kind, p_title)
$$;

create function catalog.ext(p_slug text, p_source text, p_id text) returns void
language sql as $$
  insert into catalog.external_id (source, external_id, node_id) values (p_source, p_id, catalog.id(p_slug))
$$;

create function catalog.member(p_group text, variadic p_members text[]) returns void
language sql as $$
  insert into catalog.membership (group_id, member_id)
  select catalog.id(p_group), catalog.id(m) from unnest(p_members) m
$$;

-- a member's role in its group, eg "source material"
create function catalog.set_role(p_group text, p_member text, p_role text) returns void
language plpgsql as $$
begin
  update catalog.membership set role = p_role
  where group_id = catalog.id(p_group) and member_id = catalog.id(p_member);
  if not found then
    raise exception '% is not in %', p_member, p_group;
  end if;
end $$;

-- Symmetric kinds are stored once, smaller id first, so callers never pick a
-- direction for them.
create function catalog.relate(p_from text, p_kind text, p_to text,
    p_from_unit text default null, p_to_unit text default null,
    p_note text default null, p_spoiler boolean default false,
    p_established_on date default null, p_established_by text default null)
returns void language plpgsql as $$
declare
  f   uuid := catalog.id(p_from);
  t   uuid := catalog.id(p_to);
  fu  uuid := case when p_from_unit is not null then catalog.unit_id(p_from, p_from_unit) end;
  tu  uuid := case when p_to_unit is not null then catalog.unit_id(p_to, p_to_unit) end;
  sym boolean;
begin
  select is_symmetric into strict sym from catalog.relation_kind where kind = p_kind;
  if sym and f > t then
    select t, f, tu, fu into f, t, fu, tu;
  end if;
  insert into catalog.relation (from_node, kind, to_node, is_symmetric, from_unit, to_unit,
                                note, spoiler, established_on, established_by)
  values (f, p_kind, t, sym, fu, tu, p_note, p_spoiler, p_established_on,
          case when p_established_by is not null then catalog.id(p_established_by) end);
end $$;

create function catalog.add_ordering(p_group text, p_kind text, p_name text, p_rationale text default null)
returns uuid language sql as $$
  insert into catalog.ordering (group_id, kind, name, rationale)
  values (catalog.id(p_group), p_kind, p_name, p_rationale)
  returning id
$$;

create function catalog.add_step(p_ordering uuid, p_position int, p_entry text,
    p_from text default null, p_to text default null,
    p_necessity text default 'required', p_note text default null,
    p_branch int default 1, p_step int default 1, p_section text default null)
returns void language sql as $$
  insert into catalog.ordering_item (ordering_id, position, branch, step, entry_id, unit_from, unit_to, necessity, section, note)
  values (p_ordering, p_position, p_branch, p_step, catalog.id(p_entry),
          case when p_from is not null then catalog.unit_id(p_entry, p_from) end,
          case when p_to is not null then catalog.unit_id(p_entry, p_to) end,
          p_necessity, p_section, p_note)
$$;

create function catalog.tag_with(p_kind text, p_name text, variadic p_slugs text[]) returns void
language sql as $$
  insert into catalog.tag (kind, name) values (p_kind, p_name) on conflict do nothing;
  insert into catalog.node_tag (node_id, tag_id)
  select catalog.id(s), t.id from unnest(p_slugs) s, catalog.tag t where t.kind = p_kind and t.name = p_name;
$$;

-- credits a company with an entry, adding the company the first time
create function catalog.credit(p_entry text, p_company text, p_role text, p_country char(2) default null)
returns void language plpgsql as $$
declare v uuid;
begin
  select id into v from catalog.company where name = p_company;
  if v is null then
    insert into catalog.company (name, country) values (p_company, p_country) returning id into v;
  end if;
  insert into catalog.entry_company (entry_id, company_id, role) values (catalog.id(p_entry), v, p_role);
end $$;
