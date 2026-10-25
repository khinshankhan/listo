-- Psych: series, sequel films, specials and licensed companion works.
-- https://en.wikipedia.org/wiki/List_of_Psych_episodes
-- https://www.usanetwork.com/usa-insider/psych-movies-what-to-know-where-to-watch

select catalog.add_group('psych', 'franchise', 'Psych', p_public_id => 827);
select catalog.add_entry('psych-tv', 'tv', 'series', 'Psych', '2006-07-07', 'US', 'episode', 121, p_public_id => 828);
select catalog.add_entry('psych-the-movie', 'film', 'tv_movie', 'Psych: The Movie', '2017-12-07', 'US', p_public_id => 829);
select catalog.add_entry('psych-2-lassie-come-home', 'film', 'tv_movie', 'Psych 2: Lassie Come Home', '2020-07-15', 'US', p_public_id => 830);
select catalog.add_entry('psych-3-this-is-gus', 'film', 'tv_movie', 'Psych 3: This Is Gus', '2021-11-18', 'US', p_public_id => 831);
select catalog.member('psych', 'psych-tv', 'psych-the-movie', 'psych-2-lassie-come-home', 'psych-3-this-is-gus');
select catalog.set_role('psych', 'psych-tv', 'original series');
select catalog.set_role('psych', 'psych-the-movie', 'sequel film');
select catalog.set_role('psych', 'psych-2-lassie-come-home', 'sequel film');
select catalog.set_role('psych', 'psych-3-this-is-gus', 'sequel film');

select catalog.ext('psych-tv', 'imdb', 'tt0491738');
select catalog.ext('psych-the-movie', 'imdb', 'tt6868216');
select catalog.ext('psych-2-lassie-come-home', 'imdb', 'tt9792884');
select catalog.ext('psych-3-this-is-gus', 'imdb', 'tt14641648');

-- The musical occupies two episodes in the 121-episode count.
select catalog.add_units('psych-tv', 121);
update catalog.unit u set part = d.season, number = u.ordinal - d.preceding,
  label = 'S' || d.season || ' E' || (u.ordinal - d.preceding)
from (values (1, 0, 15), (2, 15, 16), (3, 31, 16), (4, 47, 16),
  (5, 63, 16), (6, 79, 16), (7, 95, 16), (8, 111, 10)) as d (season, preceding, episodes)
where u.entry_id = catalog.id('psych-tv') and u.ordinal between d.preceding + 1 and d.preceding + d.episodes;
update catalog.unit set title = 'Psych: The Musical, Part 1' where entry_id = catalog.id('psych-tv') and ordinal = 110;
update catalog.unit set title = 'Psych: The Musical, Part 2' where entry_id = catalog.id('psych-tv') and ordinal = 111;

select catalog.relate('psych-the-movie', 'sequel_of', 'psych-tv', p_to_unit => '121');
select catalog.relate('psych-2-lassie-come-home', 'sequel_of', 'psych-the-movie');
select catalog.relate('psych-3-this-is-gus', 'sequel_of', 'psych-2-lassie-come-home');

update catalog.entry set status = 'finished', source_type = 'original'
where id in (catalog.id('psych-tv'), catalog.id('psych-the-movie'),
  catalog.id('psych-2-lassie-come-home'), catalog.id('psych-3-this-is-gus'));
update catalog.entry set end_date = '2014-03-26', end_precision = 'day', unit_minutes = 42
where id = catalog.id('psych-tv');

-- Streaming edition runtimes.
-- https://www.peacocktv.com/watch-online/movies/psych-the-movie/ff4a71e8-e75d-31fc-97a4-09fa91ff3010
-- https://www.peacocktv.com/watch-online/movies/psych-2-lassie-come-home/10633c67-7e9f-3cdb-9017-a30781015ef7
-- https://www.peacocktv.com/watch-online/movies/psych-3-this-is-gus/4695aa38-6632-3b3f-aa14-9edf3b70d24a
update catalog.entry e set length_minutes = d.minutes
from (values ('psych-the-movie', 88), ('psych-2-lassie-come-home', 88), ('psych-3-this-is-gus', 96)) as d (slug, minutes)
where e.id = catalog.id(d.slug);

