-- Detective Conan, from Wikidata (CC0) on 2026-10-18: titles, first release
-- dates, and AniList / MAL ids. "Case Closed" is kept as the english title
-- where Wikidata uses it.
--
-- Not yet: TV specials, recaps, the spin-off anime, Magic Kaito.

select catalog.add_group('detective-conan', 'franchise', 'Detective Conan');
select catalog.alias('detective-conan', 'native', '名探偵コナン');
select catalog.alias('detective-conan', 'synonym', 'Case Closed');
select catalog.ext('detective-conan', 'wikidata', 'Q185143');

-- The manga and the TV anime are both ongoing, so neither has a unit total.

select catalog.add_entry('detective-conan-manga', 'manga', 'serial', 'Detective Conan', '1994-01', 'JP', 'chapter');
select catalog.alias('detective-conan-manga', 'native', '名探偵コナン');
select catalog.alias('detective-conan-manga', 'english', 'Case Closed');
select catalog.ext('detective-conan-manga', 'wikidata', 'Q3853685');
select catalog.ext('detective-conan-manga', 'anilist_manga', '31061');
select catalog.ext('detective-conan-manga', 'mal_manga', '1061');

select catalog.add_entry('detective-conan-anime', 'anime', 'tv', 'Detective Conan', '1996-01-08', 'JP', 'episode');
select catalog.alias('detective-conan-anime', 'native', '名探偵コナン');
select catalog.alias('detective-conan-anime', 'english', 'Case Closed');
select catalog.ext('detective-conan-anime', 'wikidata', 'Q5363072');
select catalog.ext('detective-conan-anime', 'anilist_anime', '235');
select catalog.ext('detective-conan-anime', 'mal_anime', '235');
select catalog.relate('detective-conan-anime', 'adaptation_of', 'detective-conan-manga');

select catalog.member('detective-conan', 'detective-conan-manga', 'detective-conan-anime');

-- The theatrical films: one a year since 1997, numbered by release.

select catalog.add_group('detective-conan-films', 'series', 'Detective Conan films');
select catalog.ext('detective-conan-films', 'wikidata', 'Q12294798');
select catalog.member('detective-conan', 'detective-conan-films');

select catalog.add_entry('detective-conan-the-time-bombed-skyscraper', 'anime', 'movie', 'Detective Conan: The Time Bombed Skyscraper', '1997-04-19', 'JP');
select catalog.alias('detective-conan-the-time-bombed-skyscraper', 'native', '名探偵コナン 時計じかけの摩天楼');
select catalog.alias('detective-conan-the-time-bombed-skyscraper', 'english', 'Case Closed: The Time Bombed Skyscraper');
select catalog.ext('detective-conan-the-time-bombed-skyscraper', 'wikidata', 'Q1195711');
select catalog.ext('detective-conan-the-time-bombed-skyscraper', 'anilist_anime', '779');
select catalog.ext('detective-conan-the-time-bombed-skyscraper', 'mal_anime', '779');

select catalog.add_entry('detective-conan-the-fourteenth-target', 'anime', 'movie', 'Detective Conan: The Fourteenth Target', '1998-04-18', 'JP');
select catalog.alias('detective-conan-the-fourteenth-target', 'native', '名探偵コナン 14番目の標的');
select catalog.alias('detective-conan-the-fourteenth-target', 'english', 'Case Closed: The Fourteenth Target');
select catalog.ext('detective-conan-the-fourteenth-target', 'wikidata', 'Q1193711');
select catalog.ext('detective-conan-the-fourteenth-target', 'anilist_anime', '780');
select catalog.ext('detective-conan-the-fourteenth-target', 'mal_anime', '780');

select catalog.add_entry('detective-conan-the-last-wizard-of-the-century', 'anime', 'movie', 'Detective Conan: The Last Wizard of the Century', '1999-04-17', 'JP');
select catalog.alias('detective-conan-the-last-wizard-of-the-century', 'native', '名探偵コナン 世紀末の魔術師');
select catalog.alias('detective-conan-the-last-wizard-of-the-century', 'english', 'Case Closed: The Last Wizard of the Century');
select catalog.ext('detective-conan-the-last-wizard-of-the-century', 'wikidata', 'Q260821');
select catalog.ext('detective-conan-the-last-wizard-of-the-century', 'anilist_anime', '781');
select catalog.ext('detective-conan-the-last-wizard-of-the-century', 'mal_anime', '781');

