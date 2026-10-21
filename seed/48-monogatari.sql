-- Monogatari: release and story order disagree, and a first watch wants
-- release. Exercises skip-first-time, sections, a step of several films, and
-- a chronological order beside the recommended one.

select catalog.add_group('monogatari', 'series', 'Monogatari', p_public_id => 103);
select catalog.add_entry('bakemonogatari', 'anime', 'tv', 'Bakemonogatari', '2009-07', 'JP', 'episode', 15, p_public_id => 104);
select catalog.add_entry('nisemonogatari', 'anime', 'tv', 'Nisemonogatari', '2012-01', 'JP', 'episode', 11, p_public_id => 105);
select catalog.add_entry('nekomonogatari-black', 'anime', 'tv', 'Nekomonogatari (Black)', '2012-12', 'JP', 'episode', 4, p_public_id => 106);
select catalog.add_entry('monogatari-second-season', 'anime', 'tv', 'Monogatari Series: Second Season', '2013-07', 'JP', 'episode', 26, p_public_id => 107);
select catalog.add_entry('hanamonogatari', 'anime', 'tv', 'Hanamonogatari', '2014-08', 'JP', 'episode', 5, p_public_id => 108);
select catalog.add_entry('tsukimonogatari', 'anime', 'tv', 'Tsukimonogatari', '2014-12', 'JP', 'episode', 4, p_public_id => 109);
select catalog.add_entry('owarimonogatari', 'anime', 'tv', 'Owarimonogatari', '2015-10', 'JP', 'episode', 12, p_public_id => 110);
select catalog.add_entry('koyomimonogatari', 'anime', 'ona', 'Koyomimonogatari', '2016-01', 'JP', 'episode', 12, p_public_id => 111);
select catalog.add_entry('kizumonogatari-1', 'anime', 'movie', 'Kizumonogatari I: Tekketsu', '2016-01', 'JP', p_public_id => 112);
select catalog.add_entry('kizumonogatari-2', 'anime', 'movie', 'Kizumonogatari II: Nekketsu', '2016-08', 'JP', p_public_id => 113);
select catalog.add_entry('kizumonogatari-3', 'anime', 'movie', 'Kizumonogatari III: Reiketsu', '2017-01', 'JP', p_public_id => 114);
select catalog.add_entry('owarimonogatari-2', 'anime', 'tv', 'Owarimonogatari Second Season', '2017-08', 'JP', 'episode', 7, p_public_id => 115);
select catalog.add_entry('zoku-owarimonogatari', 'anime', 'tv', 'Zoku Owarimonogatari', '2018-11', 'JP', 'episode', 6, p_public_id => 116);
select catalog.add_entry('monogatari-off-and-monster', 'anime', 'ona', 'Monogatari Series: Off & Monster Season', '2024-07', 'JP', 'episode', 14, p_public_id => 117);

select catalog.add_units(e, n) from (values
  ('bakemonogatari', 15), ('nisemonogatari', 11), ('nekomonogatari-black', 4),
  ('monogatari-second-season', 26), ('hanamonogatari', 5), ('tsukimonogatari', 4),
  ('owarimonogatari', 12), ('koyomimonogatari', 12), ('owarimonogatari-2', 7),
  ('zoku-owarimonogatari', 6), ('monogatari-off-and-monster', 14)) as u (e, n);

select catalog.member('monogatari', 'bakemonogatari', 'nisemonogatari', 'nekomonogatari-black',
  'monogatari-second-season', 'hanamonogatari', 'tsukimonogatari', 'owarimonogatari', 'koyomimonogatari',
  'kizumonogatari-1', 'kizumonogatari-2', 'kizumonogatari-3', 'owarimonogatari-2', 'zoku-owarimonogatari',
  'monogatari-off-and-monster');
select catalog.set_role('monogatari', 'kizumonogatari-1', 'prequel');
select catalog.set_role('monogatari', 'nekomonogatari-black', 'prequel');
select catalog.set_role('monogatari', 'koyomimonogatari', 'side stories');

