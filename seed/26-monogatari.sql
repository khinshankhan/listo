-- Monogatari: a recommended path and the actual anime release order.
-- https://www.crunchyroll.com/news/guides/2024/3/10/a-guide-to-all-monogatari-arcs-in-order
-- https://en.wikipedia.org/wiki/List_of_Monogatari_episodes

select catalog.add_group('monogatari', 'series', 'Monogatari', p_public_id => 767);
select catalog.add_entry('bakemonogatari', 'anime', 'tv', 'Bakemonogatari', '2009-07', 'JP', 'episode', 15, p_public_id => 768);
select catalog.add_entry('nisemonogatari', 'anime', 'tv', 'Nisemonogatari', '2012-01', 'JP', 'episode', 11, p_public_id => 769);
select catalog.add_entry('nekomonogatari-black', 'anime', 'tv', 'Nekomonogatari (Black)', '2012-12', 'JP', 'episode', 4, p_public_id => 770);
select catalog.add_entry('monogatari-second-season', 'anime', 'tv', 'Monogatari Series: Second Season', '2013-07', 'JP', 'episode', 26, p_public_id => 771);
select catalog.add_entry('hanamonogatari', 'anime', 'tv', 'Hanamonogatari', '2014-08', 'JP', 'episode', 5, p_public_id => 772);
select catalog.add_entry('tsukimonogatari', 'anime', 'tv', 'Tsukimonogatari', '2014-12', 'JP', 'episode', 4, p_public_id => 773);
-- Episode 1 was a double broadcast; home video and international editions split it.
select catalog.add_entry('owarimonogatari', 'anime', 'tv', 'Owarimonogatari', '2015-10', 'JP', 'episode', 13, p_public_id => 774);
select catalog.add_entry('koyomimonogatari', 'anime', 'ona', 'Koyomimonogatari', '2016-01', 'JP', 'episode', 12, p_public_id => 775);
select catalog.add_entry('kizumonogatari-1', 'anime', 'movie', 'Kizumonogatari I: Tekketsu', '2016-01', 'JP', p_public_id => 776);
select catalog.add_entry('kizumonogatari-2', 'anime', 'movie', 'Kizumonogatari II: Nekketsu', '2016-08', 'JP', p_public_id => 777);
select catalog.add_entry('kizumonogatari-3', 'anime', 'movie', 'Kizumonogatari III: Reiketsu', '2017-01', 'JP', p_public_id => 778);
select catalog.add_entry('owarimonogatari-2', 'anime', 'tv', 'Owarimonogatari Second Season', '2017-08', 'JP', 'episode', 7, p_public_id => 779);
-- The six TV episodes premiered in May 2019, after a 2018 theatrical preview.
-- https://www.monogatari-series.com/zokuowarimonogatari/
select catalog.add_entry('zoku-owarimonogatari', 'anime', 'tv', 'Zoku Owarimonogatari', '2019-05', 'JP', 'episode', 6, p_public_id => 780);
select catalog.add_entry('monogatari-off-and-monster', 'anime', 'ona', 'Monogatari Series: Off & Monster Season', '2024-07', 'JP', 'episode', 14, p_public_id => 781);

