-- Madoka: the recap films retell the series, so the first stop is a choice.

select catalog.add_group('madoka', 'series', 'Puella Magi Madoka Magica', p_public_id => 59);
select catalog.alias('madoka', 'romanized', 'Mahou Shoujo Madoka Magica');
select catalog.add_entry('madoka-tv', 'anime', 'tv', 'Puella Magi Madoka Magica', '2011-01', 'JP', 'episode', 12, p_public_id => 60);
select catalog.add_units('madoka-tv', 12);
select catalog.add_entry('madoka-beginnings', 'anime', 'movie', 'Puella Magi Madoka Magica the Movie Part I: Beginnings', '2012-10', 'JP', p_public_id => 61);
select catalog.add_entry('madoka-eternal', 'anime', 'movie', 'Puella Magi Madoka Magica the Movie Part II: Eternal', '2012-10', 'JP', p_public_id => 62);
select catalog.add_entry('madoka-rebellion', 'anime', 'movie', 'Puella Magi Madoka Magica the Movie Part III: Rebellion', '2013-10', 'JP', p_public_id => 63);
select catalog.member('madoka', 'madoka-tv', 'madoka-beginnings', 'madoka-eternal', 'madoka-rebellion');

select catalog.ext('madoka-tv', 'wikidata', 'Q53353');
select catalog.ext('madoka-beginnings', 'wikidata', 'Q3842376');
select catalog.ext('madoka-eternal', 'wikidata', 'Q3842375');
select catalog.ext('madoka-rebellion', 'wikidata', 'Q17479150');

select catalog.relate('madoka-beginnings', 'recap_of', 'madoka-tv', p_note => 'film retelling, first half');
select catalog.relate('madoka-eternal', 'recap_of', 'madoka-tv', p_note => 'film retelling, second half');
select catalog.relate('madoka-rebellion', 'sequel_of', 'madoka-tv');
select catalog.relate('madoka-rebellion', 'sequel_of', 'madoka-eternal');

select catalog.add_ordering('madoka', 'recommended', 'first-timer order',
  'The series or the two recap films, then Rebellion.') as o \gset
select catalog.add_step(:'o', 1, 'madoka-tv', p_branch => 1);
select catalog.add_step(:'o', 1, 'madoka-beginnings', p_branch => 2, p_step => 1, p_note => 'retells the series');
select catalog.add_step(:'o', 1, 'madoka-eternal', p_branch => 2, p_step => 2);
select catalog.add_step(:'o', 2, 'madoka-rebellion');
