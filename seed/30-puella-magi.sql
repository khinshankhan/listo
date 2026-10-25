-- Madoka: the recap films retell the series, so the first stop is a choice.

select catalog.add_group('madoka', 'series', 'Puella Magi Madoka Magica', p_public_id => 900);
select catalog.alias('madoka', 'romanized', 'Mahou Shoujo Madoka Magica');
select catalog.add_entry('madoka-tv', 'anime', 'tv', 'Puella Magi Madoka Magica', '2011-01', 'JP', 'episode', 12, p_public_id => 901);
select catalog.add_units('madoka-tv', 12);
select catalog.add_entry('madoka-beginnings', 'anime', 'movie', 'Puella Magi Madoka Magica the Movie Part I: Beginnings', '2012-10', 'JP', p_public_id => 902);
select catalog.add_entry('madoka-eternal', 'anime', 'movie', 'Puella Magi Madoka Magica the Movie Part II: Eternal', '2012-10', 'JP', p_public_id => 903);
select catalog.add_entry('madoka-rebellion', 'anime', 'movie', 'Puella Magi Madoka Magica the Movie Part III: Rebellion', '2013-10', 'JP', p_public_id => 904);
select catalog.member('madoka', 'madoka-tv', 'madoka-beginnings', 'madoka-eternal', 'madoka-rebellion');

select catalog.ext('madoka-tv', 'wikidata', 'Q53353');
select catalog.ext('madoka-tv', 'anilist_anime', '9756');
select catalog.ext('madoka-tv', 'mal_anime', '9756');
select catalog.ext('madoka-beginnings', 'wikidata', 'Q3842376');
select catalog.ext('madoka-eternal', 'wikidata', 'Q3842375');
select catalog.ext('madoka-rebellion', 'wikidata', 'Q17479150');
select catalog.ext('madoka-rebellion', 'mal_anime', '11981');

select catalog.relate('madoka-beginnings', 'recap_of', 'madoka-tv', p_note => 'Retells episodes 1-8.');
select catalog.relate('madoka-eternal', 'recap_of', 'madoka-tv', p_note => 'Retells episodes 9-12.');
select catalog.relate('madoka-eternal', 'sequel_of', 'madoka-beginnings');
select catalog.relate('madoka-rebellion', 'sequel_of', 'madoka-tv');
select catalog.relate('madoka-rebellion', 'sequel_of', 'madoka-eternal');

