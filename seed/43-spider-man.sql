-- Spider-Man: crossovers point at whole continuities, and the MCU one sits
-- in both the franchise and the universe.

select catalog.add_group('mcu', 'universe', 'Marvel Cinematic Universe', p_public_id => 74);
select catalog.alias('mcu', 'synonym', 'MCU');
select catalog.add_group('spider-man', 'franchise', 'Spider-Man', p_public_id => 75);
select catalog.add_group('spider-man-raimi', 'continuity', 'Raimi trilogy', 'concluded', p_public_id => 76);
select catalog.add_group('spider-man-webb', 'continuity', 'The Amazing Spider-Man', 'abandoned', p_public_id => 77);
select catalog.add_group('spider-man-mcu', 'continuity', 'MCU Spider-Man', 'active', p_public_id => 78);
select catalog.member('spider-man', 'spider-man-raimi', 'spider-man-webb', 'spider-man-mcu');
select catalog.member('mcu', 'spider-man-mcu');

select catalog.add_entry('spider-man-2002', 'film', 'feature', 'Spider-Man', '2002', 'US', p_public_id => 79);
select catalog.add_entry('spider-man-2', 'film', 'feature', 'Spider-Man 2', '2004', 'US', p_public_id => 80);
select catalog.add_entry('spider-man-3', 'film', 'feature', 'Spider-Man 3', '2007', 'US', p_public_id => 81);
select catalog.add_entry('amazing-spider-man', 'film', 'feature', 'The Amazing Spider-Man', '2012', 'US', p_public_id => 82);
select catalog.add_entry('amazing-spider-man-2', 'film', 'feature', 'The Amazing Spider-Man 2', '2014', 'US', p_public_id => 83);
select catalog.add_entry('spider-man-homecoming', 'film', 'feature', 'Spider-Man: Homecoming', '2017', 'US', p_public_id => 84);
select catalog.add_entry('spider-man-far-from-home', 'film', 'feature', 'Spider-Man: Far From Home', '2019', 'US', p_public_id => 85);
select catalog.add_entry('spider-man-no-way-home', 'film', 'feature', 'Spider-Man: No Way Home', '2021', 'US', p_public_id => 86);
select catalog.member('spider-man-raimi', 'spider-man-2002', 'spider-man-2', 'spider-man-3');
select catalog.member('spider-man-webb', 'amazing-spider-man', 'amazing-spider-man-2');
select catalog.member('spider-man-mcu', 'spider-man-homecoming', 'spider-man-far-from-home', 'spider-man-no-way-home');

select catalog.relate('spider-man-2', 'sequel_of', 'spider-man-2002');
select catalog.relate('spider-man-3', 'sequel_of', 'spider-man-2');
select catalog.relate('amazing-spider-man', 'reboot_of', 'spider-man-2002');
select catalog.relate('amazing-spider-man-2', 'sequel_of', 'amazing-spider-man');
select catalog.relate('spider-man-homecoming', 'reboot_of', 'amazing-spider-man');
select catalog.relate('spider-man-far-from-home', 'sequel_of', 'spider-man-homecoming');
select catalog.relate('spider-man-no-way-home', 'sequel_of', 'spider-man-far-from-home');
select catalog.relate('spider-man-no-way-home', 'crossover_with', 'spider-man-raimi', p_spoiler => true,
                      p_established_on => '2021-12-17', p_established_by => 'spider-man-no-way-home');
select catalog.relate('spider-man-no-way-home', 'crossover_with', 'spider-man-webb', p_spoiler => true,
                      p_established_on => '2021-12-17', p_established_by => 'spider-man-no-way-home');
