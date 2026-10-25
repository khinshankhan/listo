-- What's known beyond titles and dates: status, how long things run, what
-- they came from, who made them, and a few descriptions. Only what's sure;
-- the rest stays null.
-- https://en.wikipedia.org/wiki/List_of_Monogatari_episodes
-- https://en.wikipedia.org/wiki/List_of_Case_Closed_films
-- https://toaru-project.com/item/onair/

-- finished series, with the month they ended where it's known
update catalog.entry e set status = 'finished', end_date = d.ended::date, end_precision = 'month'
from (values
  ('steins-gate-anime', '2011-09-01'), ('steins-gate-0-anime', '2018-09-01'), ('madoka-tv', '2011-04-01'),
  ('index-s1', '2009-03-01'), ('index-ii', '2011-04-01'), ('railgun-s1', '2010-03-01'), ('railgun-s', '2013-09-01'),
  ('naruto', '2007-02-01'), ('naruto-shippuden', '2017-03-01'),
  ('konosuba-s1', '2016-03-01'), ('konosuba-s2', '2017-03-01'), ('konosuba-s3', '2024-06-01'),
  ('shameless-uk', '2013-05-01'), ('shameless-us', '2021-04-01'), ('bizim-hikaye', '2019-05-01'),
  ('robotics-notes-anime', '2013-03-01'), ('chaos-head-anime', '2008-12-01'),
  ('chaos-child-anime', '2017-03-01'), ('occultic-nine-anime', '2016-12-01'),
  ('index-iii', '2019-04-01'), ('railgun-t', '2020-09-01'), ('accelerator', '2019-09-01'),
  ('bakemonogatari', '2010-06-01'), ('nisemonogatari', '2012-03-01'),
  ('nekomonogatari-black', '2012-12-01'), ('monogatari-second-season', '2013-12-01'),
  ('hanamonogatari', '2014-08-01'), ('tsukimonogatari', '2014-12-01'),
  ('owarimonogatari', '2015-12-01'), ('koyomimonogatari', '2016-03-01'),
  ('owarimonogatari-2', '2017-08-01'), ('zoku-owarimonogatari', '2019-06-01'),
  ('monogatari-off-and-monster', '2024-10-01'),
  ('ben-10-2005', '2008-04-01'), ('ben-10-alien-force', '2010-03-01'),
  ('ben-10-ultimate-alien', '2012-03-01'), ('ben-10-omniverse', '2014-11-01'),
  ('ben-10-2016', '2021-04-01'), ('detective-conan-zeros-tea-time', '2022-05-01'),
  ('detective-conan-wild-police-story', '2020-11-01')) as d (slug, ended)
where e.id = catalog.id(d.slug);

-- one-off games, films and OVAs that are already out
update catalog.entry e set status = 'finished'
from catalog.node n
where n.id = e.id
  and (e.medium in ('game', 'visual_novel', 'film') or (e.medium = 'anime' and e.format in ('movie', 'ova')))
  and e.status is null and e.start_date <= date '2026-10-24'
  and n.slug not in ('spider-man-homecoming', 'spider-man-far-from-home', 'spider-man-no-way-home');

-- still going
update catalog.entry set status = 'releasing'
where id in (catalog.id('detective-conan-anime'), catalog.id('detective-conan-manga'),
            catalog.id('detective-conan-the-culprit-hanzawa'), catalog.id('item'));

update catalog.entry set status = 'announced', unit_kind = 'episode' where id = catalog.id('railgun-4');

-- Three weekly episodes released by October 24; no confirmed final count.
-- https://www.crunchyroll.com/news/latest/2026/10/17/a-certain-dark-item-anime-creditless-opening-ending
-- https://www.primevideo.com/-/es/detail/0S5XNB4ABKPRZ36WKL5INQX5WO
update catalog.entry set unit_available = 3, unit_minutes = 24 where id = catalog.id('item');
select catalog.add_units('item', 3);

