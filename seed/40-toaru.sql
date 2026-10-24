-- Toaru: two series that overlap at odd episode ranges, and the spin-offs
-- and side stories around them (the anime, as AniDB relates it).

select catalog.add_group('toaru', 'franchise', 'Toaru series', p_public_id => 52);
select catalog.alias('toaru', 'synonym', 'A Certain Magical Index');
select catalog.alias('toaru', 'native', 'とあるシリーズ');
select catalog.add_group('index-anime', 'series', 'A Certain Magical Index (anime)', p_public_id => 53);
select catalog.add_group('railgun-anime', 'series', 'A Certain Scientific Railgun (anime)', p_public_id => 54);
select catalog.add_group('accelerator-anime', 'series', 'A Certain Scientific Accelerator (anime)', p_public_id => 130);
select catalog.add_group('item-anime', 'series', 'A Certain Item of Dark Side (anime)', p_public_id => 131);
select catalog.member('toaru', 'index-anime', 'railgun-anime', 'accelerator-anime', 'item-anime');

select catalog.add_entry('index-s1', 'anime', 'tv', 'A Certain Magical Index', '2008-10', 'JP', 'episode', 24, p_public_id => 55);
select catalog.add_entry('index-ii', 'anime', 'tv', 'A Certain Magical Index II', '2010-10', 'JP', 'episode', 24, p_public_id => 56);
select catalog.add_entry('railgun-s1', 'anime', 'tv', 'A Certain Scientific Railgun', '2009-10', 'JP', 'episode', 24, p_public_id => 57);
select catalog.add_entry('railgun-s', 'anime', 'tv', 'A Certain Scientific Railgun S', '2013-04', 'JP', 'episode', 24, p_public_id => 58);
select catalog.add_units(s, 24) from unnest(array['index-s1', 'index-ii', 'railgun-s1', 'railgun-s']) s;
select catalog.add_entry('index-iii', 'anime', 'tv', 'A Certain Magical Index III', '2018-10', 'JP', 'episode', 26, p_public_id => 132);
select catalog.add_entry('index-endymion', 'anime', 'movie', 'A Certain Magical Index: The Movie - The Miracle of Endymion', '2013-02', 'JP', p_public_id => 133);
select catalog.add_entry('railgun-t', 'anime', 'tv', 'A Certain Scientific Railgun T', '2020-01', 'JP', 'episode', 25, p_public_id => 134);
-- announced in 2025, no date yet
select catalog.add_entry('railgun-4', 'anime', 'tv', 'A Certain Scientific Railgun season 4', null, 'JP', p_public_id => 135);
select catalog.add_entry('railgun-ova-model', 'anime', 'ova', 'Railgun: Entenka no Satsuei Model mo Raku ja Arimasen wa ne', '2010-07', 'JP', p_public_id => 136);
select catalog.add_entry('railgun-ova-misaka', 'anime', 'ova', 'Railgun: Misaka-san wa Ima Chuumoku no Mato desukara', '2010', 'JP', p_public_id => 137);
select catalog.add_entry('railgun-s-ova', 'anime', 'ova', 'Railgun S: Daiji na Koto wa Zenbu Sentou ni Osowatta', '2014-03', 'JP', p_public_id => 138);
select catalog.add_entry('accelerator', 'anime', 'tv', 'A Certain Scientific Accelerator', '2019-07', 'JP', 'episode', 12, p_public_id => 139);
select catalog.add_entry('item', 'anime', 'tv', 'A Certain Item of Dark Side', '2026', 'JP', p_public_id => 140);
select catalog.alias('index-s1', 'romanized', 'Toaru Majutsu no Index');
select catalog.alias('railgun-s1', 'romanized', 'Toaru Kagaku no Railgun');
select catalog.member('index-anime', 'index-s1', 'index-ii', 'index-endymion', 'index-iii');
select catalog.member('railgun-anime', 'railgun-s1', 'railgun-ova-model', 'railgun-ova-misaka', 'railgun-s', 'railgun-s-ova', 'railgun-t', 'railgun-4');
select catalog.member('accelerator-anime', 'accelerator');
select catalog.member('item-anime', 'item');

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
select catalog.relate('index-iii', 'sequel_of', 'index-ii');
select catalog.relate('index-endymion', 'side_story_of', 'index-ii');
select catalog.relate('railgun-t', 'sequel_of', 'railgun-s');
select catalog.relate('railgun-4', 'sequel_of', 'railgun-t');
select catalog.relate('railgun-ova-model', 'side_story_of', 'railgun-s1');
select catalog.relate('railgun-ova-misaka', 'side_story_of', 'railgun-s1');
select catalog.relate('railgun-s-ova', 'side_story_of', 'railgun-s');
select catalog.relate('accelerator', 'spin_off_of', 'index-s1');
select catalog.relate('item', 'spin_off_of', 'index-s1');

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