update catalog.entry e set description = d.description
from (values
  ('psych-tv', 'Shawn Spencer turns his sharp observation skills into a fake psychic detective act, solving cases for the Santa Barbara police with his skeptical best friend Gus.'),
  ('psych-the-movie', 'Three years after leaving Santa Barbara, Shawn and Gus reunite with their friends when an attack on someone close to them pulls the group into a new case in San Francisco.'),
  ('psych-2-lassie-come-home', 'Shawn and Gus return to Santa Barbara after Lassiter is seriously injured, investigating strange events around his recovery clinic.'),
  ('psych-3-this-is-gus', 'With Gus preparing for marriage and fatherhood, he and Shawn investigate the hidden past of his fiancee and try to track down her estranged husband.')) as d (slug, description)
where e.id = catalog.id(d.slug);

select catalog.credit(s, 'Universal Cable Productions', 'producer', 'US') from unnest(array['psych-tv', 'psych-the-movie']) s;
select catalog.credit(s, 'Universal Content Productions', 'producer', 'US')
from unnest(array['psych-2-lassie-come-home', 'psych-3-this-is-gus']) s;
select catalog.tag_with('genre', 'comedy', 'psych-tv', 'psych-the-movie', 'psych-2-lassie-come-home', 'psych-3-this-is-gus');
select catalog.tag_with('genre', 'mystery', 'psych-tv', 'psych-the-movie', 'psych-2-lassie-come-home', 'psych-3-this-is-gus');