select catalog.relate('nisemonogatari', 'sequel_of', 'bakemonogatari');
select catalog.relate('kizumonogatari-2', 'sequel_of', 'kizumonogatari-1');
select catalog.relate('kizumonogatari-3', 'sequel_of', 'kizumonogatari-2');
select catalog.relate('bakemonogatari', 'sequel_of', 'kizumonogatari-3', p_note => 'in story order');
select catalog.relate('monogatari-second-season', 'sequel_of', 'nisemonogatari');
select catalog.relate('tsukimonogatari', 'sequel_of', 'monogatari-second-season');
select catalog.relate('owarimonogatari', 'sequel_of', 'tsukimonogatari');
select catalog.relate('owarimonogatari-2', 'sequel_of', 'owarimonogatari');
select catalog.relate('zoku-owarimonogatari', 'sequel_of', 'owarimonogatari-2');
select catalog.relate('monogatari-off-and-monster', 'sequel_of', 'zoku-owarimonogatari');

select catalog.add_ordering('monogatari', 'recommended', 'first-timer order',
  'Release and story order disagree, and the series runs on withheld information: Kizumonogatari is a prequel, but watching it first spoils what Bakemonogatari spends fifteen episodes earning. Watch roughly as it aired; the chronological order is for rewatchers.') as o \gset
select catalog.add_step(:'o', 1, 'bakemonogatari');
select catalog.add_step(:'o', 2, 'nisemonogatari');
select catalog.add_step(:'o', 3, 'kizumonogatari-1', p_step => 1,
  p_note => 'A prequel, but watch it here, not first: Bakemonogatari holds back how Araragi became what he is, and Kizu is the payoff.');
select catalog.add_step(:'o', 3, 'kizumonogatari-2', p_step => 2);
select catalog.add_step(:'o', 3, 'kizumonogatari-3', p_step => 3);
select catalog.add_step(:'o', 4, 'nekomonogatari-black');
select catalog.add_step(:'o', 5, 'monogatari-second-season');
select catalog.add_step(:'o', 6, 'hanamonogatari', p_necessity => 'recommended', p_note => 'set after the rest, but aired here');
select catalog.add_step(:'o', 7, 'tsukimonogatari');
select catalog.add_step(:'o', 8, 'koyomimonogatari', p_necessity => 'skip_first_time', p_section => 'side stories');
select catalog.add_step(:'o', 9, 'owarimonogatari');
select catalog.add_step(:'o', 10, 'owarimonogatari-2');
select catalog.add_step(:'o', 11, 'zoku-owarimonogatari', p_necessity => 'recommended');
select catalog.add_step(:'o', 12, 'monogatari-off-and-monster', p_necessity => 'recommended');

select catalog.add_ordering('monogatari', 'chronological', 'chronological',
  'Story order, by where each arc mostly takes place. Several jump around inside themselves, so it''s approximate.') as o \gset
select catalog.add_step(:'o', 1, 'kizumonogatari-1', p_step => 1);
select catalog.add_step(:'o', 1, 'kizumonogatari-2', p_step => 2);
select catalog.add_step(:'o', 1, 'kizumonogatari-3', p_step => 3);
select catalog.add_step(:'o', 2, 'nekomonogatari-black');
select catalog.add_step(:'o', 3, 'bakemonogatari');
select catalog.add_step(:'o', 4, 'nisemonogatari');
select catalog.add_step(:'o', 5, 'monogatari-second-season');
select catalog.add_step(:'o', 6, 'tsukimonogatari');
select catalog.add_step(:'o', 7, 'koyomimonogatari', p_note => 'spans the whole year');
select catalog.add_step(:'o', 8, 'owarimonogatari');
select catalog.add_step(:'o', 9, 'owarimonogatari-2');
select catalog.add_step(:'o', 10, 'zoku-owarimonogatari');
select catalog.add_step(:'o', 11, 'hanamonogatari');
select catalog.add_step(:'o', 12, 'monogatari-off-and-monster');
