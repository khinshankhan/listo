-- Steins;Gate and Science Adventure.
-- Facts are real where an external id is attached (checked against Wikidata
-- on 2026-10-17).

select catalog.add_group('science-adventure', 'universe', 'Science Adventure');
select catalog.alias('science-adventure', 'native', '科学アドベンチャーシリーズ');
select catalog.add_group('steins-gate', 'franchise', 'Steins;Gate');
select catalog.add_group('steins-gate-main', 'continuity', 'Main timeline', 'primary');
select catalog.add_group('steins-gate-zero', 'continuity', 'Divergence (0)', 'concluded');
select catalog.add_group('robotics-notes', 'franchise', 'Robotics;Notes');
select catalog.member('science-adventure', 'steins-gate', 'robotics-notes');
select catalog.member('steins-gate', 'steins-gate-main', 'steins-gate-zero');

select catalog.add_entry('steins-gate-vn', 'visual_novel', 'full', 'Steins;Gate', '2009', 'JP');
select catalog.add_entry('steins-gate-anime', 'anime', 'tv', 'Steins;Gate', '2011-04', 'JP', 'episode', 24);
select catalog.add_units('steins-gate-anime', 24);
select catalog.add_special('steins-gate-anime', '23β');
select catalog.alias('steins-gate-anime', 'native', 'シュタインズ・ゲート');
select catalog.alias('steins-gate-anime', 'synonym', 'S;G');
select catalog.add_entry('steins-gate-movie', 'anime', 'movie', 'Steins;Gate: The Movie - Load Region of Déjà Vu', '2013-04', 'JP');
select catalog.add_entry('steins-gate-elite', 'visual_novel', 'full', 'Steins;Gate Elite', '2018', 'JP');
select catalog.add_entry('steins-gate-0-vn', 'visual_novel', 'full', 'Steins;Gate 0', '2015', 'JP');
select catalog.add_entry('steins-gate-0-anime', 'anime', 'tv', 'Steins;Gate 0', '2018-04', 'JP', 'episode', 23);
select catalog.add_units('steins-gate-0-anime', 23);
select catalog.add_entry('robotics-notes-vn', 'visual_novel', 'full', 'Robotics;Notes', '2012', 'JP');
select catalog.add_entry('robotics-notes-dash', 'visual_novel', 'full', 'Robotics;Notes DaSH', '2019', 'JP');

select catalog.member('steins-gate-main', 'steins-gate-vn', 'steins-gate-anime', 'steins-gate-movie', 'steins-gate-elite');
select catalog.member('steins-gate-zero', 'steins-gate-0-vn', 'steins-gate-0-anime');
select catalog.member('robotics-notes', 'robotics-notes-vn', 'robotics-notes-dash');

select catalog.ext('steins-gate-vn', 'wikidata', 'Q482494');
select catalog.ext('steins-gate-vn', 'vndb', 'v2002');
select catalog.ext('steins-gate-anime', 'wikidata', 'Q20590069');
select catalog.ext('steins-gate-anime', 'anilist_anime', '9253');
select catalog.ext('steins-gate-anime', 'mal_anime', '9253');
select catalog.ext('steins-gate-movie', 'wikidata', 'Q7606919');
select catalog.ext('steins-gate-movie', 'anilist_anime', '11577');
select catalog.ext('steins-gate-movie', 'mal_anime', '11577');
select catalog.ext('steins-gate-elite', 'wikidata', 'Q40887622');
select catalog.ext('steins-gate-0-vn', 'wikidata', 'Q21775901');
select catalog.ext('steins-gate-0-vn', 'vndb', 'v17102');
select catalog.ext('steins-gate-0-anime', 'wikidata', 'Q48835532');
select catalog.ext('steins-gate-0-anime', 'anilist_anime', '21127');
select catalog.ext('steins-gate-0-anime', 'mal_anime', '30484');

select catalog.relate('steins-gate-anime', 'adaptation_of', 'steins-gate-vn');
select catalog.relate('steins-gate-movie', 'sequel_of', 'steins-gate-anime');
select catalog.relate('steins-gate-elite', 'remake_of', 'steins-gate-vn');
select catalog.relate('steins-gate-0-anime', 'adaptation_of', 'steins-gate-0-vn');
select catalog.relate('steins-gate-0-anime', 'sequel_of', 'steins-gate-anime', p_to_unit => '23',
                      p_note => 'branches off episode 23');
select catalog.relate('robotics-notes-dash', 'crossover_with', 'steins-gate-vn',
                      p_spoiler => true, p_established_on => '2019-01-31', p_established_by => 'robotics-notes-dash');

select catalog.add_ordering('steins-gate', 'recommended', 'first-timer order') as o \gset
select catalog.add_step(:'o', 1, 'steins-gate-anime', '1', '22');
select catalog.add_step(:'o', 2, 'steins-gate-anime', '23β', '23β', p_note => 'required for 0');
select catalog.add_step(:'o', 3, 'steins-gate-0-anime');
select catalog.add_step(:'o', 4, 'steins-gate-movie', p_necessity => 'optional', p_note => 'epilogue');

