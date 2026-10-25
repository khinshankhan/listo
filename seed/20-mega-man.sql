-- Mega Man: the games, their screen adaptations, and the manga and comics.

-- Series groups include games and their adaptations, not a single shared continuity.
select catalog.add_group('mega-man', 'franchise', 'Mega Man', p_public_id => 1);
select catalog.add_group('mega-man-classic', 'series', 'Mega Man Classic', p_public_id => 2);
select catalog.add_group('mega-man-x', 'series', 'Mega Man X', p_public_id => 3);
select catalog.add_group('mega-man-zero', 'series', 'Mega Man Zero', p_public_id => 4);
select catalog.add_group('mega-man-zx', 'series', 'Mega Man ZX', p_public_id => 5);
select catalog.add_group('mega-man-legends', 'series', 'Mega Man Legends', p_public_id => 6);
select catalog.add_group('mega-man-battle-network', 'series', 'Mega Man Battle Network', p_public_id => 7);
select catalog.add_group('mega-man-star-force', 'series', 'Mega Man Star Force', p_public_id => 8);
select catalog.add_group('mega-man-fan-games', 'series', 'Mega Man fan games', p_public_id => 9);
select catalog.alias('mega-man', 'native', 'ロックマン');
select catalog.alias('mega-man', 'synonym', 'Rockman');
select catalog.member('mega-man', 'mega-man-classic', 'mega-man-x', 'mega-man-zero', 'mega-man-zx', 'mega-man-legends', 'mega-man-battle-network', 'mega-man-star-force', 'mega-man-fan-games');

