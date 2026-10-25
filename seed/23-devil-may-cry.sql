-- Devil May Cry: a reboot in its own continuity.

select catalog.add_group('devil-may-cry-franchise', 'franchise', 'Devil May Cry', p_public_id => 477);
select catalog.add_group('dmc-mainline', 'continuity', 'Mainline', 'primary', p_public_id => 478);
select catalog.add_group('dmc-reboot', 'continuity', 'DmC', 'abandoned', p_public_id => 479);
select catalog.member('devil-may-cry-franchise', 'dmc-mainline', 'dmc-reboot');
select catalog.add_entry('devil-may-cry', 'game', 'full', 'Devil May Cry', '2001', 'JP', p_public_id => 480);
select catalog.add_entry('dmc-devil-may-cry', 'game', 'full', 'DmC: Devil May Cry', '2013-01-15', 'GB', p_public_id => 481);
select catalog.member('dmc-mainline', 'devil-may-cry');
select catalog.member('dmc-reboot', 'dmc-devil-may-cry');
select catalog.ext('devil-may-cry', 'wikidata', 'Q840889');
select catalog.ext('dmc-devil-may-cry', 'wikidata', 'Q1826997');
select catalog.relate('dmc-devil-may-cry', 'reboot_of', 'devil-may-cry');

