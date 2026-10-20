-- Devil May Cry: a reboot in its own continuity.

select catalog.add_group('devil-may-cry-franchise', 'franchise', 'Devil May Cry', p_public_id => 87);
select catalog.add_group('dmc-mainline', 'continuity', 'Mainline', 'primary', p_public_id => 88);
select catalog.add_group('dmc-reboot', 'continuity', 'DmC', 'abandoned', p_public_id => 89);
select catalog.member('devil-may-cry-franchise', 'dmc-mainline', 'dmc-reboot');
select catalog.add_entry('devil-may-cry', 'game', 'full', 'Devil May Cry', '2001', 'JP', p_public_id => 90);
select catalog.add_entry('dmc-devil-may-cry', 'game', 'full', 'DmC: Devil May Cry', '2013', p_public_id => 91);
select catalog.member('dmc-mainline', 'devil-may-cry');
select catalog.member('dmc-reboot', 'dmc-devil-may-cry');
select catalog.ext('devil-may-cry', 'wikidata', 'Q840889');
select catalog.ext('dmc-devil-may-cry', 'wikidata', 'Q1826997');
select catalog.relate('dmc-devil-may-cry', 'reboot_of', 'devil-may-cry');