-- Official games.
-- https://www.mobygames.com/group/586/mega-man-rockman-franchise/
-- https://www.mobygames.com/group/193/mega-man-rockman-series/
-- https://en.wikipedia.org/wiki/List_of_Mega_Man_video_games
select catalog.add_entry('mega-man-1', 'game', 'full', 'Mega Man', '1987', 'JP', 'none', null, p_public_id => 10);
select catalog.member('mega-man-classic', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-1');
select catalog.ext('mega-man-1', 'mobygames', '12768');
select catalog.add_entry('mega-man-2', 'game', 'full', 'Mega Man 2', '1988', 'JP', 'none', null, p_public_id => 11);
select catalog.member('mega-man-classic', 'mega-man-2');
select catalog.relate('mega-man-2', 'sequel_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-2');
select catalog.add_entry('mega-man-3', 'game', 'full', 'Mega Man 3', '1990', 'JP', 'none', null, p_public_id => 12);
select catalog.member('mega-man-classic', 'mega-man-3');
select catalog.relate('mega-man-3', 'sequel_of', 'mega-man-2');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-3');
select catalog.add_entry('mega-man-4', 'game', 'full', 'Mega Man 4', '1991', 'JP', 'none', null, p_public_id => 13);
select catalog.member('mega-man-classic', 'mega-man-4');
select catalog.relate('mega-man-4', 'sequel_of', 'mega-man-3');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-4');
select catalog.add_entry('mega-man-5', 'game', 'full', 'Mega Man 5', '1992', 'JP', 'none', null, p_public_id => 14);
select catalog.member('mega-man-classic', 'mega-man-5');
select catalog.relate('mega-man-5', 'sequel_of', 'mega-man-4');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-5');
select catalog.add_entry('mega-man-6', 'game', 'full', 'Mega Man 6', '1993', 'JP', 'none', null, p_public_id => 15);
select catalog.member('mega-man-classic', 'mega-man-6');
select catalog.relate('mega-man-6', 'sequel_of', 'mega-man-5');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-6');
select catalog.add_entry('mega-man-7', 'game', 'full', 'Mega Man 7', '1995', 'JP', 'none', null, p_public_id => 16);
select catalog.member('mega-man-classic', 'mega-man-7');
select catalog.relate('mega-man-7', 'sequel_of', 'mega-man-6');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-7');
select catalog.add_entry('mega-man-8', 'game', 'full', 'Mega Man 8', '1996', 'JP', 'none', null, p_public_id => 17);
select catalog.member('mega-man-classic', 'mega-man-8');
select catalog.relate('mega-man-8', 'sequel_of', 'mega-man-7');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-8');
select catalog.add_entry('mega-man-9', 'game', 'full', 'Mega Man 9', '2008', 'JP', 'none', null, p_public_id => 18);
select catalog.member('mega-man-classic', 'mega-man-9');
select catalog.relate('mega-man-9', 'sequel_of', 'mega-man-8');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-9');
select catalog.add_entry('mega-man-10', 'game', 'full', 'Mega Man 10', '2010', 'JP', 'none', null, p_public_id => 19);
select catalog.member('mega-man-classic', 'mega-man-10');
select catalog.relate('mega-man-10', 'sequel_of', 'mega-man-9');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-10');
select catalog.add_entry('mega-man-11', 'game', 'full', 'Mega Man 11', '2018', 'JP', 'none', null, p_public_id => 20);
select catalog.member('mega-man-classic', 'mega-man-11');
select catalog.relate('mega-man-11', 'sequel_of', 'mega-man-10');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-11');
select catalog.ext('mega-man-11', 'steam', '742320');
select catalog.add_entry('mega-man-dual-override', 'game', 'full', 'Mega Man: Dual Override', '2027', 'JP', 'none', null, p_public_id => 21);
select catalog.member('mega-man-classic', 'mega-man-dual-override');
select catalog.relate('mega-man-dual-override', 'sequel_of', 'mega-man-11');
update catalog.entry set status = 'announced', description = 'The announced twelfth mainline Mega Man game, planned for 2027. The exact release date is unconfirmed.', source_type = 'original' where id = catalog.id('mega-man-dual-override');
select catalog.add_entry('mega-man-wilys-revenge', 'game', 'full', 'Mega Man: Dr. Wily''s Revenge', '1991', 'JP', 'none', null, p_public_id => 22);
select catalog.member('mega-man-classic', 'mega-man-wilys-revenge');
select catalog.relate('mega-man-wilys-revenge', 'spin_off_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-wilys-revenge');
select catalog.add_entry('mega-man-ii-gb', 'game', 'full', 'Mega Man II (Game Boy)', '1991', 'JP', 'none', null, p_public_id => 23);
select catalog.member('mega-man-classic', 'mega-man-ii-gb');
select catalog.relate('mega-man-ii-gb', 'sequel_of', 'mega-man-wilys-revenge');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-ii-gb');
select catalog.add_entry('mega-man-iii-gb', 'game', 'full', 'Mega Man III (Game Boy)', '1992', 'JP', 'none', null, p_public_id => 24);
select catalog.member('mega-man-classic', 'mega-man-iii-gb');
select catalog.relate('mega-man-iii-gb', 'sequel_of', 'mega-man-ii-gb');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-iii-gb');
select catalog.add_entry('mega-man-iv-gb', 'game', 'full', 'Mega Man IV (Game Boy)', '1993', 'JP', 'none', null, p_public_id => 25);
select catalog.member('mega-man-classic', 'mega-man-iv-gb');
select catalog.relate('mega-man-iv-gb', 'sequel_of', 'mega-man-iii-gb');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-iv-gb');
select catalog.add_entry('mega-man-v-gb', 'game', 'full', 'Mega Man V (Game Boy)', '1994', 'JP', 'none', null, p_public_id => 26);
select catalog.member('mega-man-classic', 'mega-man-v-gb');
select catalog.relate('mega-man-v-gb', 'sequel_of', 'mega-man-iv-gb');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-v-gb');
select catalog.add_entry('mega-man-wily-wars', 'game', 'full', 'Mega Man: The Wily Wars', '1994', 'JP', 'none', null, p_public_id => 27);
select catalog.member('mega-man-classic', 'mega-man-wily-wars');
select catalog.relate('mega-man-wily-wars', 'remake_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-wily-wars');
select catalog.add_entry('mega-man-bass', 'game', 'full', 'Mega Man & Bass', '1998', 'JP', 'none', null, p_public_id => 28);
select catalog.member('mega-man-classic', 'mega-man-bass');
select catalog.relate('mega-man-bass', 'spin_off_of', 'mega-man-8');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bass');
select catalog.add_entry('mega-man-bass-challenger', 'game', 'full', 'Rockman & Forte: Challenger from the Future', '1999', 'JP', 'none', null, p_public_id => 29);
select catalog.member('mega-man-classic', 'mega-man-bass-challenger');
select catalog.relate('mega-man-bass-challenger', 'sequel_of', 'mega-man-bass');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bass-challenger');
select catalog.add_entry('mega-man-game-gear', 'game', 'full', 'Mega Man (Game Gear)', '1995', 'JP', 'none', null, p_public_id => 30);
select catalog.member('mega-man-classic', 'mega-man-game-gear');
select catalog.relate('mega-man-game-gear', 'spin_off_of', 'mega-man-4');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-game-gear');
select catalog.add_entry('mega-man-dos', 'game', 'full', 'Mega Man (DOS)', '1990', 'US', 'none', null, p_public_id => 31);
select catalog.member('mega-man-classic', 'mega-man-dos');
select catalog.relate('mega-man-dos', 'adaptation_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-dos');
select catalog.credit('mega-man-dos', 'Hi-Tech Expressions', 'publisher', 'US');
select catalog.add_entry('mega-man-iii-dos', 'game', 'full', 'Mega Man III (DOS)', '1992', 'US', 'none', null, p_public_id => 32);
select catalog.member('mega-man-classic', 'mega-man-iii-dos');
select catalog.relate('mega-man-iii-dos', 'sequel_of', 'mega-man-dos');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-iii-dos');
select catalog.credit('mega-man-iii-dos', 'Hi-Tech Expressions', 'publisher', 'US');
select catalog.add_entry('mega-man-powered-up', 'game', 'full', 'Mega Man Powered Up', '2006', 'JP', 'none', null, p_public_id => 33);
select catalog.member('mega-man-classic', 'mega-man-powered-up');
select catalog.relate('mega-man-powered-up', 'remake_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-powered-up');
select catalog.add_entry('mega-man-rockboard', 'game', 'full', 'Wily & Right no RockBoard: That''s Paradise', '1993', 'JP', 'none', null, p_public_id => 34);
select catalog.member('mega-man-classic', 'mega-man-rockboard');
select catalog.relate('mega-man-rockboard', 'spin_off_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-rockboard');
select catalog.add_entry('mega-man-soccer', 'game', 'full', 'Mega Man Soccer', '1994', 'JP', 'none', null, p_public_id => 35);
select catalog.member('mega-man-classic', 'mega-man-soccer');
select catalog.relate('mega-man-soccer', 'spin_off_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-soccer');
select catalog.add_entry('mega-man-battle-chase', 'game', 'full', 'Mega Man: Battle & Chase', '1997', 'JP', 'none', null, p_public_id => 36);
select catalog.member('mega-man-classic', 'mega-man-battle-chase');
select catalog.relate('mega-man-battle-chase', 'spin_off_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-battle-chase');
select catalog.add_entry('mega-man-power-battle', 'game', 'full', 'Mega Man: The Power Battle', '1995', 'JP', 'none', null, p_public_id => 37);
select catalog.member('mega-man-classic', 'mega-man-power-battle');
select catalog.relate('mega-man-power-battle', 'spin_off_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-power-battle');
select catalog.add_entry('mega-man-power-fighters', 'game', 'full', 'Mega Man 2: The Power Fighters', '1996', 'JP', 'none', null, p_public_id => 38);
select catalog.member('mega-man-classic', 'mega-man-power-fighters');
select catalog.relate('mega-man-power-fighters', 'sequel_of', 'mega-man-power-battle');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-power-fighters');
select catalog.add_entry('mega-man-battle-fighters', 'game', 'full', 'Mega Man: Battle & Fighters', '2000', 'JP', 'none', null, p_public_id => 39);
select catalog.member('mega-man-classic', 'mega-man-battle-fighters');
select catalog.relate('mega-man-battle-fighters', 'remake_of', 'mega-man-power-battle');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-battle-fighters');
select catalog.add_entry('super-adventure-rockman', 'game', 'full', 'Super Adventure Rockman', '1998', 'JP', 'none', null, p_public_id => 40);
select catalog.member('mega-man-classic', 'super-adventure-rockman');
select catalog.relate('super-adventure-rockman', 'spin_off_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('super-adventure-rockman');
select catalog.add_entry('rockman-strategy', 'game', 'full', 'Rockman Strategy', '2001', 'TW', 'none', null, p_public_id => 41);
select catalog.member('mega-man-classic', 'rockman-strategy');
select catalog.relate('rockman-strategy', 'spin_off_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('rockman-strategy');
select catalog.add_entry('street-fighter-x-mega-man', 'game', 'full', 'Street Fighter X Mega Man', '2012-12-17', 'SG', 'none', null, p_public_id => 42);
select catalog.member('mega-man-classic', 'street-fighter-x-mega-man');
select catalog.relate('street-fighter-x-mega-man', 'spin_off_of', 'mega-man-1');
update catalog.entry set status = 'finished', description = 'A fan-origin crossover developed by Seow Zong Hui and officially published by Capcom. It is cataloged with the licensed games.', source_type = 'original' where id = catalog.id('street-fighter-x-mega-man');
select catalog.add_entry('mega-man-x-game', 'game', 'full', 'Mega Man X', '1993', 'JP', 'none', null, p_public_id => 43);
select catalog.member('mega-man-x', 'mega-man-x-game');
select catalog.relate('mega-man-x-game', 'spin_off_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-x-game');
select catalog.add_entry('mega-man-x2', 'game', 'full', 'Mega Man X2', '1994', 'JP', 'none', null, p_public_id => 44);
select catalog.member('mega-man-x', 'mega-man-x2');
select catalog.relate('mega-man-x2', 'sequel_of', 'mega-man-x-game');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-x2');
select catalog.add_entry('mega-man-x3', 'game', 'full', 'Mega Man X3', '1995', 'JP', 'none', null, p_public_id => 45);
select catalog.member('mega-man-x', 'mega-man-x3');
select catalog.relate('mega-man-x3', 'sequel_of', 'mega-man-x2');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-x3');
select catalog.add_entry('mega-man-x4', 'game', 'full', 'Mega Man X4', '1997', 'JP', 'none', null, p_public_id => 46);
select catalog.member('mega-man-x', 'mega-man-x4');
select catalog.relate('mega-man-x4', 'sequel_of', 'mega-man-x3');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-x4');
select catalog.add_entry('mega-man-x5', 'game', 'full', 'Mega Man X5', '2000', 'JP', 'none', null, p_public_id => 47);
select catalog.member('mega-man-x', 'mega-man-x5');
select catalog.relate('mega-man-x5', 'sequel_of', 'mega-man-x4');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-x5');
select catalog.add_entry('mega-man-x6', 'game', 'full', 'Mega Man X6', '2001', 'JP', 'none', null, p_public_id => 48);
select catalog.member('mega-man-x', 'mega-man-x6');
select catalog.relate('mega-man-x6', 'sequel_of', 'mega-man-x5');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-x6');
select catalog.add_entry('mega-man-x7', 'game', 'full', 'Mega Man X7', '2003', 'JP', 'none', null, p_public_id => 49);
select catalog.member('mega-man-x', 'mega-man-x7');
select catalog.relate('mega-man-x7', 'sequel_of', 'mega-man-x6');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-x7');
select catalog.add_entry('mega-man-x8', 'game', 'full', 'Mega Man X8', '2004', 'JP', 'none', null, p_public_id => 50);
select catalog.member('mega-man-x', 'mega-man-x8');
select catalog.relate('mega-man-x8', 'sequel_of', 'mega-man-x7');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-x8');
select catalog.add_entry('mega-man-xtreme', 'game', 'full', 'Mega Man Xtreme', '2000', 'JP', 'none', null, p_public_id => 51);
select catalog.member('mega-man-x', 'mega-man-xtreme');
select catalog.relate('mega-man-xtreme', 'spin_off_of', 'mega-man-x-game');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-xtreme');
select catalog.add_entry('mega-man-xtreme-2', 'game', 'full', 'Mega Man Xtreme 2', '2001', 'JP', 'none', null, p_public_id => 52);
select catalog.member('mega-man-x', 'mega-man-xtreme-2');
select catalog.relate('mega-man-xtreme-2', 'sequel_of', 'mega-man-xtreme');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-xtreme-2');
select catalog.add_entry('mega-man-maverick-hunter-x', 'game', 'full', 'Mega Man Maverick Hunter X', '2005', 'JP', 'none', null, p_public_id => 53);
select catalog.member('mega-man-x', 'mega-man-maverick-hunter-x');
select catalog.relate('mega-man-maverick-hunter-x', 'remake_of', 'mega-man-x-game');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-maverick-hunter-x');
select catalog.add_entry('mega-man-x-dive', 'game', 'full', 'Rockman X DiVE', '2020', 'TW', 'none', null, p_public_id => 54);
select catalog.member('mega-man-x', 'mega-man-x-dive');
select catalog.relate('mega-man-x-dive', 'spin_off_of', 'mega-man-x-game');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-x-dive');
select catalog.add_entry('mega-man-x-dive-offline', 'game', 'full', 'Mega Man X DiVE Offline', '2023-08-31', 'TW', 'none', null, p_public_id => 55);
select catalog.member('mega-man-x', 'mega-man-x-dive-offline');
select catalog.relate('mega-man-x-dive-offline', 'alternative_version_of', 'mega-man-x-dive');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-x-dive-offline');
select catalog.ext('mega-man-x-dive-offline', 'steam', '2183650');
select catalog.add_entry('rockman-xover', 'game', 'full', 'Rockman Xover', '2012', 'JP', 'none', null, p_public_id => 56);
select catalog.member('mega-man-x', 'rockman-xover');
select catalog.relate('rockman-xover', 'spin_off_of', 'mega-man');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('rockman-xover');
select catalog.add_entry('mega-man-zero-1', 'game', 'full', 'Mega Man Zero', '2002', 'JP', 'none', null, p_public_id => 57);
select catalog.member('mega-man-zero', 'mega-man-zero-1');
select catalog.relate('mega-man-zero-1', 'spin_off_of', 'mega-man-x-game');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-zero-1');
select catalog.add_entry('mega-man-zero-2', 'game', 'full', 'Mega Man Zero 2', '2003', 'JP', 'none', null, p_public_id => 58);
select catalog.member('mega-man-zero', 'mega-man-zero-2');
select catalog.relate('mega-man-zero-2', 'sequel_of', 'mega-man-zero-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-zero-2');
select catalog.add_entry('mega-man-zero-3', 'game', 'full', 'Mega Man Zero 3', '2004', 'JP', 'none', null, p_public_id => 59);
select catalog.member('mega-man-zero', 'mega-man-zero-3');
select catalog.relate('mega-man-zero-3', 'sequel_of', 'mega-man-zero-2');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-zero-3');
select catalog.add_entry('mega-man-zero-4', 'game', 'full', 'Mega Man Zero 4', '2005', 'JP', 'none', null, p_public_id => 60);
select catalog.member('mega-man-zero', 'mega-man-zero-4');
select catalog.relate('mega-man-zero-4', 'sequel_of', 'mega-man-zero-3');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-zero-4');
select catalog.add_entry('mega-man-zx-game', 'game', 'full', 'Mega Man ZX', '2006', 'JP', 'none', null, p_public_id => 61);
select catalog.member('mega-man-zx', 'mega-man-zx-game');
select catalog.relate('mega-man-zx-game', 'spin_off_of', 'mega-man-zero-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-zx-game');
select catalog.add_entry('mega-man-zx-advent', 'game', 'full', 'Mega Man ZX Advent', '2007', 'JP', 'none', null, p_public_id => 62);
select catalog.member('mega-man-zx', 'mega-man-zx-advent');
select catalog.relate('mega-man-zx-advent', 'sequel_of', 'mega-man-zx-game');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-zx-advent');
select catalog.add_entry('mega-man-legends-game', 'game', 'full', 'Mega Man Legends', '1997', 'JP', 'none', null, p_public_id => 63);
select catalog.member('mega-man-legends', 'mega-man-legends-game');
select catalog.relate('mega-man-legends-game', 'spin_off_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-legends-game');
select catalog.add_entry('tron-bonne', 'game', 'full', 'The Misadventures of Tron Bonne', '1999', 'JP', 'none', null, p_public_id => 64);
select catalog.member('mega-man-legends', 'tron-bonne');
select catalog.relate('tron-bonne', 'side_story_of', 'mega-man-legends-game', p_note => 'A prequel set before the first Mega Man Legends.');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('tron-bonne');
select catalog.add_entry('mega-man-legends-2', 'game', 'full', 'Mega Man Legends 2', '2000', 'JP', 'none', null, p_public_id => 65);
select catalog.member('mega-man-legends', 'mega-man-legends-2');
select catalog.relate('mega-man-legends-2', 'sequel_of', 'mega-man-legends-game');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-legends-2');
select catalog.add_entry('mega-man-legends-five-islands', 'game', 'full', 'Rockman DASH: 5 Islands Adventure', '2008', 'JP', 'none', null, p_public_id => 66);
select catalog.member('mega-man-legends', 'mega-man-legends-five-islands');
select catalog.relate('mega-man-legends-five-islands', 'side_story_of', 'mega-man-legends-game');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-legends-five-islands');
select catalog.add_entry('mega-man-bn1', 'game', 'full', 'Mega Man Battle Network', '2001', 'JP', 'none', null, p_public_id => 67);
select catalog.member('mega-man-battle-network', 'mega-man-bn1');
select catalog.relate('mega-man-bn1', 'spin_off_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bn1');
select catalog.add_entry('mega-man-bn2', 'game', 'full', 'Mega Man Battle Network 2', '2001', 'JP', 'none', null, p_public_id => 68);
select catalog.member('mega-man-battle-network', 'mega-man-bn2');
select catalog.relate('mega-man-bn2', 'sequel_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bn2');
select catalog.add_entry('mega-man-bn3-white', 'game', 'full', 'Mega Man Battle Network 3: White', '2002', 'JP', 'none', null, p_public_id => 69);
select catalog.member('mega-man-battle-network', 'mega-man-bn3-white');
select catalog.relate('mega-man-bn3-white', 'sequel_of', 'mega-man-bn2');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bn3-white');
select catalog.add_entry('mega-man-bn3-blue', 'game', 'full', 'Mega Man Battle Network 3: Blue', '2003', 'JP', 'none', null, p_public_id => 70);
select catalog.member('mega-man-battle-network', 'mega-man-bn3-blue');
select catalog.relate('mega-man-bn3-blue', 'alternative_version_of', 'mega-man-bn3-white');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bn3-blue');
select catalog.add_entry('mega-man-bn4-red-sun', 'game', 'full', 'Mega Man Battle Network 4: Red Sun', '2003', 'JP', 'none', null, p_public_id => 71);
select catalog.member('mega-man-battle-network', 'mega-man-bn4-red-sun');
select catalog.relate('mega-man-bn4-red-sun', 'sequel_of', 'mega-man-bn3-white');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bn4-red-sun');
select catalog.add_entry('mega-man-bn4-blue-moon', 'game', 'full', 'Mega Man Battle Network 4: Blue Moon', '2003', 'JP', 'none', null, p_public_id => 72);
select catalog.member('mega-man-battle-network', 'mega-man-bn4-blue-moon');
select catalog.relate('mega-man-bn4-blue-moon', 'alternative_version_of', 'mega-man-bn4-red-sun');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bn4-blue-moon');
select catalog.add_entry('mega-man-bn5-protoman', 'game', 'full', 'Mega Man Battle Network 5: Team ProtoMan', '2004', 'JP', 'none', null, p_public_id => 73);
select catalog.member('mega-man-battle-network', 'mega-man-bn5-protoman');
select catalog.relate('mega-man-bn5-protoman', 'sequel_of', 'mega-man-bn4-red-sun');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bn5-protoman');
select catalog.add_entry('mega-man-bn5-colonel', 'game', 'full', 'Mega Man Battle Network 5: Team Colonel', '2005', 'JP', 'none', null, p_public_id => 74);
select catalog.member('mega-man-battle-network', 'mega-man-bn5-colonel');
select catalog.relate('mega-man-bn5-colonel', 'alternative_version_of', 'mega-man-bn5-protoman');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bn5-colonel');
select catalog.add_entry('mega-man-bn6-gregar', 'game', 'full', 'Mega Man Battle Network 6: Cybeast Gregar', '2005', 'JP', 'none', null, p_public_id => 75);
select catalog.member('mega-man-battle-network', 'mega-man-bn6-gregar');
select catalog.relate('mega-man-bn6-gregar', 'sequel_of', 'mega-man-bn5-protoman');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bn6-gregar');
select catalog.add_entry('mega-man-bn6-falzar', 'game', 'full', 'Mega Man Battle Network 6: Cybeast Falzar', '2005', 'JP', 'none', null, p_public_id => 76);
select catalog.member('mega-man-battle-network', 'mega-man-bn6-falzar');
select catalog.relate('mega-man-bn6-falzar', 'alternative_version_of', 'mega-man-bn6-gregar');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bn6-falzar');
select catalog.add_entry('mega-man-network-transmission', 'game', 'full', 'Mega Man Network Transmission', '2003', 'JP', 'none', null, p_public_id => 77);
select catalog.member('mega-man-battle-network', 'mega-man-network-transmission');
select catalog.relate('mega-man-network-transmission', 'side_story_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-network-transmission');
select catalog.add_entry('mega-man-battle-chip-challenge', 'game', 'full', 'Mega Man Battle Chip Challenge', '2003', 'JP', 'none', null, p_public_id => 78);
select catalog.member('mega-man-battle-network', 'mega-man-battle-chip-challenge');
select catalog.relate('mega-man-battle-chip-challenge', 'spin_off_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-battle-chip-challenge');
select catalog.add_entry('rockman-exe-n1-battle', 'game', 'full', 'Rockman EXE N1 Battle', '2003', 'JP', 'none', null, p_public_id => 79);
select catalog.member('mega-man-battle-network', 'rockman-exe-n1-battle');
select catalog.relate('rockman-exe-n1-battle', 'alternative_version_of', 'mega-man-battle-chip-challenge');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('rockman-exe-n1-battle');
select catalog.add_entry('rockman-exe-ws', 'game', 'full', 'Rockman EXE WS', '2003', 'JP', 'none', null, p_public_id => 80);
select catalog.member('mega-man-battle-network', 'rockman-exe-ws');
select catalog.relate('rockman-exe-ws', 'spin_off_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('rockman-exe-ws');
select catalog.add_entry('rockman-exe-real-operation', 'game', 'full', 'Rockman EXE 4.5: Real Operation', '2004', 'JP', 'none', null, p_public_id => 81);
select catalog.member('mega-man-battle-network', 'rockman-exe-real-operation');
select catalog.relate('rockman-exe-real-operation', 'spin_off_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('rockman-exe-real-operation');
select catalog.add_entry('mega-man-bn5-ds', 'game', 'full', 'Mega Man Battle Network 5: Double Team DS', '2005', 'JP', 'none', null, p_public_id => 82);
select catalog.member('mega-man-battle-network', 'mega-man-bn5-ds');
select catalog.relate('mega-man-bn5-ds', 'remaster_of', 'mega-man-bn5-protoman');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-bn5-ds');
select catalog.add_entry('rockman-exe-phantom-network', 'game', 'full', 'Rockman EXE: Phantom of Network', '2004', 'JP', 'none', null, p_public_id => 83);
select catalog.member('mega-man-battle-network', 'rockman-exe-phantom-network');
select catalog.relate('rockman-exe-phantom-network', 'side_story_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('rockman-exe-phantom-network');
select catalog.add_entry('rockman-exe-legend-network', 'game', 'full', 'Rockman EXE: Legend of Network', '2006', 'JP', 'none', null, p_public_id => 84);
select catalog.member('mega-man-battle-network', 'rockman-exe-legend-network');
select catalog.relate('rockman-exe-legend-network', 'sequel_of', 'rockman-exe-phantom-network');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('rockman-exe-legend-network');
select catalog.add_entry('rockman-exe-operate-shooting-star', 'game', 'full', 'Rockman EXE: Operate Shooting Star', '2009', 'JP', 'none', null, p_public_id => 85);
select catalog.member('mega-man-battle-network', 'rockman-exe-operate-shooting-star');
select catalog.relate('rockman-exe-operate-shooting-star', 'remake_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('rockman-exe-operate-shooting-star');
select catalog.add_entry('mega-man-sf-pegasus', 'game', 'full', 'Mega Man Star Force: Pegasus', '2006', 'JP', 'none', null, p_public_id => 86);
select catalog.member('mega-man-star-force', 'mega-man-sf-pegasus');
select catalog.relate('mega-man-sf-pegasus', 'spin_off_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-sf-pegasus');
select catalog.add_entry('mega-man-sf-leo', 'game', 'full', 'Mega Man Star Force: Leo', '2006', 'JP', 'none', null, p_public_id => 87);
select catalog.member('mega-man-star-force', 'mega-man-sf-leo');
select catalog.relate('mega-man-sf-leo', 'alternative_version_of', 'mega-man-sf-pegasus');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-sf-leo');
select catalog.add_entry('mega-man-sf-dragon', 'game', 'full', 'Mega Man Star Force: Dragon', '2006', 'JP', 'none', null, p_public_id => 88);
select catalog.member('mega-man-star-force', 'mega-man-sf-dragon');
select catalog.relate('mega-man-sf-dragon', 'alternative_version_of', 'mega-man-sf-pegasus');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-sf-dragon');
select catalog.add_entry('mega-man-sf2-ninja', 'game', 'full', 'Mega Man Star Force 2: Zerker x Ninja', '2007', 'JP', 'none', null, p_public_id => 89);
select catalog.member('mega-man-star-force', 'mega-man-sf2-ninja');
select catalog.relate('mega-man-sf2-ninja', 'sequel_of', 'mega-man-sf-pegasus');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-sf2-ninja');
select catalog.add_entry('mega-man-sf2-saurian', 'game', 'full', 'Mega Man Star Force 2: Zerker x Saurian', '2007', 'JP', 'none', null, p_public_id => 90);
select catalog.member('mega-man-star-force', 'mega-man-sf2-saurian');
select catalog.relate('mega-man-sf2-saurian', 'alternative_version_of', 'mega-man-sf2-ninja');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-sf2-saurian');
select catalog.add_entry('mega-man-sf3-black-ace', 'game', 'full', 'Mega Man Star Force 3: Black Ace', '2008', 'JP', 'none', null, p_public_id => 91);
select catalog.member('mega-man-star-force', 'mega-man-sf3-black-ace');
select catalog.relate('mega-man-sf3-black-ace', 'sequel_of', 'mega-man-sf2-ninja');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-sf3-black-ace');
select catalog.add_entry('mega-man-sf3-red-joker', 'game', 'full', 'Mega Man Star Force 3: Red Joker', '2008', 'JP', 'none', null, p_public_id => 92);
select catalog.member('mega-man-star-force', 'mega-man-sf3-red-joker');
select catalog.relate('mega-man-sf3-red-joker', 'alternative_version_of', 'mega-man-sf3-black-ace');
update catalog.entry set status = 'finished', source_type = 'original' where id = catalog.id('mega-man-sf3-red-joker');

select catalog.relate('mega-man-wily-wars', 'remake_of', 'mega-man-2');
select catalog.relate('mega-man-wily-wars', 'remake_of', 'mega-man-3');
select catalog.relate('mega-man-battle-fighters', 'remake_of', 'mega-man-power-fighters');
select catalog.relate('mega-man-bn5-ds', 'remaster_of', 'mega-man-bn5-colonel');
select catalog.relate('rockman-exe-operate-shooting-star', 'crossover_with', 'mega-man-star-force');
select catalog.alias('mega-man-legends-game', 'synonym', 'Mega Man 64');
select catalog.alias('mega-man-legends-game', 'romanized', 'Rockman DASH');
select catalog.alias('mega-man-bn3-blue', 'romanized', 'Rockman EXE 3 Black');

-- Compilations.
-- https://store.steampowered.com/app/363440/
-- https://store.steampowered.com/app/495050/
-- https://store.steampowered.com/app/743890/
-- https://store.steampowered.com/app/743900/
-- https://store.steampowered.com/app/999020/
-- https://store.steampowered.com/app/1798010/
-- https://store.steampowered.com/app/1798020/
-- https://store.steampowered.com/app/3500390/
select catalog.add_entry('mega-man-anniversary-collection', 'game', 'full', 'Mega Man Anniversary Collection', '2004', 'JP', 'none', null, p_public_id => 93);
select catalog.member('mega-man-classic', 'mega-man-anniversary-collection');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man, Mega Man 2, Mega Man 3, Mega Man 4, Mega Man 5, Mega Man 6, Mega Man 7, Mega Man 8, Mega Man: The Power Battle, Mega Man 2: The Power Fighters.', source_type = 'original' where id = catalog.id('mega-man-anniversary-collection');
select catalog.add_entry('mega-man-legacy-collection', 'game', 'full', 'Mega Man Legacy Collection', '2015', 'JP', 'none', null, p_public_id => 94);
select catalog.member('mega-man-classic', 'mega-man-legacy-collection');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man, Mega Man 2, Mega Man 3, Mega Man 4, Mega Man 5, Mega Man 6.', source_type = 'original' where id = catalog.id('mega-man-legacy-collection');
select catalog.ext('mega-man-legacy-collection', 'steam', '363440');
select catalog.add_entry('mega-man-legacy-collection-2', 'game', 'full', 'Mega Man Legacy Collection 2', '2017', 'JP', 'none', null, p_public_id => 95);
select catalog.member('mega-man-classic', 'mega-man-legacy-collection-2');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man 7, Mega Man 8, Mega Man 9, Mega Man 10.', source_type = 'original' where id = catalog.id('mega-man-legacy-collection-2');
select catalog.ext('mega-man-legacy-collection-2', 'steam', '495050');
select catalog.add_entry('mega-man-power-battle-fighters', 'game', 'full', 'Rockman Power Battle Fighters', '2004', 'JP', 'none', null, p_public_id => 96);
select catalog.member('mega-man-classic', 'mega-man-power-battle-fighters');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man: The Power Battle, Mega Man 2: The Power Fighters.', source_type = 'original' where id = catalog.id('mega-man-power-battle-fighters');
select catalog.add_entry('mega-man-x-collection', 'game', 'full', 'Mega Man X Collection', '2006', 'JP', 'none', null, p_public_id => 97);
select catalog.member('mega-man-x', 'mega-man-x-collection');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man X, Mega Man X2, Mega Man X3, Mega Man X4, Mega Man X5, Mega Man X6, Mega Man: Battle & Chase.', source_type = 'original' where id = catalog.id('mega-man-x-collection');
select catalog.add_entry('mega-man-x-legacy-collection', 'game', 'full', 'Mega Man X Legacy Collection', '2018', 'JP', 'none', null, p_public_id => 98);
select catalog.member('mega-man-x', 'mega-man-x-legacy-collection');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man X, Mega Man X2, Mega Man X3, Mega Man X4.', source_type = 'original' where id = catalog.id('mega-man-x-legacy-collection');
select catalog.ext('mega-man-x-legacy-collection', 'steam', '743890');
select catalog.add_entry('mega-man-x-legacy-collection-2', 'game', 'full', 'Mega Man X Legacy Collection 2', '2018', 'JP', 'none', null, p_public_id => 99);
select catalog.member('mega-man-x', 'mega-man-x-legacy-collection-2');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man X5, Mega Man X6, Mega Man X7, Mega Man X8.', source_type = 'original' where id = catalog.id('mega-man-x-legacy-collection-2');
select catalog.ext('mega-man-x-legacy-collection-2', 'steam', '743900');
select catalog.add_entry('mega-man-zero-collection', 'game', 'full', 'Mega Man Zero Collection', '2010', 'JP', 'none', null, p_public_id => 100);
select catalog.member('mega-man-zero', 'mega-man-zero-collection');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man Zero, Mega Man Zero 2, Mega Man Zero 3, Mega Man Zero 4.', source_type = 'original' where id = catalog.id('mega-man-zero-collection');
select catalog.add_entry('mega-man-zero-zx-collection', 'game', 'full', 'Mega Man Zero/ZX Legacy Collection', '2020', 'JP', 'none', null, p_public_id => 101);
select catalog.member('mega-man-zero', 'mega-man-zero-zx-collection');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man Zero, Mega Man Zero 2, Mega Man Zero 3, Mega Man Zero 4, Mega Man ZX, Mega Man ZX Advent.', source_type = 'original' where id = catalog.id('mega-man-zero-zx-collection');
select catalog.ext('mega-man-zero-zx-collection', 'steam', '999020');
select catalog.add_entry('mega-man-bn-legacy-vol1', 'game', 'full', 'Mega Man Battle Network Legacy Collection Vol. 1', '2023', 'JP', 'none', null, p_public_id => 102);
select catalog.member('mega-man-battle-network', 'mega-man-bn-legacy-vol1');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man Battle Network, Mega Man Battle Network 2, Mega Man Battle Network 3: White, Mega Man Battle Network 3: Blue.', source_type = 'original' where id = catalog.id('mega-man-bn-legacy-vol1');
select catalog.ext('mega-man-bn-legacy-vol1', 'steam', '1798010');
select catalog.add_entry('mega-man-bn-legacy-vol2', 'game', 'full', 'Mega Man Battle Network Legacy Collection Vol. 2', '2023', 'JP', 'none', null, p_public_id => 103);
select catalog.member('mega-man-battle-network', 'mega-man-bn-legacy-vol2');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man Battle Network 4: Red Sun, Mega Man Battle Network 4: Blue Moon, Mega Man Battle Network 5: Team ProtoMan, Mega Man Battle Network 5: Team Colonel, Mega Man Battle Network 6: Cybeast Gregar, Mega Man Battle Network 6: Cybeast Falzar.', source_type = 'original' where id = catalog.id('mega-man-bn-legacy-vol2');
select catalog.ext('mega-man-bn-legacy-vol2', 'steam', '1798020');
select catalog.add_entry('mega-man-sf-legacy-collection', 'game', 'full', 'Mega Man Star Force Legacy Collection', '2026-03-26', 'JP', 'none', null, p_public_id => 104);
select catalog.member('mega-man-star-force', 'mega-man-sf-legacy-collection');
update catalog.entry set status = 'finished', description = 'A compilation of Mega Man Star Force: Pegasus, Mega Man Star Force: Leo, Mega Man Star Force: Dragon, Mega Man Star Force 2: Zerker x Ninja, Mega Man Star Force 2: Zerker x Saurian, Mega Man Star Force 3: Black Ace, Mega Man Star Force 3: Red Joker.', source_type = 'original' where id = catalog.id('mega-man-sf-legacy-collection');
select catalog.ext('mega-man-sf-legacy-collection', 'steam', '3500390');
select catalog.relate('mega-man-anniversary-collection', 'compilation_of', 'mega-man-1');
select catalog.relate('mega-man-anniversary-collection', 'compilation_of', 'mega-man-2');
select catalog.relate('mega-man-anniversary-collection', 'compilation_of', 'mega-man-3');
select catalog.relate('mega-man-anniversary-collection', 'compilation_of', 'mega-man-4');
select catalog.relate('mega-man-anniversary-collection', 'compilation_of', 'mega-man-5');
select catalog.relate('mega-man-anniversary-collection', 'compilation_of', 'mega-man-6');
select catalog.relate('mega-man-anniversary-collection', 'compilation_of', 'mega-man-7');
select catalog.relate('mega-man-anniversary-collection', 'compilation_of', 'mega-man-8');
select catalog.relate('mega-man-anniversary-collection', 'compilation_of', 'mega-man-power-battle');
select catalog.relate('mega-man-anniversary-collection', 'compilation_of', 'mega-man-power-fighters');
select catalog.relate('mega-man-legacy-collection', 'compilation_of', 'mega-man-1');
select catalog.relate('mega-man-legacy-collection', 'compilation_of', 'mega-man-2');
select catalog.relate('mega-man-legacy-collection', 'compilation_of', 'mega-man-3');
select catalog.relate('mega-man-legacy-collection', 'compilation_of', 'mega-man-4');
select catalog.relate('mega-man-legacy-collection', 'compilation_of', 'mega-man-5');
select catalog.relate('mega-man-legacy-collection', 'compilation_of', 'mega-man-6');
select catalog.relate('mega-man-legacy-collection-2', 'compilation_of', 'mega-man-7');
select catalog.relate('mega-man-legacy-collection-2', 'compilation_of', 'mega-man-8');
select catalog.relate('mega-man-legacy-collection-2', 'compilation_of', 'mega-man-9');
select catalog.relate('mega-man-legacy-collection-2', 'compilation_of', 'mega-man-10');
select catalog.relate('mega-man-power-battle-fighters', 'compilation_of', 'mega-man-power-battle');
select catalog.relate('mega-man-power-battle-fighters', 'compilation_of', 'mega-man-power-fighters');
select catalog.relate('mega-man-x-collection', 'compilation_of', 'mega-man-x-game');
select catalog.relate('mega-man-x-collection', 'compilation_of', 'mega-man-x2');
select catalog.relate('mega-man-x-collection', 'compilation_of', 'mega-man-x3');
select catalog.relate('mega-man-x-collection', 'compilation_of', 'mega-man-x4');
select catalog.relate('mega-man-x-collection', 'compilation_of', 'mega-man-x5');
select catalog.relate('mega-man-x-collection', 'compilation_of', 'mega-man-x6');
select catalog.relate('mega-man-x-collection', 'compilation_of', 'mega-man-battle-chase');
select catalog.relate('mega-man-x-legacy-collection', 'compilation_of', 'mega-man-x-game');
select catalog.relate('mega-man-x-legacy-collection', 'compilation_of', 'mega-man-x2');
select catalog.relate('mega-man-x-legacy-collection', 'compilation_of', 'mega-man-x3');
select catalog.relate('mega-man-x-legacy-collection', 'compilation_of', 'mega-man-x4');
select catalog.relate('mega-man-x-legacy-collection-2', 'compilation_of', 'mega-man-x5');
select catalog.relate('mega-man-x-legacy-collection-2', 'compilation_of', 'mega-man-x6');
select catalog.relate('mega-man-x-legacy-collection-2', 'compilation_of', 'mega-man-x7');
select catalog.relate('mega-man-x-legacy-collection-2', 'compilation_of', 'mega-man-x8');
select catalog.relate('mega-man-zero-collection', 'compilation_of', 'mega-man-zero-1');
select catalog.relate('mega-man-zero-collection', 'compilation_of', 'mega-man-zero-2');
select catalog.relate('mega-man-zero-collection', 'compilation_of', 'mega-man-zero-3');
select catalog.relate('mega-man-zero-collection', 'compilation_of', 'mega-man-zero-4');
select catalog.relate('mega-man-zero-zx-collection', 'compilation_of', 'mega-man-zero-1');
select catalog.relate('mega-man-zero-zx-collection', 'compilation_of', 'mega-man-zero-2');
select catalog.relate('mega-man-zero-zx-collection', 'compilation_of', 'mega-man-zero-3');
select catalog.relate('mega-man-zero-zx-collection', 'compilation_of', 'mega-man-zero-4');
select catalog.relate('mega-man-zero-zx-collection', 'compilation_of', 'mega-man-zx-game');
select catalog.relate('mega-man-zero-zx-collection', 'compilation_of', 'mega-man-zx-advent');
select catalog.relate('mega-man-bn-legacy-vol1', 'compilation_of', 'mega-man-bn1');
select catalog.relate('mega-man-bn-legacy-vol1', 'compilation_of', 'mega-man-bn2');
select catalog.relate('mega-man-bn-legacy-vol1', 'compilation_of', 'mega-man-bn3-white');
select catalog.relate('mega-man-bn-legacy-vol1', 'compilation_of', 'mega-man-bn3-blue');
select catalog.relate('mega-man-bn-legacy-vol2', 'compilation_of', 'mega-man-bn4-red-sun');
select catalog.relate('mega-man-bn-legacy-vol2', 'compilation_of', 'mega-man-bn4-blue-moon');
select catalog.relate('mega-man-bn-legacy-vol2', 'compilation_of', 'mega-man-bn5-protoman');
select catalog.relate('mega-man-bn-legacy-vol2', 'compilation_of', 'mega-man-bn5-colonel');
select catalog.relate('mega-man-bn-legacy-vol2', 'compilation_of', 'mega-man-bn6-gregar');
select catalog.relate('mega-man-bn-legacy-vol2', 'compilation_of', 'mega-man-bn6-falzar');
select catalog.relate('mega-man-sf-legacy-collection', 'compilation_of', 'mega-man-sf-pegasus');
select catalog.relate('mega-man-sf-legacy-collection', 'compilation_of', 'mega-man-sf-leo');
select catalog.relate('mega-man-sf-legacy-collection', 'compilation_of', 'mega-man-sf-dragon');
select catalog.relate('mega-man-sf-legacy-collection', 'compilation_of', 'mega-man-sf2-ninja');
select catalog.relate('mega-man-sf-legacy-collection', 'compilation_of', 'mega-man-sf2-saurian');
select catalog.relate('mega-man-sf-legacy-collection', 'compilation_of', 'mega-man-sf3-black-ace');
select catalog.relate('mega-man-sf-legacy-collection', 'compilation_of', 'mega-man-sf3-red-joker');

-- Main games only. Remakes and compilations are alternatives, not additional sequels.
select catalog.add_ordering('mega-man-classic', 'recommended', 'main story order', 'Play one version at each numbered position. Adaptations and fan games have independent stories.') as o \gset
select catalog.add_step(:'o', 1, 'mega-man-1', p_branch => 1);
select catalog.add_step(:'o', 2, 'mega-man-2', p_branch => 1);
select catalog.add_step(:'o', 3, 'mega-man-3', p_branch => 1);
select catalog.add_step(:'o', 4, 'mega-man-4', p_branch => 1);
select catalog.add_step(:'o', 5, 'mega-man-5', p_branch => 1);
select catalog.add_step(:'o', 6, 'mega-man-6', p_branch => 1);
select catalog.add_step(:'o', 7, 'mega-man-7', p_branch => 1);
select catalog.add_step(:'o', 8, 'mega-man-8', p_branch => 1);
select catalog.add_step(:'o', 9, 'mega-man-9', p_branch => 1);
select catalog.add_step(:'o', 10, 'mega-man-10', p_branch => 1);
select catalog.add_step(:'o', 11, 'mega-man-11', p_branch => 1);
select catalog.add_ordering('mega-man-x', 'recommended', 'main story order', 'Play one version at each numbered position. Adaptations and fan games have independent stories.') as o \gset
select catalog.add_step(:'o', 1, 'mega-man-x-game', p_branch => 1);
select catalog.add_step(:'o', 2, 'mega-man-x2', p_branch => 1);
select catalog.add_step(:'o', 3, 'mega-man-x3', p_branch => 1);
select catalog.add_step(:'o', 4, 'mega-man-x4', p_branch => 1);
select catalog.add_step(:'o', 5, 'mega-man-x5', p_branch => 1);
select catalog.add_step(:'o', 6, 'mega-man-x6', p_branch => 1);
select catalog.add_step(:'o', 7, 'mega-man-x7', p_branch => 1);
select catalog.add_step(:'o', 8, 'mega-man-x8', p_branch => 1);
select catalog.add_ordering('mega-man-zero', 'recommended', 'main story order', 'Play one version at each numbered position. Adaptations and fan games have independent stories.') as o \gset
select catalog.add_step(:'o', 1, 'mega-man-zero-1', p_branch => 1);
select catalog.add_step(:'o', 2, 'mega-man-zero-2', p_branch => 1);
select catalog.add_step(:'o', 3, 'mega-man-zero-3', p_branch => 1);
select catalog.add_step(:'o', 4, 'mega-man-zero-4', p_branch => 1);
select catalog.add_ordering('mega-man-zx', 'recommended', 'main story order', 'Play one version at each numbered position. Adaptations and fan games have independent stories.') as o \gset
select catalog.add_step(:'o', 1, 'mega-man-zx-game', p_branch => 1);
select catalog.add_step(:'o', 2, 'mega-man-zx-advent', p_branch => 1);
select catalog.add_ordering('mega-man-legends', 'chronological', 'main story order', 'Play one version at each numbered position. Adaptations and fan games have independent stories.') as o \gset
select catalog.add_step(:'o', 1, 'tron-bonne', p_branch => 1);
select catalog.add_step(:'o', 2, 'mega-man-legends-game', p_branch => 1);
select catalog.add_step(:'o', 3, 'mega-man-legends-2', p_branch => 1);
select catalog.add_ordering('mega-man-battle-network', 'recommended', 'main story order', 'Play one version at each numbered position. Adaptations and fan games have independent stories.') as o \gset
select catalog.add_step(:'o', 1, 'mega-man-bn1', p_branch => 1);
select catalog.add_step(:'o', 2, 'mega-man-bn2', p_branch => 1);
select catalog.add_step(:'o', 3, 'mega-man-bn3-white', p_branch => 1);
select catalog.add_step(:'o', 3, 'mega-man-bn3-blue', p_branch => 2);
select catalog.add_step(:'o', 4, 'mega-man-bn4-red-sun', p_branch => 1);
select catalog.add_step(:'o', 4, 'mega-man-bn4-blue-moon', p_branch => 2);
select catalog.add_step(:'o', 5, 'mega-man-bn5-protoman', p_branch => 1);
select catalog.add_step(:'o', 5, 'mega-man-bn5-colonel', p_branch => 2);
select catalog.add_step(:'o', 6, 'mega-man-bn6-gregar', p_branch => 1);
select catalog.add_step(:'o', 6, 'mega-man-bn6-falzar', p_branch => 2);
select catalog.add_ordering('mega-man-star-force', 'recommended', 'main story order', 'Play one version at each numbered position. Adaptations and fan games have independent stories.') as o \gset
select catalog.add_step(:'o', 1, 'mega-man-sf-pegasus', p_branch => 1);
select catalog.add_step(:'o', 1, 'mega-man-sf-leo', p_branch => 2);
select catalog.add_step(:'o', 1, 'mega-man-sf-dragon', p_branch => 3);
select catalog.add_step(:'o', 2, 'mega-man-sf2-ninja', p_branch => 1);
select catalog.add_step(:'o', 2, 'mega-man-sf2-saurian', p_branch => 2);
select catalog.add_step(:'o', 3, 'mega-man-sf3-black-ace', p_branch => 1);
select catalog.add_step(:'o', 3, 'mega-man-sf3-red-joker', p_branch => 2);

-- Screen adaptations.
-- https://anilist.co/anime/2108/RockmanEXE-Axess
-- https://anilist.co/anime/2109/RockmanEXE-Stream/
-- https://anilist.co/anime/2110/RockmanEXE-BEAST/
-- https://anilist.co/anime/2713/RockmanEXE-Beast
-- https://anilist.co/anime/1835/Ryuusei-no-Rockman/
-- https://anilist.co/anime/3493/Ryuusei-no-Rockman-Tribe/
-- https://www.tvmaze.com/shows/35560/mega-man-fully-charged
-- https://en.wikipedia.org/wiki/Mega_Man:_Upon_a_Star
-- https://en.wikipedia.org/wiki/Mega_Man_(1994_TV_series)
-- https://en.wikipedia.org/wiki/List_of_MegaMan_NT_Warrior_media
select catalog.add_entry('mega-man-ruby-spears', 'tv', 'series', 'Mega Man (1994 TV Series)', '1994-09-11', 'US', 'episode', 27, p_public_id => 105);
select catalog.member('mega-man-classic', 'mega-man-ruby-spears');
select catalog.relate('mega-man-ruby-spears', 'adaptation_of', 'mega-man-1');
update catalog.entry set status = 'finished', description = 'The Ruby-Spears animated adaptation, with its own versions of Mega Man, Dr. Light and Dr. Wily. Two seasons contain 27 episodes.', source_type = 'game', unit_minutes = 22, end_date = '1996-01-19', end_precision = 'day' where id = catalog.id('mega-man-ruby-spears');
select catalog.add_entry('mega-man-upon-a-star', 'anime', 'ova', 'Mega Man: Upon a Star', '2002-09-20', 'JP', 'episode', 3, p_public_id => 106);
select catalog.member('mega-man-classic', 'mega-man-upon-a-star');
select catalog.relate('mega-man-upon-a-star', 'adaptation_of', 'mega-man-5');
update catalog.entry set status = 'finished', description = 'Three educational OVAs in which Mega Man and his friends visit Japan. Produced in the 1990s; the first documented home-video release is the Japanese DVD in 2002.', source_type = 'game' where id = catalog.id('mega-man-upon-a-star');
select catalog.add_entry('mega-man-fully-charged', 'tv', 'series', 'Mega Man: Fully Charged', '2018-08-05', 'US', 'episode', 52, p_public_id => 107);
select catalog.member('mega-man-classic', 'mega-man-fully-charged');
select catalog.relate('mega-man-fully-charged', 'adaptation_of', 'mega-man-1');
update catalog.entry set status = 'finished', description = 'An independent animated continuity starring Aki Light. The 52 units count individual eleven-minute stories, rather than paired broadcast blocks.', source_type = 'game', unit_minutes = 11, end_date = '2019-05-23', end_precision = 'day' where id = catalog.id('mega-man-fully-charged');
select catalog.add_entry('mega-man-day-of-sigma', 'anime', 'ova', 'Mega Man X: The Day of Sigma', '2005-12-15', 'JP', 'none', null, p_public_id => 108);
select catalog.member('mega-man-x', 'mega-man-day-of-sigma');
select catalog.relate('mega-man-day-of-sigma', 'adaptation_of', 'mega-man-maverick-hunter-x');
update catalog.entry set status = 'finished', description = 'The animated prologue bundled with Mega Man Maverick Hunter X, later included in Mega Man X Legacy Collection.', source_type = 'game', length_minutes = 25 where id = catalog.id('mega-man-day-of-sigma');
select catalog.add_entry('rockman-exe-anime', 'anime', 'tv', 'MegaMan NT Warrior', '2002-03-04', 'JP', 'episode', 56, p_public_id => 109);
select catalog.member('mega-man-battle-network', 'rockman-exe-anime');
select catalog.relate('rockman-exe-anime', 'adaptation_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', source_type = 'game', unit_minutes = 24, end_date = '2003-03-31', end_precision = 'day' where id = catalog.id('rockman-exe-anime');
select catalog.add_entry('rockman-exe-axess', 'anime', 'tv', 'MegaMan NT Warrior: Axess', '2003-10-04', 'JP', 'episode', 51, p_public_id => 110);
select catalog.member('mega-man-battle-network', 'rockman-exe-axess');
select catalog.relate('rockman-exe-axess', 'sequel_of', 'rockman-exe-anime');
update catalog.entry set status = 'finished', source_type = 'game', unit_minutes = 24, end_date = '2004-09-25', end_precision = 'day' where id = catalog.id('rockman-exe-axess');
select catalog.ext('rockman-exe-axess', 'anilist_anime', '2108');
select catalog.add_entry('rockman-exe-stream', 'anime', 'tv', 'Rockman EXE Stream', '2004-10-02', 'JP', 'episode', 51, p_public_id => 111);
select catalog.member('mega-man-battle-network', 'rockman-exe-stream');
select catalog.relate('rockman-exe-stream', 'sequel_of', 'rockman-exe-axess');
update catalog.entry set status = 'finished', source_type = 'game', unit_minutes = 25, end_date = '2005-09-24', end_precision = 'day' where id = catalog.id('rockman-exe-stream');
select catalog.ext('rockman-exe-stream', 'anilist_anime', '2109');
select catalog.add_entry('rockman-exe-beast', 'anime', 'tv', 'Rockman EXE Beast', '2005-10-01', 'JP', 'episode', 25, p_public_id => 112);
select catalog.member('mega-man-battle-network', 'rockman-exe-beast');
select catalog.relate('rockman-exe-beast', 'sequel_of', 'rockman-exe-stream');
update catalog.entry set status = 'finished', source_type = 'game', unit_minutes = 25, end_date = '2006-04-01', end_precision = 'day' where id = catalog.id('rockman-exe-beast');
select catalog.ext('rockman-exe-beast', 'anilist_anime', '2110');
select catalog.add_entry('rockman-exe-beast-plus', 'anime', 'tv', 'Rockman EXE Beast+', '2006-04-08', 'JP', 'episode', 26, p_public_id => 113);
select catalog.member('mega-man-battle-network', 'rockman-exe-beast-plus');
select catalog.relate('rockman-exe-beast-plus', 'sequel_of', 'rockman-exe-beast');
update catalog.entry set status = 'finished', source_type = 'game', unit_minutes = 10, end_date = '2006-09-30', end_precision = 'day' where id = catalog.id('rockman-exe-beast-plus');
select catalog.ext('rockman-exe-beast-plus', 'anilist_anime', '2713');
select catalog.add_entry('rockman-exe-movie', 'anime', 'movie', 'Rockman EXE: The Program of Light and Darkness', '2005-03-12', 'JP', 'none', null, p_public_id => 114);
select catalog.member('mega-man-battle-network', 'rockman-exe-movie');
select catalog.relate('rockman-exe-movie', 'side_story_of', 'rockman-exe-stream');
update catalog.entry set status = 'finished', source_type = 'game', length_minutes = 50 where id = catalog.id('rockman-exe-movie');
select catalog.add_entry('mega-man-star-force-anime', 'anime', 'tv', 'Mega Man Star Force', '2006-10-07', 'JP', 'episode', 55, p_public_id => 115);
select catalog.member('mega-man-star-force', 'mega-man-star-force-anime');
select catalog.relate('mega-man-star-force-anime', 'adaptation_of', 'mega-man-sf-pegasus');
update catalog.entry set status = 'finished', source_type = 'game', unit_minutes = 10, end_date = '2007-10-27', end_precision = 'day' where id = catalog.id('mega-man-star-force-anime');
select catalog.ext('mega-man-star-force-anime', 'anilist_anime', '1835');
select catalog.add_entry('mega-man-star-force-tribe', 'anime', 'tv', 'Mega Man Star Force Tribe', '2007-11-03', 'JP', 'episode', 21, p_public_id => 116);
select catalog.member('mega-man-star-force', 'mega-man-star-force-tribe');
select catalog.relate('mega-man-star-force-tribe', 'sequel_of', 'mega-man-star-force-anime');
update catalog.entry set status = 'finished', source_type = 'game', unit_minutes = 10, end_date = '2008-03-29', end_precision = 'day' where id = catalog.id('mega-man-star-force-tribe');
select catalog.ext('mega-man-star-force-tribe', 'anilist_anime', '3493');

select catalog.alias('rockman-exe-anime', 'romanized', 'Rockman.EXE');
select catalog.alias('mega-man-upon-a-star', 'romanized', 'Rockman: Hoshi ni Negai wo');
select catalog.alias('mega-man-day-of-sigma', 'romanized', 'Irregular Hunter X: The Day of Sigma');
select catalog.alias('mega-man-star-force-anime', 'romanized', 'Ryuusei no Rockman');
select catalog.alias('mega-man-star-force-tribe', 'romanized', 'Ryuusei no Rockman Tribe');
select catalog.credit('mega-man-ruby-spears', 'Ruby-Spears Productions', 'studio', 'US');
select catalog.credit('mega-man-ruby-spears', 'Ashi Productions', 'studio', 'JP');
select catalog.credit('mega-man-upon-a-star', 'Ashi Productions', 'studio', 'JP');
select catalog.credit('mega-man-fully-charged', 'DHX Studios Vancouver', 'studio', 'CA');
select catalog.credit(s, 'Xebec', 'studio', 'JP') from unnest(array[
  'rockman-exe-anime','rockman-exe-axess','rockman-exe-stream','rockman-exe-beast',
  'rockman-exe-beast-plus','rockman-exe-movie','mega-man-star-force-anime','mega-man-star-force-tribe']) s;

select catalog.add_ordering('mega-man-battle-network', 'release', 'anime order',
  'The anime is independent from the games. The theatrical film is a side story during Stream.') as o \gset
select catalog.add_step(:'o', 1, 'rockman-exe-anime');
select catalog.add_step(:'o', 2, 'rockman-exe-axess');
select catalog.add_step(:'o', 3, 'rockman-exe-stream');
select catalog.add_step(:'o', 4, 'rockman-exe-movie', p_necessity => 'optional');
select catalog.add_step(:'o', 5, 'rockman-exe-beast');
select catalog.add_step(:'o', 6, 'rockman-exe-beast-plus');
select catalog.add_ordering('mega-man-star-force', 'release', 'anime order',
  'Two animated series in their own adaptation continuity.') as o \gset
select catalog.add_step(:'o', 1, 'mega-man-star-force-anime');
select catalog.add_step(:'o', 2, 'mega-man-star-force-tribe');

-- Licensed manga.
-- https://anilist.co/manga/49438
-- https://anilist.co/manga/54302/Rockman-Gigamix
-- https://anilist.co/manga/38525/Rockman-X/
-- https://anilist.co/manga/51684/Rockman-X4
-- https://anilist.co/manga/38008/Rockman-ZX/social
-- https://anilist.co/manga/30761
-- https://www.fukkan.com/list/isearch?tag=%E6%B1%A0%E5%8E%9F%E3%81%97%E3%81%92%E3%81%A8
-- https://www.fukkan.com/fk/CartSearchDetail?i_no=68319096
-- https://www.fukkan.com/fk/CartSearchDetail?i_no=37556921
select catalog.add_entry('rockman-manga-ikehara', 'manga', 'serial', 'Rockman (Shigeto Ikehara)', null, 'JP', 'volume', null, p_public_id => 117);
select catalog.member('mega-man-classic', 'rockman-manga-ikehara');
select catalog.relate('rockman-manga-ikehara', 'adaptation_of', 'mega-man-1');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-manga-ikehara');
select catalog.add_entry('rockman-2-manga', 'manga', 'serial', 'Rockman 2 (Shigeto Ikehara)', null, 'JP', 'volume', null, p_public_id => 118);
select catalog.member('mega-man-classic', 'rockman-2-manga');
select catalog.relate('rockman-2-manga', 'adaptation_of', 'mega-man-2');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-2-manga');
select catalog.add_entry('rockman-3-manga', 'manga', 'serial', 'Rockman 3 (Shigeto Ikehara)', null, 'JP', 'volume', null, p_public_id => 119);
select catalog.member('mega-man-classic', 'rockman-3-manga');
select catalog.relate('rockman-3-manga', 'adaptation_of', 'mega-man-3');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-3-manga');
select catalog.add_entry('rockman-4-manga', 'manga', 'serial', 'Rockman 4 (Shigeto Ikehara)', null, 'JP', 'volume', null, p_public_id => 120);
select catalog.member('mega-man-classic', 'rockman-4-manga');
select catalog.relate('rockman-4-manga', 'adaptation_of', 'mega-man-4');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-4-manga');
select catalog.add_entry('rockman-5-manga', 'manga', 'serial', 'Rockman 5 (Shigeto Ikehara)', null, 'JP', 'volume', null, p_public_id => 121);
select catalog.member('mega-man-classic', 'rockman-5-manga');
select catalog.relate('rockman-5-manga', 'adaptation_of', 'mega-man-5');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-5-manga');
select catalog.add_entry('rockman-6-manga', 'manga', 'serial', 'Rockman 6 (Shigeto Ikehara)', null, 'JP', 'volume', null, p_public_id => 122);
select catalog.member('mega-man-classic', 'rockman-6-manga');
select catalog.relate('rockman-6-manga', 'adaptation_of', 'mega-man-6');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-6-manga');
select catalog.add_entry('rockman-7-manga', 'manga', 'serial', 'Rockman 7 (Shigeto Ikehara)', null, 'JP', 'volume', 2, p_public_id => 123);
select catalog.member('mega-man-classic', 'rockman-7-manga');
select catalog.relate('rockman-7-manga', 'adaptation_of', 'mega-man-7');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-7-manga');
select catalog.add_units('rockman-7-manga', 2);
select catalog.add_entry('rockman-8-manga', 'manga', 'serial', 'Rockman 8 (Koji Izuki)', null, 'JP', 'volume', null, p_public_id => 124);
select catalog.member('mega-man-classic', 'rockman-8-manga');
select catalog.relate('rockman-8-manga', 'adaptation_of', 'mega-man-8');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-8-manga');
select catalog.add_entry('rockman-bass-manga', 'manga', 'serial', 'Rockman & Forte (Koji Izuki)', null, 'JP', 'volume', null, p_public_id => 125);
select catalog.member('mega-man-classic', 'rockman-bass-manga');
select catalog.relate('rockman-bass-manga', 'adaptation_of', 'mega-man-bass');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-bass-manga');
select catalog.add_entry('rockman-x-manga', 'manga', 'serial', 'Mega Man X (Yoshihiro Iwamoto)', null, 'JP', 'volume', 3, p_public_id => 126);
select catalog.member('mega-man-x', 'rockman-x-manga');
select catalog.relate('rockman-x-manga', 'adaptation_of', 'mega-man-x-game');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-x-manga');
select catalog.add_units('rockman-x-manga', 3);
select catalog.ext('rockman-x-manga', 'anilist_manga', '38525');
select catalog.add_entry('rockman-x2-manga', 'manga', 'serial', 'Mega Man X2 (Yoshihiro Iwamoto)', null, 'JP', 'volume', null, p_public_id => 127);
select catalog.member('mega-man-x', 'rockman-x2-manga');
select catalog.relate('rockman-x2-manga', 'adaptation_of', 'mega-man-x2');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-x2-manga');
select catalog.add_entry('rockman-x3-manga', 'manga', 'serial', 'Mega Man X3 (Yoshihiro Iwamoto)', null, 'JP', 'volume', null, p_public_id => 128);
select catalog.member('mega-man-x', 'rockman-x3-manga');
select catalog.relate('rockman-x3-manga', 'adaptation_of', 'mega-man-x3');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-x3-manga');
select catalog.add_entry('rockman-x4-manga', 'manga', 'serial', 'Mega Man X4 (Yoshihiro Iwamoto)', null, 'JP', 'volume', 2, p_public_id => 129);
select catalog.member('mega-man-x', 'rockman-x4-manga');
select catalog.relate('rockman-x4-manga', 'adaptation_of', 'mega-man-x4');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-x4-manga');
select catalog.add_units('rockman-x4-manga', 2);
select catalog.ext('rockman-x4-manga', 'anilist_manga', '51684');
select catalog.add_entry('irregular-hunter-x-manga', 'manga', 'serial', 'Irregular Hunter Rockman X (Shigeto Ikehara)', null, 'JP', 'volume', 2, p_public_id => 130);
select catalog.member('mega-man-x', 'irregular-hunter-x-manga');
select catalog.relate('irregular-hunter-x-manga', 'adaptation_of', 'mega-man-x-game');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('irregular-hunter-x-manga');
select catalog.add_units('irregular-hunter-x-manga', 2);
select catalog.add_entry('rockman-zero-manga', 'manga', 'serial', 'Rockman Zero (Hideto Kajima)', null, 'JP', 'volume', null, p_public_id => 131);
select catalog.member('mega-man-zero', 'rockman-zero-manga');
select catalog.relate('rockman-zero-manga', 'adaptation_of', 'mega-man-zero-1');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-zero-manga');
select catalog.add_entry('rockman-zx-manga', 'manga', 'serial', 'Mega Man ZX (Shin Ogino)', null, 'JP', 'volume', 2, p_public_id => 132);
select catalog.member('mega-man-zx', 'rockman-zx-manga');
select catalog.relate('rockman-zx-manga', 'adaptation_of', 'mega-man-zx-game');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-zx-manga');
select catalog.add_units('rockman-zx-manga', 2);
select catalog.ext('rockman-zx-manga', 'anilist_manga', '38008');
select catalog.add_entry('rockman-zx-advent-manga', 'manga', 'serial', 'Rockman ZX Advent (Shin Ogino)', null, 'JP', 'volume', null, p_public_id => 133);
select catalog.member('mega-man-zx', 'rockman-zx-advent-manga');
select catalog.relate('rockman-zx-advent-manga', 'adaptation_of', 'mega-man-zx-advent');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-zx-advent-manga');
select catalog.add_entry('rockman-exe-manga', 'manga', 'serial', 'MegaMan NT Warrior (Ryo Takamisaki)', null, 'JP', 'volume', 13, p_public_id => 134);
select catalog.member('mega-man-battle-network', 'rockman-exe-manga');
select catalog.relate('rockman-exe-manga', 'adaptation_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-exe-manga');
select catalog.add_units('rockman-exe-manga', 13);
select catalog.ext('rockman-exe-manga', 'anilist_manga', '30761');
select catalog.add_entry('rockman-exe-battle-story', 'manga', 'serial', 'Battle Story Rockman EXE', null, 'JP', 'volume', null, p_public_id => 135);
select catalog.member('mega-man-battle-network', 'rockman-exe-battle-story');
select catalog.relate('rockman-exe-battle-story', 'adaptation_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-exe-battle-story');
select catalog.add_entry('rockman-star-force-manga', 'manga', 'serial', 'Ryuusei no Rockman (Masaya Itagaki)', null, 'JP', 'volume', null, p_public_id => 136);
select catalog.member('mega-man-star-force', 'rockman-star-force-manga');
select catalog.relate('rockman-star-force-manga', 'adaptation_of', 'mega-man-sf-pegasus');
update catalog.entry set status = 'finished', description = 'A licensed manga adaptation. The original serialization date is unverified; later collected editions are kept in this work.', source_type = 'game' where id = catalog.id('rockman-star-force-manga');
select catalog.add_entry('rockman-megamix', 'manga', 'serial', 'Mega Man Megamix', null, 'JP', 'volume', 3, p_public_id => 137);
select catalog.member('mega-man-classic', 'rockman-megamix');
select catalog.relate('rockman-megamix', 'adaptation_of', 'mega-man-1');
update catalog.entry set status = 'finished', description = 'Hitoshi Ariga''s connected short stories. The three-volume collected edition includes material previously published as Rockman Remix.', source_type = 'game' where id = catalog.id('rockman-megamix');
select catalog.add_units('rockman-megamix', 3);
select catalog.ext('rockman-megamix', 'anilist_manga', '49438');
select catalog.add_entry('rockman-gigamix', 'manga', 'serial', 'Mega Man Gigamix', '2009-11-30', 'JP', 'volume', 3, p_public_id => 138);
select catalog.member('mega-man-classic', 'rockman-gigamix');
select catalog.relate('rockman-gigamix', 'side_story_of', 'rockman-megamix');
update catalog.entry set status = 'finished', description = 'Hitoshi Ariga''s companion manga to Megamix, containing both earlier and later stories.', source_type = 'game', end_date = '2010-11-27', end_precision = 'day' where id = catalog.id('rockman-gigamix');
select catalog.add_units('rockman-gigamix', 3);
select catalog.ext('rockman-gigamix', 'anilist_manga', '54302');

-- Western comics.
-- https://archiecomics.com/preview-the-new-archie-comics-on-sale-today-including-mega-man-55/
-- https://www.boom-studios.com/archives/meet-the-all-new-mega-man-in-a-fully-charged-comic-book-series-from-boom-studios/
-- https://comicvine.gamespot.com/mega-man-mastermix/4050-108109/
-- https://store.udonentertainment.com/blogs/news/new-classic-mega-man-comics-from-udon
-- https://store.udonentertainment.com/collections/mega-man
-- https://www.midtowncomics.com/weekly-release-home/print?wdate=02%2F18%2F2026
-- https://store.udonentertainment.com/collections/comics/products/mega-man-gear-masters-1-cvr-a-ryan-jampole
select catalog.add_entry('mega-man-archie', 'comic', 'serial', 'Mega Man (Archie Comics)', '2011', 'US', 'chapter', 55, p_public_id => 139);
select catalog.member('mega-man-classic', 'mega-man-archie');
select catalog.relate('mega-man-archie', 'adaptation_of', 'mega-man-1');
update catalog.entry set status = 'hiatus', description = 'The 55-issue Archie comic, including its own adaptations of the games and its Sonic crossover installments. Publication went on indefinite hiatus after issue 55.', source_type = 'game', end_date = '2015-12-23', end_precision = 'day' where id = catalog.id('mega-man-archie');
select catalog.credit('mega-man-archie', 'Archie Comics', 'publisher', 'US');
select catalog.add_entry('mega-man-mastermix', 'comic', 'serial', 'Mega Man Mastermix', '2018', 'US', 'chapter', 4, p_public_id => 140);
select catalog.member('mega-man-classic', 'mega-man-mastermix');
select catalog.relate('mega-man-mastermix', 'alternative_version_of', 'rockman-megamix');
update catalog.entry set status = 'finished', description = 'A four-issue, full-color presentation of stories from Hitoshi Ariga''s Megamix manga.', source_type = 'manga' where id = catalog.id('mega-man-mastermix');
select catalog.credit('mega-man-mastermix', 'UDON Entertainment', 'publisher', 'US');
select catalog.add_entry('mega-man-fully-charged-comic', 'comic', 'serial', 'Mega Man: Fully Charged (Comic)', '2020', 'US', 'chapter', 6, p_public_id => 141);
select catalog.member('mega-man-classic', 'mega-man-fully-charged-comic');
select catalog.relate('mega-man-fully-charged-comic', 'sequel_of', 'mega-man-fully-charged');
update catalog.entry set status = 'finished', description = 'The six-issue comic continuation of the Fully Charged animated continuity.', source_type = 'other', end_date = '2021-01-01', end_precision = 'year' where id = catalog.id('mega-man-fully-charged-comic');
select catalog.credit('mega-man-fully-charged-comic', 'BOOM! Studios', 'publisher', 'US');
select catalog.add_entry('mega-man-udon-fcbd-2025', 'comic', 'one_shot', 'Mega Man #0 (Free Comic Book Day 2025)', '2025-05-03', 'US', 'none', null, p_public_id => 142);
select catalog.member('mega-man-classic', 'mega-man-udon-fcbd-2025');
select catalog.relate('mega-man-udon-fcbd-2025', 'adaptation_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'game' where id = catalog.id('mega-man-udon-fcbd-2025');
select catalog.credit('mega-man-udon-fcbd-2025', 'UDON Entertainment', 'publisher', 'US');
select catalog.add_entry('mega-man-timelines', 'comic', 'one_shot', 'Mega Man Timelines #1', '2025-05-28', 'US', 'none', null, p_public_id => 143);
select catalog.member('mega-man-classic', 'mega-man-timelines');
select catalog.relate('mega-man-timelines', 'adaptation_of', 'mega-man-1');
update catalog.entry set status = 'finished', source_type = 'game' where id = catalog.id('mega-man-timelines');
select catalog.credit('mega-man-timelines', 'UDON Entertainment', 'publisher', 'US');
select catalog.add_entry('mega-man-x-timelines', 'comic', 'one_shot', 'Mega Man X Timelines #1', '2025-06', 'US', 'none', null, p_public_id => 144);
select catalog.member('mega-man-x', 'mega-man-x-timelines');
select catalog.relate('mega-man-x-timelines', 'adaptation_of', 'mega-man-x-game');
update catalog.entry set status = 'finished', source_type = 'game' where id = catalog.id('mega-man-x-timelines');
select catalog.credit('mega-man-x-timelines', 'UDON Entertainment', 'publisher', 'US');
select catalog.add_entry('mega-man-zero-timelines', 'comic', 'one_shot', 'Mega Man Zero Timelines #1', '2025-07', 'US', 'none', null, p_public_id => 145);
select catalog.member('mega-man-zero', 'mega-man-zero-timelines');
select catalog.relate('mega-man-zero-timelines', 'adaptation_of', 'mega-man-zero-1');
update catalog.entry set status = 'finished', source_type = 'game' where id = catalog.id('mega-man-zero-timelines');
select catalog.credit('mega-man-zero-timelines', 'UDON Entertainment', 'publisher', 'US');
select catalog.add_entry('mega-man-zx-timelines', 'comic', 'one_shot', 'Mega Man ZX Timelines #1', '2025-08', 'US', 'none', null, p_public_id => 146);
select catalog.member('mega-man-zx', 'mega-man-zx-timelines');
select catalog.relate('mega-man-zx-timelines', 'adaptation_of', 'mega-man-zx-game');
update catalog.entry set status = 'finished', source_type = 'game' where id = catalog.id('mega-man-zx-timelines');
select catalog.credit('mega-man-zx-timelines', 'UDON Entertainment', 'publisher', 'US');
select catalog.add_entry('mega-man-legends-timelines', 'comic', 'one_shot', 'Mega Man Legends Timelines #1', '2025-09', 'US', 'none', null, p_public_id => 147);
select catalog.member('mega-man-legends', 'mega-man-legends-timelines');
select catalog.relate('mega-man-legends-timelines', 'adaptation_of', 'mega-man-legends-game');
update catalog.entry set status = 'finished', source_type = 'game' where id = catalog.id('mega-man-legends-timelines');
select catalog.credit('mega-man-legends-timelines', 'UDON Entertainment', 'publisher', 'US');
select catalog.add_entry('mega-man-x-udon', 'comic', 'serial', 'Mega Man X (UDON Comics)', '2026-02', 'US', 'chapter', null, p_public_id => 148);
select catalog.member('mega-man-x', 'mega-man-x-udon');
select catalog.relate('mega-man-x-udon', 'adaptation_of', 'mega-man-x-game');
update catalog.entry set status = 'releasing', description = 'The ongoing UDON Mega Man X comic. This series is separate from the 2025 Timelines one-shot and the Japanese manga.', source_type = 'game' where id = catalog.id('mega-man-x-udon');
select catalog.credit('mega-man-x-udon', 'UDON Entertainment', 'publisher', 'US');
select catalog.add_entry('mega-man-gear-masters', 'comic', 'serial', 'Mega Man: Gear Masters', null, 'US', 'chapter', null, p_public_id => 149);
select catalog.member('mega-man-classic', 'mega-man-gear-masters');
select catalog.relate('mega-man-gear-masters', 'adaptation_of', 'mega-man-11');
update catalog.entry set status = 'announced', description = 'The announced UDON miniseries about the Double Gear system. The publisher still advertises issue 1 as a preorder; its actual release date is unverified.', source_type = 'game' where id = catalog.id('mega-man-gear-masters');
select catalog.credit('mega-man-gear-masters', 'UDON Entertainment', 'publisher', 'US');

-- Released fan games.
-- https://megaphilx.com/index.php/home/games/mega-man-unlimited/
-- https://www.mobygames.com/game/105307/mega-man-unlimited/
-- https://www.mobygames.com/game/105860/mega-man-rock-force/
-- https://dennisengelhard.com/rock-n-roll/
-- https://dennisengelhard.com/
-- https://petersjostrand.com/
-- https://petersjostrand.com/faq.html
-- https://megamanmaker.com/
-- https://mm8bdm.net/
-- https://www.siliconera.com/fan-made-game-mega-man-super-fighting-robot-is-now-out/
select catalog.add_entry('mega-man-unlimited', 'game', 'full', 'Mega Man Unlimited', '2013-07-14', null, 'none', null, p_public_id => 150);
select catalog.member('mega-man-fan-games', 'mega-man-unlimited');
select catalog.relate('mega-man-unlimited', 'fan_work_of', 'mega-man');
update catalog.entry set status = 'finished', description = 'An unofficial NES-style fan game by Philippe Poulin, Jean-Simon Brochu and Gabriel Leblanc.', source_type = 'game' where id = catalog.id('mega-man-unlimited');
select catalog.ext('mega-man-unlimited', 'mobygames', '105307');
select catalog.add_entry('mega-man-rock-force', 'game', 'full', 'Mega Man: Rock Force', '2013-09-03', null, 'none', null, p_public_id => 151);
select catalog.member('mega-man-fan-games', 'mega-man-rock-force');
select catalog.relate('mega-man-rock-force', 'fan_work_of', 'mega-man');
update catalog.entry set status = 'finished', description = 'An unofficial fan game in which Mega Man and allied Robot Masters form the Rock Force.', source_type = 'game' where id = catalog.id('mega-man-rock-force');
select catalog.ext('mega-man-rock-force', 'mobygames', '105860');
select catalog.add_entry('mega-man-rock-n-roll', 'game', 'full', 'Mega Man: Rock N Roll', '2019', null, 'none', null, p_public_id => 152);
select catalog.member('mega-man-fan-games', 'mega-man-rock-n-roll');
select catalog.relate('mega-man-rock-n-roll', 'fan_work_of', 'mega-man');
update catalog.entry set status = 'finished', description = 'An unofficial fan game by Dennis Engelhard, with playable Mega Man, Roll and Proto Man.', source_type = 'game' where id = catalog.id('mega-man-rock-n-roll');
select catalog.add_entry('mega-man-25d', 'game', 'full', 'Mega Man 2.5D', '2017', null, 'none', null, p_public_id => 153);
select catalog.member('mega-man-fan-games', 'mega-man-25d');
select catalog.relate('mega-man-25d', 'fan_work_of', 'mega-man');
update catalog.entry set status = 'finished', description = 'An unofficial fan game by Peter Sjostrand and collaborators, featuring single-player and cooperative campaigns.', source_type = 'game' where id = catalog.id('mega-man-25d');
select catalog.add_entry('mega-man-maker', 'game', 'full', 'Mega Man Maker', '2017-07-15', null, 'none', null, p_public_id => 154);
select catalog.member('mega-man-fan-games', 'mega-man-maker');
select catalog.relate('mega-man-maker', 'fan_work_of', 'mega-man');
update catalog.entry set status = 'releasing', description = 'An unofficial fan game for creating, sharing and playing custom Mega Man levels.', source_type = 'game' where id = catalog.id('mega-man-maker');
select catalog.add_entry('mega-man-8bit-deathmatch', 'game', 'full', 'Mega Man 8-Bit Deathmatch', '2010', null, 'none', null, p_public_id => 155);
select catalog.member('mega-man-fan-games', 'mega-man-8bit-deathmatch');
select catalog.relate('mega-man-8bit-deathmatch', 'fan_work_of', 'mega-man');
update catalog.entry set status = 'releasing', description = 'An unofficial first-person fan game by CutmanMike and Team MM8BDM, with a campaign and multiplayer modes.', source_type = 'game' where id = catalog.id('mega-man-8bit-deathmatch');
select catalog.add_entry('mega-man-super-fighting-robot', 'game', 'full', 'Mega Man: Super Fighting Robot', '2015', null, 'none', null, p_public_id => 156);
select catalog.member('mega-man-fan-games', 'mega-man-super-fighting-robot');
select catalog.relate('mega-man-super-fighting-robot', 'fan_work_of', 'mega-man');
update catalog.entry set status = 'finished', description = 'An unofficial fan platform game by MrWeirdGuy, featuring upgrades and secret bosses.', source_type = 'game' where id = catalog.id('mega-man-super-fighting-robot');
select catalog.set_role('mega-man-fan-games', 'mega-man-unlimited', 'unofficial fan game');
select catalog.set_role('mega-man-fan-games', 'mega-man-rock-force', 'unofficial fan game');
select catalog.set_role('mega-man-fan-games', 'mega-man-rock-n-roll', 'unofficial fan game');
select catalog.set_role('mega-man-fan-games', 'mega-man-25d', 'unofficial fan game');
select catalog.set_role('mega-man-fan-games', 'mega-man-maker', 'unofficial fan game');
select catalog.set_role('mega-man-fan-games', 'mega-man-8bit-deathmatch', 'unofficial fan game');
select catalog.set_role('mega-man-fan-games', 'mega-man-super-fighting-robot', 'unofficial fan game');

select catalog.credit(n.slug, 'Capcom', 'publisher', 'JP')
from catalog.entry e join catalog.node n on n.id = e.id
where e.medium = 'game' and exists (
  select 1 from catalog.membership m join catalog.node g on g.id = m.group_id
  where m.member_id = n.id and g.slug in ('mega-man-classic','mega-man-x',
    'mega-man-zero','mega-man-zx','mega-man-legends','mega-man-battle-network','mega-man-star-force'))
  and n.slug not in ('mega-man-dos','mega-man-iii-dos','rockman-strategy');
select catalog.credit(s, 'Inti Creates', 'developer', 'JP') from unnest(array[
  'mega-man-9','mega-man-10','mega-man-zero-1','mega-man-zero-2',
  'mega-man-zero-3','mega-man-zero-4','mega-man-zx-game','mega-man-zx-advent']) s;
select catalog.tag_with('theme', 'fan work', variadic array(
  select n.slug from catalog.membership m join catalog.node n on n.id = m.member_id
  where m.group_id = catalog.id('mega-man-fan-games')));

-- The Game Gear game was developed in the United States.
update catalog.entry set country = 'US' where id = catalog.id('mega-man-game-gear');

-- Announced fan projects and a fan remake.
-- https://megaphilx.com/index.php/home/games/mega-man-unlimited-2/
-- https://www.igdb.com/games/mega-man-x-corrupted/presskit
-- https://github.com/AlyssonDaPaz/Mega-Man-X8-16-bit
select catalog.add_entry('mega-man-unlimited-2', 'game', 'full', 'Mega Man Unlimited 2 (Working Title)', null, null, 'none', null, p_public_id => 157);
select catalog.member('mega-man-fan-games', 'mega-man-unlimited-2');
select catalog.relate('mega-man-unlimited-2', 'fan_work_of', 'mega-man');
update catalog.entry set status = 'announced', description = 'An unofficial sequel to Mega Man Unlimited by Philippe Poulin. The developer lists it as in development, without a release date.', source_type = 'game' where id = catalog.id('mega-man-unlimited-2');
select catalog.add_entry('mega-man-x-corrupted', 'game', 'full', 'Mega Man X: Corrupted', null, null, 'none', null, p_public_id => 158);
select catalog.member('mega-man-fan-games', 'mega-man-x-corrupted');
select catalog.relate('mega-man-x-corrupted', 'fan_work_of', 'mega-man-x');
update catalog.entry set status = 'announced', description = 'An unofficial Mega Man X fan game by JKB Games, featuring an interconnected world and character upgrades. It remains in development with no confirmed release date.', source_type = 'game' where id = catalog.id('mega-man-x-corrupted');
select catalog.add_entry('mega-man-x8-16bit', 'game', 'full', 'Mega Man X8 16-bit', null, null, 'none', null, p_public_id => 159);
select catalog.member('mega-man-fan-games', 'mega-man-x8-16bit');
select catalog.relate('mega-man-x8-16bit', 'fan_work_of', 'mega-man-x8');
update catalog.entry set status = 'finished', description = 'An unofficial reimagining of Mega Man X8 by Alysson da Paz, using a style inspired by the Super NES games. The original release date is unverified.', source_type = 'game' where id = catalog.id('mega-man-x8-16bit');

select catalog.relate('mega-man-unlimited-2','sequel_of','mega-man-unlimited');
select catalog.relate('mega-man-x8-16bit','remake_of','mega-man-x8');
select catalog.set_role('mega-man-fan-games','mega-man-unlimited-2','unofficial fan game');
select catalog.set_role('mega-man-fan-games','mega-man-x-corrupted','unofficial fan game');
select catalog.set_role('mega-man-fan-games','mega-man-x8-16bit','unofficial fan game');
select catalog.tag_with('theme','fan work','mega-man-unlimited-2','mega-man-x-corrupted','mega-man-x8-16bit');

-- Additional licensed manga.
-- https://www.fukkan.com/list/isearch?c=1&o=send&s=send&z=2
-- https://anilist.co/manga/133300/Zenin-Shuugou-Gag-tte-Rockman-
select catalog.add_entry('rockman-world-manga', 'manga', 'serial', 'Rockman World (Shigeto Ikehara)', null, 'JP', 'volume', null, p_public_id => 160);
select catalog.member('mega-man-classic', 'rockman-world-manga');
select catalog.relate('rockman-world-manga', 'adaptation_of', 'mega-man-wilys-revenge');
update catalog.entry set status = 'finished', source_type = 'game' where id = catalog.id('rockman-world-manga');
select catalog.add_entry('rockman-world-3-manga', 'manga', 'serial', 'Rockman World 3 (Shigeto Ikehara)', null, 'JP', 'volume', null, p_public_id => 161);
select catalog.member('mega-man-classic', 'rockman-world-3-manga');
select catalog.relate('rockman-world-3-manga', 'adaptation_of', 'mega-man-iii-gb');
update catalog.entry set status = 'finished', source_type = 'game' where id = catalog.id('rockman-world-3-manga');
select catalog.add_entry('rockman-gag-manga', 'manga', 'serial', 'Zenin Shuugou: Gag tte!! Rockman!', null, 'JP', 'chapter', 30, p_public_id => 162);
select catalog.member('mega-man-battle-network', 'rockman-gag-manga');
select catalog.relate('rockman-gag-manga', 'spin_off_of', 'mega-man-bn1');
update catalog.entry set status = 'finished', source_type = 'game' where id = catalog.id('rockman-gag-manga');
select catalog.ext('rockman-gag-manga', 'anilist_manga', '133300');

-- These totals describe the publisher's collected edition.
update catalog.entry set unit_total = 2 where id = catalog.id('rockman-5-manga');
select catalog.add_units('rockman-5-manga', 2);
update catalog.entry set unit_total = 3 where id = catalog.id('rockman-6-manga');
select catalog.add_units('rockman-6-manga', 3);
