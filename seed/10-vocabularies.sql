-- Vocabularies.

-- Visual novels are read, not played, hence their own medium.
insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating, part_label) values
  ('anime',        'anime',        'watching', 'plan to watch', 'rewatching', null),
  ('tv',           'TV',           'watching', 'plan to watch', 'rewatching', 'season'),
  ('film',         'film',         'watching', 'plan to watch', 'rewatching', null),
  ('manga',        'manga',        'reading',  'plan to read',  'rereading',  'volume'),
  ('comic',        'comic',        'reading',  'plan to read',  'rereading',  'volume'),
  ('book',         'book',         'reading',  'plan to read',  'rereading',  null),
  ('visual_novel', 'visual novel', 'reading',  'plan to read',  'rereading',  null),
  ('game',         'game',         'playing',  'plan to play',  'replaying',  null),
  ('audio',        'audio',        'listening', 'plan to listen', 'relistening', null);

insert into catalog.medium_format (medium, format, label) values
  ('anime', 'tv', 'anime'), ('anime', 'movie', 'film'), ('anime', 'ova', 'OVA'), ('anime', 'ona', 'ONA'),
  ('tv', 'series', 'TV'), ('film', 'feature', 'film'), ('film', 'tv_movie', 'TV film'),
  ('manga', 'serial', 'manga'), ('manga', 'one_shot', 'one-shot'),
  ('book', 'novel', 'novel'), ('book', 'light_novel', 'light novel'),
  ('visual_novel', 'full', 'visual novel'), ('game', 'full', 'game'),
  ('anime', 'special', 'special'), ('tv', 'special', 'special'),
  ('comic', 'serial', 'comic'), ('book', 'short_story', 'short story'),
  ('comic', 'one_shot', 'one-shot'), ('comic', 'graphic_novel', 'graphic novel'),
  ('book', 'nonfiction', 'nonfiction'), ('book', 'anthology', 'anthology'),
  ('book', 'gamebook', 'gamebook'), ('game', 'expansion', 'expansion'),
  ('film', 'stage', 'recorded performance'), ('film', 'short', 'short film'),
  ('audio', 'drama', 'audio drama'), ('audio', 'podcast', 'podcast'),
  ('audio', 'audiobook', 'audiobook');

insert into catalog.relation_kind (kind, label, inverse_label, is_symmetric, family) values
  ('sequel_of',              'sequel of',              'prequel of',              false, 'sequence'),
  ('side_story_of',          'side story of',          'has side story',          false, 'sequence'),
  ('spin_off_of',            'spin-off of',            'has spin-off',            false, 'sequence'),
  ('adaptation_of',          'adaptation of',          'adapted as',              false, 'derivative'),
  ('remake_of',              'remake of',              'remade as',               false, 'derivative'),
  ('reboot_of',              'reboot of',              'rebooted as',             false, 'derivative'),
  ('remaster_of',            'remaster of',            'remastered as',           false, 'version'),
  ('recap_of',               'recap of',               'recapped by',             false, 'version'),
  ('compilation_of',         'compilation of',         'included in compilation', false, 'version'),
  ('fan_work_of',            'fan work of',            'has fan work',            false, 'derivative'),
  ('alternative_version_of', 'alternative version of', null,                      true,  'version'),
  ('crossover_with',         'crossover with',         null,                      true,  'setting'),
  ('same_setting_as',        'same setting as',        null,                      true,  'setting'),
  ('references',             'references',             'referenced by',           false, 'setting');

insert into catalog.id_source (source, site, url_template) values
  ('wikidata',      'Wikidata',    'https://www.wikidata.org/wiki/{id}'),
  ('anilist_anime', 'AniList',     'https://anilist.co/anime/{id}'),
  ('anilist_manga', 'AniList',     'https://anilist.co/manga/{id}'),
  ('mal_anime',     'MyAnimeList', 'https://myanimelist.net/anime/{id}'),
  ('mal_manga',     'MyAnimeList', 'https://myanimelist.net/manga/{id}'),
  ('vndb',          'VNDB',        'https://vndb.org/{id}'),
  ('imdb',          'IMDb',        'https://www.imdb.com/title/{id}/'),
  ('tmdb_movie',    'TMDB',        'https://www.themoviedb.org/movie/{id}'),
  ('tmdb_tv',       'TMDB',        'https://www.themoviedb.org/tv/{id}'),
  ('isbn',          'ISBN',        'https://openlibrary.org/isbn/{id}'),
  ('steam',         'Steam',       'https://store.steampowered.com/app/{id}/'),
  ('mobygames',     'MobyGames',   'https://www.mobygames.com/game/{id}/'),
  ('anidb',         'AniDB',       'https://anidb.net/anime/{id}'),
  ('comicvine',     'Comic Vine',  'https://comicvine.gamespot.com/wd/{id}/'),
  ('igdb',          'IGDB',        'https://www.igdb.com/games/{id}'),
  ('vgmdb',         'VGMdb',       'https://vgmdb.net/product/{id}');

