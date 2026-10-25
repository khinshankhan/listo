-- Shameless: a remake of a remake, across countries.

select catalog.add_group('shameless', 'franchise', 'Shameless', p_public_id => 1440);
-- https://en.wikipedia.org/wiki/Shameless_(British_TV_series)
-- https://en.wikipedia.org/wiki/Shameless_(American_TV_series)
-- https://en.wikipedia.org/wiki/Bizim_Hikaye
select catalog.add_entry('shameless-uk', 'tv', 'series', 'Shameless', '2004-01-13', 'GB', 'episode', 139, p_public_id => 1441);
select catalog.add_entry('shameless-us', 'tv', 'series', 'Shameless', '2011-01-09', 'US', 'episode', 134, p_public_id => 1442);
select catalog.add_entry('bizim-hikaye', 'tv', 'series', 'Bizim Hikaye', '2017-09-14', 'TR', 'episode', 70, p_title_kind => 'native', p_public_id => 1443);
select catalog.alias('bizim-hikaye', 'synonym', 'Our Story');
select catalog.member('shameless', 'shameless-uk', 'shameless-us', 'bizim-hikaye');
select catalog.ext('shameless-uk', 'wikidata', 'Q23605');
select catalog.ext('shameless-us', 'wikidata', 'Q23599');
select catalog.ext('bizim-hikaye', 'wikidata', 'Q47010718');
select catalog.relate('shameless-us', 'remake_of', 'shameless-uk');
select catalog.relate('bizim-hikaye', 'remake_of', 'shameless-us');

-- https://www.paramountpressexpress.com/showtime/releases/?view=56586
-- https://en.wikipedia.org/wiki/List_of_Shameless_(American_TV_series)_episodes
-- https://www.tvmaze.com/shows/31186/besstydniki
-- https://en.wikipedia.org/wiki/Itti_Si_Khushi
-- https://www.sonyliv.com/shows/title-1790006675
-- https://www.youtube.com/watch?v=GGCTK_32NY0
select catalog.add_entry('shameless-hall-of-shame', 'tv', 'series', 'Shameless Hall of Shame', '2020-12-27', 'US', 'episode', 6, p_public_id => 1444);
select catalog.member('shameless', 'shameless-hall-of-shame');
select catalog.relate('shameless-hall-of-shame', 'recap_of', 'shameless-us');
update catalog.entry set status = 'finished', description = 'Six companion episodes combine newly filmed scenes with retrospectives on the American series'' characters. They aired during the final season and are outside its 134 regular episodes.', end_date = '2021-02-28', end_precision = 'day' where id = catalog.id('shameless-hall-of-shame');
select catalog.add_entry('shameless-russia', 'tv', 'series', 'Besstydniki', '2017-09-24', 'RU', 'episode', 24, p_public_id => 1445);
select catalog.member('shameless', 'shameless-russia');
select catalog.relate('shameless-russia', 'remake_of', 'shameless-uk');
update catalog.entry set status = 'finished', description = 'A Russian adaptation relocates the dysfunctional family story to the Gruzdev household near Moscow.', unit_minutes = 45, end_date = '2017-12-10', end_precision = 'day' where id = catalog.id('shameless-russia');
select catalog.add_entry('itti-si-khushi', 'tv', 'series', 'Itti Si Khushi', '2025-08-18', 'IN', 'episode', 234, p_public_id => 1446);
select catalog.member('shameless', 'itti-si-khushi');
select catalog.relate('itti-si-khushi', 'remake_of', 'shameless-uk');
update catalog.entry set status = 'finished', description = 'The official Hindi adaptation follows Anvita Diwekar as she cares for her siblings in Mumbai while their parents fail them.', end_date = '2026-05-16', end_precision = 'day' where id = catalog.id('itti-si-khushi');

-- https://www.comedy.co.uk/tv/shameless/episodes/1/8/
-- https://www.comedy.co.uk/tv/shameless/episodes/3/1/
-- https://en.wikipedia.org/wiki/List_of_Shameless_(British_TV_series)_episodes
select catalog.add_entry('shameless-uk-christmas-special', 'tv', 'special', 'Shameless: Christmas in Chatsworth', '2004-12-23', 'GB', 'none', null, p_public_id => 1447);
select catalog.member('shameless', 'shameless-uk-christmas-special');
select catalog.relate('shameless-uk-christmas-special', 'side_story_of', 'shameless-uk', p_note => 'The Christmas special is already included at ordinal 8 in the original series'' 139-episode count.');
update catalog.entry set status = 'finished' where id = catalog.id('shameless-uk-christmas-special');
select catalog.add_entry('shameless-uk-new-year-special', 'tv', 'special', 'Shameless: New Year Special', '2006-01-03', 'GB', 'none', null, p_public_id => 1448);
select catalog.member('shameless', 'shameless-uk-new-year-special');
select catalog.relate('shameless-uk-new-year-special', 'side_story_of', 'shameless-uk', p_note => 'The first episode of series 3, included at ordinal 19 in the original series'' episode count.');
update catalog.entry set status = 'finished' where id = catalog.id('shameless-uk-new-year-special');

select catalog.add_units('shameless-uk', 139);
update catalog.unit set title = 'Christmas in Chatsworth' where entry_id = catalog.id('shameless-uk') and ordinal = 8;
update catalog.unit set title = 'New Year Special' where entry_id = catalog.id('shameless-uk') and ordinal = 19;
