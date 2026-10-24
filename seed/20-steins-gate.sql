-- Steins;Gate and Science Adventure.
-- Facts are real where an external id is attached (checked against Wikidata
-- on 2026-10-17).

select catalog.add_group('science-adventure', 'universe', 'Science Adventure', p_public_id => 1);
select catalog.alias('science-adventure', 'native', '科学アドベンチャーシリーズ');
select catalog.add_group('steins-gate', 'franchise', 'Steins;Gate', p_public_id => 2);
select catalog.add_group('steins-gate-main', 'continuity', 'Main timeline', 'primary', p_public_id => 3);
select catalog.add_group('steins-gate-zero', 'continuity', 'Divergence (0)', 'concluded', p_public_id => 4);
select catalog.add_group('robotics-notes', 'franchise', 'Robotics;Notes', p_public_id => 5);
select catalog.add_group('chaos-head', 'franchise', 'Chaos;Head', p_public_id => 118);
select catalog.add_group('occultic-nine', 'franchise', 'Occultic;Nine', p_public_id => 119);
select catalog.add_group('anonymous-code', 'franchise', 'Anonymous;Code', p_public_id => 120);
select catalog.member('science-adventure', 'chaos-head', 'steins-gate', 'robotics-notes', 'occultic-nine', 'anonymous-code');
select catalog.member('steins-gate', 'steins-gate-main', 'steins-gate-zero');

select catalog.add_entry('steins-gate-vn', 'visual_novel', 'full', 'Steins;Gate', '2009', 'JP', p_public_id => 6);
select catalog.add_entry('steins-gate-anime', 'anime', 'tv', 'Steins;Gate', '2011-04', 'JP', 'episode', 24, p_public_id => 7);
select catalog.add_units('steins-gate-anime', 24);
select catalog.add_special('steins-gate-anime', '23β');
select catalog.alias('steins-gate-anime', 'native', 'シュタインズ・ゲート');
select catalog.alias('steins-gate-anime', 'synonym', 'S;G');
select catalog.add_entry('steins-gate-movie', 'anime', 'movie', 'Steins;Gate: The Movie - Load Region of Déjà Vu', '2013-04', 'JP', p_public_id => 8);
select catalog.add_entry('steins-gate-elite', 'visual_novel', 'full', 'Steins;Gate Elite', '2018', 'JP', p_public_id => 9);
select catalog.add_entry('steins-gate-0-vn', 'visual_novel', 'full', 'Steins;Gate 0', '2015', 'JP', p_public_id => 10);
select catalog.add_entry('steins-gate-0-anime', 'anime', 'tv', 'Steins;Gate 0', '2018-04', 'JP', 'episode', 23, p_public_id => 11);
select catalog.add_units('steins-gate-0-anime', 23);
select catalog.add_entry('robotics-notes-vn', 'visual_novel', 'full', 'Robotics;Notes', '2012', 'JP', p_public_id => 12);
select catalog.add_entry('robotics-notes-dash', 'visual_novel', 'full', 'Robotics;Notes DaSH', '2019', 'JP', p_public_id => 13);
select catalog.add_entry('robotics-notes-anime', 'anime', 'tv', 'Robotics;Notes', '2012-10', 'JP', 'episode', 22, p_public_id => 121);
select catalog.add_entry('chaos-head-vn', 'visual_novel', 'full', 'Chaos;Head', '2008-04', 'JP', p_public_id => 122);
select catalog.add_entry('chaos-head-anime', 'anime', 'tv', 'Chaos;Head', '2008-10', 'JP', 'episode', 12, p_public_id => 123);
select catalog.add_entry('chaos-child-vn', 'visual_novel', 'full', 'Chaos;Child', '2014-12', 'JP', p_public_id => 124);
select catalog.add_entry('chaos-child-anime', 'anime', 'tv', 'Chaos;Child', '2017-01', 'JP', 'episode', 12, p_public_id => 125);
select catalog.add_entry('occultic-nine-ln', 'book', 'light_novel', 'Occultic;Nine', '2014-08', 'JP', p_public_id => 126);
select catalog.add_entry('occultic-nine-anime', 'anime', 'tv', 'Occultic;Nine', '2016-10', 'JP', 'episode', 12, p_public_id => 127);
select catalog.add_entry('occultic-nine-vn', 'visual_novel', 'full', 'Occultic;Nine', '2017-11', 'JP', p_public_id => 128);
select catalog.add_entry('anonymous-code-vn', 'visual_novel', 'full', 'Anonymous;Code', '2022-07', 'JP', p_public_id => 129);

select catalog.member('steins-gate-main', 'steins-gate-vn', 'steins-gate-anime', 'steins-gate-movie', 'steins-gate-elite');
select catalog.member('steins-gate-zero', 'steins-gate-0-vn', 'steins-gate-0-anime');
select catalog.set_role('steins-gate-main', 'steins-gate-vn', 'source material');
select catalog.set_role('steins-gate-main', 'steins-gate-anime', 'adaptation of the VN');
select catalog.set_role('steins-gate-main', 'steins-gate-movie', 'sequel to the anime');
select catalog.set_role('steins-gate-main', 'steins-gate-elite', 'remake of the VN');
select catalog.set_role('steins-gate-zero', 'steins-gate-0-vn', 'sequel from another worldline');
select catalog.set_role('steins-gate-zero', 'steins-gate-0-anime', 'adaptation of 0');
select catalog.member('robotics-notes', 'robotics-notes-vn', 'robotics-notes-anime', 'robotics-notes-dash');
select catalog.member('chaos-head', 'chaos-head-vn', 'chaos-head-anime', 'chaos-child-vn', 'chaos-child-anime');
select catalog.member('occultic-nine', 'occultic-nine-ln', 'occultic-nine-anime', 'occultic-nine-vn');
select catalog.member('anonymous-code', 'anonymous-code-vn');

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
select catalog.relate('robotics-notes-anime', 'adaptation_of', 'robotics-notes-vn');
select catalog.relate('robotics-notes-dash', 'sequel_of', 'robotics-notes-vn');
select catalog.relate('chaos-head-anime', 'adaptation_of', 'chaos-head-vn');
select catalog.relate('chaos-child-vn', 'sequel_of', 'chaos-head-vn');
select catalog.relate('chaos-child-anime', 'adaptation_of', 'chaos-child-vn');
select catalog.relate('occultic-nine-anime', 'adaptation_of', 'occultic-nine-ln');
select catalog.relate('occultic-nine-vn', 'adaptation_of', 'occultic-nine-ln');
select catalog.relate('robotics-notes-dash', 'crossover_with', 'steins-gate-vn',
                      p_spoiler => true, p_established_on => '2019-01-31', p_established_by => 'robotics-notes-dash');

select catalog.add_ordering('steins-gate', 'recommended', 'first-timer order') as o \gset
select catalog.add_step(:'o', 1, 'steins-gate-anime', '1', '22');
select catalog.add_step(:'o', 2, 'steins-gate-anime', '23β', '23β', p_note => 'required for 0');
select catalog.add_step(:'o', 3, 'steins-gate-0-anime');
select catalog.add_step(:'o', 4, 'steins-gate-movie', p_necessity => 'optional', p_note => 'epilogue');

