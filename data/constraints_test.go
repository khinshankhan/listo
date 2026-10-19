package data_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// One bad write per schema rule. Each must be rejected by its own
// constraint, so dropping or weakening one fails here.
func TestSchemaRejectsBadWrites(t *testing.T) {
	tx := testTx(t)
	ctx := context.Background()

	// per-run names; relation kinds allow letters and underscores only
	run := fmt.Sprint(time.Now().UnixNano())
	letters := strings.Map(func(r rune) rune { return 'a' + (r - '0') }, run)
	medium, source := "test_medium_"+run, "test_source_"+run
	directed, symmetric := "test_sequel_"+letters, "test_crossover_"+letters

	exec := func(q string, args ...any) {
		t.Helper()
		if _, err := tx.Exec(ctx, q, args...); err != nil {
			t.Fatal(err)
		}
	}
	id := func(q string, args ...any) uuid.UUID {
		t.Helper()
		var v uuid.UUID
		if err := tx.QueryRow(ctx, q, args...).Scan(&v); err != nil {
			t.Fatal(err)
		}
		return v
	}
	entry := func(name string) uuid.UUID {
		t.Helper()
		v := id(`insert into catalog.node (type, slug) values ('entry', $1) returning id`, fmt.Sprintf("test-%s-%s", name, run))
		exec(`insert into catalog.entry (id, medium, format, title, unit_kind) values ($1, $2, 'tv', $3, 'episode')`, v, medium, name)
		return v
	}

	exec(`insert into catalog.medium (medium, label, verb_current, verb_planning, verb_repeating)
	      values ($1, 'test medium', 'watching', 'plan to watch', 'rewatching')`, medium)
	exec(`insert into catalog.medium_format (medium, format, label) values ($1, 'tv', 'tv')`, medium)
	exec(`insert into catalog.relation_kind (kind, label, inverse_label, is_symmetric, family) values
	      ($1, 'sequel of', 'prequel of', false, 'sequence'), ($2, 'crossover with', null, true, 'setting')`, directed, symmetric)
	exec(`insert into catalog.id_source (source, site, url_template) values ($1, 'Test', 'https://example.test/{id}')`, source)

	first, second, other, tracked := entry("first"), entry("second"), entry("other"), entry("tracked")
	otherUnit := id(`insert into catalog.unit (entry_id, ordinal, kind, label) values ($1, 1, 'episode', 'Ep. 1') returning id`, other)
	exec(`insert into catalog.external_id (source, external_id, node_id) values ($1, '9253', $2)`, source, first)
	group := id(`insert into catalog.node (type, slug) values ('group', $1) returning id`, "test-group-"+run)
	exec(`insert into catalog.grp (id, kind, name) values ($1, 'franchise', 'Test Franchise')`, group)
	account := id(`insert into tracking.account (username) values ($1) returning id`, "test-"+run)
	exec(`insert into tracking.list_entry (account_id, entry_id, status) values ($1, $2, 'current')`, account, tracked)

	// a symmetric edge is stored smaller id first; this is the other way round
	hi, lo := first, second
	if strings.Compare(hi.String(), lo.String()) < 0 {
		hi, lo = lo, hi
	}

	cases := []struct {
		name       string
		sql        string
		args       []any
		code       string // SQLSTATE: 23505 unique, 23503 foreign key, 23514 check
		constraint string
	}{
		{"a. an external id belongs to one node",
			`insert into catalog.external_id (source, external_id, node_id) values ($1, '9253', $2)`,
			[]any{source, second}, "23505", "external_id_dedupe"},
		{"b. a reference to a node that doesn't exist is an error, not a new node",
			`insert into catalog.relation (from_node, kind, to_node) values ($1, $2, $3)`,
			[]any{uuid.New(), directed, first}, "23503", "relation_from_node_fkey"},
		{"c. an anchor must be a unit of the entry it anchors",
			`insert into catalog.relation (from_node, kind, to_node, to_unit) values ($1, $2, $3, $4)`,
			[]any{first, directed, second, otherUnit}, "23503", "relation_to_unit_on_to_node"},
		{"d. a symmetric edge is stored once, smaller id first",
			`insert into catalog.relation (from_node, kind, to_node, is_symmetric) values ($1, $2, $3, true)`,
			[]any{hi, symmetric, lo}, "23514", "relation_symmetric_stored_once"},
		{"e. a symmetric kind can't be filed as directed",
			`insert into catalog.relation (from_node, kind, to_node) values ($1, $2, $3)`,
			[]any{lo, symmetric, hi}, "23503", "relation_kind_symmetry"},
		{"f. standing belongs to continuities only",
			`with n as (insert into catalog.node (type, slug) values ('group', $1) returning id)
			 insert into catalog.grp (id, kind, name, standing) select id, 'franchise', 'Bad', 'primary' from n`,
			[]any{"test-bad-group-" + run}, "23514", "grp_standing_only_on_continuity"},
		{"g. a format must exist for its medium",
			`with n as (insert into catalog.node (type, slug) values ('entry', $1) returning id)
			 insert into catalog.entry (id, medium, format, title) select id, $2, 'feature', 'Bad' from n`,
			[]any{"test-bad-entry-" + run, medium}, "23503", "entry_format_fits_medium"},
		{"h. a tracked entry can't be deleted out from under a list",
			`delete from catalog.node where id = $1`,
			[]any{tracked}, "23503", "list_entry_entry_id_fkey"},
		{"i. a group can't contain itself",
			`insert into catalog.membership (group_id, member_id) values ($1, $1)`,
			[]any{group}, "23514", "membership_not_in_self"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// a savepoint: the rejected write aborts only itself
			sp, err := tx.Begin(ctx)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				_ = sp.Rollback(ctx)
			}()

			_, err = sp.Exec(ctx, tc.sql, tc.args...)
			var pgErr *pgconn.PgError
			if !errors.As(err, &pgErr) {
				t.Fatalf("write was accepted (err = %v), want %s from %s", err, tc.code, tc.constraint)
			}
			if pgErr.Code != tc.code || pgErr.ConstraintName != tc.constraint {
				t.Errorf("rejected by %s %q (%s), want %s %q", pgErr.Code, pgErr.ConstraintName, pgErr.Message, tc.code, tc.constraint)
			}
		})
	}
}