select catalog.add_entry('detective-conan-captured-in-her-eyes', 'anime', 'movie', 'Detective Conan: Captured in Her Eyes', '2000-04-22', 'JP');
select catalog.alias('detective-conan-captured-in-her-eyes', 'native', '名探偵コナン 瞳の中の暗殺者');
select catalog.alias('detective-conan-captured-in-her-eyes', 'english', 'Case Closed: Captured in Her Eyes');
select catalog.ext('detective-conan-captured-in-her-eyes', 'wikidata', 'Q2840937');
select catalog.ext('detective-conan-captured-in-her-eyes', 'anilist_anime', '1363');
select catalog.ext('detective-conan-captured-in-her-eyes', 'mal_anime', '1363');

select catalog.add_entry('detective-conan-countdown-to-heaven', 'anime', 'movie', 'Detective Conan: Countdown to Heaven', '2001-04-21', 'JP');
select catalog.alias('detective-conan-countdown-to-heaven', 'native', '名探偵コナン 天国へのカウントダウン');
select catalog.alias('detective-conan-countdown-to-heaven', 'english', 'Case Closed: Countdown to Heaven');
select catalog.ext('detective-conan-countdown-to-heaven', 'wikidata', 'Q287096');
select catalog.ext('detective-conan-countdown-to-heaven', 'anilist_anime', '1364');
select catalog.ext('detective-conan-countdown-to-heaven', 'mal_anime', '1364');

select catalog.add_entry('detective-conan-the-phantom-of-baker-street', 'anime', 'movie', 'Detective Conan: The Phantom of Baker Street', '2002-04-20', 'JP');
select catalog.alias('detective-conan-the-phantom-of-baker-street', 'native', '名探偵コナン ベイカー街の亡霊');
select catalog.alias('detective-conan-the-phantom-of-baker-street', 'english', 'Case Closed: The Phantom of Baker Street');
select catalog.ext('detective-conan-the-phantom-of-baker-street', 'wikidata', 'Q1196712');
select catalog.ext('detective-conan-the-phantom-of-baker-street', 'anilist_anime', '1365');
select catalog.ext('detective-conan-the-phantom-of-baker-street', 'mal_anime', '1365');

select catalog.add_entry('detective-conan-crossroad-in-the-ancient-capital', 'anime', 'movie', 'Detective Conan: Crossroad in the Ancient Capital', '2003-04-19', 'JP');
select catalog.alias('detective-conan-crossroad-in-the-ancient-capital', 'native', '名探偵コナン 迷宮の十字路');
select catalog.ext('detective-conan-crossroad-in-the-ancient-capital', 'wikidata', 'Q1195701');
select catalog.ext('detective-conan-crossroad-in-the-ancient-capital', 'anilist_anime', '1366');
select catalog.ext('detective-conan-crossroad-in-the-ancient-capital', 'mal_anime', '1366');

select catalog.add_entry('detective-conan-magician-of-the-silver-sky', 'anime', 'movie', 'Detective Conan: Magician of the Silver Sky', '2004-04-17', 'JP');
select catalog.alias('detective-conan-magician-of-the-silver-sky', 'native', '名探偵コナン 銀翼の奇術師');
select catalog.ext('detective-conan-magician-of-the-silver-sky', 'wikidata', 'Q1195724');
select catalog.ext('detective-conan-magician-of-the-silver-sky', 'anilist_anime', '1367');
select catalog.ext('detective-conan-magician-of-the-silver-sky', 'mal_anime', '1367');

select catalog.add_entry('detective-conan-strategy-above-the-depths', 'anime', 'movie', 'Detective Conan: Strategy Above the Depths', '2005-04-09', 'JP');
select catalog.alias('detective-conan-strategy-above-the-depths', 'native', '名探偵コナン 水平線上の陰謀');
select catalog.ext('detective-conan-strategy-above-the-depths', 'wikidata', 'Q1207384');
select catalog.ext('detective-conan-strategy-above-the-depths', 'anilist_anime', '1505');
select catalog.ext('detective-conan-strategy-above-the-depths', 'mal_anime', '1505');

select catalog.add_entry('detective-conan-the-private-eyes-requiem', 'anime', 'movie', 'Detective Conan: The Private Eyes'' Requiem', '2006-04-15', 'JP');
select catalog.alias('detective-conan-the-private-eyes-requiem', 'native', '名探偵コナン 探偵たちの鎮魂歌');
select catalog.ext('detective-conan-the-private-eyes-requiem', 'wikidata', 'Q1317462');
select catalog.ext('detective-conan-the-private-eyes-requiem', 'anilist_anime', '1506');
select catalog.ext('detective-conan-the-private-eyes-requiem', 'mal_anime', '1506');

