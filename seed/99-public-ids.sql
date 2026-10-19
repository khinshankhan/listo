-- seeds pin public ids, so move the identity past them
select setval(pg_get_serial_sequence('catalog.node', 'public_id'), (select max(public_id) from catalog.node));
