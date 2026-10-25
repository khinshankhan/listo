-- Tags, once every entry they name exists.

select catalog.tag_with('genre', 'action', 'index-s1', 'index-ii', 'railgun-s1', 'railgun-s', 'naruto',
                        'naruto-shippuden', 'devil-may-cry', 'dmc-devil-may-cry');
select catalog.tag_with('genre', 'comedy', 'konosuba-s1', 'konosuba-s2', 'konosuba-s3');
select catalog.tag_with('genre', 'sci-fi', 'steins-gate-anime', 'steins-gate-0-anime', 'index-s1', 'railgun-s1');
select catalog.tag_with('theme', 'time travel', 'steins-gate-vn', 'steins-gate-anime', 'steins-gate-0-anime');

select catalog.tag_with('genre', 'mystery', 'detective-conan-manga', 'detective-conan-anime',
  'detective-conan-films', 'detective-conan-short-stories', 'detective-conan-zeros-tea-time',
  'detective-conan-wild-police-story', 'chaos-head-vn', 'chaos-head-anime', 'chaos-child-vn',
  'chaos-child-anime', 'occultic-nine-ln', 'occultic-nine-anime', 'occultic-nine-vn');
select catalog.tag_with('genre', 'mystery', variadic array(
  select n.slug from catalog.membership m join catalog.node n on n.id = m.member_id
  where m.group_id = catalog.id('detective-conan-films')));
select catalog.tag_with('genre', 'comedy', 'detective-conan-the-culprit-hanzawa', 'shameless-uk',
  'shameless-us', 'bizim-hikaye', 'konosuba-legend-of-crimson');
select catalog.tag_with('genre', 'drama', 'shameless-uk', 'shameless-us', 'bizim-hikaye');
select catalog.tag_with('genre', 'fantasy', 'madoka-tv', 'madoka-beginnings', 'madoka-eternal',
  'madoka-rebellion', 'konosuba-s1', 'konosuba-s2', 'konosuba-s3', 'konosuba-legend-of-crimson');
select catalog.tag_with('genre', 'supernatural', 'monogatari', 'bakemonogatari', 'nisemonogatari',
  'nekomonogatari-black', 'monogatari-second-season', 'hanamonogatari', 'tsukimonogatari',
  'owarimonogatari', 'koyomimonogatari', 'kizumonogatari-1', 'kizumonogatari-2',
  'kizumonogatari-3', 'owarimonogatari-2', 'zoku-owarimonogatari', 'monogatari-off-and-monster');
select catalog.tag_with('genre', 'action', 'index-iii', 'railgun-t', 'accelerator', 'item',
  'ben-10-2005', 'ben-10-alien-force', 'ben-10-ultimate-alien', 'ben-10-omniverse', 'ben-10-2016',
  'spider-man-2002', 'spider-man-2', 'spider-man-3', 'amazing-spider-man', 'amazing-spider-man-2');