select catalog.add_entry('detective-conan-jolly-roger-in-the-deep-azure', 'anime', 'movie', 'Detective Conan: Jolly Roger in the Deep Azure', '2007-04-21', 'JP');
select catalog.alias('detective-conan-jolly-roger-in-the-deep-azure', 'native', '名探偵コナン 紺碧の棺');
select catalog.ext('detective-conan-jolly-roger-in-the-deep-azure', 'wikidata', 'Q1300513');
select catalog.ext('detective-conan-jolly-roger-in-the-deep-azure', 'anilist_anime', '2171');
select catalog.ext('detective-conan-jolly-roger-in-the-deep-azure', 'mal_anime', '2171');

select catalog.add_entry('detective-conan-full-score-of-fear', 'anime', 'movie', 'Detective Conan: Full Score of Fear', '2008-04-19', 'JP');
select catalog.alias('detective-conan-full-score-of-fear', 'native', '名探偵コナン 戦慄の楽譜');
select catalog.ext('detective-conan-full-score-of-fear', 'wikidata', 'Q1300488');
select catalog.ext('detective-conan-full-score-of-fear', 'anilist_anime', '4447');
select catalog.ext('detective-conan-full-score-of-fear', 'mal_anime', '4447');

select catalog.add_entry('detective-conan-the-raven-chaser', 'anime', 'movie', 'Detective Conan: The Raven Chaser', '2009-04-18', 'JP');
select catalog.alias('detective-conan-the-raven-chaser', 'native', '名探偵コナン 漆黒の追跡者');
select catalog.ext('detective-conan-the-raven-chaser', 'wikidata', 'Q1207394');
select catalog.ext('detective-conan-the-raven-chaser', 'anilist_anime', '5460');
select catalog.ext('detective-conan-the-raven-chaser', 'mal_anime', '5460');

select catalog.add_entry('detective-conan-the-lost-ship-in-the-sky', 'anime', 'movie', 'Detective Conan: The Lost Ship in the Sky', '2010-04-17', 'JP');
select catalog.alias('detective-conan-the-lost-ship-in-the-sky', 'native', '名探偵コナン 天空の難破船');
select catalog.ext('detective-conan-the-lost-ship-in-the-sky', 'wikidata', 'Q1204889');
select catalog.ext('detective-conan-the-lost-ship-in-the-sky', 'anilist_anime', '6467');
select catalog.ext('detective-conan-the-lost-ship-in-the-sky', 'mal_anime', '6467');

select catalog.add_entry('detective-conan-quarter-of-silence', 'anime', 'movie', 'Detective Conan: Quarter of Silence', '2011-04-16', 'JP');
select catalog.alias('detective-conan-quarter-of-silence', 'native', '名探偵コナン 沈黙の15分');
select catalog.ext('detective-conan-quarter-of-silence', 'wikidata', 'Q512262');
select catalog.ext('detective-conan-quarter-of-silence', 'anilist_anime', '9963');
select catalog.ext('detective-conan-quarter-of-silence', 'mal_anime', '9963');

select catalog.add_entry('detective-conan-the-eleventh-striker', 'anime', 'movie', 'Detective Conan: The Eleventh Striker', '2012-04-14', 'JP');
select catalog.alias('detective-conan-the-eleventh-striker', 'native', '名探偵コナン 11人目のストライカー');
select catalog.ext('detective-conan-the-eleventh-striker', 'wikidata', 'Q717510');
select catalog.ext('detective-conan-the-eleventh-striker', 'anilist_anime', '12117');
select catalog.ext('detective-conan-the-eleventh-striker', 'mal_anime', '12117');

select catalog.add_entry('detective-conan-private-eye-in-the-distant-sea', 'anime', 'movie', 'Detective Conan: Private Eye in the Distant Sea', '2013-04-20', 'JP');
select catalog.alias('detective-conan-private-eye-in-the-distant-sea', 'native', '名探偵コナン 絶海の探偵');
select catalog.ext('detective-conan-private-eye-in-the-distant-sea', 'wikidata', 'Q218516');
select catalog.ext('detective-conan-private-eye-in-the-distant-sea', 'anilist_anime', '14735');
select catalog.ext('detective-conan-private-eye-in-the-distant-sea', 'mal_anime', '14735');

