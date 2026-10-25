-- Ben 10: a film from an alternate timeline that a later series picks back
-- up, so it watches better later than its release slot.

select catalog.add_group('ben-10', 'franchise', 'Ben 10', p_public_id => 421);
select catalog.add_group('ben-10-classic', 'continuity', 'Classic', 'concluded', p_public_id => 422);
select catalog.add_group('ben-10-alt', 'continuity', 'Race Against Time timeline', 'concluded', p_public_id => 423);
select catalog.add_group('ben-10-reboot', 'continuity', '2016 reboot', 'concluded', p_public_id => 424);
select catalog.member('ben-10', 'ben-10-classic', 'ben-10-alt', 'ben-10-reboot');

-- Episode totals exclude separately released films. The reboot debuted abroad.
-- https://en.wikipedia.org/wiki/Ben_10_(2005_TV_series)
-- https://en.wikipedia.org/wiki/Ben_10_(2016_TV_series)
select catalog.add_entry('ben-10-2005', 'tv', 'series', 'Ben 10', '2005-12-27', 'US', 'episode', 49, p_public_id => 425);
select catalog.add_entry('ben-10-race-against-time', 'film', 'tv_movie', 'Ben 10: Race Against Time', '2007-11-21', 'US', p_public_id => 426);
select catalog.add_entry('ben-10-alien-force', 'tv', 'series', 'Ben 10: Alien Force', '2008-04-18', 'US', 'episode', 46, p_public_id => 427);
select catalog.add_entry('ben-10-ultimate-alien', 'tv', 'series', 'Ben 10: Ultimate Alien', '2010-04-23', 'US', 'episode', 52, p_public_id => 428);
select catalog.add_entry('ben-10-omniverse', 'tv', 'series', 'Ben 10: Omniverse', '2012-08', 'US', 'episode', 80, p_public_id => 429);
select catalog.add_entry('ben-10-2016', 'tv', 'series', 'Ben 10', '2016-10-01', 'US', 'episode', 166, p_public_id => 430);
select catalog.member('ben-10-classic', 'ben-10-2005', 'ben-10-alien-force', 'ben-10-ultimate-alien', 'ben-10-omniverse');
select catalog.member('ben-10-alt', 'ben-10-race-against-time');
select catalog.member('ben-10-reboot', 'ben-10-2016');

select catalog.ext('ben-10-2005', 'wikidata', 'Q207746');
select catalog.ext('ben-10-omniverse', 'wikidata', 'Q113346');
select catalog.ext('ben-10-2016', 'wikidata', 'Q27061887');
select catalog.ext('ben-10', 'wikidata', 'Q5594233');

select catalog.relate('ben-10-alien-force', 'sequel_of', 'ben-10-2005');
select catalog.relate('ben-10-ultimate-alien', 'sequel_of', 'ben-10-alien-force');
select catalog.relate('ben-10-omniverse', 'sequel_of', 'ben-10-ultimate-alien');
select catalog.relate('ben-10-2016', 'reboot_of', 'ben-10-2005');
select catalog.relate('ben-10-ultimate-alien', 'references', 'ben-10-race-against-time', p_spoiler => true,
                      p_note => 'the third series picks the film back up');

