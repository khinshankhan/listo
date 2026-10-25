-- KonoSuba: seasons as entries, with a romanized title and a short synonym.

select catalog.add_group('konosuba', 'series', 'KonoSuba', p_public_id => 850);
select catalog.add_entry('konosuba-s1', 'anime', 'tv', 'KonoSuba: God''s Blessing on This Wonderful World!', '2016-01', 'JP', 'episode', 10, p_public_id => 851);
select catalog.add_entry('konosuba-s2', 'anime', 'tv', 'KonoSuba: God''s Blessing on This Wonderful World! 2', '2017-01', 'JP', 'episode', 10, p_public_id => 852);
select catalog.add_entry('konosuba-s3', 'anime', 'tv', 'KonoSuba: God''s Blessing on This Wonderful World! 3', '2024-04', 'JP', 'episode', 11, p_public_id => 853);
-- https://konosuba.com/legend_of_crimson/news/?id=20190716_2&mode=detail
select catalog.add_entry('konosuba-legend-of-crimson', 'anime', 'movie', 'KonoSuba: God''s Blessing on This Wonderful World! Legend of Crimson', '2019-08-30', 'JP', p_public_id => 854);
select catalog.alias('konosuba-s2', 'romanized', 'Kono Subarashii Sekai ni Shukufuku wo! 2');
select catalog.alias('konosuba-s2', 'synonym', 'KonoSuba S2');
select catalog.member('konosuba', 'konosuba-s1', 'konosuba-s2', 'konosuba-legend-of-crimson', 'konosuba-s3');
select catalog.ext('konosuba', 'wikidata', 'Q86838208');
select catalog.relate('konosuba-s2', 'sequel_of', 'konosuba-s1');
select catalog.relate('konosuba-legend-of-crimson', 'sequel_of', 'konosuba-s2');
select catalog.relate('konosuba-s3', 'sequel_of', 'konosuba-legend-of-crimson');