-- The light novels have three published volumes, with no verified final total.
-- https://en.wikipedia.org/wiki/Occultic;Nine
update catalog.entry set unit_kind = 'volume', unit_available = 3 where id = catalog.id('occultic-nine-ln');
update catalog.entry e set unit_total = d.total
from (values ('detective-conan-zeros-tea-time', 60), ('detective-conan-wild-police-story', 13)) as d (slug, total)
where e.id = catalog.id(d.slug);

-- Dates of the first Japanese editions, rather than their later localizations.
update catalog.entry e set start_date = d.released::date, start_precision = 'day'
from (values
  ('steins-gate-vn', '2009-10-15'), ('steins-gate-elite', '2018-09-20'),
  ('steins-gate-0-vn', '2015-12-10'), ('robotics-notes-vn', '2012-06-28'),
  ('robotics-notes-dash', '2019-01-31'), ('chaos-head-vn', '2008-04-25'),
  ('chaos-child-vn', '2014-12-18'), ('occultic-nine-ln', '2014-08-25'),
  ('occultic-nine-vn', '2017-11-09'), ('anonymous-code-vn', '2022-07-28'),
  ('detective-conan-manga', '1994-01-05'), ('devil-may-cry', '2001-08-23'),
  ('railgun-ova-model', '2010-07-24'), ('railgun-s-ova', '2014-03-27')) as d (slug, released)
where e.id = catalog.id(d.slug);

update catalog.entry set unit_kind = 'episode', unit_total = 1
where id in (catalog.id('railgun-ova-model'), catalog.id('railgun-ova-misaka'), catalog.id('railgun-s-ova'));

update catalog.entry set unit_minutes = 24
where id in (select catalog.id(s) from unnest(array['steins-gate-anime', 'steins-gate-0-anime', 'madoka-tv',
  'konosuba-s1', 'konosuba-s2', 'konosuba-s3', 'index-s1', 'index-ii', 'railgun-s1', 'railgun-s',
  'index-iii', 'railgun-t', 'accelerator', 'robotics-notes-anime', 'chaos-head-anime',
  'chaos-child-anime', 'occultic-nine-anime', 'bakemonogatari', 'nisemonogatari',
  'nekomonogatari-black', 'monogatari-second-season', 'hanamonogatari', 'tsukimonogatari',
  'owarimonogatari', 'owarimonogatari-2', 'zoku-owarimonogatari']) s);

update catalog.entry e set unit_minutes = d.minutes
from (values ('naruto', 23), ('naruto-shippuden', 23), ('detective-conan-anime', 25),
  ('ben-10-2005', 22), ('ben-10-alien-force', 23), ('ben-10-ultimate-alien', 22),
  ('ben-10-2016', 11), ('koyomimonogatari', 14), ('bizim-hikaye', 120)) as d (slug, minutes)
where e.id = catalog.id(d.slug);

-- films, by runtime
update catalog.entry e set length_minutes = d.minutes
from (values
  ('steins-gate-movie', 90), ('madoka-beginnings', 130), ('madoka-eternal', 109), ('madoka-rebellion', 116),
  ('index-endymion', 90), ('ben-10-race-against-time', 67), ('konosuba-legend-of-crimson', 90),
  ('kizumonogatari-1', 64), ('kizumonogatari-2', 69), ('kizumonogatari-3', 83),
  ('spider-man-2002', 121), ('spider-man-2', 127), ('spider-man-3', 139),
  ('amazing-spider-man', 136), ('amazing-spider-man-2', 142),
  ('spider-man-homecoming', 133), ('spider-man-far-from-home', 129), ('spider-man-no-way-home', 148)) as d (slug, minutes)
where e.id = catalog.id(d.slug);