-- https://en.wikipedia.org/wiki/Ben_10
-- https://en.wikipedia.org/wiki/List_of_Ben_10_(2016_TV_series)_episodes
-- https://en.wikipedia.org/wiki/Ben_10_(comic_series)
-- https://www.dynamite.com/htmlfiles/catalog/DynamiteDispatches-Mar4May26.pdf
-- https://www.comics.org/feature/2087/issues/?display=grid&sort=issue
select catalog.add_entry('ben-10-secret-of-the-omnitrix', 'film', 'tv_movie', 'Ben 10: Secret of the Omnitrix', '2007-08-10', 'US', 'none', null, p_public_id => 431);
select catalog.member('ben-10-classic', 'ben-10-secret-of-the-omnitrix');
select catalog.relate('ben-10-secret-of-the-omnitrix', 'side_story_of', 'ben-10-2005');
update catalog.entry set status = 'finished', length_minutes = 67 where id = catalog.id('ben-10-secret-of-the-omnitrix');
select catalog.add_entry('ben-10-alien-swarm', 'film', 'tv_movie', 'Ben 10: Alien Swarm', '2009-11-25', 'US', 'none', null, p_public_id => 432);
select catalog.member('ben-10-classic', 'ben-10-alien-swarm');
select catalog.relate('ben-10-alien-swarm', 'side_story_of', 'ben-10-alien-force');
update catalog.entry set status = 'finished', length_minutes = 69 where id = catalog.id('ben-10-alien-swarm');
select catalog.add_entry('ben-10-destroy-all-aliens', 'film', 'tv_movie', 'Ben 10: Destroy All Aliens', '2012-03-23', 'US', 'none', null, p_public_id => 433);
select catalog.member('ben-10-classic', 'ben-10-destroy-all-aliens');
select catalog.relate('ben-10-destroy-all-aliens', 'side_story_of', 'ben-10-2005');
update catalog.entry set status = 'finished', length_minutes = 69 where id = catalog.id('ben-10-destroy-all-aliens');
select catalog.add_entry('ben-10-versus-the-universe', 'film', 'tv_movie', 'Ben 10 Versus the Universe: The Movie', '2020-10-10', 'US', 'none', null, p_public_id => 434);
select catalog.member('ben-10-reboot', 'ben-10-versus-the-universe');
select catalog.relate('ben-10-versus-the-universe', 'side_story_of', 'ben-10-2016');
update catalog.entry set status = 'finished', length_minutes = 72 where id = catalog.id('ben-10-versus-the-universe');
select catalog.add_entry('ben-10-generator-rex-heroes-united', 'tv', 'special', 'Ben 10/Generator Rex: Heroes United', '2011-11-25', 'US', 'none', null, p_public_id => 435);
select catalog.member('ben-10-classic', 'ben-10-generator-rex-heroes-united');
select catalog.relate('ben-10-generator-rex-heroes-united', 'crossover_with', 'ben-10-ultimate-alien');
update catalog.entry set status = 'finished', length_minutes = 44 where id = catalog.id('ben-10-generator-rex-heroes-united');
select catalog.add_entry('ben-10-010', 'tv', 'special', 'Ben 10,010', '2021-04-09', 'US', 'none', null, p_public_id => 436);
select catalog.member('ben-10-reboot', 'ben-10-010');
select catalog.relate('ben-10-010', 'side_story_of', 'ben-10-2016', p_note => 'Also recorded as an extra unit on the 2016 series; not one of its 166 regular episodes.');
update catalog.entry set status = 'finished', length_minutes = 44 where id = catalog.id('ben-10-010');
select catalog.add_entry('ben-gen-10', 'tv', 'special', 'Ben Gen 10', '2021-04-10', 'US', 'none', null, p_public_id => 437);
select catalog.member('ben-10-reboot', 'ben-gen-10');
select catalog.relate('ben-gen-10', 'crossover_with', 'ben-10-2016', p_note => 'Also recorded as an extra unit on the reboot. Features a reboot version of Generator Rex.');
update catalog.entry set status = 'finished', length_minutes = 44 where id = catalog.id('ben-gen-10');
select catalog.add_entry('ben-10-alien-x-tinction', 'tv', 'special', 'Alien X-Tinction', '2021-04-11', 'US', 'none', null, p_public_id => 438);
select catalog.member('ben-10-reboot', 'ben-10-alien-x-tinction');
select catalog.relate('ben-10-alien-x-tinction', 'crossover_with', 'ben-10-2016', p_note => 'Also recorded as an extra unit on the reboot; crosses over multiple incarnations of Ben.');
update catalog.entry set status = 'finished', length_minutes = 44 where id = catalog.id('ben-10-alien-x-tinction');
select catalog.add_entry('ben-10-hyperscan', 'game', 'full', 'Ben 10 (HyperScan)', '2006-10-23', 'US', 'none', null, p_public_id => 439);
select catalog.member('ben-10-classic', 'ben-10-hyperscan');
select catalog.relate('ben-10-hyperscan', 'adaptation_of', 'ben-10-2005');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-hyperscan');
select catalog.add_entry('ben-10-protector-of-earth', 'game', 'full', 'Ben 10: Protector of Earth', '2007-10-30', 'US', 'none', null, p_public_id => 440);
select catalog.member('ben-10-classic', 'ben-10-protector-of-earth');
select catalog.relate('ben-10-protector-of-earth', 'adaptation_of', 'ben-10-2005');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-protector-of-earth');
select catalog.add_entry('ben-10-alien-force-game', 'game', 'full', 'Ben 10: Alien Force', '2008-10-28', 'US', 'none', null, p_public_id => 441);
select catalog.member('ben-10-classic', 'ben-10-alien-force-game');
select catalog.relate('ben-10-alien-force-game', 'adaptation_of', 'ben-10-alien-force');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-alien-force-game');
select catalog.add_entry('ben-10-vilgax-attacks', 'game', 'full', 'Ben 10 Alien Force: Vilgax Attacks', '2009-10-27', 'US', 'none', null, p_public_id => 442);
select catalog.member('ben-10-classic', 'ben-10-vilgax-attacks');
select catalog.relate('ben-10-vilgax-attacks', 'adaptation_of', 'ben-10-alien-force');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-vilgax-attacks');
select catalog.add_entry('ben-10-rise-of-hex', 'game', 'full', 'Ben 10 Alien Force: The Rise of Hex', '2010-05-31', 'US', 'none', null, p_public_id => 443);
select catalog.member('ben-10-classic', 'ben-10-rise-of-hex');
select catalog.relate('ben-10-rise-of-hex', 'adaptation_of', 'ben-10-alien-force');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-rise-of-hex');
select catalog.add_entry('ben-10-cosmic-destruction', 'game', 'full', 'Ben 10 Ultimate Alien: Cosmic Destruction', '2010-10-05', 'US', 'none', null, p_public_id => 444);
select catalog.member('ben-10-classic', 'ben-10-cosmic-destruction');
select catalog.relate('ben-10-cosmic-destruction', 'adaptation_of', 'ben-10-ultimate-alien');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-cosmic-destruction');
select catalog.add_entry('ben-10-galactic-racing', 'game', 'full', 'Ben 10: Galactic Racing', '2011-10-18', 'US', 'none', null, p_public_id => 445);
select catalog.member('ben-10-classic', 'ben-10-galactic-racing');
select catalog.relate('ben-10-galactic-racing', 'spin_off_of', 'ben-10-ultimate-alien');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-galactic-racing');
select catalog.add_entry('ben-10-omniverse-game', 'game', 'full', 'Ben 10: Omniverse', '2012-11-13', 'US', 'none', null, p_public_id => 446);
select catalog.member('ben-10-classic', 'ben-10-omniverse-game');
select catalog.relate('ben-10-omniverse-game', 'adaptation_of', 'ben-10-omniverse');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-omniverse-game');
select catalog.add_entry('ben-10-omniverse-2-game', 'game', 'full', 'Ben 10: Omniverse 2', '2013-11', 'US', 'none', null, p_public_id => 447);
select catalog.member('ben-10-classic', 'ben-10-omniverse-2-game');
select catalog.relate('ben-10-omniverse-2-game', 'sequel_of', 'ben-10-omniverse-game');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-omniverse-2-game');
select catalog.add_entry('ben-10-all-out-attack', 'game', 'full', 'Ben 10: All Out Attack!', '2008-11-04', 'US', 'none', null, p_public_id => 448);
select catalog.member('ben-10-classic', 'ben-10-all-out-attack');
select catalog.relate('ben-10-all-out-attack', 'adaptation_of', 'ben-10-2005');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-all-out-attack');
select catalog.add_entry('ben-10-game-2017', 'game', 'full', 'Ben 10', '2017-11-10', 'US', 'none', null, p_public_id => 449);
select catalog.member('ben-10-reboot', 'ben-10-game-2017');
select catalog.relate('ben-10-game-2017', 'adaptation_of', 'ben-10-2016');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-game-2017');
select catalog.add_entry('ben-10-power-trip', 'game', 'full', 'Ben 10: Power Trip', '2020-10-09', 'US', 'none', null, p_public_id => 450);
select catalog.member('ben-10-reboot', 'ben-10-power-trip');
select catalog.relate('ben-10-power-trip', 'adaptation_of', 'ben-10-2016');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-power-trip');
select catalog.add_entry('ben-10-minecraft', 'game', 'expansion', 'Minecraft: Ben 10', '2021-06-02', 'US', 'none', null, p_public_id => 451);
select catalog.member('ben-10-reboot', 'ben-10-minecraft');
select catalog.relate('ben-10-minecraft', 'adaptation_of', 'ben-10-2016');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-minecraft');
select catalog.add_entry('ben-10-idw-comic', 'comic', 'serial', 'Ben 10 (IDW)', '2013-11', 'US', 'chapter', 4, p_public_id => 452);
select catalog.member('ben-10-classic', 'ben-10-idw-comic');
select catalog.relate('ben-10-idw-comic', 'side_story_of', 'ben-10-2005');
update catalog.entry set status = 'finished', end_date = '2014-02-01', end_precision = 'month' where id = catalog.id('ben-10-idw-comic');
select catalog.credit('ben-10-idw-comic', 'IDW Publishing', 'publisher', 'US');
select catalog.add_entry('ben-10-dynamite-comic', 'comic', 'serial', 'Ben 10 (Dynamite)', '2026-05-06', 'US', 'chapter', null, p_public_id => 453);
select catalog.member('ben-10', 'ben-10-dynamite-comic');
select catalog.relate('ben-10-dynamite-comic', 'reboot_of', 'ben-10-2005');
update catalog.entry set status = 'releasing', description = 'A new comic continuity reimagines the original road-trip premise, published by Dynamite Entertainment.' where id = catalog.id('ben-10-dynamite-comic');
select catalog.credit('ben-10-dynamite-comic', 'Dynamite Entertainment', 'publisher', 'US');
select catalog.add_entry('cartoon-network-action-pack', 'comic', 'serial', 'Cartoon Network Action Pack', '2006', 'US', 'chapter', null, p_public_id => 454);
select catalog.member('ben-10', 'cartoon-network-action-pack');
select catalog.relate('cartoon-network-action-pack', 'crossover_with', 'ben-10-2005');
update catalog.entry set status = 'finished', description = 'A Cartoon Network anthology comic containing Ben 10 stories alongside other animated properties.', end_date = '2012-01-01', end_precision = 'year' where id = catalog.id('cartoon-network-action-pack');
select catalog.credit('cartoon-network-action-pack', 'DC Comics', 'publisher', 'US');