-- https://en.wikipedia.org/wiki/KonoSuba
-- https://en.wikipedia.org/wiki/KonoSuba:_An_Explosion_on_This_Wonderful_World!
-- https://en.wikipedia.org/wiki/List_of_KonoSuba_volumes
-- https://konosuba.com/10th_anniv/news/?id=20260726&mode=detail
-- https://konosuba.com/3rd-bonus/
-- https://isekai-quartet.com/
-- https://kadokawa-anime.com/titles/entry-196.html
select catalog.add_entry('konosuba-web-novel', 'book', 'novel', 'KonoSuba (original web novel)', '2012-12-20', 'JP', 'none', null, p_public_id => 855);
select catalog.member('konosuba', 'konosuba-web-novel');
update catalog.entry set status = 'finished', source_type = 'original', end_date = '2013-10-21', end_precision = 'day' where id = catalog.id('konosuba-web-novel');
select catalog.add_entry('konosuba-ln', 'book', 'light_novel', 'KonoSuba: God''s Blessing on This Wonderful World!', '2013-10-01', 'JP', 'volume', 17, p_public_id => 856);
select catalog.member('konosuba', 'konosuba-ln');
select catalog.relate('konosuba-ln', 'adaptation_of', 'konosuba-web-novel');
update catalog.entry set status = 'finished', end_date = '2020-05-01', end_precision = 'day' where id = catalog.id('konosuba-ln');
select catalog.add_units('konosuba-ln', 17);
select catalog.credit('konosuba-ln', 'Kadokawa Shoten', 'publisher', 'JP');
select catalog.add_entry('konosuba-explosion-ln', 'book', 'light_novel', 'KonoSuba: An Explosion on This Wonderful World!', '2014-07-01', 'JP', 'volume', 3, p_public_id => 857);
select catalog.member('konosuba', 'konosuba-explosion-ln');
select catalog.relate('konosuba-explosion-ln', 'side_story_of', 'konosuba-ln');
update catalog.entry set status = 'finished', end_date = '2015-06-01', end_precision = 'day' where id = catalog.id('konosuba-explosion-ln');
select catalog.add_units('konosuba-explosion-ln', 3);
select catalog.add_entry('konosuba-explosion-bonus-ln', 'book', 'light_novel', 'KonoSuba: An Explosion on This Wonderful World! Bonus Story', '2016-12-28', 'JP', 'volume', 2, p_public_id => 858);
select catalog.member('konosuba', 'konosuba-explosion-bonus-ln');
select catalog.relate('konosuba-explosion-bonus-ln', 'sequel_of', 'konosuba-explosion-ln');
update catalog.entry set status = 'finished', end_date = '2019-03-01', end_precision = 'day' where id = catalog.id('konosuba-explosion-bonus-ln');
select catalog.add_units('konosuba-explosion-bonus-ln', 2);
select catalog.add_entry('konosuba-masked-devil-ln', 'book', 'light_novel', 'KonoSuba: Consulting with This Masked Devil!', '2016-04-01', 'JP', 'volume', 1, p_public_id => 859);
select catalog.member('konosuba', 'konosuba-masked-devil-ln');
select catalog.relate('konosuba-masked-devil-ln', 'side_story_of', 'konosuba-ln');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-masked-devil-ln');
select catalog.add_units('konosuba-masked-devil-ln', 1);
select catalog.add_entry('konosuba-dust-ln', 'book', 'light_novel', 'KonoSuba: God''s Blessing on This Wonderful World! Extra Attention to That Wonderful Fool!', '2017-08-01', 'JP', 'volume', null, p_public_id => 860);
select catalog.member('konosuba', 'konosuba-dust-ln');
select catalog.relate('konosuba-dust-ln', 'side_story_of', 'konosuba-ln');
update catalog.entry set status = null, unit_available = 8 where id = catalog.id('konosuba-dust-ln');
select catalog.add_units('konosuba-dust-ln', 8);
select catalog.add_entry('konosuba-detour-ln', 'book', 'light_novel', 'KonoSuba: God''s Blessing on This Wonderful World! Detour!', '2020-01-01', 'JP', 'volume', 4, p_public_id => 861);
select catalog.member('konosuba', 'konosuba-detour-ln');
select catalog.relate('konosuba-detour-ln', 'side_story_of', 'konosuba-ln');
update catalog.entry set status = 'finished', end_date = '2026-02-28', end_precision = 'day' where id = catalog.id('konosuba-detour-ln');
select catalog.add_units('konosuba-detour-ln', 4);
select catalog.add_entry('konosuba-manga', 'manga', 'serial', 'KonoSuba: God''s Blessing on This Wonderful World!', '2014-09-09', 'JP', 'volume', null, p_public_id => 862);
select catalog.member('konosuba', 'konosuba-manga');
select catalog.relate('konosuba-manga', 'adaptation_of', 'konosuba-ln');
update catalog.entry set status = 'releasing', unit_available = 23 where id = catalog.id('konosuba-manga');
select catalog.add_units('konosuba-manga', 23);
select catalog.add_entry('konosuba-explosion-manga', 'manga', 'serial', 'KonoSuba: An Explosion on This Wonderful World!', '2016-05-27', 'JP', 'volume', 5, p_public_id => 863);
select catalog.member('konosuba', 'konosuba-explosion-manga');
select catalog.relate('konosuba-explosion-manga', 'adaptation_of', 'konosuba-explosion-ln');
update catalog.entry set status = 'finished', end_date = '2017-12-27', end_precision = 'day' where id = catalog.id('konosuba-explosion-manga');
select catalog.add_units('konosuba-explosion-manga', 5);
select catalog.add_entry('konosuba-explosion-bonus-manga', 'manga', 'serial', 'KonoSuba: Even More Explosions on This Wonderful World!', '2018-02-27', 'JP', 'volume', 4, p_public_id => 864);
select catalog.member('konosuba', 'konosuba-explosion-bonus-manga');
select catalog.relate('konosuba-explosion-bonus-manga', 'adaptation_of', 'konosuba-explosion-bonus-ln');
update catalog.entry set status = 'finished', end_date = '2020-05-27', end_precision = 'day' where id = catalog.id('konosuba-explosion-bonus-manga');
select catalog.add_units('konosuba-explosion-bonus-manga', 4);
select catalog.add_entry('konosuba-dust-manga', 'manga', 'serial', 'KonoSuba: Extra Attention to That Wonderful Fool!', '2018-01-26', 'JP', 'volume', 2, p_public_id => 865);
select catalog.member('konosuba', 'konosuba-dust-manga');
select catalog.relate('konosuba-dust-manga', 'adaptation_of', 'konosuba-dust-ln');
update catalog.entry set status = 'finished', end_date = '2020-01-24', end_precision = 'day' where id = catalog.id('konosuba-dust-manga');
select catalog.add_units('konosuba-dust-manga', 2);
select catalog.add_entry('konosuba-masked-devil-manga', 'manga', 'serial', 'KonoSuba: Consulting with This Masked Devil!', '2018-02-27', 'JP', 'volume', 2, p_public_id => 866);
select catalog.member('konosuba', 'konosuba-masked-devil-manga');
select catalog.relate('konosuba-masked-devil-manga', 'adaptation_of', 'konosuba-masked-devil-ln');
update catalog.entry set status = 'finished', end_date = '2019-01-26', end_precision = 'day' where id = catalog.id('konosuba-masked-devil-manga');
select catalog.add_units('konosuba-masked-devil-manga', 2);
select catalog.add_entry('konosuba-ova-1', 'anime', 'ova', 'KonoSuba: God''s Blessings on This Wonderful Choker!', '2016-06', 'JP', 'episode', 1, p_public_id => 867);
select catalog.member('konosuba', 'konosuba-ova-1');
select catalog.relate('konosuba-ova-1', 'side_story_of', 'konosuba-s1');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-ova-1');
select catalog.add_entry('konosuba-ova-2', 'anime', 'ova', 'KonoSuba 2: God''s Blessings on These Wonderful Works of Art!', '2017-07-24', 'JP', 'episode', 1, p_public_id => 868);
select catalog.member('konosuba', 'konosuba-ova-2');
select catalog.relate('konosuba-ova-2', 'side_story_of', 'konosuba-s2');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-ova-2');
select catalog.add_entry('konosuba-3-bonus-stage', 'anime', 'ova', 'KonoSuba: God''s Blessing on This Wonderful World! 3 - Bonus Stage', '2025-03-14', 'JP', 'episode', 2, p_public_id => 869);
select catalog.member('konosuba', 'konosuba-3-bonus-stage');
select catalog.relate('konosuba-3-bonus-stage', 'side_story_of', 'konosuba-s3', p_note => 'Two bonus episodes premiered theatrically before the April 25 home-video release.');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-3-bonus-stage');
select catalog.add_entry('konosuba-explosion-anime', 'anime', 'tv', 'KonoSuba: An Explosion on This Wonderful World!', '2023-04-06', 'JP', 'episode', 12, p_public_id => 870);
select catalog.member('konosuba', 'konosuba-explosion-anime');
select catalog.relate('konosuba-explosion-anime', 'adaptation_of', 'konosuba-explosion-ln');
update catalog.entry set status = 'finished', end_date = '2023-06-22', end_precision = 'day' where id = catalog.id('konosuba-explosion-anime');
select catalog.add_entry('konosuba-s4', 'anime', 'tv', 'KonoSuba: God''s Blessing on This Wonderful World! 4', '2027', 'JP', 'episode', null, p_public_id => 871);
select catalog.member('konosuba', 'konosuba-s4');
select catalog.relate('konosuba-s4', 'sequel_of', 'konosuba-s3');
update catalog.entry set status = 'announced' where id = catalog.id('konosuba-s4');
select catalog.add_entry('konosuba-in-the-life', 'game', 'full', 'KonoSuba: In the Life', '2016-03-25', 'JP', 'none', null, p_public_id => 872);
select catalog.member('konosuba', 'konosuba-in-the-life');
select catalog.relate('konosuba-in-the-life', 'side_story_of', 'konosuba-s1');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-in-the-life');
select catalog.add_entry('konosuba-revival-of-beldia', 'game', 'full', 'KonoSuba: Revival of Beldia', '2017-04-28', 'JP', 'none', null, p_public_id => 873);
select catalog.member('konosuba', 'konosuba-revival-of-beldia');
select catalog.relate('konosuba-revival-of-beldia', 'side_story_of', 'konosuba-s2');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-revival-of-beldia');
select catalog.add_entry('konosuba-judgment-greedy-game', 'visual_novel', 'full', 'KonoSuba: Judgment on This Greedy Game!', '2017-09-07', 'JP', 'none', null, p_public_id => 874);
select catalog.member('konosuba', 'konosuba-judgment-greedy-game');
select catalog.relate('konosuba-judgment-greedy-game', 'side_story_of', 'konosuba-ln');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-judgment-greedy-game');
select catalog.add_entry('konosuba-love-clothes-desire', 'visual_novel', 'full', 'KonoSuba: Love for These Clothes of Desire!', '2020-09-24', 'JP', 'none', null, p_public_id => 875);
select catalog.member('konosuba', 'konosuba-love-clothes-desire');
select catalog.relate('konosuba-love-clothes-desire', 'side_story_of', 'konosuba-ln');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-love-clothes-desire');
select catalog.add_entry('konosuba-labyrinth-hope', 'game', 'full', 'KonoSuba: Labyrinth of Hope and the Gathering of Adventurers!', '2019-06-27', 'JP', 'none', null, p_public_id => 876);
select catalog.member('konosuba', 'konosuba-labyrinth-hope');
select catalog.relate('konosuba-labyrinth-hope', 'side_story_of', 'konosuba-ln');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-labyrinth-hope');
select catalog.add_entry('konosuba-labyrinth-hope-plus', 'game', 'full', 'KonoSuba: Labyrinth of Hope and the Gathering of Adventurers! Plus', '2020-08-27', 'JP', 'none', null, p_public_id => 877);
select catalog.member('konosuba', 'konosuba-labyrinth-hope-plus');
select catalog.relate('konosuba-labyrinth-hope-plus', 'remake_of', 'konosuba-labyrinth-hope');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-labyrinth-hope-plus');
select catalog.add_entry('konosuba-cursed-relic', 'game', 'full', 'KonoSuba: Cursed Relic and the Perplexed Adventurers', '2022', 'JP', 'none', null, p_public_id => 878);
select catalog.member('konosuba', 'konosuba-cursed-relic');
select catalog.relate('konosuba-cursed-relic', 'side_story_of', 'konosuba-ln');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-cursed-relic');
select catalog.add_entry('konosuba-fantastic-days', 'game', 'full', 'KonoSuba: Fantastic Days', '2020-02-27', 'JP', 'none', null, p_public_id => 879);
select catalog.member('konosuba', 'konosuba-fantastic-days');
select catalog.relate('konosuba-fantastic-days', 'side_story_of', 'konosuba-ln');
update catalog.entry set status = 'finished', end_date = '2025-01-30', end_precision = 'day' where id = catalog.id('konosuba-fantastic-days');
select catalog.add_entry('konosuba-fantastic-days-ln', 'book', 'light_novel', 'KonoSuba: Fantastic Days', '2022-03-01', 'JP', 'volume', 1, p_public_id => 880);
select catalog.member('konosuba', 'konosuba-fantastic-days-ln');
select catalog.relate('konosuba-fantastic-days-ln', 'adaptation_of', 'konosuba-fantastic-days');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-fantastic-days-ln');
select catalog.add_units('konosuba-fantastic-days-ln', 1);
select catalog.add_entry('konosuba-fantastic-days-manga', 'manga', 'serial', 'KonoSuba: Fantastic Days', '2022-03-26', 'JP', 'volume', 4, p_public_id => 881);
select catalog.member('konosuba', 'konosuba-fantastic-days-manga');
select catalog.relate('konosuba-fantastic-days-manga', 'adaptation_of', 'konosuba-fantastic-days');
update catalog.entry set status = 'finished', end_date = '2024-05-27', end_precision = 'day' where id = catalog.id('konosuba-fantastic-days-manga');
select catalog.add_units('konosuba-fantastic-days-manga', 4);
select catalog.add_entry('konosuba-drama-cd', 'audio', 'drama', 'KonoSuba: God''s Blessing on This Wonderful World! Drama CD', '2015-03-01', 'JP', 'none', null, p_public_id => 882);
select catalog.member('konosuba', 'konosuba-drama-cd');
select catalog.relate('konosuba-drama-cd', 'side_story_of', 'konosuba-ln');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-drama-cd');
select catalog.add_entry('konosuba-drama-cd-2', 'audio', 'drama', 'KonoSuba: God''s Blessing on This Wonderful World! 2 Drama CD', '2017-03', 'JP', 'none', null, p_public_id => 883);
select catalog.member('konosuba', 'konosuba-drama-cd-2');
select catalog.relate('konosuba-drama-cd-2', 'side_story_of', 'konosuba-s2');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-drama-cd-2');
select catalog.add_entry('isekai-quartet-s1', 'anime', 'tv', 'Isekai Quartet', '2019-04-09', 'JP', 'episode', 12, p_public_id => 884);
select catalog.member('konosuba', 'isekai-quartet-s1');
select catalog.relate('isekai-quartet-s1', 'crossover_with', 'konosuba-s1');
update catalog.entry set status = 'finished', end_date = '2019-06-25', end_precision = 'day' where id = catalog.id('isekai-quartet-s1');
select catalog.add_entry('isekai-quartet-s2', 'anime', 'tv', 'Isekai Quartet 2', '2020-01-14', 'JP', 'episode', 12, p_public_id => 885);
select catalog.member('konosuba', 'isekai-quartet-s2');
select catalog.relate('isekai-quartet-s2', 'sequel_of', 'isekai-quartet-s1');
update catalog.entry set status = 'finished', end_date = '2020-03-31', end_precision = 'day' where id = catalog.id('isekai-quartet-s2');
select catalog.add_entry('isekai-quartet-movie', 'anime', 'movie', 'Isekai Quartet: The Movie - Another World', '2022-06-10', 'JP', 'none', null, p_public_id => 886);
select catalog.member('konosuba', 'isekai-quartet-movie');
select catalog.relate('isekai-quartet-movie', 'sequel_of', 'isekai-quartet-s2');
update catalog.entry set status = 'finished' where id = catalog.id('isekai-quartet-movie');
select catalog.add_entry('isekai-quartet-s3', 'anime', 'tv', 'Isekai Quartet 3', '2025-10-13', 'JP', 'episode', 11, p_public_id => 887);
select catalog.member('konosuba', 'isekai-quartet-s3');
select catalog.relate('isekai-quartet-s3', 'sequel_of', 'isekai-quartet-movie');
update catalog.entry set status = 'finished', end_date = '2025-12-22', end_precision = 'day' where id = catalog.id('isekai-quartet-s3');