-- https://en.wikipedia.org/wiki/Devil_May_Cry
-- https://www.netflix.com/tudum/articles/devil-may-cry-release-date-photos
-- https://en.wikipedia.org/wiki/Devil_May_Cry_(novels)
-- https://prtimes.jp/main/html/rd/p/000001158.000013450.html
select catalog.add_entry('devil-may-cry-2', 'game', 'full', 'Devil May Cry 2', '2003', 'JP', 'none', null, p_public_id => 482);
select catalog.member('dmc-mainline', 'devil-may-cry-2');
select catalog.relate('devil-may-cry-2', 'sequel_of', 'devil-may-cry');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-2');
select catalog.add_entry('devil-may-cry-3', 'game', 'full', 'Devil May Cry 3: Dante''s Awakening', '2005', 'JP', 'none', null, p_public_id => 483);
select catalog.member('dmc-mainline', 'devil-may-cry-3');
select catalog.relate('devil-may-cry-3', 'same_setting_as', 'devil-may-cry');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-3');
select catalog.add_entry('devil-may-cry-4', 'game', 'full', 'Devil May Cry 4', '2008', 'JP', 'none', null, p_public_id => 484);
select catalog.member('dmc-mainline', 'devil-may-cry-4');
select catalog.relate('devil-may-cry-4', 'sequel_of', 'devil-may-cry-2');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-4');
select catalog.add_entry('devil-may-cry-5', 'game', 'full', 'Devil May Cry 5', '2019-03-08', 'JP', 'none', null, p_public_id => 485);
select catalog.member('dmc-mainline', 'devil-may-cry-5');
select catalog.relate('devil-may-cry-5', 'sequel_of', 'devil-may-cry-4');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-5');
select catalog.add_entry('devil-may-cry-3-special-edition', 'game', 'full', 'Devil May Cry 3: Special Edition', '2006', 'JP', 'none', null, p_public_id => 486);
select catalog.member('dmc-mainline', 'devil-may-cry-3-special-edition');
select catalog.relate('devil-may-cry-3-special-edition', 'remaster_of', 'devil-may-cry-3');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-3-special-edition');
select catalog.add_entry('devil-may-cry-4-special-edition', 'game', 'full', 'Devil May Cry 4: Special Edition', '2015', 'JP', 'none', null, p_public_id => 487);
select catalog.member('dmc-mainline', 'devil-may-cry-4-special-edition');
select catalog.relate('devil-may-cry-4-special-edition', 'remaster_of', 'devil-may-cry-4');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-4-special-edition');
select catalog.add_entry('devil-may-cry-5-special-edition', 'game', 'full', 'Devil May Cry 5: Special Edition', '2020', 'JP', 'none', null, p_public_id => 488);
select catalog.member('dmc-mainline', 'devil-may-cry-5-special-edition');
select catalog.relate('devil-may-cry-5-special-edition', 'remaster_of', 'devil-may-cry-5');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-5-special-edition');
select catalog.add_entry('devil-may-cry-hd-collection', 'game', 'full', 'Devil May Cry HD Collection', '2012', 'JP', 'none', null, p_public_id => 489);
select catalog.member('dmc-mainline', 'devil-may-cry-hd-collection');
select catalog.relate('devil-may-cry-hd-collection', 'remaster_of', 'devil-may-cry');
update catalog.entry set status = 'finished', description = 'A remastered compilation of Devil May Cry, Devil May Cry 2 and Devil May Cry 3: Special Edition.' where id = catalog.id('devil-may-cry-hd-collection');
select catalog.add_entry('dmc-definitive-edition', 'game', 'full', 'DmC: Devil May Cry - Definitive Edition', '2015', 'GB', 'none', null, p_public_id => 490);
select catalog.member('dmc-reboot', 'dmc-definitive-edition');
select catalog.relate('dmc-definitive-edition', 'remaster_of', 'dmc-devil-may-cry');
update catalog.entry set status = 'finished' where id = catalog.id('dmc-definitive-edition');
select catalog.add_entry('dmc-vergils-downfall', 'game', 'expansion', 'DmC: Devil May Cry - Vergil''s Downfall', '2013-03', 'GB', 'none', null, p_public_id => 491);
select catalog.member('dmc-reboot', 'dmc-vergils-downfall');
select catalog.relate('dmc-vergils-downfall', 'sequel_of', 'dmc-devil-may-cry');
update catalog.entry set status = 'finished' where id = catalog.id('dmc-vergils-downfall');
select catalog.add_entry('devil-may-cry-deadshot', 'game', 'full', 'Devil May Cry: Deadshot', '2003', 'JP', 'none', null, p_public_id => 492);
select catalog.member('dmc-mainline', 'devil-may-cry-deadshot');
select catalog.relate('devil-may-cry-deadshot', 'adaptation_of', 'devil-may-cry');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-deadshot');
select catalog.add_entry('devil-may-cry-dante-x-vergil', 'game', 'full', 'Devil May Cry: Dante x Vergil', '2007', 'JP', 'none', null, p_public_id => 493);
select catalog.member('dmc-mainline', 'devil-may-cry-dante-x-vergil');
select catalog.relate('devil-may-cry-dante-x-vergil', 'adaptation_of', 'devil-may-cry-3');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-dante-x-vergil');
select catalog.add_entry('devil-may-cry-mobile-2008', 'game', 'full', 'Devil May Cry (mobile)', '2008', 'JP', 'none', null, p_public_id => 494);
select catalog.member('dmc-mainline', 'devil-may-cry-mobile-2008');
select catalog.relate('devil-may-cry-mobile-2008', 'adaptation_of', 'devil-may-cry-3');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-mobile-2008');
select catalog.add_entry('devil-may-cry-peak-of-combat', 'game', 'full', 'Devil May Cry: Peak of Combat', '2021', 'CN', 'none', null, p_public_id => 495);
select catalog.member('dmc-mainline', 'devil-may-cry-peak-of-combat');
select catalog.relate('devil-may-cry-peak-of-combat', 'spin_off_of', 'devil-may-cry');
update catalog.entry set status = 'finished', description = 'A licensed mobile action game developed by NebulaJoy, originally released in China before its international launch.' where id = catalog.id('devil-may-cry-peak-of-combat');
select catalog.add_entry('devil-may-cry-anime', 'anime', 'tv', 'Devil May Cry: The Animated Series', '2007-06-12', 'JP', 'episode', 12, p_public_id => 496);
select catalog.member('dmc-mainline', 'devil-may-cry-anime');
select catalog.relate('devil-may-cry-anime', 'side_story_of', 'devil-may-cry');
update catalog.entry set status = 'finished', end_date = '2007-09-01', end_precision = 'month' where id = catalog.id('devil-may-cry-anime');
select catalog.add_entry('devil-may-cry-netflix-s1', 'tv', 'series', 'Devil May Cry: Season 1', '2025-04-03', 'US', 'episode', 8, p_public_id => 497);
select catalog.member('devil-may-cry-franchise', 'devil-may-cry-netflix-s1');
select catalog.relate('devil-may-cry-netflix-s1', 'adaptation_of', 'devil-may-cry');
update catalog.entry set status = 'finished', description = 'Adi Shankar''s animated adaptation reimagines Dante and the demon-hunting world in a separate continuity, animated by Studio Mir.', end_date = '2025-04-03', end_precision = 'day' where id = catalog.id('devil-may-cry-netflix-s1');
select catalog.add_entry('devil-may-cry-netflix-s2', 'tv', 'series', 'Devil May Cry: Season 2', '2026-05-12', 'US', 'episode', 8, p_public_id => 498);
select catalog.member('devil-may-cry-franchise', 'devil-may-cry-netflix-s2');
select catalog.relate('devil-may-cry-netflix-s2', 'sequel_of', 'devil-may-cry-netflix-s1');
update catalog.entry set status = 'finished', end_date = '2026-05-12', end_precision = 'day' where id = catalog.id('devil-may-cry-netflix-s2');
select catalog.add_entry('devil-may-cry-netflix-s3', 'tv', 'series', 'Devil May Cry: Season 3', null, 'US', 'episode', null, p_public_id => 499);
select catalog.member('devil-may-cry-franchise', 'devil-may-cry-netflix-s3');
select catalog.relate('devil-may-cry-netflix-s3', 'sequel_of', 'devil-may-cry-netflix-s2');
update catalog.entry set status = 'announced', description = 'Netflix has renewed the animated adaptation for a third and final season; no release date is confirmed.' where id = catalog.id('devil-may-cry-netflix-s3');
select catalog.add_entry('devil-may-cry-novel-1', 'book', 'light_novel', 'Devil May Cry: Volume 1', '2002', 'JP', 'volume', 1, p_public_id => 500);
select catalog.member('dmc-mainline', 'devil-may-cry-novel-1');
select catalog.relate('devil-may-cry-novel-1', 'side_story_of', 'devil-may-cry');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-novel-1');
select catalog.add_units('devil-may-cry-novel-1', 1);
select catalog.add_entry('devil-may-cry-novel-2', 'book', 'light_novel', 'Devil May Cry: Volume 2', '2003', 'JP', 'volume', 1, p_public_id => 501);
select catalog.member('dmc-mainline', 'devil-may-cry-novel-2');
select catalog.relate('devil-may-cry-novel-2', 'side_story_of', 'devil-may-cry-2');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-novel-2');
select catalog.add_units('devil-may-cry-novel-2', 1);
select catalog.add_entry('devil-may-cry-deadly-fortune', 'book', 'light_novel', 'Devil May Cry 4: Deadly Fortune', '2009', 'JP', 'volume', 2, p_public_id => 502);
select catalog.member('dmc-mainline', 'devil-may-cry-deadly-fortune');
select catalog.relate('devil-may-cry-deadly-fortune', 'adaptation_of', 'devil-may-cry-4');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-deadly-fortune');
select catalog.add_units('devil-may-cry-deadly-fortune', 2);
select catalog.add_entry('devil-may-cry-before-the-nightmare', 'book', 'light_novel', 'Devil May Cry 5: Before the Nightmare', '2019-03-01', 'JP', 'volume', 1, p_public_id => 503);
select catalog.member('dmc-mainline', 'devil-may-cry-before-the-nightmare');
select catalog.relate('devil-may-cry-before-the-nightmare', 'side_story_of', 'devil-may-cry-5');
update catalog.entry set status = 'finished' where id = catalog.id('devil-may-cry-before-the-nightmare');
select catalog.add_units('devil-may-cry-before-the-nightmare', 1);
select catalog.add_entry('devil-may-cry-3-manga', 'manga', 'serial', 'Devil May Cry 3', '2005', 'JP', 'volume', 2, p_public_id => 504);
select catalog.member('dmc-mainline', 'devil-may-cry-3-manga');
select catalog.relate('devil-may-cry-3-manga', 'side_story_of', 'devil-may-cry-3');
update catalog.entry set status = 'cancelled', description = 'A prequel manga with two published volumes; the planned third volume was never released.' where id = catalog.id('devil-may-cry-3-manga');
select catalog.add_units('devil-may-cry-3-manga', 2);
select catalog.add_entry('devil-may-cry-visions-of-v', 'manga', 'serial', 'Devil May Cry 5: Visions of V', '2019-03-07', 'JP', 'volume', 5, p_public_id => 505);
select catalog.member('dmc-mainline', 'devil-may-cry-visions-of-v');
select catalog.relate('devil-may-cry-visions-of-v', 'side_story_of', 'devil-may-cry-5', p_note => 'The prologue was released March 7; regular serialization began April 27.');
update catalog.entry set status = 'finished', end_date = '2022-03-19', end_precision = 'day' where id = catalog.id('devil-may-cry-visions-of-v');
select catalog.add_units('devil-may-cry-visions-of-v', 5);
select catalog.add_entry('devil-may-cry-dreamwave-comic', 'comic', 'serial', 'Devil May Cry', '2004', 'CA', 'chapter', 3, p_public_id => 506);
select catalog.member('dmc-mainline', 'devil-may-cry-dreamwave-comic');
select catalog.relate('devil-may-cry-dreamwave-comic', 'adaptation_of', 'devil-may-cry');
update catalog.entry set status = 'cancelled', description = 'Dreamwave''s comic adaptation ended after three issues without completing its story.' where id = catalog.id('devil-may-cry-dreamwave-comic');
select catalog.add_entry('devil-may-cry-live-hacker', 'film', 'stage', 'Devil May Cry: The Live Hacker', '2019-08', 'JP', 'none', null, p_public_id => 507);
select catalog.member('dmc-mainline', 'devil-may-cry-live-hacker');
select catalog.relate('devil-may-cry-live-hacker', 'side_story_of', 'devil-may-cry');
update catalog.entry set status = 'finished', description = 'The released DVD recording of Jun Yoriko''s stage musical, performed in March 2019.' where id = catalog.id('devil-may-cry-live-hacker');