select catalog.relate('ben-10-alien-x-tinction', 'crossover_with', 'ben-10-2005');
select catalog.relate('ben-10-alien-x-tinction', 'crossover_with', 'ben-10-alien-force');
select catalog.relate('ben-10-alien-x-tinction', 'crossover_with', 'ben-10-ultimate-alien');
select catalog.relate('ben-10-alien-x-tinction', 'crossover_with', 'ben-10-omniverse');
-- Alternate opening sequences of the same film, not three additional films.
select catalog.add_special('ben-10-secret-of-the-omnitrix', 'Red opening: Heatblast');
select catalog.add_special('ben-10-secret-of-the-omnitrix', 'Blue opening: XLR8');
select catalog.add_special('ben-10-secret-of-the-omnitrix', 'Gold opening: Eye Guy');

select catalog.add_units('ben-10-2016', 166);
select catalog.add_ordering('ben-10', 'release', 'release order',
  'Series and separately released films/specials. The three 2021 specials are shown separately from the reboot''s regular episodes; each also has a bonus-unit record on that series.') as o \gset
select catalog.add_step(:'o', 1, 'ben-10-2005');
select catalog.add_step(:'o', 2, 'ben-10-secret-of-the-omnitrix');
select catalog.add_step(:'o', 3, 'ben-10-race-against-time');
select catalog.add_step(:'o', 4, 'ben-10-alien-force');
select catalog.add_step(:'o', 5, 'ben-10-alien-swarm');
select catalog.add_step(:'o', 6, 'ben-10-ultimate-alien');
select catalog.add_step(:'o', 7, 'ben-10-generator-rex-heroes-united', p_necessity => 'optional');
select catalog.add_step(:'o', 8, 'ben-10-destroy-all-aliens');
select catalog.add_step(:'o', 9, 'ben-10-omniverse');
select catalog.add_step(:'o', 10, 'ben-10-2016', '1', '166', p_necessity => 'optional');
select catalog.add_step(:'o', 11, 'ben-10-versus-the-universe', p_necessity => 'optional');
select catalog.add_step(:'o', 12, 'ben-10-010', p_necessity => 'optional');
select catalog.add_step(:'o', 13, 'ben-gen-10', p_necessity => 'optional');
select catalog.add_step(:'o', 14, 'ben-10-alien-x-tinction', p_necessity => 'optional');