select catalog.add_units(e, n) from (values
  ('bakemonogatari', 15), ('nisemonogatari', 11), ('nekomonogatari-black', 4),
  ('monogatari-second-season', 26), ('hanamonogatari', 5), ('tsukimonogatari', 4),
  ('owarimonogatari', 13), ('koyomimonogatari', 12), ('owarimonogatari-2', 7),
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
select catalog.relate('nekomonogatari-black', 'sequel_of', 'kizumonogatari-3', p_note => 'Golden Week, after the spring-break prequel.');
select catalog.relate('bakemonogatari', 'sequel_of', 'nekomonogatari-black', p_note => 'The main story follows the Golden Week prequel.');
select catalog.relate('monogatari-second-season', 'sequel_of', 'nisemonogatari');
select catalog.relate('hanamonogatari', 'side_story_of', 'monogatari-second-season', p_note => 'Suruga Devil takes place after Araragi graduates.');
select catalog.relate('koyomimonogatari', 'side_story_of', 'bakemonogatari', p_note => 'Monthly stories spanning the school year; the ending leads into Owarimonogatari Second Season.');
select catalog.relate('tsukimonogatari', 'sequel_of', 'monogatari-second-season');
select catalog.relate('owarimonogatari', 'sequel_of', 'tsukimonogatari');
select catalog.relate('owarimonogatari-2', 'sequel_of', 'owarimonogatari');
select catalog.relate('owarimonogatari-2', 'sequel_of', 'koyomimonogatari');
select catalog.relate('zoku-owarimonogatari', 'sequel_of', 'owarimonogatari-2');
select catalog.relate('monogatari-off-and-monster', 'sequel_of', 'zoku-owarimonogatari');

select catalog.add_ordering('monogatari', 'recommended', 'first-timer order',
  'Start with Bakemonogatari, then its Kizumonogatari prequel, and follow the anime seasons. Koyomimonogatari is required before Owarimonogatari Second Season.') as o \gset
select catalog.add_step(:'o', 1, 'bakemonogatari');
select catalog.add_step(:'o', 2, 'kizumonogatari-1', p_step => 1,
  p_note => 'The spring-break prequel, watched after the introductory series.');
select catalog.add_step(:'o', 2, 'kizumonogatari-2', p_step => 2);
select catalog.add_step(:'o', 2, 'kizumonogatari-3', p_step => 3);
select catalog.add_step(:'o', 3, 'nisemonogatari');
select catalog.add_step(:'o', 4, 'nekomonogatari-black');
select catalog.add_step(:'o', 5, 'monogatari-second-season');
select catalog.add_step(:'o', 6, 'hanamonogatari', p_note => 'Set after graduation, but released here.');
select catalog.add_step(:'o', 7, 'tsukimonogatari');
select catalog.add_step(:'o', 8, 'owarimonogatari');
select catalog.add_step(:'o', 9, 'koyomimonogatari', p_note => 'The final episodes lead directly into the next season.');
select catalog.add_step(:'o', 10, 'owarimonogatari-2');
select catalog.add_step(:'o', 11, 'zoku-owarimonogatari', p_necessity => 'recommended');
select catalog.add_step(:'o', 12, 'monogatari-off-and-monster', p_necessity => 'recommended');

select catalog.add_ordering('monogatari', 'release', 'release order',
  'First release of each seeded anime entry. Koyomimonogatari began between the first two Kizumonogatari films; Zoku uses its six-episode TV release, after the theatrical preview.') as o \gset
select catalog.add_step(:'o', 1, 'bakemonogatari');
select catalog.add_step(:'o', 2, 'nisemonogatari');
select catalog.add_step(:'o', 3, 'nekomonogatari-black');
select catalog.add_step(:'o', 4, 'monogatari-second-season');
select catalog.add_step(:'o', 5, 'hanamonogatari');
select catalog.add_step(:'o', 6, 'tsukimonogatari');
select catalog.add_step(:'o', 7, 'owarimonogatari');
select catalog.add_step(:'o', 8, 'kizumonogatari-1');
select catalog.add_step(:'o', 9, 'koyomimonogatari');
select catalog.add_step(:'o', 10, 'kizumonogatari-2');
select catalog.add_step(:'o', 11, 'kizumonogatari-3');
select catalog.add_step(:'o', 12, 'owarimonogatari-2');
select catalog.add_step(:'o', 13, 'zoku-owarimonogatari');
select catalog.add_step(:'o', 14, 'monogatari-off-and-monster');

-- https://en.wikipedia.org/wiki/List_of_Monogatari_novels
-- https://en.wikipedia.org/wiki/Monogatari_(series)
-- https://www.monogatari-series.com/oms/
-- https://www.kodansha.co.jp/book/products/0000427391
select catalog.add_entry('bakemonogatari-ln', 'book', 'light_novel', 'Bakemonogatari', '2006-11-02', 'JP', 'volume', 2, p_public_id => 782);
select catalog.member('monogatari', 'bakemonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2006-12-05', end_precision = 'day' where id = catalog.id('bakemonogatari-ln');
select catalog.add_units('bakemonogatari-ln', 2);
select catalog.credit('bakemonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('kizumonogatari-ln', 'book', 'light_novel', 'Kizumonogatari', '2008-05-08', 'JP', 'volume', 1, p_public_id => 783);
select catalog.member('monogatari', 'kizumonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2008-05-08', end_precision = 'day' where id = catalog.id('kizumonogatari-ln');
select catalog.add_units('kizumonogatari-ln', 1);
select catalog.credit('kizumonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('nisemonogatari-ln', 'book', 'light_novel', 'Nisemonogatari', '2008-09-02', 'JP', 'volume', 2, p_public_id => 784);
select catalog.member('monogatari', 'nisemonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2009-06-11', end_precision = 'day' where id = catalog.id('nisemonogatari-ln');
select catalog.add_units('nisemonogatari-ln', 2);
select catalog.credit('nisemonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('nekomonogatari-black-ln', 'book', 'light_novel', 'Nekomonogatari (Black)', '2010-07-29', 'JP', 'volume', 1, p_public_id => 785);
select catalog.member('monogatari', 'nekomonogatari-black-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2010-07-29', end_precision = 'day' where id = catalog.id('nekomonogatari-black-ln');
select catalog.add_units('nekomonogatari-black-ln', 1);
select catalog.credit('nekomonogatari-black-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('nekomonogatari-white-ln', 'book', 'light_novel', 'Nekomonogatari (White)', '2010-10-28', 'JP', 'volume', 1, p_public_id => 786);
select catalog.member('monogatari', 'nekomonogatari-white-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2010-10-28', end_precision = 'day' where id = catalog.id('nekomonogatari-white-ln');
select catalog.add_units('nekomonogatari-white-ln', 1);
select catalog.credit('nekomonogatari-white-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('kabukimonogatari-ln', 'book', 'light_novel', 'Kabukimonogatari', '2010-12-25', 'JP', 'volume', 1, p_public_id => 787);
select catalog.member('monogatari', 'kabukimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2010-12-25', end_precision = 'day' where id = catalog.id('kabukimonogatari-ln');
select catalog.add_units('kabukimonogatari-ln', 1);
select catalog.credit('kabukimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('hanamonogatari-ln', 'book', 'light_novel', 'Hanamonogatari', '2011-03-30', 'JP', 'volume', 1, p_public_id => 788);
select catalog.member('monogatari', 'hanamonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2011-03-30', end_precision = 'day' where id = catalog.id('hanamonogatari-ln');
select catalog.add_units('hanamonogatari-ln', 1);
select catalog.credit('hanamonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('otorimonogatari-ln', 'book', 'light_novel', 'Otorimonogatari', '2011-06-29', 'JP', 'volume', 1, p_public_id => 789);
select catalog.member('monogatari', 'otorimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2011-06-29', end_precision = 'day' where id = catalog.id('otorimonogatari-ln');
select catalog.add_units('otorimonogatari-ln', 1);
select catalog.credit('otorimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('onimonogatari-ln', 'book', 'light_novel', 'Onimonogatari', '2011-09-29', 'JP', 'volume', 1, p_public_id => 790);
select catalog.member('monogatari', 'onimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2011-09-29', end_precision = 'day' where id = catalog.id('onimonogatari-ln');
select catalog.add_units('onimonogatari-ln', 1);
select catalog.credit('onimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('koimonogatari-ln', 'book', 'light_novel', 'Koimonogatari', '2011-12-21', 'JP', 'volume', 1, p_public_id => 791);
select catalog.member('monogatari', 'koimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2011-12-21', end_precision = 'day' where id = catalog.id('koimonogatari-ln');
select catalog.add_units('koimonogatari-ln', 1);
select catalog.credit('koimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('tsukimonogatari-ln', 'book', 'light_novel', 'Tsukimonogatari', '2012-09-27', 'JP', 'volume', 1, p_public_id => 792);
select catalog.member('monogatari', 'tsukimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2012-09-27', end_precision = 'day' where id = catalog.id('tsukimonogatari-ln');
select catalog.add_units('tsukimonogatari-ln', 1);
select catalog.credit('tsukimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('koyomimonogatari-ln', 'book', 'light_novel', 'Koyomimonogatari', '2013-05-21', 'JP', 'volume', 1, p_public_id => 793);
select catalog.member('monogatari', 'koyomimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2013-05-21', end_precision = 'day' where id = catalog.id('koyomimonogatari-ln');
select catalog.add_units('koyomimonogatari-ln', 1);
select catalog.credit('koyomimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('owarimonogatari-ln', 'book', 'light_novel', 'Owarimonogatari', '2013-10-22', 'JP', 'volume', 3, p_public_id => 794);
select catalog.member('monogatari', 'owarimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2014-04-02', end_precision = 'day' where id = catalog.id('owarimonogatari-ln');
select catalog.add_units('owarimonogatari-ln', 3);
select catalog.credit('owarimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('zoku-owarimonogatari-ln', 'book', 'light_novel', 'Zoku Owarimonogatari', '2014-09-18', 'JP', 'volume', 1, p_public_id => 795);
select catalog.member('monogatari', 'zoku-owarimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2014-09-18', end_precision = 'day' where id = catalog.id('zoku-owarimonogatari-ln');
select catalog.add_units('zoku-owarimonogatari-ln', 1);
select catalog.credit('zoku-owarimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('orokamonogatari-ln', 'book', 'light_novel', 'Orokamonogatari', '2015-10-06', 'JP', 'volume', 1, p_public_id => 796);
select catalog.member('monogatari', 'orokamonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2015-10-06', end_precision = 'day' where id = catalog.id('orokamonogatari-ln');
select catalog.add_units('orokamonogatari-ln', 1);
select catalog.credit('orokamonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('wazamonogatari-ln', 'book', 'light_novel', 'Wazamonogatari', '2016-01-14', 'JP', 'volume', 1, p_public_id => 797);
select catalog.member('monogatari', 'wazamonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2016-01-14', end_precision = 'day' where id = catalog.id('wazamonogatari-ln');
select catalog.add_units('wazamonogatari-ln', 1);
select catalog.credit('wazamonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('nademonogatari-ln', 'book', 'light_novel', 'Nademonogatari', '2016-07-28', 'JP', 'volume', 1, p_public_id => 798);
select catalog.member('monogatari', 'nademonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2016-07-28', end_precision = 'day' where id = catalog.id('nademonogatari-ln');
select catalog.add_units('nademonogatari-ln', 1);
select catalog.credit('nademonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('musubimonogatari-ln', 'book', 'light_novel', 'Musubimonogatari', '2017-01-12', 'JP', 'volume', 1, p_public_id => 799);
select catalog.member('monogatari', 'musubimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2017-01-12', end_precision = 'day' where id = catalog.id('musubimonogatari-ln');
select catalog.add_units('musubimonogatari-ln', 1);
select catalog.credit('musubimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('shinobumonogatari-ln', 'book', 'light_novel', 'Shinobumonogatari', '2017-07-20', 'JP', 'volume', 1, p_public_id => 800);
select catalog.member('monogatari', 'shinobumonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2017-07-20', end_precision = 'day' where id = catalog.id('shinobumonogatari-ln');
select catalog.add_units('shinobumonogatari-ln', 1);
select catalog.credit('shinobumonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('yoimonogatari-ln', 'book', 'light_novel', 'Yoimonogatari', '2018-06-14', 'JP', 'volume', 1, p_public_id => 801);
select catalog.member('monogatari', 'yoimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2018-06-14', end_precision = 'day' where id = catalog.id('yoimonogatari-ln');
select catalog.add_units('yoimonogatari-ln', 1);
select catalog.credit('yoimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('amarimonogatari-ln', 'book', 'light_novel', 'Amarimonogatari', '2019-04-17', 'JP', 'volume', 1, p_public_id => 802);
select catalog.member('monogatari', 'amarimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2019-04-17', end_precision = 'day' where id = catalog.id('amarimonogatari-ln');
select catalog.add_units('amarimonogatari-ln', 1);
select catalog.credit('amarimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('ogimonogatari-ln', 'book', 'light_novel', 'Ogimonogatari', '2020-10-28', 'JP', 'volume', 1, p_public_id => 803);
select catalog.member('monogatari', 'ogimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2020-10-28', end_precision = 'day' where id = catalog.id('ogimonogatari-ln');
select catalog.add_units('ogimonogatari-ln', 1);
select catalog.credit('ogimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('shinomonogatari-ln', 'book', 'light_novel', 'Shinomonogatari', '2021-08-19', 'JP', 'volume', 2, p_public_id => 804);
select catalog.member('monogatari', 'shinomonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2021-08-19', end_precision = 'day' where id = catalog.id('shinomonogatari-ln');
select catalog.add_units('shinomonogatari-ln', 2);
select catalog.credit('shinomonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('ikusamonogatari-ln', 'book', 'light_novel', 'Ikusamonogatari', '2023-05-17', 'JP', 'volume', 1, p_public_id => 805);
select catalog.member('monogatari', 'ikusamonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2023-05-17', end_precision = 'day' where id = catalog.id('ikusamonogatari-ln');
select catalog.add_units('ikusamonogatari-ln', 1);
select catalog.credit('ikusamonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('tsugimonogatari-ln', 'book', 'light_novel', 'Tsugimonogatari', '2025-10-16', 'JP', 'volume', 1, p_public_id => 806);
select catalog.member('monogatari', 'tsugimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2025-10-16', end_precision = 'day' where id = catalog.id('tsugimonogatari-ln');
select catalog.add_units('tsugimonogatari-ln', 1);
select catalog.credit('tsugimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('torimonogatari-ln', 'book', 'light_novel', 'Torimonogatari', '2026-09-04', 'JP', 'volume', 1, p_public_id => 807);
select catalog.member('monogatari', 'torimonogatari-ln');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2026-09-04', end_precision = 'day' where id = catalog.id('torimonogatari-ln');
select catalog.add_units('torimonogatari-ln', 1);
select catalog.credit('torimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('ootorimonogatari-ln', 'book', 'light_novel', 'Ootorimonogatari', null, 'JP', 'volume', null, p_public_id => 808);
select catalog.member('monogatari', 'ootorimonogatari-ln');
update catalog.entry set status = 'announced', source_type = 'original' where id = catalog.id('ootorimonogatari-ln');
select catalog.credit('ootorimonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('mijikanamonogatari-ln', 'book', 'light_novel', 'Mijikanamonogatari', '2024-09-11', 'JP', 'volume', 1, p_public_id => 809);
select catalog.member('monogatari', 'mijikanamonogatari-ln');
update catalog.entry set status = 'finished', description = 'A published collection of short Monogatari stories originally issued across assorted booklets and promotional publications.' where id = catalog.id('mijikanamonogatari-ln');
select catalog.add_units('mijikanamonogatari-ln', 1);
select catalog.credit('mijikanamonogatari-ln', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('bakemonogatari-manga', 'manga', 'serial', 'Bakemonogatari', '2018-03-14', 'JP', 'volume', 22, p_public_id => 810);
select catalog.member('monogatari', 'bakemonogatari-manga');
select catalog.relate('bakemonogatari-manga', 'adaptation_of', 'bakemonogatari-ln');
update catalog.entry set status = 'finished', end_date = '2023-03-15', end_precision = 'day' where id = catalog.id('bakemonogatari-manga');
select catalog.add_units('bakemonogatari-manga', 22);
select catalog.credit('bakemonogatari-manga', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('bakemonogatari-portable', 'visual_novel', 'full', 'Bakemonogatari Portable', '2012-08-23', 'JP', 'none', null, p_public_id => 811);
select catalog.member('monogatari', 'bakemonogatari-portable');
select catalog.relate('bakemonogatari-portable', 'adaptation_of', 'bakemonogatari-ln');
update catalog.entry set status = 'finished' where id = catalog.id('bakemonogatari-portable');
select catalog.add_entry('monogatari-puc-puc', 'game', 'full', 'Monogatari Series Puc Puc', '2018-08-21', 'JP', 'none', null, p_public_id => 812);
select catalog.member('monogatari', 'monogatari-puc-puc');
select catalog.relate('monogatari-puc-puc', 'spin_off_of', 'bakemonogatari');
update catalog.entry set status = 'finished', end_date = '2020-03-31', end_precision = 'day' where id = catalog.id('monogatari-puc-puc');
select catalog.add_entry('kizumonogatari-koyomi-vamp', 'anime', 'movie', 'Kizumonogatari: Koyomi Vamp', '2024-01-12', 'JP', 'none', null, p_public_id => 813);
select catalog.member('monogatari', 'kizumonogatari-koyomi-vamp');
select catalog.relate('kizumonogatari-koyomi-vamp', 'recap_of', 'kizumonogatari-1');
update catalog.entry set status = 'finished', description = 'A recut theatrical presentation condenses the three Kizumonogatari films into one feature.', length_minutes = 144 where id = catalog.id('kizumonogatari-koyomi-vamp');
select catalog.add_entry('zoku-owarimonogatari-theatrical', 'anime', 'movie', 'Zoku Owarimonogatari (theatrical edition)', '2018-11-10', 'JP', 'none', null, p_public_id => 814);
select catalog.member('monogatari', 'zoku-owarimonogatari-theatrical');
select catalog.relate('zoku-owarimonogatari-theatrical', 'alternative_version_of', 'zoku-owarimonogatari');
update catalog.entry set status = 'finished', description = 'The theatrical presentation released before the six-episode television version in 2019.', length_minutes = 148 where id = catalog.id('zoku-owarimonogatari-theatrical');
select catalog.add_entry('monogatari-off-and-monster-2', 'anime', 'ona', 'Monogatari Series: Off & Monster Season - Part 2', '2026', 'JP', 'episode', null, p_public_id => 815);
select catalog.member('monogatari', 'monogatari-off-and-monster-2');
select catalog.relate('monogatari-off-and-monster-2', 'sequel_of', 'monogatari-off-and-monster');
update catalog.entry set status = 'announced', description = 'A continuation announced for winter 2026, including Karen Ogre. The official site has not supplied a precise premiere date in the checked material.' where id = catalog.id('monogatari-off-and-monster-2');
select catalog.add_entry('monogatari-hyakumonogatari', 'audio', 'drama', 'Hyakumonogatari', '2009-08-03', 'JP', 'none', null, p_public_id => 816);
select catalog.member('monogatari', 'monogatari-hyakumonogatari');
select catalog.relate('monogatari-hyakumonogatari', 'spin_off_of', 'bakemonogatari-ln');
update catalog.entry set status = 'finished' where id = catalog.id('monogatari-hyakumonogatari');

select catalog.relate('bakemonogatari', 'adaptation_of', 'bakemonogatari-ln');
select catalog.relate('nisemonogatari', 'adaptation_of', 'nisemonogatari-ln');
select catalog.relate('nekomonogatari-black', 'adaptation_of', 'nekomonogatari-black-ln');
select catalog.relate('monogatari-second-season', 'adaptation_of', 'nekomonogatari-white-ln');
select catalog.relate('monogatari-second-season', 'adaptation_of', 'kabukimonogatari-ln');
select catalog.relate('monogatari-second-season', 'adaptation_of', 'otorimonogatari-ln');
select catalog.relate('monogatari-second-season', 'adaptation_of', 'onimonogatari-ln');
select catalog.relate('monogatari-second-season', 'adaptation_of', 'koimonogatari-ln');
select catalog.relate('hanamonogatari', 'adaptation_of', 'hanamonogatari-ln');
select catalog.relate('tsukimonogatari', 'adaptation_of', 'tsukimonogatari-ln');
select catalog.relate('koyomimonogatari', 'adaptation_of', 'koyomimonogatari-ln');
select catalog.relate('owarimonogatari', 'adaptation_of', 'owarimonogatari-ln');
select catalog.relate('owarimonogatari-2', 'adaptation_of', 'owarimonogatari-ln');
select catalog.relate('zoku-owarimonogatari', 'adaptation_of', 'zoku-owarimonogatari-ln');
select catalog.relate('kizumonogatari-1', 'adaptation_of', 'kizumonogatari-ln');
select catalog.relate('kizumonogatari-2', 'adaptation_of', 'kizumonogatari-ln');
select catalog.relate('kizumonogatari-3', 'adaptation_of', 'kizumonogatari-ln');
select catalog.relate('monogatari-off-and-monster', 'adaptation_of', 'orokamonogatari-ln');
select catalog.relate('monogatari-off-and-monster', 'adaptation_of', 'nademonogatari-ln');
select catalog.relate('monogatari-off-and-monster', 'adaptation_of', 'wazamonogatari-ln');
select catalog.relate('monogatari-off-and-monster', 'adaptation_of', 'shinobumonogatari-ln');
select catalog.relate('kizumonogatari-koyomi-vamp', 'recap_of', 'kizumonogatari-2');
select catalog.relate('kizumonogatari-koyomi-vamp', 'recap_of', 'kizumonogatari-3');
select catalog.credit('kizumonogatari-koyomi-vamp', 'Shaft', 'studio', 'JP');
select catalog.credit('zoku-owarimonogatari-theatrical', 'Shaft', 'studio', 'JP');

select catalog.add_entry('monogatari-koyomi-history', 'anime', 'ova', 'Koyomi History', '2016-09-22', 'JP', 'episode', 1, p_public_id => 817);
select catalog.member('monogatari', 'monogatari-koyomi-history');
select catalog.relate('monogatari-koyomi-history', 'side_story_of', 'bakemonogatari', p_note => 'Home-video release of the short shown at the earlier Madogatari exhibition.');
update catalog.entry set status = 'finished' where id = catalog.id('monogatari-koyomi-history');
select catalog.add_entry('monogatari-beautiful-princess', 'anime', 'ona', 'Monogatari: A Cruel Fairy Tale - The Beautiful Princess', '2024-08-17', 'JP', 'episode', 1, p_public_id => 818);
select catalog.member('monogatari', 'monogatari-beautiful-princess');
select catalog.relate('monogatari-beautiful-princess', 'side_story_of', 'monogatari-off-and-monster', p_note => 'Episode 6.5 is extra to the fourteen numbered Off & Monster Season episodes.');
update catalog.entry set status = 'finished' where id = catalog.id('monogatari-beautiful-princess');
select catalog.relate('monogatari-beautiful-princess', 'adaptation_of', 'wazamonogatari-ln');

-- https://www.kodansha.co.jp/book/products/0000188411
-- https://www.kodansha.co.jp/book/products/0000189571
-- https://en.wikipedia.org/wiki/Monogatari_(series)#Novel_series
-- https://en.wikipedia.org/wiki/List_of_Monogatari_novels#Audiobook
-- https://www.audible.co.jp/series/B08WH8LJF5
select catalog.add_entry('monogatari-heroine-books', 'book', 'nonfiction', 'Monogatari Series Heroine Books', '2013-09-20', 'JP', 'volume', 8, p_public_id => 819);
select catalog.member('monogatari', 'monogatari-heroine-books');
select catalog.relate('monogatari-heroine-books', 'references', 'monogatari');
update catalog.entry set status = 'finished', description = 'Eight character books combining illustrations, interviews and original short fiction: Tsubasa, Mayoi, Shinobu, Nadeko, Hitagi, Suruga, the Fire Sisters and Yotsugi.', end_date = '2015-12-23', end_precision = 'day' where id = catalog.id('monogatari-heroine-books');
select catalog.add_units('monogatari-heroine-books', 8);
select catalog.credit('monogatari-heroine-books', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('monogatari-japanese-audiobooks', 'audio', 'audiobook', 'Monogatari Series (Japanese Audiobooks)', '2021-02-17', 'JP', 'volume', null, p_public_id => 820);
select catalog.member('monogatari', 'monogatari-japanese-audiobooks');
select catalog.relate('monogatari-japanese-audiobooks', 'adaptation_of', 'monogatari');
update catalog.entry set status = null, description = 'Kodansha and Audible editions narrated by the anime cast. The series had reached the second Shinomonogatari volume by June 2023; later publication coverage is unconfirmed.' where id = catalog.id('monogatari-japanese-audiobooks');
select catalog.add_entry('kizumonogatari-english-audiobook', 'audio', 'audiobook', 'Kizumonogatari (English Audiobook)', '2016-05-25', 'US', 'none', null, p_public_id => 821);
select catalog.member('monogatari', 'kizumonogatari-english-audiobook');
select catalog.relate('kizumonogatari-english-audiobook', 'adaptation_of', 'kizumonogatari-ln');
update catalog.entry set status = 'finished' where id = catalog.id('kizumonogatari-english-audiobook');
select catalog.add_entry('nekomonogatari-white-english-audiobook', 'audio', 'audiobook', 'Nekomonogatari: White (English Audiobook)', '2019-04-09', 'US', 'none', null, p_public_id => 822);
select catalog.member('monogatari', 'nekomonogatari-white-english-audiobook');
select catalog.relate('nekomonogatari-white-english-audiobook', 'adaptation_of', 'nekomonogatari-white-ln');
update catalog.entry set status = 'finished' where id = catalog.id('nekomonogatari-white-english-audiobook');
select catalog.add_entry('bakemonogatari-english-audiobook', 'audio', 'audiobook', 'Bakemonogatari (English Audiobooks)', '2020-03-24', 'US', 'volume', 3, p_public_id => 823);
select catalog.member('monogatari', 'bakemonogatari-english-audiobook');
select catalog.relate('bakemonogatari-english-audiobook', 'adaptation_of', 'bakemonogatari-ln');
update catalog.entry set status = 'finished' where id = catalog.id('bakemonogatari-english-audiobook');
select catalog.add_units('bakemonogatari-english-audiobook', 3);

-- https://www.kodansha.co.jp/r/book/product?isbn=9784065132920
-- https://x.com/NISIOISIN_info/status/1988200422503538835
-- https://realsound.jp/book/2026/01/post-2271336.html
select catalog.add_entry('mazemonogatari', 'book', 'anthology', 'Mazemonogatari', '2019-02-06', 'JP', 'none', null, p_public_id => 824);
select catalog.member('monogatari', 'mazemonogatari');
select catalog.relate('mazemonogatari', 'crossover_with', 'monogatari');
update catalog.entry set status = 'finished', description = 'Fifteen crossover short stories pairing Monogatari characters with characters from other Nisio Isin works. Twelve first appeared as Kizumonogatari theater handouts; three were added for the collection.' where id = catalog.id('mazemonogatari');
select catalog.ext('mazemonogatari', 'isbn', '9784065132920');
select catalog.credit('mazemonogatari', 'Kodansha', 'publisher', 'JP');
select catalog.add_entry('monogatari-yotsugi-blank', 'book', 'short_story', 'Yotsugi Blank', '2025-11-11', 'JP', 'none', null, p_public_id => 825);
select catalog.member('monogatari', 'monogatari-yotsugi-blank');
select catalog.relate('monogatari-yotsugi-blank', 'side_story_of', 'monogatari');
update catalog.entry set status = 'finished' where id = catalog.id('monogatari-yotsugi-blank');
select catalog.add_entry('monogatari-koyomi-fifties', 'book', 'short_story', 'Koyomi Fifties', '2026-01-08', 'JP', 'none', null, p_public_id => 826);
select catalog.member('monogatari', 'monogatari-koyomi-fifties');
select catalog.relate('monogatari-koyomi-fifties', 'side_story_of', 'monogatari', p_note => 'Distributed at the Shaft 50th Anniversary Exhibition.');
update catalog.entry set status = 'finished' where id = catalog.id('monogatari-koyomi-fifties');