select catalog.relate('konosuba-s1', 'adaptation_of', 'konosuba-ln');
select catalog.relate('konosuba-s2', 'adaptation_of', 'konosuba-ln');
select catalog.relate('konosuba-s3', 'adaptation_of', 'konosuba-ln');
select catalog.relate('konosuba-legend-of-crimson', 'adaptation_of', 'konosuba-ln');
select catalog.credit('konosuba-explosion-anime', 'Drive', 'studio', 'JP');
select catalog.credit('konosuba-3-bonus-stage', 'Drive', 'studio', 'JP');
select catalog.credit('konosuba-s4', 'ENGI', 'studio', 'JP');

select catalog.add_ordering('konosuba', 'recommended', 'first-timer order',
  'The main TV story, with optional bonus episodes and the Megumin prequel. Legend of Crimson bridges seasons 2 and 3.') as o \gset
select catalog.add_step(:'o', 1, 'konosuba-s1');
select catalog.add_step(:'o', 2, 'konosuba-ova-1', p_necessity => 'optional');
select catalog.add_step(:'o', 3, 'konosuba-s2');
select catalog.add_step(:'o', 4, 'konosuba-ova-2', p_necessity => 'optional');
select catalog.add_step(:'o', 5, 'konosuba-explosion-anime', p_necessity => 'optional', p_note => 'A prequel centered on Megumin.');
select catalog.add_step(:'o', 6, 'konosuba-legend-of-crimson');
select catalog.add_step(:'o', 7, 'konosuba-s3');
select catalog.add_step(:'o', 8, 'konosuba-3-bonus-stage', p_necessity => 'optional');