-- https://www.universalpicturesathome.com/tv/psych-the-musical
-- https://www.jjsedelmaier.com/bc3
-- https://psychusa.fandom.com/wiki/Psych:_The_Big_Adventures_of_Little_Shawn_and_Gus
-- https://www.thewrap.com/attack-shows-kevin-pereira-host-psych-finale-show/
-- https://www.youtube.com/watch?v=4ORfpdI7JR8
-- https://catalog.freelibrary.org/Record/1608997
-- https://www.tvguide.com/news/psych-fans-vote-clue-episode-1062707/
-- https://en.wikipedia.org/wiki/Psych
-- https://www.penguin.com.au/books/psych-a-fatal-frame-of-mind-9780451231598
-- https://www.hachettebookgroup.com/titles/shawn-spencer/psychs-guide-to-crime-fighting-for-the-totally-unqualified/9781455512867/
-- https://www.igdb.com/games/psych/presskit
-- https://worldscreen.com/usa-network-debuts-new-animated-shorts/
select catalog.add_entry('psych-the-musical', 'tv', 'special', 'Psych: The Musical', '2013-12-15', 'US', 'none', null, p_public_id => 832);
select catalog.member('psych', 'psych-the-musical');
select catalog.relate('psych-the-musical', 'side_story_of', 'psych-tv', p_to_unit => '110', p_note => 'The standalone edition contains the same story as TV episodes 110-111.');
update catalog.entry set status = 'finished', description = 'Shawn and Gus investigate a case involving a psychiatric patient and a burned theater in a two-part musical story. The same story is included as season 7 episodes 15 and 16 of the TV series.', source_type = 'original', length_minutes = 88 where id = catalog.id('psych-the-musical');
select catalog.add_entry('psych-little-shawn-and-gus', 'tv', 'series', 'The Big Adventures of Little Shawn and Gus', '2008-01-11', 'US', 'episode', 14, p_public_id => 833);
select catalog.member('psych', 'psych-little-shawn-and-gus');
select catalog.relate('psych-little-shawn-and-gus', 'spin_off_of', 'psych-tv');
update catalog.entry set status = 'finished', description = 'Animated shorts follow Shawn and Gus as children. Two seasons were produced for television breaks, the web and home video; individual shorts run from under a minute to two minutes.', source_type = 'original' where id = catalog.id('psych-little-shawn-and-gus');
select catalog.ext('psych-little-shawn-and-gus', 'imdb', 'tt1654413');
select catalog.add_entry('psych-after-pshow', 'tv', 'special', 'Psych After Pshow', '2014-03-26', 'US', 'none', null, p_public_id => 834);
select catalog.member('psych', 'psych-after-pshow');
select catalog.relate('psych-after-pshow', 'references', 'psych-tv');
update catalog.entry set status = 'finished', description = 'Kevin Pereira hosts the cast for a live farewell aftershow following the original series finale.', length_minutes = 44 where id = catalog.id('psych-after-pshow');
select catalog.ext('psych-after-pshow', 'imdb', 'tt4381292');
select catalog.add_entry('psych-cast-reunion-2020', 'tv', 'special', 'Psych Cast Reunites To Talk Classic Moments', '2020-05-04', 'US', 'none', null, p_public_id => 835);
select catalog.member('psych', 'psych-cast-reunion-2020');
select catalog.relate('psych-cast-reunion-2020', 'references', 'psych-tv');
update catalog.entry set status = 'finished', description = 'The cast reunites remotely during the Psychtacular Binge-a-thon to discuss favorite moments from the series.' where id = catalog.id('psych-cast-reunion-2020');
select catalog.add_entry('psych-international-pilot', 'tv', 'special', 'Psych: International Pilot', '2007', 'US', 'none', null, p_public_id => 836);
select catalog.member('psych', 'psych-international-pilot');
select catalog.relate('psych-international-pilot', 'alternative_version_of', 'psych-tv', p_to_unit => '1');
update catalog.entry set status = 'finished', description = 'The longer international cut of the pilot, also included among the first-season DVD extras.', length_minutes = 75 where id = catalog.id('psych-international-pilot');
select catalog.add_entry('psych-100-clues-alternate-ending', 'tv', 'special', 'Psych: 100 Clues - West Coast Ending', '2013-03-27', 'US', 'none', null, p_public_id => 837);
select catalog.member('psych', 'psych-100-clues-alternate-ending');
select catalog.relate('psych-100-clues-alternate-ending', 'alternative_version_of', 'psych-tv', p_to_unit => '100');
update catalog.entry set status = 'finished', description = 'The alternate West Coast broadcast ending of the viewer-voted 100th episode, 100 Clues.' where id = catalog.id('psych-100-clues-alternate-ending');
select catalog.add_entry('psych-a-mind-is-a-terrible-thing-to-read', 'book', 'novel', 'Psych: A Mind Is a Terrible Thing to Read', '2009-01-06', 'US', 'none', null, p_public_id => 838);
select catalog.member('psych', 'psych-a-mind-is-a-terrible-thing-to-read');
select catalog.relate('psych-a-mind-is-a-terrible-thing-to-read', 'side_story_of', 'psych-tv');
update catalog.entry set status = 'finished', description = 'The first original Psych tie-in novel by William Rabkin, following Shawn and Gus in a case outside the television episodes.' where id = catalog.id('psych-a-mind-is-a-terrible-thing-to-read');
select catalog.credit('psych-a-mind-is-a-terrible-thing-to-read', 'Berkley', 'publisher', 'US');
select catalog.add_entry('psych-mind-over-magic', 'book', 'novel', 'Psych: Mind Over Magic', '2009-07-07', 'US', 'none', null, p_public_id => 839);
select catalog.member('psych', 'psych-mind-over-magic');
select catalog.relate('psych-mind-over-magic', 'side_story_of', 'psych-tv');
update catalog.entry set status = 'finished', description = 'William Rabkin sends Shawn and Gus into a murder investigation involving a stage magician in the second original tie-in novel.' where id = catalog.id('psych-mind-over-magic');
select catalog.credit('psych-mind-over-magic', 'Berkley', 'publisher', 'US');
select catalog.add_entry('psych-the-call-of-the-mild', 'book', 'novel', 'Psych: The Call of the Mild', '2010-01-05', 'US', 'none', null, p_public_id => 840);
select catalog.member('psych', 'psych-the-call-of-the-mild');
select catalog.relate('psych-the-call-of-the-mild', 'side_story_of', 'psych-tv');
update catalog.entry set status = 'finished', description = 'The third original Psych tie-in novel by William Rabkin takes Shawn and Gus away from their usual Santa Barbara surroundings.' where id = catalog.id('psych-the-call-of-the-mild');
select catalog.credit('psych-the-call-of-the-mild', 'Berkley', 'publisher', 'US');
select catalog.add_entry('psych-a-fatal-frame-of-mind', 'book', 'novel', 'Psych: A Fatal Frame of Mind', '2010-08-03', 'US', 'none', null, p_public_id => 841);
select catalog.member('psych', 'psych-a-fatal-frame-of-mind');
select catalog.relate('psych-a-fatal-frame-of-mind', 'side_story_of', 'psych-tv');
update catalog.entry set status = 'finished', description = 'Shawn and Gus investigate a death at an art museum in William Rabkin''s fourth original Psych tie-in novel.' where id = catalog.id('psych-a-fatal-frame-of-mind');
select catalog.credit('psych-a-fatal-frame-of-mind', 'Berkley', 'publisher', 'US');
select catalog.add_entry('psych-mind-altering-murder', 'book', 'novel', 'Psych: Mind-Altering Murder', '2011-02-01', 'US', 'none', null, p_public_id => 842);
select catalog.member('psych', 'psych-mind-altering-murder');
select catalog.relate('psych-mind-altering-murder', 'side_story_of', 'psych-tv');
update catalog.entry set status = 'finished', description = 'William Rabkin''s fifth original Psych tie-in novel follows Shawn and Gus into another murder investigation.' where id = catalog.id('psych-mind-altering-murder');
select catalog.credit('psych-mind-altering-murder', 'Berkley', 'publisher', 'US');
select catalog.add_entry('psych-guide-to-crime-fighting', 'book', 'nonfiction', 'Psych''s Guide to Crime Fighting for the Totally Unqualified', '2013-05-07', 'US', 'none', null, p_public_id => 843);
select catalog.member('psych', 'psych-guide-to-crime-fighting');
select catalog.relate('psych-guide-to-crime-fighting', 'references', 'psych-tv');
update catalog.entry set status = 'finished', description = 'Chad Gervich writes an in-character companion guide attributed to Shawn Spencer and Burton Guster, combining crime-solving advice with the show''s humor.' where id = catalog.id('psych-guide-to-crime-fighting');
select catalog.credit('psych-guide-to-crime-fighting', 'Grand Central Publishing', 'publisher', 'US');
select catalog.add_entry('psych-game-2014', 'game', 'full', 'Psych', '2014-03-26', 'US', 'none', null, p_public_id => 844);
select catalog.member('psych', 'psych-game-2014');
select catalog.relate('psych-game-2014', 'adaptation_of', 'psych-tv');
update catalog.entry set status = 'finished', description = 'A licensed hidden-object mystery game developed by Legacy Games, featuring Shawn and Gus.' where id = catalog.id('psych-game-2014');
select catalog.add_entry('psych-pogo-adventure', 'game', 'full', 'Pogo Adventure', '2008', 'US', 'none', null, p_public_id => 845);
select catalog.member('psych', 'psych-pogo-adventure');
select catalog.relate('psych-pogo-adventure', 'spin_off_of', 'psych-little-shawn-and-gus');
update catalog.entry set status = 'finished', description = 'An official browser game released alongside the Little Shawn and Gus animated shorts.' where id = catalog.id('psych-pogo-adventure');
select catalog.add_entry('psych-hotshot-shootout', 'game', 'full', 'Hotshot Shootout', '2008', 'US', 'none', null, p_public_id => 846);
select catalog.member('psych', 'psych-hotshot-shootout');
select catalog.relate('psych-hotshot-shootout', 'spin_off_of', 'psych-little-shawn-and-gus');
update catalog.entry set status = 'finished', description = 'An official browser game accompanying the Little Shawn and Gus animated shorts.' where id = catalog.id('psych-hotshot-shootout');
select catalog.add_entry('psychologists-are-in', 'audio', 'podcast', 'The Psychologists Are In: A Psych Rewatch', '2021-11', 'US', 'episode', null, p_public_id => 847);
select catalog.member('psych', 'psychologists-are-in');
select catalog.relate('psychologists-are-in', 'references', 'psych-tv');
update catalog.entry set status = 'releasing', description = 'Maggie Lawson and Timothy Omundson revisit Psych episodes and films, interviewing fellow cast members and crew about making the series.' where id = catalog.id('psychologists-are-in');

