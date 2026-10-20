-- Shameless: a remake of a remake, across countries.

select catalog.add_group('shameless', 'franchise', 'Shameless', p_public_id => 92);
select catalog.add_entry('shameless-uk', 'tv', 'series', 'Shameless', '2004-01', 'GB', p_public_id => 93);
select catalog.add_entry('shameless-us', 'tv', 'series', 'Shameless', '2011-01', 'US', p_public_id => 94);
select catalog.add_entry('bizim-hikaye', 'tv', 'series', 'Bizim Hikaye', '2017', 'TR', p_title_kind => 'native', p_public_id => 95);
select catalog.alias('bizim-hikaye', 'synonym', 'Our Story');
select catalog.member('shameless', 'shameless-uk', 'shameless-us', 'bizim-hikaye');
select catalog.ext('shameless-uk', 'wikidata', 'Q23605');
select catalog.ext('shameless-us', 'wikidata', 'Q23599');
select catalog.ext('bizim-hikaye', 'wikidata', 'Q47010718');
select catalog.relate('shameless-us', 'remake_of', 'shameless-uk');
select catalog.relate('bizim-hikaye', 'remake_of', 'shameless-us');