select catalog.add_entry('detective-conan-dimensional-sniper', 'anime', 'movie', 'Detective Conan: Dimensional Sniper', '2014-04-19', 'JP');
select catalog.alias('detective-conan-dimensional-sniper', 'native', '名探偵コナン 異次元の狙撃手');
select catalog.ext('detective-conan-dimensional-sniper', 'wikidata', 'Q15279283');
select catalog.ext('detective-conan-dimensional-sniper', 'anilist_anime', '20546');
select catalog.ext('detective-conan-dimensional-sniper', 'mal_anime', '21419');

select catalog.add_entry('detective-conan-sunflowers-of-inferno', 'anime', 'movie', 'Detective Conan: Sunflowers of Inferno', '2015-04-18', 'JP');
select catalog.alias('detective-conan-sunflowers-of-inferno', 'native', '名探偵コナン 業火の向日葵');
select catalog.ext('detective-conan-sunflowers-of-inferno', 'wikidata', 'Q18601583');
select catalog.ext('detective-conan-sunflowers-of-inferno', 'anilist_anime', '21646');
select catalog.ext('detective-conan-sunflowers-of-inferno', 'mal_anime', '28479');

select catalog.add_entry('detective-conan-the-darkest-nightmare', 'anime', 'movie', 'Detective Conan: The Darkest Nightmare', '2016-04-16', 'JP');
select catalog.alias('detective-conan-the-darkest-nightmare', 'native', '名探偵コナン 純黒の悪夢');
select catalog.ext('detective-conan-the-darkest-nightmare', 'wikidata', 'Q21600367');
select catalog.ext('detective-conan-the-darkest-nightmare', 'anilist_anime', '21470');
select catalog.ext('detective-conan-the-darkest-nightmare', 'mal_anime', '32005');

select catalog.add_entry('detective-conan-crimson-love-letter', 'anime', 'movie', 'Detective Conan: Crimson Love Letter', '2017-04-15', 'JP');
select catalog.alias('detective-conan-crimson-love-letter', 'native', '名探偵コナン から紅の恋歌');
select catalog.ext('detective-conan-crimson-love-letter', 'wikidata', 'Q27918415');
select catalog.ext('detective-conan-crimson-love-letter', 'anilist_anime', '98222');
select catalog.ext('detective-conan-crimson-love-letter', 'mal_anime', '34430');

select catalog.add_entry('detective-conan-zero-the-enforcer', 'anime', 'movie', 'Detective Conan: Zero the Enforcer', '2018-04-13', 'JP');
select catalog.alias('detective-conan-zero-the-enforcer', 'native', '名探偵コナン ゼロの執行人');
select catalog.ext('detective-conan-zero-the-enforcer', 'wikidata', 'Q43984092');
select catalog.ext('detective-conan-zero-the-enforcer', 'anilist_anime', '100653');
select catalog.ext('detective-conan-zero-the-enforcer', 'mal_anime', '35798');

select catalog.add_entry('detective-conan-the-fist-of-blue-sapphire', 'anime', 'movie', 'Detective Conan: The Fist of Blue Sapphire', '2019-04-12', 'JP');
select catalog.alias('detective-conan-the-fist-of-blue-sapphire', 'native', '名探偵コナン 紺青の拳');
select catalog.ext('detective-conan-the-fist-of-blue-sapphire', 'wikidata', 'Q59418140');
select catalog.ext('detective-conan-the-fist-of-blue-sapphire', 'anilist_anime', '106206');
select catalog.ext('detective-conan-the-fist-of-blue-sapphire', 'mal_anime', '38770');

select catalog.add_entry('detective-conan-the-scarlet-bullet', 'anime', 'movie', 'Detective Conan: The Scarlet Bullet', '2021-04-16', 'JP');
select catalog.alias('detective-conan-the-scarlet-bullet', 'native', '名探偵コナン 緋色の弾丸');
select catalog.ext('detective-conan-the-scarlet-bullet', 'wikidata', 'Q76769702');
select catalog.ext('detective-conan-the-scarlet-bullet', 'anilist_anime', '113653');
select catalog.ext('detective-conan-the-scarlet-bullet', 'mal_anime', '39764');

