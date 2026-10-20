-- Tags, once every entry they name exists.

select catalog.tag_with('genre', 'action', 'index-s1', 'index-ii', 'railgun-s1', 'railgun-s', 'naruto',
                        'naruto-shippuden', 'devil-may-cry', 'dmc-devil-may-cry');
select catalog.tag_with('genre', 'comedy', 'konosuba-s1', 'konosuba-s2', 'konosuba-s3');
select catalog.tag_with('genre', 'sci-fi', 'steins-gate-anime', 'steins-gate-0-anime', 'index-s1', 'railgun-s1');
select catalog.tag_with('theme', 'time travel', 'steins-gate-vn', 'steins-gate-anime', 'steins-gate-0-anime');