-- Theatrical runtimes of the numbered Conan films, plus the Lupin crossover.
update catalog.entry e set length_minutes = d.minutes
from (values
  ('detective-conan-the-time-bombed-skyscraper', 94), ('detective-conan-the-fourteenth-target', 99),
  ('detective-conan-the-last-wizard-of-the-century', 99), ('detective-conan-captured-in-her-eyes', 100),
  ('detective-conan-countdown-to-heaven', 100), ('detective-conan-the-phantom-of-baker-street', 107),
  ('detective-conan-crossroad-in-the-ancient-capital', 108), ('detective-conan-magician-of-the-silver-sky', 108),
  ('detective-conan-strategy-above-the-depths', 109), ('detective-conan-the-private-eyes-requiem', 111),
  ('detective-conan-jolly-roger-in-the-deep-azure', 107), ('detective-conan-full-score-of-fear', 115),
  ('detective-conan-the-raven-chaser', 111), ('detective-conan-the-lost-ship-in-the-sky', 102),
  ('detective-conan-quarter-of-silence', 109), ('detective-conan-the-eleventh-striker', 110),
  ('detective-conan-private-eye-in-the-distant-sea', 110), ('detective-conan-dimensional-sniper', 110),
  ('detective-conan-sunflowers-of-inferno', 112), ('detective-conan-the-darkest-nightmare', 112),
  ('detective-conan-crimson-love-letter', 112), ('detective-conan-zero-the-enforcer', 110),
  ('detective-conan-the-fist-of-blue-sapphire', 109), ('detective-conan-the-scarlet-bullet', 112),
  ('detective-conan-the-bride-of-halloween', 111), ('detective-conan-black-iron-submarine', 110),
  ('detective-conan-the-million-dollar-pentagram', 111), ('detective-conan-one-eyed-flashback', 110),
  ('detective-conan-fallen-angel-of-the-highway', 109), ('lupin-the-3rd-vs-detective-conan-the-movie', 107)) as d (slug, minutes)
where e.id = catalog.id(d.slug);

-- what they came from; an adaptation_of relation, where there is one, says it better
update catalog.entry e set source_type = d.source
from (values
  ('madoka-tv', 'original'), ('ben-10-2005', 'original'),
  ('bakemonogatari', 'light_novel'), ('nisemonogatari', 'light_novel'), ('nekomonogatari-black', 'light_novel'),
  ('monogatari-second-season', 'light_novel'), ('kizumonogatari-1', 'light_novel'),
  ('konosuba-s1', 'light_novel'), ('index-s1', 'light_novel'), ('railgun-s1', 'manga'),
  ('naruto', 'manga'), ('detective-conan-anime', 'manga'),
  ('steins-gate-anime', 'visual_novel'), ('steins-gate-0-anime', 'visual_novel'),
  ('steins-gate-vn', 'original'), ('steins-gate-0-vn', 'original'), ('steins-gate-elite', 'visual_novel'),
  ('steins-gate-movie', 'original'), ('robotics-notes-vn', 'original'), ('robotics-notes-dash', 'original'),
  ('robotics-notes-anime', 'visual_novel'), ('chaos-head-vn', 'original'), ('chaos-child-vn', 'original'),
  ('chaos-head-anime', 'visual_novel'), ('chaos-child-anime', 'visual_novel'),
  ('occultic-nine-ln', 'original'), ('occultic-nine-anime', 'light_novel'), ('occultic-nine-vn', 'light_novel'),
  ('anonymous-code-vn', 'original'), ('madoka-beginnings', 'original'), ('madoka-eternal', 'original'),
  ('madoka-rebellion', 'original'), ('index-ii', 'light_novel'), ('index-iii', 'light_novel'),
  ('index-endymion', 'original'), ('railgun-s', 'manga'), ('railgun-t', 'manga'), ('railgun-4', 'manga'),
  ('accelerator', 'manga'), ('item', 'light_novel'), ('naruto-shippuden', 'manga'),
  ('konosuba-s2', 'light_novel'), ('konosuba-s3', 'light_novel'), ('konosuba-legend-of-crimson', 'light_novel'),
  ('kizumonogatari-2', 'light_novel'), ('kizumonogatari-3', 'light_novel'), ('hanamonogatari', 'light_novel'),
  ('tsukimonogatari', 'light_novel'), ('owarimonogatari', 'light_novel'), ('koyomimonogatari', 'light_novel'),
  ('owarimonogatari-2', 'light_novel'), ('zoku-owarimonogatari', 'light_novel'),
  ('monogatari-off-and-monster', 'light_novel'), ('ben-10-alien-force', 'original'),
  ('ben-10-ultimate-alien', 'original'), ('ben-10-omniverse', 'original'), ('ben-10-2016', 'original'),
  ('ben-10-race-against-time', 'original'), ('shameless-uk', 'original'),
  ('devil-may-cry', 'original'), ('dmc-devil-may-cry', 'original'),
  ('spider-man-2002', 'comic'), ('spider-man-2', 'comic'), ('spider-man-3', 'comic'),
  ('amazing-spider-man', 'comic'), ('amazing-spider-man-2', 'comic')) as d (slug, source)