select catalog.relate('devil-may-cry', 'sequel_of', 'devil-may-cry-3', p_note => 'Dante''s Awakening is a prequel; the original game follows it in story chronology.');
select catalog.relate('devil-may-cry-hd-collection', 'remaster_of', 'devil-may-cry-2');
select catalog.relate('devil-may-cry-hd-collection', 'remaster_of', 'devil-may-cry-3-special-edition');
select catalog.credit('devil-may-cry-anime', 'Madhouse', 'studio', 'JP');
select catalog.credit('devil-may-cry-netflix-s1', 'Studio Mir', 'studio', 'KR');
select catalog.credit('devil-may-cry-netflix-s2', 'Studio Mir', 'studio', 'KR');

select catalog.add_ordering('devil-may-cry-franchise', 'recommended', 'main games in story order',
  'The current main continuity runs 3, 1, 2, 4, 5. Special editions are alternatives to the same games, and DmC belongs to a separate reboot.') as o \gset
select catalog.add_step(:'o', 1, 'devil-may-cry-3');
select catalog.add_step(:'o', 2, 'devil-may-cry');
select catalog.add_step(:'o', 3, 'devil-may-cry-2');
select catalog.add_step(:'o', 4, 'devil-may-cry-4');
select catalog.add_step(:'o', 5, 'devil-may-cry-5');

