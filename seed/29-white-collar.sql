-- White Collar: original series, licensed game and proposed revival.
-- https://en.wikipedia.org/wiki/White_Collar_(TV_series)
-- https://en.wikipedia.org/wiki/List_of_White_Collar_episodes

select catalog.add_group('white-collar', 'franchise', 'White Collar', p_public_id => 896);
select catalog.add_entry('white-collar-tv', 'tv', 'series', 'White Collar', '2009-10-23', 'US', 'episode', 81, p_public_id => 897);
select catalog.member('white-collar', 'white-collar-tv');
select catalog.ext('white-collar-tv', 'imdb', 'tt1358522');

select catalog.add_units('white-collar-tv', 81);
update catalog.unit u set part = d.season, number = u.ordinal - d.preceding,
  label = 'S' || d.season || ' E' || (u.ordinal - d.preceding)
from (values (1, 0, 14), (2, 14, 16), (3, 30, 16), (4, 46, 16),
  (5, 62, 13), (6, 75, 6)) as d (season, preceding, episodes)
where u.entry_id = catalog.id('white-collar-tv') and u.ordinal between d.preceding + 1 and d.preceding + d.episodes;

update catalog.entry set status = 'finished', end_date = '2014-12-18', end_precision = 'day',
  unit_minutes = 43, source_type = 'original',
  description = 'Con artist Neal Caffrey trades prison time for a tracking anklet and a job helping FBI agent Peter Burke catch art thieves, forgers and other white-collar criminals in New York.'
where id = catalog.id('white-collar-tv');

select catalog.credit('white-collar-tv', 'Fox Television Studios', 'producer', 'US');
select catalog.tag_with('genre', 'crime', 'white-collar-tv');
select catalog.tag_with('genre', 'drama', 'white-collar-tv');
select catalog.tag_with('genre', 'mystery', 'white-collar-tv');

-- https://www.chiefmarketer.com/ford-white-collar-chasing-the-shadow/
-- https://rbr.com/usa-network-ford-launch-game-for-white-collar/
select catalog.add_entry('white-collar-chasing-the-shadow', 'game', 'full', 'White Collar: Chasing the Shadow', '2010', 'US', 'chapter', 7, p_public_id => 898);
select catalog.member('white-collar', 'white-collar-chasing-the-shadow');
select catalog.relate('white-collar-chasing-the-shadow', 'side_story_of', 'white-collar-tv');
update catalog.entry set status = 'finished', description = 'An official interactive story extension, developed with creator Jeff Eastin, used live-action video and browser puzzles across seven weekly cases.' where id = catalog.id('white-collar-chasing-the-shadow');

-- https://www.tvline.com/news/white-collar-revival-renaissance-name-title-1235339469/
-- https://deadline.com/2025/02/white-collar-reboot-shopped-jeff-eastin-hulu-1236241010/
select catalog.add_entry('white-collar-renaissance', 'tv', 'series', 'White Collar: Renaissance', null, 'US', 'episode', null, p_public_id => 899);
select catalog.member('white-collar', 'white-collar-renaissance');
select catalog.relate('white-collar-renaissance', 'sequel_of', 'white-collar-tv', p_note => 'Proposed revival continuing the original series; a series order and final production format are unconfirmed.');
update catalog.entry set status = 'announced', description = 'A proposed continuation of White Collar, developed by original creator Jeff Eastin under the working title White Collar Renaissance. The pilot script is titled Masquerade. The project has no confirmed release date or episode count.', source_type = 'original' where id = catalog.id('white-collar-renaissance');