where e.id = catalog.id(d.slug);

-- The films use the manga's characters but tell original screen stories.
update catalog.entry e set source_type = 'manga'
from catalog.membership m
where m.member_id = e.id and m.group_id = catalog.id('detective-conan-films');

-- studios
select catalog.credit(s, 'White Fox', 'studio', 'JP')
from unnest(array['steins-gate-anime', 'steins-gate-0-anime', 'steins-gate-movie']) s;
select catalog.credit(s, 'Shaft', 'studio', 'JP')
from unnest(array['madoka-tv', 'madoka-beginnings', 'madoka-eternal', 'madoka-rebellion',
  'bakemonogatari', 'nisemonogatari', 'nekomonogatari-black', 'monogatari-second-season', 'hanamonogatari',
  'tsukimonogatari', 'owarimonogatari', 'koyomimonogatari', 'kizumonogatari-1', 'kizumonogatari-2',
  'kizumonogatari-3', 'owarimonogatari-2', 'zoku-owarimonogatari', 'monogatari-off-and-monster']) s;
select catalog.credit(s, 'J.C.Staff', 'studio', 'JP')
from unnest(array['index-s1', 'index-ii', 'index-iii', 'index-endymion', 'railgun-s1', 'railgun-s',
  'railgun-t', 'railgun-ova-model', 'railgun-ova-misaka', 'railgun-s-ova', 'accelerator', 'item',
  'konosuba-legend-of-crimson']) s;
select catalog.credit(s, 'Pierrot', 'studio', 'JP') from unnest(array['naruto', 'naruto-shippuden']) s;
select catalog.credit(s, 'Studio Deen', 'studio', 'JP') from unnest(array['konosuba-s1', 'konosuba-s2']) s;
select catalog.credit('konosuba-s3', 'Drive', 'studio', 'JP');
select catalog.credit('robotics-notes-anime', 'Production I.G', 'studio', 'JP');
select catalog.credit('chaos-head-anime', 'Madhouse', 'studio', 'JP');
select catalog.credit('chaos-child-anime', 'Silver Link', 'studio', 'JP');
select catalog.credit('occultic-nine-anime', 'A-1 Pictures', 'studio', 'JP');
select catalog.credit(s, 'Cartoon Network Studios', 'studio', 'US')
from unnest(array['ben-10-2005', 'ben-10-alien-force', 'ben-10-ultimate-alien', 'ben-10-omniverse', 'ben-10-2016']) s;
-- the TV series and every film
select catalog.credit(slug, 'TMS Entertainment', 'studio', 'JP')
from catalog.node n join catalog.entry e on e.id = n.id
where n.slug in ('detective-conan-anime', 'lupin-the-3rd-vs-detective-conan-the-movie')
   or (n.slug like 'detective-conan-%' and e.format = 'movie');