-- https://konosuba-machisuba.com/
-- https://www.entergram.co.jp/konosuba/product/
select catalog.add_entry('konosuba-prosperity-beloved-town', 'game', 'full', 'KonoSuba: Prosperity to This Beloved Town!', '2027-02', 'JP', 'none', null, p_public_id => 888);
select catalog.member('konosuba', 'konosuba-prosperity-beloved-town');
select catalog.relate('konosuba-prosperity-beloved-town', 'side_story_of', 'konosuba-ln');
update catalog.entry set status = 'announced', description = 'An announced mobile and PC role-playing game with an original story supervised by Natsume Akatsuki.' where id = catalog.id('konosuba-prosperity-beloved-town');
select catalog.credit('konosuba-prosperity-beloved-town', 'Kadokawa', 'publisher', 'JP');
select catalog.add_entry('konosuba-education-useless-goddess', 'game', 'full', 'KonoSuba: Education for This Useless Goddess!', '2019-06-27', 'JP', 'none', null, p_public_id => 889);
select catalog.member('konosuba', 'konosuba-education-useless-goddess');
select catalog.relate('konosuba-education-useless-goddess', 'spin_off_of', 'konosuba-labyrinth-hope');
update catalog.entry set status = 'finished', description = 'A separate Aqua-raising bonus game bundled with Labyrinth of Hope and the Gathering of Adventurers.' where id = catalog.id('konosuba-education-useless-goddess');
select catalog.credit('konosuba-education-useless-goddess', 'Entergram', 'publisher', 'JP');

