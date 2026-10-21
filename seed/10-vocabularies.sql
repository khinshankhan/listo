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
  ('game',         'game',         'playing',  'plan to play',  'replaying',  null);

insert into catalog.medium_format (medium, format, label) values
  ('anime', 'tv', 'anime'), ('anime', 'movie', 'film'), ('anime', 'ova', 'OVA'), ('anime', 'ona', 'ONA'),
  ('tv', 'series', 'TV'), ('film', 'feature', 'film'),
  ('manga', 'serial', 'manga'), ('book', 'novel', 'novel'), ('book', 'light_novel', 'light novel'),
  ('visual_novel', 'full', 'visual novel'), ('game', 'full', 'game');

insert into catalog.relation_kind (kind, label, inverse_label, is_symmetric, family) values
  ('sequel_of',              'sequel of',              'prequel of',              false, 'sequence'),
  ('side_story_of',          'side story of',          'has side story',          false, 'sequence'),
  ('spin_off_of',            'spin-off of',            'has spin-off',            false, 'sequence'),
  ('adaptation_of',          'adaptation of',          'adapted as',              false, 'derivative'),
  ('remake_of',              'remake of',              'remade as',               false, 'derivative'),
  ('reboot_of',              'reboot of',              'rebooted as',             false, 'derivative'),
  ('remaster_of',            'remaster of',            'remastered as',           false, 'version'),
  ('recap_of',               'recap of',               'recapped by',             false, 'version'),
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
  ('tmdb_tv',       'TMDB',        'https://www.themoviedb.org/tv/{id}');