-- publishers
select catalog.credit(s, 'Shogakukan', 'publisher', 'JP')
from unnest(array['detective-conan-manga', 'detective-conan-short-stories',
  'detective-conan-the-culprit-hanzawa', 'detective-conan-zeros-tea-time', 'detective-conan-wild-police-story']) s;
select catalog.credit('occultic-nine-ln', 'Overlap', 'publisher', 'JP');

-- games
select catalog.credit('devil-may-cry', 'Capcom', 'developer', 'JP');
select catalog.credit('devil-may-cry', 'Capcom', 'publisher', 'JP');
select catalog.credit('dmc-devil-may-cry', 'Ninja Theory', 'developer', 'GB');
select catalog.credit('dmc-devil-may-cry', 'Capcom', 'publisher', 'JP');
select catalog.credit(s, 'Mages', 'developer', 'JP')
from unnest(array['steins-gate-vn', 'steins-gate-0-vn', 'steins-gate-elite', 'robotics-notes-vn',
  'robotics-notes-dash', 'chaos-head-vn', 'chaos-child-vn', 'occultic-nine-vn', 'anonymous-code-vn']) s;
select catalog.credit(s, 'Nitroplus', 'developer', 'JP') from unnest(array['steins-gate-vn', 'chaos-head-vn']) s;
select catalog.credit('anonymous-code-vn', 'Chiyomaru Studio', 'developer', 'JP');

