-- KonoSuba: seasons as entries, with a romanized title and a short synonym.

select catalog.add_group('konosuba', 'series', 'KonoSuba', p_public_id => 99);
select catalog.add_entry('konosuba-s1', 'anime', 'tv', 'KonoSuba: God''s Blessing on This Wonderful World!', '2016-01', 'JP', 'episode', 10, p_public_id => 100);
select catalog.add_entry('konosuba-s2', 'anime', 'tv', 'KonoSuba: God''s Blessing on This Wonderful World! 2', '2017-01', 'JP', 'episode', 10, p_public_id => 101);
select catalog.add_entry('konosuba-s3', 'anime', 'tv', 'KonoSuba: God''s Blessing on This Wonderful World! 3', '2024-04', 'JP', 'episode', 11, p_public_id => 102);
select catalog.alias('konosuba-s2', 'romanized', 'Kono Subarashii Sekai ni Shukufuku wo! 2');
select catalog.alias('konosuba-s2', 'synonym', 'KonoSuba S2');
select catalog.member('konosuba', 'konosuba-s1', 'konosuba-s2', 'konosuba-s3');
select catalog.ext('konosuba', 'wikidata', 'Q86838208');
select catalog.relate('konosuba-s2', 'sequel_of', 'konosuba-s1');
select catalog.relate('konosuba-s3', 'sequel_of', 'konosuba-s2');