-- These first-season DVD features are bonus units, not new story episodes.
select catalog.add_special('psych-tv', 'Season 1: James Roday audition tape');
select catalog.add_special('psych-tv', 'Season 1: Psych Revealed');
select catalog.add_special('psych-tv', 'Season 1: Inside the Writers Room');
select catalog.add_special('psych-tv', 'Season 1: character profiles');
select catalog.add_special('psych-tv', 'Season 1: blooper reel');
select catalog.add_special('psych-tv', 'Season 1: deleted scenes');
select catalog.add_special('psych-tv', 'Season 1: audio commentaries');

select catalog.add_units('psych-little-shawn-and-gus', 14);
update catalog.unit u set title = d.title from (values
  (1, 'Far Out'), (2, 'The Big Dig'), (3, 'Pre-School Reunion'), (4, 'Mystery Meat'),
  (5, 'Teacher''s Lounge'), (6, 'Bake Sale'), (7, 'Aliens'), (8, 'Missing Lunch'),
  (9, 'Hygeia'), (10, 'We All Scream'), (11, 'There Will Be Chocolate Sauce'),
  (12, 'Roller Chase'), (13, 'The Daredevil''s in the Details'), (14, 'Deja Who?')) as d (ordinal, title)
where u.entry_id = catalog.id('psych-little-shawn-and-gus') and u.ordinal = d.ordinal;
select catalog.credit('psych-little-shawn-and-gus', 'J. J. Sedelmaier Productions', 'studio', 'US');
select catalog.credit('psych-after-pshow', 'Wilshire Studios', 'producer', 'US');
select catalog.credit('psych-game-2014', 'Legacy Games', 'developer', 'US');
select catalog.credit('psychologists-are-in', 'Cloud10', 'producer', 'US');