-- descriptions, in our own words
update catalog.entry e set description = d.description
from (values
  ('steins-gate-anime', 'A self-styled mad scientist and his friends rig a microwave that sends text messages into the past. Every message rewrites the present a little more, until one change costs far more than they meant to spend.'),
  ('madoka-tv', 'A schoolgirl is offered any wish she likes in exchange for becoming a magical girl who fights witches. The longer she weighs the offer, the clearer it gets that the bargain isn''t what it looks like.'),
  ('bakemonogatari', 'Koyomi Araragi, still half a vampire after one very bad spring break, keeps running into classmates caught up in strange afflictions, and keeps trying to help.'),
  ('detective-conan-anime', 'A teenage detective is poisoned and shrunk into a child''s body. Hiding who he is from almost everyone, he keeps solving cases while hunting the organization that did it.'),
  ('konosuba-s1', 'A shut-in dies, is reborn into a fantasy world with a goddess as the one thing he brought along, and ends up leading the least capable adventuring party in town.'),
  ('steins-gate-vn', 'A group of friends in Akihabara discovers how to send messages into the past. Their experiments draw them into a struggle over who gets to control time.'),
  ('steins-gate-0-vn', 'An alternate path follows Okabe after he abandons his attempt to save Kurisu, when an artificial intelligence brings her memories back into his life.'),
  ('steins-gate-0-anime', 'Okabe tries to live with the outcome of a failed rescue, but an artificial intelligence based on Kurisu pulls him back into the struggle over time travel.'),
  ('steins-gate-elite', 'An animated retelling of the original visual novel, combining its branching story with footage from the anime and new animation for the other routes.'),
  ('steins-gate-movie', 'Kurisu returns to Japan and discovers that Okabe is beginning to disappear, leaving her to confront the consequences of his journeys across worldlines.'),
  ('robotics-notes-vn', 'A struggling school robotics club on Tanegashima sets out to build a giant robot and uncovers a conspiracy reaching far beyond its workshop.'),
  ('robotics-notes-anime', 'Two students try to save their school robotics club by building a giant robot, while a mysterious set of reports points to a much larger threat.'),
  ('robotics-notes-dash', 'The former robotics club members reunite on Tanegashima and meet an eccentric visitor whose arrival draws them into another crisis.'),
  ('chaos-head-vn', 'A withdrawn student in Shibuya finds himself linked to a series of grotesque murders and struggles to distinguish his delusions from reality.'),
  ('chaos-head-anime', 'After witnessing a murder, a reclusive student becomes entangled in a series of killings that erodes his trust in both other people and his own senses.'),
  ('chaos-child-vn', 'Years after a disaster in Shibuya, a school newspaper club investigates murders that seem to repeat an earlier series of killings.'),
  ('chaos-child-anime', 'A student newspaper club investigates a new wave of murders in Shibuya and discovers connections to the catastrophe that changed the city.'),
  ('occultic-nine-ln', 'A teenager running an occult blog becomes connected to eight other people as apparently unrelated supernatural mysteries converge.'),
  ('occultic-nine-anime', 'An occult blogger and a circle of strangers are drawn together by a murder and a string of increasingly disturbing paranormal events.'),
  ('occultic-nine-vn', 'An occult blogger investigates strange events with a group of unlikely allies in a visual novel adaptation of the supernatural mystery series.'),
  ('anonymous-code-vn', 'In a future Tokyo scarred by a technological disaster, a young hacker meets a mysterious girl and discovers a power that can overturn his world.'),
  ('index-s1', 'In a city built around developing psychic powers, a student who can cancel supernatural abilities meets a nun carrying a library of forbidden magic in her memory.'),
  ('railgun-s1', 'A powerful teenage electromaster and her friends investigate incidents in Academy City, where scientific progress often hides dangerous experiments.'),
  ('accelerator', 'Academy City''s strongest psychic tries to protect a young clone while a new enemy threatens the fragile life they have begun to build.'),
  ('item', 'Four young women working in Academy City''s underworld take on covert jobs whose dangers reach beyond their clients'' instructions.'),
  ('ben-10-2005', 'Ten-year-old Ben Tennyson finds an alien device that lets him transform into other species, turning a summer road trip with his family into a series of adventures.'),
  ('ben-10-2016', 'A new version of ten-year-old Ben discovers the Omnitrix and uses its alien transformations during a summer road trip with Gwen and Grandpa Max.'),
  ('devil-may-cry', 'Demon hunter Dante travels to an isolated island after a mysterious woman offers him a chance to confront the power behind his family''s tragedy.'),
  ('dmc-devil-may-cry', 'A reimagined young Dante fights through a city controlled by demons and discovers his place in a rebellion against their rule.'),
  ('shameless-uk', 'On a Manchester council estate, the Gallagher children navigate work, relationships and trouble while their unreliable father Frank pursues his own schemes.'),
  ('shameless-us', 'The Gallagher siblings try to keep their family afloat on Chicago''s South Side while their father Frank creates as many problems as he avoids.'),
  ('bizim-hikaye', 'Filiz raises her younger siblings in an Istanbul neighborhood while managing her father''s drinking and a romance that complicates her responsibilities.'),
  ('naruto', 'An ostracized young ninja pursues his dream of becoming his village''s leader while training with teammates who each carry ambitions and secrets of their own.'),
  ('naruto-shippuden', 'After years of training, Naruto returns to his village and faces a growing conflict involving his missing friend Sasuke and the organization hunting the tailed beasts.'),
  ('konosuba-legend-of-crimson', 'Kazuma''s party travels to Megumin and Yunyun''s home village after a warning that the Crimson Demons are in danger.')) as d (slug, description)
where e.id = catalog.id(d.slug);

-- Fixed totals have matching progress units; specials remain extra units.
select catalog.add_units(n.slug, e.unit_total)
from catalog.entry e join catalog.node n on n.id = e.id
where e.unit_kind in ('episode', 'chapter') and e.unit_total is not null
  and not exists (select 1 from catalog.unit u where u.entry_id = e.id);

select catalog.add_units('occultic-nine-ln', 3);

-- The reboot's three 2021 specials are outside its 166 regular episodes.
-- https://en.wikipedia.org/wiki/List_of_Ben_10_(2016_TV_series)_episodes
select catalog.add_special('ben-10-2016', 'Ben 10,010');
select catalog.add_special('ben-10-2016', 'Ben Gen 10');
select catalog.add_special('ben-10-2016', 'Alien X-Tinction');