-- https://www.dmc-tv.com/goods/
-- https://www.cdjapan.co.jp/product/FCCM-214
select catalog.add_entry('devil-may-cry-drama-cd', 'audio', 'drama', 'TV Animation Devil May Cry: Drama CD', '2007-12-21', 'JP', 'volume', 2, p_public_id => 508);
select catalog.member('dmc-mainline', 'devil-may-cry-drama-cd');
select catalog.relate('devil-may-cry-drama-cd', 'side_story_of', 'devil-may-cry-anime');
update catalog.entry set status = 'finished', end_date = '2008-01-01', end_precision = 'year' where id = catalog.id('devil-may-cry-drama-cd');
select catalog.add_units('devil-may-cry-drama-cd', 2);

-- https://www.engadget.com/2013-01-24-dmc-prequel-comic-the-chronicles-of-vergil-unleashes-first-dig.html
-- https://www.4gamer.net/games/288/G028872/20151009036/
-- https://www.suruga-ya.jp/product/detail/151094965
select catalog.add_entry('dmc-vergil-chronicles', 'comic', 'serial', 'DmC: Devil May Cry - The Vergil Chronicles', '2013-01', 'FR', 'chapter', 2, p_public_id => 509);
select catalog.member('dmc-reboot', 'dmc-vergil-chronicles');
select catalog.relate('dmc-vergil-chronicles', 'side_story_of', 'dmc-devil-may-cry', p_note => 'A prequel about Vergil and Kat in the DmC reboot continuity.');
update catalog.entry set status = 'finished' where id = catalog.id('dmc-vergil-chronicles');
select catalog.credit('dmc-vergil-chronicles', 'Titan Comics', 'publisher', 'FR');
select catalog.add_entry('sengoku-basara-vs-devil-may-cry', 'film', 'stage', 'Sengoku BASARA vs Devil May Cry', '2015-11-11', 'JP', 'none', null, p_public_id => 510);
select catalog.member('devil-may-cry-franchise', 'sengoku-basara-vs-devil-may-cry');
select catalog.relate('sengoku-basara-vs-devil-may-cry', 'crossover_with', 'devil-may-cry-4', p_note => 'Recorded stage crossover; date is the home-video release, not the August performance.');
update catalog.entry set status = 'finished', length_minutes = 183 where id = catalog.id('sengoku-basara-vs-devil-may-cry');