-- https://game.mages.co.jp/konosuba/shinpan/product/
-- https://game.mages.co.jp/konosuba/product/
select catalog.add_entry('konosuba-attack-of-the-destroyer', 'game', 'full', 'KonoSuba: Attack of the Destroyer', '2017-09-07', 'JP', 'none', null, p_public_id => 890);
select catalog.member('konosuba', 'konosuba-attack-of-the-destroyer');
select catalog.relate('konosuba-attack-of-the-destroyer', 'side_story_of', 'konosuba-judgment-greedy-game', p_note => 'Bonus shooting game distributed with first-print copies of Judgment on This Greedy Game!');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-attack-of-the-destroyer');
select catalog.add_entry('konosuba-clothes-bonus-shooter', 'game', 'full', 'KonoSuba: Love for These Clothes of Desire! - Bonus Shooting Game', '2020-09-24', 'JP', 'none', null, p_public_id => 891);
select catalog.member('konosuba', 'konosuba-clothes-bonus-shooter');
select catalog.relate('konosuba-clothes-bonus-shooter', 'side_story_of', 'konosuba-love-clothes-desire', p_note => 'Six-stage bonus shooting game distributed with first-print copies; the product page names it by the parent title.');
update catalog.entry set status = 'finished' where id = catalog.id('konosuba-clothes-bonus-shooter');