-- https://en.wikipedia.org/wiki/Puella_Magi_Madoka_Magica
-- https://wiki.puella-magi.net/Manga
-- https://en.wikipedia.org/wiki/Magia_Record
-- https://www.madoka-magica.com/wr/
-- https://wiki.puella-magi.net/Null_Magical_Girl
-- https://wiki.puella-magi.net/Madoka_Magica_Online
-- https://wiki.puella-magi.net/Madoka_Mobage
select catalog.add_entry('madoka-walpurgisnacht-rising', 'anime', 'movie', 'Puella Magi Madoka Magica the Movie: Walpurgisnacht - Rising', '2026-08-28', 'JP', 'none', null, p_public_id => 905);
select catalog.member('madoka', 'madoka-walpurgisnacht-rising');
select catalog.relate('madoka-walpurgisnacht-rising', 'sequel_of', 'madoka-rebellion');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-walpurgisnacht-rising');
select catalog.add_entry('madoka-concept-movie', 'anime', 'special', 'Puella Magi Madoka Magica: Concept Movie', '2015-11-27', 'JP', 'none', null, p_public_id => 906);
select catalog.member('madoka', 'madoka-concept-movie');
select catalog.relate('madoka-concept-movie', 'references', 'madoka-rebellion');
update catalog.entry set status = 'finished', description = 'A short concept film exhibited at Madogatari, presenting images and ideas for a future Madoka project.', length_minutes = 4 where id = catalog.id('madoka-concept-movie');
select catalog.add_entry('madoka-manga', 'manga', 'serial', 'Puella Magi Madoka Magica', '2011-02-12', 'JP', 'volume', 3, p_public_id => 907);
select catalog.member('madoka', 'madoka-manga');
select catalog.relate('madoka-manga', 'adaptation_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2011-05-30', end_precision = 'day' where id = catalog.id('madoka-manga');
select catalog.add_units('madoka-manga', 3);
select catalog.add_entry('madoka-manga-new-complete-edition', 'manga', 'serial', 'Puella Magi Madoka Magica: New Complete Edition', '2021', 'JP', 'volume', null, p_public_id => 908);
select catalog.member('madoka', 'madoka-manga-new-complete-edition');
select catalog.relate('madoka-manga-new-complete-edition', 'remake_of', 'madoka-manga');
update catalog.entry set status = 'finished', description = 'Hanokage''s redrawn adaptation brings the original manga closer to the completed television production.' where id = catalog.id('madoka-manga-new-complete-edition');
select catalog.add_entry('madoka-rebellion-manga', 'manga', 'serial', 'Puella Magi Madoka Magica: Rebellion', '2013-11', 'JP', 'volume', 3, p_public_id => 909);
select catalog.member('madoka', 'madoka-rebellion-manga');
select catalog.relate('madoka-rebellion-manga', 'adaptation_of', 'madoka-rebellion');
update catalog.entry set status = 'finished', end_date = '2014-01-01', end_precision = 'month' where id = catalog.id('madoka-rebellion-manga');
select catalog.add_units('madoka-rebellion-manga', 3);
select catalog.add_entry('madoka-different-story', 'manga', 'serial', 'Puella Magi Madoka Magica: The Different Story', '2012-10-12', 'JP', 'volume', 3, p_public_id => 910);
select catalog.member('madoka', 'madoka-different-story');
select catalog.relate('madoka-different-story', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2012-11-12', end_precision = 'day' where id = catalog.id('madoka-different-story');
select catalog.add_units('madoka-different-story', 3);
select catalog.add_entry('madoka-wraith-arc', 'manga', 'serial', 'Puella Magi Madoka Magica: Wraith Arc', '2015-06-10', 'JP', 'volume', 3, p_public_id => 911);
select catalog.member('madoka', 'madoka-wraith-arc');
select catalog.relate('madoka-wraith-arc', 'sequel_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2016-10-08', end_precision = 'day' where id = catalog.id('madoka-wraith-arc');
select catalog.add_units('madoka-wraith-arc', 3);
select catalog.add_entry('madoka-homuras-revenge', 'manga', 'serial', 'Puella Magi Madoka Magica: Homura''s Revenge!', '2012-06-08', 'JP', 'volume', 2, p_public_id => 912);
select catalog.member('madoka', 'madoka-homuras-revenge');
select catalog.relate('madoka-homuras-revenge', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-homuras-revenge');
select catalog.add_units('madoka-homuras-revenge', 2);
select catalog.add_entry('madoka-homura-tamura', 'manga', 'serial', 'Puella Magi Homura Tamura', '2012-06-08', 'JP', 'volume', 3, p_public_id => 913);
select catalog.member('madoka', 'madoka-homura-tamura');
select catalog.relate('madoka-homura-tamura', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-homura-tamura');
select catalog.add_units('madoka-homura-tamura', 3);
select catalog.add_entry('madoka-kazumi-magica', 'manga', 'serial', 'Puella Magi Kazumi Magica: The Innocent Malice', '2011-01-24', 'JP', 'volume', 5, p_public_id => 914);
select catalog.member('madoka', 'madoka-kazumi-magica');
select catalog.relate('madoka-kazumi-magica', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2012-11-24', end_precision = 'day' where id = catalog.id('madoka-kazumi-magica');
select catalog.add_units('madoka-kazumi-magica', 5);
select catalog.add_entry('madoka-oriko-magica', 'manga', 'serial', 'Puella Magi Oriko Magica', '2011-05-12', 'JP', 'volume', 2, p_public_id => 915);
select catalog.member('madoka', 'madoka-oriko-magica');
select catalog.relate('madoka-oriko-magica', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2011-06-12', end_precision = 'day' where id = catalog.id('madoka-oriko-magica');
select catalog.add_units('madoka-oriko-magica', 2);
select catalog.add_entry('madoka-oriko-extra-story', 'manga', 'serial', 'Puella Magi Oriko Magica: Extra Story', '2012-08-09', 'JP', 'volume', 1, p_public_id => 916);
select catalog.member('madoka', 'madoka-oriko-extra-story');
select catalog.relate('madoka-oriko-extra-story', 'side_story_of', 'madoka-oriko-magica');
update catalog.entry set status = 'finished', end_date = '2013-06-10', end_precision = 'day' where id = catalog.id('madoka-oriko-extra-story');
select catalog.add_units('madoka-oriko-extra-story', 1);
select catalog.add_entry('madoka-oriko-sadness-prayer', 'manga', 'serial', 'Puella Magi Oriko Magica: Sadness Prayer', '2013-11-08', 'JP', 'volume', 4, p_public_id => 917);
select catalog.member('madoka', 'madoka-oriko-sadness-prayer');
select catalog.relate('madoka-oriko-sadness-prayer', 'side_story_of', 'madoka-oriko-magica');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-oriko-sadness-prayer');
select catalog.add_units('madoka-oriko-sadness-prayer', 4);
select catalog.add_entry('madoka-suzune-magica', 'manga', 'serial', 'Puella Magi Suzune Magica', '2013-11-12', 'JP', 'volume', 3, p_public_id => 918);
select catalog.member('madoka', 'madoka-suzune-magica');
select catalog.relate('madoka-suzune-magica', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2014-11-22', end_precision = 'day' where id = catalog.id('madoka-suzune-magica');
select catalog.add_units('madoka-suzune-magica', 3);
select catalog.add_entry('madoka-tart-magica', 'manga', 'serial', 'Puella Magi Tart Magica: The Legend of Jeanne d''Arc', '2013-11-08', 'JP', 'volume', 5, p_public_id => 919);
select catalog.member('madoka', 'madoka-tart-magica');
select catalog.relate('madoka-tart-magica', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-tart-magica');
select catalog.add_units('madoka-tart-magica', 5);
select catalog.add_entry('madoka-cafe-grief-seed', 'manga', 'serial', 'Welcome to Cafe Grief Seed!', '2012-06-08', 'JP', 'volume', null, p_public_id => 920);
select catalog.member('madoka', 'madoka-cafe-grief-seed');
select catalog.relate('madoka-cafe-grief-seed', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-cafe-grief-seed');
select catalog.add_entry('madoka-sute-ma', 'manga', 'serial', 'Sute Ma!', '2012-06-08', 'JP', 'volume', null, p_public_id => 921);
select catalog.member('madoka', 'madoka-sute-ma');
select catalog.relate('madoka-sute-ma', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-sute-ma');
select catalog.add_entry('madoka-magical-girl-club', 'manga', 'serial', 'Mahou Shoujobu Madoka Magica', '2012-06-08', 'JP', 'volume', null, p_public_id => 922);
select catalog.member('madoka', 'madoka-magical-girl-club');
select catalog.relate('madoka-magical-girl-club', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-magical-girl-club');
select catalog.add_entry('madoka-pomu-magi', 'manga', 'serial', 'Pomu Magi', '2012-06-08', 'JP', 'volume', null, p_public_id => 923);
select catalog.member('madoka', 'madoka-pomu-magi');
select catalog.relate('madoka-pomu-magi', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-pomu-magi');
select catalog.add_entry('madoka-veranda', 'manga', 'serial', 'The Veranda of Madoka', '2012-06-08', 'JP', 'volume', null, p_public_id => 924);
select catalog.member('madoka', 'madoka-veranda');
select catalog.relate('madoka-veranda', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-veranda');
select catalog.add_entry('madoka-mitakihara-kindergarten', 'manga', 'serial', 'Mitakihara Kindergarten Mahougumi', '2012-06-08', 'JP', 'volume', null, p_public_id => 925);
select catalog.member('madoka', 'madoka-mitakihara-kindergarten');
select catalog.relate('madoka-mitakihara-kindergarten', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-mitakihara-kindergarten');
select catalog.add_entry('madoka-mami-everyday-life', 'manga', 'serial', 'Mami Tomoe''s Everyday Life', '2012-06-08', 'JP', 'volume', null, p_public_id => 926);
select catalog.member('madoka', 'madoka-mami-everyday-life');
select catalog.relate('madoka-mami-everyday-life', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = null where id = catalog.id('madoka-mami-everyday-life');
select catalog.add_entry('madoka-mitakihara-anti-materials', 'manga', 'serial', 'Mitakihara Anti-Materials', '2012-06-08', 'JP', 'volume', null, p_public_id => 927);
select catalog.member('madoka', 'madoka-mitakihara-anti-materials');
select catalog.relate('madoka-mitakihara-anti-materials', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-mitakihara-anti-materials');
select catalog.add_entry('madoka-floating-madoka', 'manga', 'serial', 'Floating Madoka', '2012-12-10', 'JP', 'volume', null, p_public_id => 928);
select catalog.member('madoka', 'madoka-floating-madoka');
select catalog.relate('madoka-floating-madoka', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-floating-madoka');
select catalog.add_entry('madoka-kanameke', 'manga', 'serial', 'Kanameke', '2013-10-10', 'JP', 'volume', null, p_public_id => 929);
select catalog.member('madoka', 'madoka-kanameke');
select catalog.relate('madoka-kanameke', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-kanameke');
select catalog.add_entry('madoka-comic-anthology', 'manga', 'serial', 'Puella Magi Madoka Magica: Comic Anthology', '2011-09-12', 'JP', 'volume', null, p_public_id => 930);
select catalog.member('madoka', 'madoka-comic-anthology');
select catalog.relate('madoka-comic-anthology', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2014-06-01', end_precision = 'month' where id = catalog.id('madoka-comic-anthology');
select catalog.add_entry('madoka-four-panel-anthology', 'manga', 'serial', 'Puella Magi Madoka Magica: 4-Panel Anthology', '2012-08', 'JP', 'volume', null, p_public_id => 931);
select catalog.member('madoka', 'madoka-four-panel-anthology');
select catalog.relate('madoka-four-panel-anthology', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2013-10-01', end_precision = 'month' where id = catalog.id('madoka-four-panel-anthology');
select catalog.add_entry('madoka-film-memories', 'manga', 'serial', 'Puella Magi Madoka Magica: Film Memories', '2012-05-26', 'JP', 'volume', null, p_public_id => 932);
select catalog.member('madoka', 'madoka-film-memories');
select catalog.relate('madoka-film-memories', 'adaptation_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-film-memories');
select catalog.add_entry('madoka-novel', 'book', 'light_novel', 'Puella Magi Madoka Magica', '2011-08-12', 'JP', 'volume', 2, p_public_id => 933);
select catalog.member('madoka', 'madoka-novel');
select catalog.relate('madoka-novel', 'adaptation_of', 'madoka-tv', p_note => 'First sold at Comiket before the August 14 general release.');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-novel');
select catalog.add_units('madoka-novel', 2);
select catalog.add_entry('madoka-beginning-story', 'book', 'nonfiction', 'Puella Magi Madoka Magica: The Beginning Story', '2011-11', 'JP', 'none', null, p_public_id => 934);
select catalog.member('madoka', 'madoka-beginning-story');
select catalog.relate('madoka-beginning-story', 'references', 'madoka-tv');
update catalog.entry set status = 'finished', description = 'Gen Urobuchi''s original scenario and production material for the television series.' where id = catalog.id('madoka-beginning-story');
select catalog.add_entry('madoka-null-magical-girl', 'book', 'light_novel', 'Null Magical Girl', '2021-06-25', 'JP', 'chapter', 3, p_public_id => 935);
select catalog.member('madoka', 'madoka-null-magical-girl');
select catalog.relate('madoka-null-magical-girl', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished', description = 'GenGen Kusano''s anniversary-book spinoff follows a graduate student whose encounter with Kyubey takes a science-fiction turn.', end_date = '2021-12-25', end_precision = 'day' where id = catalog.id('madoka-null-magical-girl');
select catalog.add_entry('madoka-portable', 'visual_novel', 'full', 'Puella Magi Madoka Magica Portable', '2012-03-15', 'JP', 'none', null, p_public_id => 936);
select catalog.member('madoka', 'madoka-portable');
select catalog.relate('madoka-portable', 'adaptation_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-portable');
select catalog.add_entry('madoka-battle-pentagram', 'game', 'full', 'Puella Magi Madoka Magica: The Battle Pentagram', '2013-12-19', 'JP', 'none', null, p_public_id => 937);
select catalog.member('madoka', 'madoka-battle-pentagram');
select catalog.relate('madoka-battle-pentagram', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-battle-pentagram');
select catalog.add_entry('madoka-mobage', 'game', 'full', 'Puella Magi Madoka Magica Mobage', '2011-09-16', 'JP', 'none', null, p_public_id => 938);
select catalog.member('madoka', 'madoka-mobage');
select catalog.relate('madoka-mobage', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2012-11-19', end_precision = 'day' where id = catalog.id('madoka-mobage');
select catalog.add_entry('madoka-online', 'game', 'full', 'Puella Magi Madoka Magica Online', '2012-09-04', 'JP', 'none', null, p_public_id => 939);
select catalog.member('madoka', 'madoka-online');
select catalog.relate('madoka-online', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2015-05-28', end_precision = 'day' where id = catalog.id('madoka-online');
select catalog.add_entry('madoka-mami-tiro-finale', 'game', 'full', 'Mami''s Heart Pounding Tiro Finale', '2011-10-14', 'JP', 'none', null, p_public_id => 940);
select catalog.member('madoka', 'madoka-mami-tiro-finale');
select catalog.relate('madoka-mami-tiro-finale', 'adaptation_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-mami-tiro-finale');
select catalog.add_entry('madoka-tps-homura', 'game', 'full', 'Puella Magi Madoka Magica TPS Featuring Homura Akemi', '2011-12', 'JP', 'none', null, p_public_id => 941);
select catalog.member('madoka', 'madoka-tps-homura');
select catalog.relate('madoka-tps-homura', 'adaptation_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-tps-homura');
select catalog.add_entry('madoka-tps-mami', 'game', 'full', 'Puella Magi Madoka Magica TPS Featuring Mami Tomoe', '2012-08', 'JP', 'none', null, p_public_id => 942);
select catalog.member('madoka', 'madoka-tps-mami');
select catalog.relate('madoka-tps-mami', 'adaptation_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-tps-mami');
select catalog.add_entry('madoka-tps-sayaka-kyoko', 'game', 'full', 'Puella Magi Madoka Magica TPS Featuring Sayaka Miki and Kyoko Sakura', '2012-10-16', 'JP', 'none', null, p_public_id => 943);
select catalog.member('madoka', 'madoka-tps-sayaka-kyoko');
select catalog.relate('madoka-tps-sayaka-kyoko', 'adaptation_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-tps-sayaka-kyoko');
select catalog.add_entry('madoka-puzzle-of-memories', 'game', 'full', 'Puella Magi Madoka Magica: Puzzle of Memories', '2013-03-29', 'JP', 'none', null, p_public_id => 944);
select catalog.member('madoka', 'madoka-puzzle-of-memories');
select catalog.relate('madoka-puzzle-of-memories', 'adaptation_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-puzzle-of-memories');
select catalog.add_entry('magia-record-game', 'game', 'full', 'Magia Record: Puella Magi Madoka Magica Side Story', '2017-08-22', 'JP', 'none', null, p_public_id => 945);
select catalog.member('madoka', 'magia-record-game');
select catalog.relate('magia-record-game', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2024-07-31', end_precision = 'day' where id = catalog.id('magia-record-game');
select catalog.add_entry('magia-exedra', 'game', 'full', 'Puella Magi Madoka Magica: Magia Exedra', '2025-03-27', 'JP', 'none', null, p_public_id => 946);
select catalog.member('madoka', 'magia-exedra');
select catalog.relate('magia-exedra', 'spin_off_of', 'madoka-tv');
update catalog.entry set status = 'releasing' where id = catalog.id('magia-exedra');
select catalog.add_entry('magia-record-s1', 'anime', 'tv', 'Magia Record: Puella Magi Madoka Magica Side Story', '2020-01-04', 'JP', 'episode', 13, p_public_id => 947);
select catalog.member('madoka', 'magia-record-s1');
select catalog.relate('magia-record-s1', 'adaptation_of', 'magia-record-game');
update catalog.entry set status = 'finished', end_date = '2020-03-29', end_precision = 'day' where id = catalog.id('magia-record-s1');
select catalog.add_entry('magia-record-s2', 'anime', 'tv', 'Magia Record: The Eve of Awakening', '2021-08-01', 'JP', 'episode', 8, p_public_id => 948);
select catalog.member('madoka', 'magia-record-s2');
select catalog.relate('magia-record-s2', 'sequel_of', 'magia-record-s1');
update catalog.entry set status = 'finished', end_date = '2021-09-26', end_precision = 'day' where id = catalog.id('magia-record-s2');
select catalog.add_entry('magia-record-final-season', 'anime', 'tv', 'Magia Record: Dawn of a Shallow Dream', '2022-04-03', 'JP', 'episode', 4, p_public_id => 949);
select catalog.member('madoka', 'magia-record-final-season');
select catalog.relate('magia-record-final-season', 'sequel_of', 'magia-record-s2');
update catalog.entry set status = 'finished', end_date = '2022-04-03', end_precision = 'day' where id = catalog.id('magia-record-final-season');
select catalog.add_entry('magia-record-manga', 'manga', 'serial', 'Magia Record: Puella Magi Madoka Magica Side Story', '2018-08-24', 'JP', 'volume', null, p_public_id => 950);
select catalog.member('madoka', 'magia-record-manga');
select catalog.relate('magia-record-manga', 'adaptation_of', 'magia-record-game');
update catalog.entry set status = 'releasing' where id = catalog.id('magia-record-manga');
select catalog.add_entry('magia-record-another-story-manga', 'manga', 'serial', 'Magia Record: Another Story', '2019-05-25', 'JP', 'volume', null, p_public_id => 951);
select catalog.member('madoka', 'magia-record-another-story-manga');
select catalog.relate('magia-record-another-story-manga', 'adaptation_of', 'magia-record-game');
update catalog.entry set status = null where id = catalog.id('magia-record-another-story-manga');
select catalog.add_entry('magia-report', 'manga', 'serial', 'Magia Report', '2017', 'JP', 'volume', null, p_public_id => 952);
select catalog.member('madoka', 'magia-report');
select catalog.relate('magia-report', 'spin_off_of', 'magia-record-game');
update catalog.entry set status = 'releasing' where id = catalog.id('magia-report');
select catalog.add_entry('magia-record-stage', 'film', 'stage', 'Magia Record: Puella Magi Madoka Magica Side Story (stage adaptation)', '2019-02-27', 'JP', 'none', null, p_public_id => 953);
select catalog.member('madoka', 'magia-record-stage');
select catalog.relate('magia-record-stage', 'adaptation_of', 'magia-record-game');
update catalog.entry set status = 'finished', description = 'The released home-video recording of the 2018 stage production starring members of Hiragana Keyakizaka46.' where id = catalog.id('magia-record-stage');

select catalog.relate('madoka-rebellion-manga', 'sequel_of', 'madoka-wraith-arc');
select catalog.credit(s, 'Shaft', 'studio', 'JP') from unnest(array['madoka-walpurgisnacht-rising', 'madoka-concept-movie', 'magia-record-s1', 'magia-record-s2', 'magia-record-final-season']) s;
select catalog.credit(s, 'Houbunsha', 'publisher', 'JP') from catalog.node n join catalog.entry e on e.id = n.id,
  lateral (select n.slug) as x(s)
where e.medium = 'manga' and n.slug in
  (select member.slug from catalog.membership m join catalog.node member on member.id = m.member_id where m.group_id = catalog.id('madoka'));

-- https://wiki.puella-magi.net/Puella_Magi_Madoka_Magica_Plus
-- https://wiki.puella-magi.net/Puella_Magi_Madoka_Magica_Mobcast
-- https://wiki.puella-magi.net/Drama_CD_1
-- https://wiki.puella-magi.net/Drama_CD_2
-- https://wiki.puella-magi.net/Drama_CD_3
-- https://en.wikipedia.org/wiki/List_of_Puella_Magi_Madoka_Magica_episodes
-- https://wiki.puella-magi.net/Puella_Magi_Series
select catalog.add_entry('madoka-plus', 'game', 'full', 'Puella Magi Madoka Magica Plus', '2012-12-29', 'JP', 'none', null, p_public_id => 954);
select catalog.member('madoka', 'madoka-plus');
select catalog.relate('madoka-plus', 'side_story_of', 'madoka-tv');
update catalog.entry set status = 'finished', end_date = '2015-12-25', end_precision = 'day' where id = catalog.id('madoka-plus');
select catalog.add_entry('madoka-mobcast', 'game', 'full', 'Puella Magi Madoka Magica the Movie: Beginnings / Eternal (Mobcast)', '2016-12-28', 'CN', 'none', null, p_public_id => 955);
select catalog.member('madoka', 'madoka-mobcast');
select catalog.relate('madoka-mobcast', 'adaptation_of', 'madoka-beginnings');
update catalog.entry set status = 'finished', end_date = '2018-01-04', end_precision = 'day' where id = catalog.id('madoka-mobcast');
select catalog.add_entry('madoka-memories-of-you', 'audio', 'drama', 'Puella Magi Madoka Magica: Memories of You', '2011-04-27', 'JP', 'none', null, p_public_id => 956);
select catalog.member('madoka', 'madoka-memories-of-you');
select catalog.relate('madoka-memories-of-you', 'side_story_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-memories-of-you');
select catalog.add_entry('madoka-sunny-day-life', 'audio', 'drama', 'Puella Magi Madoka Magica: Sunny Day Life', '2011-06-22', 'JP', 'none', null, p_public_id => 957);
select catalog.member('madoka', 'madoka-sunny-day-life');
select catalog.relate('madoka-sunny-day-life', 'side_story_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-sunny-day-life');
select catalog.add_entry('madoka-farewell-story', 'audio', 'drama', 'Puella Magi Madoka Magica: Farewell Story', '2011-08-24', 'JP', 'none', null, p_public_id => 958);
select catalog.member('madoka', 'madoka-farewell-story');
select catalog.relate('madoka-farewell-story', 'side_story_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-farewell-story');
select catalog.add_entry('madoka-scene-zero', 'game', 'expansion', 'Puella Magi Madoka Magica: Scene 0', '2023-10-03', 'JP', 'none', null, p_public_id => 959);
select catalog.member('madoka', 'madoka-scene-zero');
select catalog.relate('madoka-scene-zero', 'side_story_of', 'madoka-tv');
update catalog.entry set status = 'finished', description = 'A prequel story first released within Magia Record, introducing Mabayu Aki. A later presentation was added to Magia Exedra. Not a separate standalone application.' where id = catalog.id('madoka-scene-zero');
select catalog.add_entry('magia-record-manga-anthology', 'manga', 'serial', 'Magia Record: Comic Anthology', '2017', 'JP', 'volume', null, p_public_id => 960);
select catalog.member('madoka', 'magia-record-manga-anthology');
select catalog.relate('magia-record-manga-anthology', 'spin_off_of', 'magia-record-game');
update catalog.entry set status = 'finished' where id = catalog.id('magia-record-manga-anthology');

select catalog.relate('madoka-different-story', 'adaptation_of', 'madoka-farewell-story', p_note => 'The opening volume adapts Farewell Story before the manga develops its own alternate timeline.');
select catalog.relate('madoka-scene-zero', 'side_story_of', 'magia-record-game');
select catalog.relate('madoka-scene-zero', 'references', 'magia-exedra', p_note => 'The story was later presented inside Magia Exedra as well.');

select catalog.add_ordering('madoka', 'recommended', 'first-timer order',
  'Choose the TV series or its two recap films, then continue with Rebellion and Walpurgisnacht - Rising. Magia Record is a separate spinoff.') as o \gset
select catalog.add_step(:'o', 1, 'madoka-tv', p_branch => 1);
select catalog.add_step(:'o', 1, 'madoka-beginnings', p_branch => 2, p_step => 1, p_note => 'Retells the opening eight episodes.');
select catalog.add_step(:'o', 1, 'madoka-eternal', p_branch => 2, p_step => 2);
select catalog.add_step(:'o', 2, 'madoka-rebellion');
select catalog.add_step(:'o', 3, 'madoka-walpurgisnacht-rising');

-- https://wiki.puella-magi.net/Drama_CD_4
select catalog.add_entry('madoka-summer-training-camp', 'audio', 'drama', 'Puella Magi Madoka Magica: Summer Magical Girl Training Camp!!', null, 'JP', 'none', null, p_public_id => 961);
select catalog.member('madoka', 'madoka-summer-training-camp');
select catalog.relate('madoka-summer-training-camp', 'side_story_of', 'madoka-tv');
update catalog.entry set status = 'finished' where id = catalog.id('madoka-summer-training-camp');

-- https://www.hanmoto.com/bd/isbn/9784832297517
-- https://www.hanmoto.com/bd/isbn/9784832297739
-- https://www.hanmoto.com/bd/isbn/9784832297111
-- https://www.hanmoto.com/bd/isbn/9784832297418
select catalog.add_entry('madoka-walpurgisnacht-rising-manga', 'manga', 'serial', 'Puella Magi Madoka Magica the Movie: Walpurgisnacht - Rising (Manga)', '2026-09-18', 'JP', 'volume', null, p_public_id => 962);
select catalog.member('madoka', 'madoka-walpurgisnacht-rising-manga');
select catalog.relate('madoka-walpurgisnacht-rising-manga', 'adaptation_of', 'madoka-walpurgisnacht-rising', p_note => 'The first collected volume is released; volume 2 is scheduled for December 23, 2026 and is not counted as available.');
update catalog.entry set status = 'releasing', unit_available = 1 where id = catalog.id('madoka-walpurgisnacht-rising-manga');
select catalog.add_units('madoka-walpurgisnacht-rising-manga', 1);
select catalog.credit('madoka-walpurgisnacht-rising-manga', 'Houbunsha', 'publisher', 'JP');
select catalog.add_entry('magia-etcetera', 'manga', 'serial', 'Magia Etcetera', '2026-04-03', 'JP', 'volume', null, p_public_id => 963);
select catalog.member('madoka', 'magia-etcetera');
select catalog.relate('magia-etcetera', 'spin_off_of', 'magia-exedra', p_note => 'Date is the first collected volume. PAPA''s webcomic succeeded Magia Report and began online earlier.');
update catalog.entry set status = null, unit_available = 1 where id = catalog.id('magia-etcetera');
select catalog.add_units('magia-etcetera', 1);
select catalog.credit('magia-etcetera', 'Houbunsha', 'publisher', 'JP');

select catalog.relate('magia-etcetera', 'references', 'magia-report');
update catalog.entry set unit_available = 13 where id = catalog.id('magia-record-manga');
select catalog.add_units('magia-record-manga', 13);