select catalog.add_ordering('ben-10', 'recommended', 'watch order',
  'The classic continuity first, including its films. Race Against Time is an optional alternate timeline. The 2016 reboot and its films/specials form a separate branch.') as o \gset
select catalog.add_step(:'o', 1, 'ben-10-2005');
select catalog.add_step(:'o', 2, 'ben-10-secret-of-the-omnitrix');
select catalog.add_step(:'o', 3, 'ben-10-destroy-all-aliens');
select catalog.add_step(:'o', 4, 'ben-10-alien-force');
select catalog.add_step(:'o', 5, 'ben-10-alien-swarm');
select catalog.add_step(:'o', 6, 'ben-10-race-against-time', p_necessity => 'optional', p_note => 'An alternate timeline introduces Eon, who returns in Ultimate Alien.');
select catalog.add_step(:'o', 7, 'ben-10-ultimate-alien');
select catalog.add_step(:'o', 8, 'ben-10-generator-rex-heroes-united', p_necessity => 'optional');
select catalog.add_step(:'o', 9, 'ben-10-omniverse');
select catalog.add_step(:'o', 10, 'ben-10-2016', '1', '166', p_necessity => 'optional');
select catalog.add_step(:'o', 11, 'ben-10-versus-the-universe', p_necessity => 'optional');
select catalog.add_step(:'o', 12, 'ben-10-010', p_necessity => 'optional');
select catalog.add_step(:'o', 13, 'ben-gen-10', p_necessity => 'optional');
select catalog.add_step(:'o', 14, 'ben-10-alien-x-tinction', p_necessity => 'optional');