select catalog.add_entry('detective-conan-the-bride-of-halloween', 'anime', 'movie', 'Detective Conan: The Bride of Halloween', '2022-04-15', 'JP');
select catalog.alias('detective-conan-the-bride-of-halloween', 'native', '名探偵コナン ハロウィンの花嫁');
select catalog.ext('detective-conan-the-bride-of-halloween', 'wikidata', 'Q109844247');
select catalog.ext('detective-conan-the-bride-of-halloween', 'anilist_anime', '142219');
select catalog.ext('detective-conan-the-bride-of-halloween', 'mal_anime', '49320');

select catalog.add_entry('detective-conan-black-iron-submarine', 'anime', 'movie', 'Detective Conan: Black Iron Submarine', '2023-04-14', 'JP');
select catalog.alias('detective-conan-black-iron-submarine', 'native', '名探偵コナン 黒鉄の魚影');
select catalog.ext('detective-conan-black-iron-submarine', 'wikidata', 'Q115182766');
select catalog.ext('detective-conan-black-iron-submarine', 'anilist_anime', '156841');
select catalog.ext('detective-conan-black-iron-submarine', 'mal_anime', '53540');

select catalog.add_entry('detective-conan-the-million-dollar-pentagram', 'anime', 'movie', 'Detective Conan: The Million-dollar Pentagram', '2024-04-12', 'JP');
select catalog.alias('detective-conan-the-million-dollar-pentagram', 'native', '名探偵コナン 100万ドルの五稜星');
select catalog.ext('detective-conan-the-million-dollar-pentagram', 'wikidata', 'Q123354230');
select catalog.ext('detective-conan-the-million-dollar-pentagram', 'anilist_anime', '169754');
select catalog.ext('detective-conan-the-million-dollar-pentagram', 'mal_anime', '56785');

select catalog.add_entry('detective-conan-one-eyed-flashback', 'anime', 'movie', 'Detective Conan: One-Eyed Flashback', '2025-04-18', 'JP');
select catalog.alias('detective-conan-one-eyed-flashback', 'native', '名探偵コナン 隻眼の残像');
select catalog.ext('detective-conan-one-eyed-flashback', 'wikidata', 'Q131358956');
select catalog.ext('detective-conan-one-eyed-flashback', 'anilist_anime', '185212');
select catalog.ext('detective-conan-one-eyed-flashback', 'mal_anime', '60445');

select catalog.add_entry('detective-conan-fallen-angel-of-the-highway', 'anime', 'movie', 'Detective Conan: Fallen Angel of the Highway', '2026-04-10', 'JP');
select catalog.alias('detective-conan-fallen-angel-of-the-highway', 'native', '名探偵コナン ハイウェイの堕天使');
select catalog.ext('detective-conan-fallen-angel-of-the-highway', 'wikidata', 'Q136049243');
select catalog.ext('detective-conan-fallen-angel-of-the-highway', 'anilist_anime', '198369');
select catalog.ext('detective-conan-fallen-angel-of-the-highway', 'mal_anime', '62387');

insert into catalog.membership (group_id, member_id, position_label)
select catalog.id('detective-conan-films'), catalog.id(slug), label from (values
  ('detective-conan-the-time-bombed-skyscraper', '#1'),
  ('detective-conan-the-fourteenth-target', '#2'),
  ('detective-conan-the-last-wizard-of-the-century', '#3'),
  ('detective-conan-captured-in-her-eyes', '#4'),
  ('detective-conan-countdown-to-heaven', '#5'),
  ('detective-conan-the-phantom-of-baker-street', '#6'),
  ('detective-conan-crossroad-in-the-ancient-capital', '#7'),
  ('detective-conan-magician-of-the-silver-sky', '#8'),
  ('detective-conan-strategy-above-the-depths', '#9'),
  ('detective-conan-the-private-eyes-requiem', '#10'),
  ('detective-conan-jolly-roger-in-the-deep-azure', '#11'),
  ('detective-conan-full-score-of-fear', '#12'),
  ('detective-conan-the-raven-chaser', '#13'),
  ('detective-conan-the-lost-ship-in-the-sky', '#14'),
  ('detective-conan-quarter-of-silence', '#15'),
  ('detective-conan-the-eleventh-striker', '#16'),
  ('detective-conan-private-eye-in-the-distant-sea', '#17'),
  ('detective-conan-dimensional-sniper', '#18'),
  ('detective-conan-sunflowers-of-inferno', '#19'),
  ('detective-conan-the-darkest-nightmare', '#20'),
  ('detective-conan-crimson-love-letter', '#21'),
  ('detective-conan-zero-the-enforcer', '#22'),
  ('detective-conan-the-fist-of-blue-sapphire', '#23'),
  ('detective-conan-the-scarlet-bullet', '#24'),
  ('detective-conan-the-bride-of-halloween', '#25'),
  ('detective-conan-black-iron-submarine', '#26'),
  ('detective-conan-the-million-dollar-pentagram', '#27'),
  ('detective-conan-one-eyed-flashback', '#28'),
  ('detective-conan-fallen-angel-of-the-highway', '#29')
) as film (slug, label);

