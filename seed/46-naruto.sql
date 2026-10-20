-- Naruto: a series and its sequel.

select catalog.add_group('naruto-franchise', 'franchise', 'Naruto', p_public_id => 96);
select catalog.add_entry('naruto', 'anime', 'tv', 'Naruto', '2002-10', 'JP', 'episode', 220, p_title_kind => 'romanized', p_public_id => 97);
select catalog.add_entry('naruto-shippuden', 'anime', 'tv', 'Naruto Shippuden', '2007-02', 'JP', 'episode', 500, p_public_id => 98);
select catalog.alias('naruto-shippuden', 'romanized', 'Naruto: Shippuuden');
select catalog.member('naruto-franchise', 'naruto', 'naruto-shippuden');
select catalog.ext('naruto-franchise', 'wikidata', 'Q642');
select catalog.ext('naruto', 'wikidata', 'Q25929253');
select catalog.ext('naruto-shippuden', 'wikidata', 'Q12500134');
select catalog.relate('naruto-shippuden', 'sequel_of', 'naruto');