-- https://www.goodreads.com/series/215672-ben-10-chapter-books
-- https://www.scholastic.com/internationalsales/sts/61311sts/ScholasticFall2010Catalog.pdf
-- https://www.boom-studios.com/archives/ben-10-for-science-first-look/
-- https://www.boom-studios.com/archives/ben-10-mecha-madness-first-look/
-- https://www.boom-studios.com/archives/vacation-with-monsters-in-ben-10-the-creature-from-serenity-shore-from-boom-studios/
-- https://en.wikipedia.org/wiki/List_of_Boom!_Studios_publications
-- https://www.comics.org/issue/2116320/
-- https://www.viz.com/press-releases/1006094
-- https://en.wikipedia.org/wiki/Ben_10:_Ultimate_Challenge
select catalog.add_entry('ben-10-book-bens-knightmare', 'book', 'novel', 'Ben 10 Alien Force: Ben''s Knightmare', '2009', 'US', 'none', null, p_public_id => 455);
select catalog.member('ben-10-classic', 'ben-10-book-bens-knightmare');
select catalog.relate('ben-10-book-bens-knightmare', 'adaptation_of', 'ben-10-alien-force');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-book-bens-knightmare');
select catalog.credit('ben-10-book-bens-knightmare', 'Scholastic', 'publisher', 'US');
select catalog.add_entry('ben-10-book-war-worlds', 'book', 'novel', 'Ben 10 Alien Force: War of the Worlds', '2009-09', 'US', 'none', null, p_public_id => 456);
select catalog.member('ben-10-classic', 'ben-10-book-war-worlds');
select catalog.relate('ben-10-book-war-worlds', 'adaptation_of', 'ben-10-alien-force');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-book-war-worlds');
select catalog.ext('ben-10-book-war-worlds', 'isbn', '9780545177153');
select catalog.credit('ben-10-book-war-worlds', 'Scholastic', 'publisher', 'US');
select catalog.add_entry('ben-10-book-triple-threat', 'book', 'novel', 'Ben 10 Alien Force: Triple Threat', '2010', 'US', 'none', null, p_public_id => 457);
select catalog.member('ben-10-classic', 'ben-10-book-triple-threat');
select catalog.relate('ben-10-book-triple-threat', 'adaptation_of', 'ben-10-alien-force');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-book-triple-threat');
select catalog.credit('ben-10-book-triple-threat', 'Scholastic', 'publisher', 'US');
select catalog.add_entry('ben-10-book-double-negative', 'book', 'novel', 'Ben 10 Alien Force: Double Negative', '2010-01', 'US', 'none', null, p_public_id => 458);
select catalog.member('ben-10-classic', 'ben-10-book-double-negative');
select catalog.relate('ben-10-book-double-negative', 'side_story_of', 'ben-10-alien-force');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-book-double-negative');
select catalog.ext('ben-10-book-double-negative', 'isbn', '9780545177160');
select catalog.credit('ben-10-book-double-negative', 'Scholastic', 'publisher', 'US');
select catalog.add_entry('ben-10-book-dark-knight', 'book', 'novel', 'Ben 10 Alien Force: The Dark of Knight', '2010-04', 'US', 'none', null, p_public_id => 459);
select catalog.member('ben-10-classic', 'ben-10-book-dark-knight');
select catalog.relate('ben-10-book-dark-knight', 'side_story_of', 'ben-10-alien-force');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-book-dark-knight');
select catalog.ext('ben-10-book-dark-knight', 'isbn', '9780545206273');
select catalog.credit('ben-10-book-dark-knight', 'Scholastic', 'publisher', 'US');
select catalog.add_entry('ben-10-book-galaxy-wars', 'book', 'novel', 'Ben 10 Alien Force: Galaxy Wars', '2010-06', 'US', 'none', null, p_public_id => 460);
select catalog.member('ben-10-classic', 'ben-10-book-galaxy-wars');
select catalog.relate('ben-10-book-galaxy-wars', 'side_story_of', 'ben-10-alien-force');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-book-galaxy-wars');
select catalog.ext('ben-10-book-galaxy-wars', 'isbn', '9780545222525');
select catalog.credit('ben-10-book-galaxy-wars', 'Scholastic', 'publisher', 'US');
select catalog.add_entry('ben-10-book-science-friction', 'book', 'novel', 'Ben 10: Science Friction', '2010', 'US', 'none', null, p_public_id => 461);
select catalog.member('ben-10-classic', 'ben-10-book-science-friction');
select catalog.relate('ben-10-book-science-friction', 'side_story_of', 'ben-10-alien-force');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-book-science-friction');
select catalog.credit('ben-10-book-science-friction', 'Scholastic', 'publisher', 'US');
select catalog.add_entry('ben-10-book-chill-day', 'book', 'novel', 'Ben 10 Ultimate Alien: Chill for a Day', '2011', 'US', 'none', null, p_public_id => 462);
select catalog.member('ben-10-classic', 'ben-10-book-chill-day');
select catalog.relate('ben-10-book-chill-day', 'side_story_of', 'ben-10-ultimate-alien');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-book-chill-day');
select catalog.credit('ben-10-book-chill-day', 'Scholastic', 'publisher', 'US');
select catalog.add_entry('ben-10-book-ultimate-guide', 'book', 'nonfiction', 'Ben 10 Ultimate Alien: The Complete Guide', '2010-09', 'US', 'none', null, p_public_id => 463);
select catalog.member('ben-10-classic', 'ben-10-book-ultimate-guide');
select catalog.relate('ben-10-book-ultimate-guide', 'references', 'ben-10-ultimate-alien');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-book-ultimate-guide');
select catalog.ext('ben-10-book-ultimate-guide', 'isbn', '9780545225380');
select catalog.credit('ben-10-book-ultimate-guide', 'Scholastic', 'publisher', 'US');
select catalog.add_entry('ben-10-omniverse-ghost-ship', 'comic', 'graphic_novel', 'Ben 10 Omniverse: Ghost Ship', '2013', 'US', 'none', null, p_public_id => 464);
select catalog.member('ben-10-classic', 'ben-10-omniverse-ghost-ship');
select catalog.relate('ben-10-omniverse-ghost-ship', 'side_story_of', 'ben-10-omniverse');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-omniverse-ghost-ship');
select catalog.ext('ben-10-omniverse-ghost-ship', 'isbn', '9781421557410');
select catalog.credit('ben-10-omniverse-ghost-ship', 'VIZ Media', 'publisher', 'US');
select catalog.add_entry('ben-10-omniverse-joyrides', 'comic', 'graphic_novel', 'Ben 10 Omniverse: Joyrides', '2014-01', 'US', 'none', null, p_public_id => 465);
select catalog.member('ben-10-classic', 'ben-10-omniverse-joyrides');
select catalog.relate('ben-10-omniverse-joyrides', 'side_story_of', 'ben-10-omniverse');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-omniverse-joyrides');
select catalog.credit('ben-10-omniverse-joyrides', 'VIZ Media', 'publisher', 'US');
select catalog.add_entry('ben-10-omniverse-parallel-paradox', 'comic', 'graphic_novel', 'Ben 10 Omniverse: Parallel Paradox', '2014-04-01', 'US', 'none', null, p_public_id => 466);
select catalog.member('ben-10-classic', 'ben-10-omniverse-parallel-paradox');
select catalog.relate('ben-10-omniverse-parallel-paradox', 'side_story_of', 'ben-10-omniverse');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-omniverse-parallel-paradox');
select catalog.credit('ben-10-omniverse-parallel-paradox', 'VIZ Media', 'publisher', 'US');
select catalog.add_entry('ben-10-truth-is-out-there', 'comic', 'graphic_novel', 'Ben 10: The Truth Is Out There', '2019-03', 'US', 'none', null, p_public_id => 467);
select catalog.member('ben-10-reboot', 'ben-10-truth-is-out-there');
select catalog.relate('ben-10-truth-is-out-there', 'side_story_of', 'ben-10-2016');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-truth-is-out-there');
select catalog.credit('ben-10-truth-is-out-there', 'BOOM! Studios', 'publisher', 'US');
select catalog.add_entry('ben-10-for-science', 'comic', 'graphic_novel', 'Ben 10: For Science!', '2019-07', 'US', 'none', null, p_public_id => 468);
select catalog.member('ben-10-reboot', 'ben-10-for-science');
select catalog.relate('ben-10-for-science', 'side_story_of', 'ben-10-2016');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-for-science');
select catalog.credit('ben-10-for-science', 'BOOM! Studios', 'publisher', 'US');
select catalog.add_entry('ben-10-mecha-madness', 'comic', 'graphic_novel', 'Ben 10: Mecha Madness', '2019-11-13', 'US', 'none', null, p_public_id => 469);
select catalog.member('ben-10-reboot', 'ben-10-mecha-madness');
select catalog.relate('ben-10-mecha-madness', 'side_story_of', 'ben-10-2016');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-mecha-madness');
select catalog.credit('ben-10-mecha-madness', 'BOOM! Studios', 'publisher', 'US');
select catalog.add_entry('ben-10-manchester-mystery', 'comic', 'graphic_novel', 'Ben 10: The Manchester Mystery', '2020-03', 'US', 'none', null, p_public_id => 470);
select catalog.member('ben-10-reboot', 'ben-10-manchester-mystery');
select catalog.relate('ben-10-manchester-mystery', 'side_story_of', 'ben-10-2016');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-manchester-mystery');
select catalog.credit('ben-10-manchester-mystery', 'BOOM! Studios', 'publisher', 'US');
select catalog.add_entry('ben-10-creature-serenity-shore', 'comic', 'graphic_novel', 'Ben 10: The Creature from Serenity Shore', '2020', 'US', 'none', null, p_public_id => 471);
select catalog.member('ben-10-reboot', 'ben-10-creature-serenity-shore');
select catalog.relate('ben-10-creature-serenity-shore', 'side_story_of', 'ben-10-2016', p_note => 'Sources disagree on the July solicitation and August retail dates; only the release year is established here.');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-creature-serenity-shore');
select catalog.ext('ben-10-creature-serenity-shore', 'isbn', '9781684155323');
select catalog.credit('ben-10-creature-serenity-shore', 'BOOM! Studios', 'publisher', 'US');
select catalog.add_entry('ben-10-ultimate-challenge', 'tv', 'series', 'Ben 10: Ultimate Challenge (UK)', '2011', 'GB', 'episode', null, p_public_id => 472);
select catalog.member('ben-10', 'ben-10-ultimate-challenge');
select catalog.relate('ben-10-ultimate-challenge', 'spin_off_of', 'ben-10-classic');
update catalog.entry set status = 'finished', description = 'A licensed competition game show for children, produced by Twenty Twenty. The UK version is distinct from its regional editions.', unit_minutes = 30 where id = catalog.id('ben-10-ultimate-challenge');
select catalog.add_entry('ben-10-ultimate-challenge-asia', 'tv', 'series', 'Ben 10: Ultimate Challenge Asia', '2013', 'MY', 'episode', null, p_public_id => 473);
select catalog.member('ben-10', 'ben-10-ultimate-challenge-asia');
select catalog.relate('ben-10-ultimate-challenge-asia', 'alternative_version_of', 'ben-10-ultimate-challenge');
update catalog.entry set status = 'finished' where id = catalog.id('ben-10-ultimate-challenge-asia');
select catalog.add_entry('ben-10-challenge', 'tv', 'series', 'Ben 10 Challenge (UK)', '2017-10-13', 'GB', 'episode', null, p_public_id => 474);
select catalog.member('ben-10', 'ben-10-challenge');
select catalog.relate('ben-10-challenge', 'spin_off_of', 'ben-10-2016');
update catalog.entry set status = 'finished', description = 'A licensed live-action competition based on the reboot, produced in Madrid by La Competencia Productions for Turner EMEA.' where id = catalog.id('ben-10-challenge');

