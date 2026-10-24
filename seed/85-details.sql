-- What's known beyond titles and dates: status, how long things run, what
-- they came from, who made them, and a few descriptions. Only what's sure;
-- the rest stays null.

-- finished series, with the month they ended where it's known
update catalog.entry e set status = 'finished', end_date = d.ended::date, end_precision = 'month'
from (values
  ('steins-gate-anime', '2011-09-01'), ('steins-gate-0-anime', '2018-09-01'), ('madoka-tv', '2011-04-01'),
  ('index-s1', '2009-03-01'), ('index-ii', '2011-04-01'), ('railgun-s1', '2010-03-01'), ('railgun-s', '2013-09-01'),
  ('naruto', '2007-02-01'), ('naruto-shippuden', '2017-03-01'),
  ('konosuba-s1', '2016-03-01'), ('konosuba-s2', '2017-03-01'), ('konosuba-s3', '2024-06-01'),
  ('shameless-uk', '2013-05-01'), ('shameless-us', '2021-04-01')) as d (slug, ended)
where e.id = catalog.id(d.slug);

update catalog.entry set status = 'finished'
where id in (select catalog.id(s) from unnest(array['bakemonogatari', 'nisemonogatari', 'monogatari-second-season',
  'ben-10-2005', 'ben-10-alien-force', 'ben-10-ultimate-alien', 'ben-10-omniverse']) s);

-- still going
update catalog.entry set status = 'releasing'
where id in (catalog.id('detective-conan-anime'), catalog.id('detective-conan-manga'));

update catalog.entry set unit_minutes = 24
where id in (select catalog.id(s) from unnest(array['steins-gate-anime', 'steins-gate-0-anime', 'madoka-tv',
  'konosuba-s1', 'konosuba-s2', 'konosuba-s3', 'index-s1', 'index-ii', 'railgun-s1', 'railgun-s']) s);

-- films, by runtime
update catalog.entry e set length_minutes = d.minutes
from (values
  ('steins-gate-movie', 90), ('madoka-rebellion', 116),
  ('kizumonogatari-1', 64), ('kizumonogatari-2', 69), ('kizumonogatari-3', 83),
  ('spider-man-2002', 121), ('spider-man-2', 127), ('spider-man-3', 139),
  ('amazing-spider-man', 136), ('amazing-spider-man-2', 142),
  ('spider-man-homecoming', 133), ('spider-man-far-from-home', 129), ('spider-man-no-way-home', 148)) as d (slug, minutes)
where e.id = catalog.id(d.slug);

-- what they came from; an adaptation_of relation, where there is one, says it better
update catalog.entry e set source_type = d.source
from (values
  ('madoka-tv', 'original'), ('ben-10-2005', 'original'),
  ('bakemonogatari', 'light_novel'), ('nisemonogatari', 'light_novel'), ('nekomonogatari-black', 'light_novel'),
  ('monogatari-second-season', 'light_novel'), ('kizumonogatari-1', 'light_novel'),
  ('konosuba-s1', 'light_novel'), ('index-s1', 'light_novel'), ('railgun-s1', 'manga'),
  ('naruto', 'manga'), ('detective-conan-anime', 'manga'),
  ('steins-gate-anime', 'visual_novel'), ('steins-gate-0-anime', 'visual_novel')) as d (slug, source)
where e.id = catalog.id(d.slug);

-- studios
select catalog.credit(s, 'White Fox', 'studio', 'JP')
from unnest(array['steins-gate-anime', 'steins-gate-0-anime', 'steins-gate-movie']) s;
select catalog.credit(s, 'Shaft', 'studio', 'JP')
from unnest(array['madoka-tv', 'madoka-beginnings', 'madoka-eternal', 'madoka-rebellion',
  'bakemonogatari', 'nisemonogatari', 'nekomonogatari-black', 'monogatari-second-season', 'hanamonogatari',
  'tsukimonogatari', 'owarimonogatari', 'koyomimonogatari', 'kizumonogatari-1', 'kizumonogatari-2',
  'kizumonogatari-3', 'owarimonogatari-2', 'zoku-owarimonogatari', 'monogatari-off-and-monster']) s;
select catalog.credit(s, 'J.C.Staff', 'studio', 'JP') from unnest(array['index-s1', 'index-ii', 'railgun-s1', 'railgun-s']) s;
select catalog.credit(s, 'Pierrot', 'studio', 'JP') from unnest(array['naruto', 'naruto-shippuden']) s;
select catalog.credit(s, 'Studio Deen', 'studio', 'JP') from unnest(array['konosuba-s1', 'konosuba-s2']) s;
select catalog.credit('konosuba-s3', 'Drive', 'studio', 'JP');
-- the TV series and every film
select catalog.credit(slug, 'TMS Entertainment', 'studio', 'JP')
from catalog.node n join catalog.entry e on e.id = n.id
where n.slug = 'detective-conan-anime' or (n.slug like 'detective-conan-%' and e.format = 'movie');

-- games
select catalog.credit('devil-may-cry', 'Capcom', 'developer', 'JP');
select catalog.credit('dmc-devil-may-cry', 'Ninja Theory', 'developer', 'GB');
select catalog.credit('dmc-devil-may-cry', 'Capcom', 'publisher', 'JP');

-- descriptions, in our own words
update catalog.entry e set description = d.description
from (values
  ('steins-gate-anime', 'A self-styled mad scientist and his friends rig a microwave that sends text messages into the past. Every message rewrites the present a little more, until one change costs far more than they meant to spend.'),
  ('madoka-tv', 'A schoolgirl is offered any wish she likes in exchange for becoming a magical girl who fights witches. The longer she weighs the offer, the clearer it gets that the bargain isn''t what it looks like.'),
  ('bakemonogatari', 'Koyomi Araragi, still half a vampire after one very bad spring break, keeps running into classmates caught up in strange afflictions, and keeps trying to help.'),
  ('detective-conan-anime', 'A teenage detective is poisoned and shrunk into a child''s body. Hiding who he is from almost everyone, he keeps solving cases while hunting the organization that did it.'),
  ('konosuba-s1', 'A shut-in dies, is reborn into a fantasy world with a goddess as the one thing he brought along, and ends up leading the least capable adventuring party in town.')) as d (slug, description)
where e.id = catalog.id(d.slug);