-- A crossover film, outside the numbered run.

select catalog.add_entry('lupin-the-3rd-vs-detective-conan-the-movie', 'anime', 'movie', 'Lupin the 3rd vs. Detective Conan: The Movie', '2013-12-07', 'JP');
select catalog.alias('lupin-the-3rd-vs-detective-conan-the-movie', 'native', 'ルパン三世VS名探偵コナン THE MOVIE');
select catalog.ext('lupin-the-3rd-vs-detective-conan-the-movie', 'wikidata', 'Q13417248');
select catalog.ext('lupin-the-3rd-vs-detective-conan-the-movie', 'anilist_anime', '18429');
select catalog.ext('lupin-the-3rd-vs-detective-conan-the-movie', 'mal_anime', '18429');
select catalog.member('detective-conan', 'lupin-the-3rd-vs-detective-conan-the-movie');

-- Spin-off manga. Start dates only where Wikidata has one.

select catalog.add_entry('detective-conan-short-stories', 'manga', 'serial', 'Detective Conan Short Stories', '1997', 'JP', 'chapter');
select catalog.alias('detective-conan-short-stories', 'native', '名探偵コナン特別編');
select catalog.ext('detective-conan-short-stories', 'wikidata', 'Q1201105');
select catalog.ext('detective-conan-short-stories', 'anilist_manga', '33202');
select catalog.ext('detective-conan-short-stories', 'mal_manga', '3202');
select catalog.relate('detective-conan-short-stories', 'spin_off_of', 'detective-conan-manga');
select catalog.member('detective-conan', 'detective-conan-short-stories');

select catalog.add_entry('detective-conan-the-culprit-hanzawa', 'manga', 'serial', 'Detective Conan: The Culprit Hanzawa', null, 'JP', 'chapter');
select catalog.alias('detective-conan-the-culprit-hanzawa', 'native', '名探偵コナン 犯人の犯沢さん');
select catalog.ext('detective-conan-the-culprit-hanzawa', 'wikidata', 'Q97220908');
select catalog.ext('detective-conan-the-culprit-hanzawa', 'anilist_manga', '106015');
select catalog.ext('detective-conan-the-culprit-hanzawa', 'mal_manga', '112744');
select catalog.relate('detective-conan-the-culprit-hanzawa', 'spin_off_of', 'detective-conan-manga');
select catalog.member('detective-conan', 'detective-conan-the-culprit-hanzawa');

select catalog.add_entry('detective-conan-zeros-tea-time', 'manga', 'serial', 'Detective Conan: Zero''s Tea Time', '2018-05-09', 'JP', 'chapter');
select catalog.alias('detective-conan-zeros-tea-time', 'native', '名探偵コナン ゼロの日常');
select catalog.ext('detective-conan-zeros-tea-time', 'wikidata', 'Q64797224');
select catalog.ext('detective-conan-zeros-tea-time', 'anilist_manga', '102750');
select catalog.ext('detective-conan-zeros-tea-time', 'mal_manga', '114379');
select catalog.relate('detective-conan-zeros-tea-time', 'spin_off_of', 'detective-conan-manga');
select catalog.member('detective-conan', 'detective-conan-zeros-tea-time');

select catalog.add_entry('detective-conan-wild-police-story', 'manga', 'serial', 'Detective Conan: Wild Police Story', '2019-10-02', 'JP', 'chapter');
select catalog.alias('detective-conan-wild-police-story', 'native', '名探偵コナン 警察学校編 Wild Police Story');
select catalog.ext('detective-conan-wild-police-story', 'wikidata', 'Q104181129');
select catalog.ext('detective-conan-wild-police-story', 'anilist_manga', '112094');
select catalog.ext('detective-conan-wild-police-story', 'mal_manga', '121951');
select catalog.relate('detective-conan-wild-police-story', 'spin_off_of', 'detective-conan-manga');
select catalog.member('detective-conan', 'detective-conan-wild-police-story');