-- https://openlibrary.org/books/OL26860702M/Ben_10_Comic_Story_Book_Permanent_Retirement_(Ben_10_Comic_Book)
-- https://www.comicsreview.co.uk/nowreadthis/2009/01/06/ben-10-comic-book-vol-1-and-then-there-were-10/
select catalog.add_entry('ben-10-egmont-comic-books', 'comic', 'serial', 'Ben 10 Comic Story Books (Egmont)', '2008', 'GB', 'volume', null, p_public_id => 475);
select catalog.member('ben-10-classic', 'ben-10-egmont-comic-books');
select catalog.relate('ben-10-egmont-comic-books', 'adaptation_of', 'ben-10-2005');
update catalog.entry set status = null, description = 'The licensed comic-storybook adaptations include And Then There Were 10 and Permanent Retirement. Individual volumes retell television episodes.' where id = catalog.id('ben-10-egmont-comic-books');
select catalog.credit('ben-10-egmont-comic-books', 'Egmont', 'publisher', 'GB');
select catalog.add_entry('ben-10-alien-worlds', 'tv', 'series', 'Ben 10: Alien Worlds', '2017', 'US', 'episode', null, p_public_id => 476);
select catalog.member('ben-10-reboot', 'ben-10-alien-worlds');
select catalog.relate('ben-10-alien-worlds', 'spin_off_of', 'ben-10-2016');
update catalog.entry set status = 'finished', description = 'Animated shorts explore the home planets, history and species of the reboot aliens. The collection ran across multiple batches from 2017 to 2019.', end_date = '2019-01-01', end_precision = 'year' where id = catalog.id('ben-10-alien-worlds');
