-- Ben 10: a film from an alternate timeline that a later series picks back
-- up, so it watches better later than its release slot.

select catalog.add_group('ben-10', 'franchise', 'Ben 10', p_public_id => 64);
select catalog.add_group('ben-10-classic', 'continuity', 'Classic', 'concluded', p_public_id => 65);
select catalog.add_group('ben-10-alt', 'continuity', 'Race Against Time timeline', 'concluded', p_public_id => 66);
select catalog.add_group('ben-10-reboot', 'continuity', '2016 reboot', 'concluded', p_public_id => 67);
select catalog.member('ben-10', 'ben-10-classic', 'ben-10-alt', 'ben-10-reboot');

select catalog.add_entry('ben-10-2005', 'tv', 'series', 'Ben 10', '2005', 'US', p_public_id => 68);
select catalog.add_entry('ben-10-race-against-time', 'film', 'feature', 'Ben 10: Race Against Time', '2007', 'US', p_public_id => 69);
select catalog.add_entry('ben-10-alien-force', 'tv', 'series', 'Ben 10: Alien Force', '2008', 'US', p_public_id => 70);
select catalog.add_entry('ben-10-ultimate-alien', 'tv', 'series', 'Ben 10: Ultimate Alien', '2010', 'US', p_public_id => 71);
select catalog.add_entry('ben-10-omniverse', 'tv', 'series', 'Ben 10: Omniverse', '2012', 'US', p_public_id => 72);
select catalog.add_entry('ben-10-2016', 'tv', 'series', 'Ben 10', '2016', 'US', p_public_id => 73);
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

select catalog.add_ordering('ben-10', 'release', 'release order') as o \gset
select catalog.add_step(:'o', 1, 'ben-10-2005');
select catalog.add_step(:'o', 2, 'ben-10-race-against-time');
select catalog.add_step(:'o', 3, 'ben-10-alien-force');
select catalog.add_step(:'o', 4, 'ben-10-ultimate-alien');
select catalog.add_step(:'o', 5, 'ben-10-omniverse');

select catalog.add_ordering('ben-10', 'recommended', 'watch order') as o \gset
select catalog.add_step(:'o', 1, 'ben-10-2005');
select catalog.add_step(:'o', 2, 'ben-10-alien-force');
select catalog.add_step(:'o', 3, 'ben-10-race-against-time', p_necessity => 'optional',
                        p_note => 'fits by age after the first series; matters once the third picks it up');
select catalog.add_step(:'o', 4, 'ben-10-ultimate-alien');
select catalog.add_step(:'o', 5, 'ben-10-omniverse');