-- Split around the musical's existing episode slots so each story is watched once.
select catalog.add_ordering('psych', 'release', 'release order',
  'TV episodes in broadcast order with the standalone musical replacing episodes 110-111, optional companion specials, then the sequel films. The animated shorts are an optional collection, not a required interruption between TV episodes.') as o \gset
select catalog.add_step(:'o', 1, 'psych-tv', '1', '15');
select catalog.add_step(:'o', 2, 'psych-international-pilot', p_necessity => 'optional', p_note => 'An alternate cut of the pilot, documented on the 2007 season-one DVD.');
select catalog.add_step(:'o', 3, 'psych-tv', '16', '25');
select catalog.add_step(:'o', 4, 'psych-little-shawn-and-gus', p_necessity => 'optional', p_note => 'The shorts began airing during season 2; the collection also includes later shorts.');
select catalog.add_step(:'o', 5, 'psych-tv', '26', '100');
select catalog.add_step(:'o', 6, 'psych-100-clues-alternate-ending', p_necessity => 'optional');
select catalog.add_step(:'o', 7, 'psych-tv', '101', '109');
select catalog.add_step(:'o', 8, 'psych-the-musical', p_note => 'Replaces the same musical story at TV episode slots 110-111.');
select catalog.add_step(:'o', 9, 'psych-tv', '112', '121');
select catalog.add_step(:'o', 10, 'psych-after-pshow', p_necessity => 'optional');
select catalog.add_step(:'o', 11, 'psych-the-movie');
select catalog.add_step(:'o', 12, 'psych-cast-reunion-2020', p_necessity => 'optional');
select catalog.add_step(:'o', 13, 'psych-2-lassie-come-home');
select catalog.add_step(:'o', 14, 'psych-3-this-is-gus');

-- https://www.usanetwork.com/usa-insider/will-there-be-psych-4-more-movies
-- https://www.tvinsider.com/1118094/psych-4-movie-timothy-omundson-carlton-lassiter/
-- https://en.wikipedia.org/wiki/Psych_season_7
-- https://psychusa.fandom.com/wiki/100_Clues
select catalog.add_entry('psych-4', 'film', 'tv_movie', 'Psych 4 (Working Title)', null, 'US', 'none', null, p_public_id => 848);
select catalog.member('psych', 'psych-4');
select catalog.relate('psych-4', 'sequel_of', 'psych-3-this-is-gus');
update catalog.entry set status = 'announced', description = 'A fourth Psych film has a completed script and returning-cast interest. Production has not been greenlit, and no final title or release date is confirmed.', source_type = 'original' where id = catalog.id('psych-4');
select catalog.add_entry('psych-100-clues-dvd-ending', 'tv', 'special', 'Psych: 100 Clues - DVD Ending', '2013-10-08', 'US', 'none', null, p_public_id => 849);
select catalog.member('psych', 'psych-100-clues-dvd-ending');
select catalog.relate('psych-100-clues-dvd-ending', 'alternative_version_of', 'psych-tv', p_to_unit => '100', p_note => 'The third ending, issued on the season-seven DVD rather than either original coast broadcast.');
update catalog.entry set status = 'finished' where id = catalog.id('psych-100-clues-dvd-ending');