-- https://store.kadokawa.co.jp/shop/g/g321611000533/
-- https://natalie.mu/comic/news/216801
-- https://www.kadokawa.co.jp/product/321608000309/
-- https://store.kadokawa.co.jp/shop/g/g321605000145/
select catalog.add_entry('konosuba-everyday-life', 'manga', 'serial', 'KonoSuba: Everyday Life in This Wonderful World!', '2017', 'JP', 'none', null, p_public_id => 892);
select catalog.member('konosuba', 'konosuba-everyday-life');
select catalog.relate('konosuba-everyday-life', 'spin_off_of', 'konosuba-ln', p_note => 'The first collected edition was published in 2017; this is a comedy spinoff by Somemiya Suzume.');
update catalog.entry set status = null where id = catalog.id('konosuba-everyday-life');
select catalog.credit('konosuba-everyday-life', 'Kadokawa', 'publisher', 'JP');
select catalog.add_entry('konosuba-kappore', 'manga', 'serial', 'KonoSuba: Kappore!', '2016', 'JP', 'volume', null, p_public_id => 893);
select catalog.member('konosuba', 'konosuba-kappore');
select catalog.relate('konosuba-kappore', 'spin_off_of', 'konosuba-ln');
update catalog.entry set status = null where id = catalog.id('konosuba-kappore');
select catalog.credit('konosuba-kappore', 'Kadokawa', 'publisher', 'JP');
select catalog.add_entry('konosuba-megumin-anthology', 'manga', 'serial', 'KonoSuba: Megumin Anthology', '2016-12-23', 'JP', 'volume', null, p_public_id => 894);
select catalog.member('konosuba', 'konosuba-megumin-anthology');
select catalog.relate('konosuba-megumin-anthology', 'spin_off_of', 'konosuba-ln');
update catalog.entry set status = null where id = catalog.id('konosuba-megumin-anthology');
select catalog.credit('konosuba-megumin-anthology', 'Kadokawa', 'publisher', 'JP');
select catalog.add_entry('konosuba-comic-anthology', 'manga', 'serial', 'KonoSuba: Comic Anthology', '2017-01-07', 'JP', 'volume', null, p_public_id => 895);
select catalog.member('konosuba', 'konosuba-comic-anthology');
select catalog.relate('konosuba-comic-anthology', 'spin_off_of', 'konosuba-ln');
update catalog.entry set status = null where id = catalog.id('konosuba-comic-anthology');
select catalog.credit('konosuba-comic-anthology', 'Kadokawa', 'publisher', 'JP');
