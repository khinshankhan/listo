-- Toaru: two series that overlap at odd episode ranges.

select catalog.add_group('toaru', 'franchise', 'Toaru series', p_public_id => 52);
select catalog.alias('toaru', 'synonym', 'A Certain Magical Index');
select catalog.alias('toaru', 'native', 'とあるシリーズ');
select catalog.add_group('index-anime', 'series', 'A Certain Magical Index (anime)', p_public_id => 53);
select catalog.add_group('railgun-anime', 'series', 'A Certain Scientific Railgun (anime)', p_public_id => 54);
select catalog.member('toaru', 'index-anime', 'railgun-anime');

select catalog.add_entry('index-s1', 'anime', 'tv', 'A Certain Magical Index', '2008-10', 'JP', 'episode', 24, p_public_id => 55);
select catalog.add_entry('index-ii', 'anime', 'tv', 'A Certain Magical Index II', '2010-10', 'JP', 'episode', 24, p_public_id => 56);
select catalog.add_entry('railgun-s1', 'anime', 'tv', 'A Certain Scientific Railgun', '2009-10', 'JP', 'episode', 24, p_public_id => 57);
select catalog.add_entry('railgun-s', 'anime', 'tv', 'A Certain Scientific Railgun S', '2013-04', 'JP', 'episode', 24, p_public_id => 58);
select catalog.add_units(s, 24) from unnest(array['index-s1', 'index-ii', 'railgun-s1', 'railgun-s']) s;
select catalog.alias('index-s1', 'romanized', 'Toaru Majutsu no Index');
select catalog.alias('railgun-s1', 'romanized', 'Toaru Kagaku no Railgun');
select catalog.member('index-anime', 'index-s1', 'index-ii');
select catalog.member('railgun-anime', 'railgun-s1', 'railgun-s');

-- Wikidata's anime items are whole shows and ours are seasons, so the show
-- ids go on the series groups.
select catalog.ext('index-anime', 'wikidata', 'Q17219196');
select catalog.ext('railgun-anime', 'wikidata', 'Q66105329');
select catalog.ext('index-s1', 'wikidata', 'Q65083155');
select catalog.ext('index-s1', 'anilist_anime', '4654');
select catalog.ext('index-s1', 'mal_anime', '4654');
select catalog.ext('index-ii', 'wikidata', 'Q65083202');
select catalog.ext('index-ii', 'anilist_anime', '8937');
select catalog.ext('index-ii', 'mal_anime', '8937');
select catalog.ext('railgun-s1', 'wikidata', 'Q65083461');
select catalog.ext('railgun-s', 'wikidata', 'Q65083462');
select catalog.ext('railgun-s', 'anilist_anime', '16049');
select catalog.ext('railgun-s', 'mal_anime', '16049');

select catalog.relate('index-ii', 'sequel_of', 'index-s1');
select catalog.relate('railgun-s', 'sequel_of', 'railgun-s1');
select catalog.relate('railgun-s1', 'spin_off_of', 'index-s1');

-- NOTE: the ranges are placeholders, not a real guide.
select catalog.add_ordering('toaru', 'chronological', 'chronological (shape test)',
  'Interleaves two series at episode granularity. The ranges are placeholders, not a real guide.') as o \gset
select catalog.add_step(:'o', 1, 'index-s1', '1', '8');
select catalog.add_step(:'o', 2, 'railgun-s1', '1', '12');
select catalog.add_step(:'o', 3, 'index-s1', '9', '13');
select catalog.add_step(:'o', 4, 'railgun-s', '1', '16', p_note => 'the same arc from the other side');
select catalog.add_step(:'o', 5, 'index-s1', '14', '24');
select catalog.add_step(:'o', 6, 'railgun-s1', '13', '24');
select catalog.add_step(:'o', 7, 'index-ii', '1', '24');
select catalog.add_step(:'o', 8, 'railgun-s', '17', '24');
