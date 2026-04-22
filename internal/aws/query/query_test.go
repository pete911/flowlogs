package query

import (
	"strings"
	"testing"
)

func TestNewQuery(t *testing.T) {
	q := NewQuery(200, 30)
	if q.GetLimit() != 200 {
		t.Errorf("limit: got %d, want 200", q.GetLimit())
	}
	if q.GetSinceMinutes() != 30 {
		t.Errorf("sinceMinutes: got %d, want 30", q.GetSinceMinutes())
	}
	if !strings.HasPrefix(q.GetQuery(), "fields ") {
		t.Errorf("query should start with 'fields ', got: %q", q.GetQuery())
	}
}

func TestQueryFilters(t *testing.T) {
	tests := []struct {
		name string
		fn   func(Query) Query
		want string
	}{
		{"NoNoData", func(q Query) Query { return q.NoNoData() }, `| filter logStatus != "NODATA"`},
		{"NoSkipData", func(q Query) Query { return q.NoSkipData() }, `| filter logStatus != "SKIPDATA"`},
		{"InterfaceId", func(q Query) Query { return q.InterfaceId("eni-abc") }, `| filter interfaceId == "eni-abc"`},
		{"Ingress", func(q Query) Query { return q.Ingress() }, `| filter flowDirection == "ingress"`},
		{"Egress", func(q Query) Query { return q.Egress() }, `| filter flowDirection == "egress"`},
		{"Accept", func(q Query) Query { return q.Accept() }, `| filter action == "ACCEPT"`},
		{"Reject", func(q Query) Query { return q.Reject() }, `| filter action == "REJECT"`},
		{"Protocol TCP", func(q Query) Query { return q.Protocol("TCP") }, `| filter protocol == "6"`},
		{"Protocol udp lowercase", func(q Query) Query { return q.Protocol("udp") }, `| filter protocol == "17"`},
		{"Port", func(q Query) Query { return q.Port(443) }, `| filter srcPort == "443" or dstPort == "443"`},
		{"SourcePort", func(q Query) Query { return q.SourcePort(22) }, `| filter srcPort == "22"`},
		{"DestinationPort", func(q Query) Query { return q.DestinationPort(80) }, `| filter dstPort == "80"`},
		{"Address", func(q Query) Query { return q.Address("10.0.0.1") }, `| filter srcAddr == "10.0.0.1" or pktSrcAddr == "10.0.0.1" or dstAddr == "10.0.0.1" or pktDstAddr == "10.0.0.1"`},
		{"SourceAddress", func(q Query) Query { return q.SourceAddress("10.0.0.2") }, `| filter srcAddr == "10.0.0.2"`},
		{"PktSourceAddress", func(q Query) Query { return q.PktSourceAddress("10.0.0.3") }, `| filter pktSrcAddr == "10.0.0.3"`},
		{"DestinationAddress", func(q Query) Query { return q.DestinationAddress("10.0.0.4") }, `| filter dstAddr == "10.0.0.4"`},
		{"PktDestinationAddress", func(q Query) Query { return q.PktDestinationAddress("10.0.0.5") }, `| filter pktDstAddr == "10.0.0.5"`},
		{"Sort", func(q Query) Query { return q.Sort() }, `| sort @timestamp desc`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fn(NewQuery(100, 60)).GetQuery()
			if !strings.Contains(got, tc.want) {
				t.Errorf("query missing clause\n  got:  %q\n  want: %q", got, tc.want)
			}
		})
	}
}

func TestProtocolUnknownKeywordIsNoOp(t *testing.T) {
	base := NewQuery(100, 60)
	got := base.Protocol("not-a-real-protocol").GetQuery()
	if got != base.GetQuery() {
		t.Errorf("unknown protocol should leave query unchanged\n  base: %q\n  got:  %q", base.GetQuery(), got)
	}
}

func TestChainedFiltersJoinedWithNewlines(t *testing.T) {
	got := NewQuery(100, 60).NoNoData().Accept().SourcePort(22).Sort().GetQuery()
	wantClauses := []string{
		"fields ",
		`| filter logStatus != "NODATA"`,
		`| filter action == "ACCEPT"`,
		`| filter srcPort == "22"`,
		`| sort @timestamp desc`,
	}
	for _, c := range wantClauses {
		if !strings.Contains(got, c) {
			t.Errorf("chain missing clause %q\n  got: %q", c, got)
		}
	}
	if lines := strings.Split(got, "\n"); len(lines) != len(wantClauses) {
		t.Errorf("expected %d lines, got %d: %q", len(wantClauses), len(lines), got)
	}
}

func TestQueryBranchingDoesNotAlias(t *testing.T) {
	base := NewQuery(100, 60).NoNoData()
	a := base.Accept()
	b := base.Reject()

	if strings.Contains(a.GetQuery(), "REJECT") {
		t.Errorf("branch A leaked REJECT from branch B: %q", a.GetQuery())
	}
	if strings.Contains(b.GetQuery(), "ACCEPT") {
		t.Errorf("branch B leaked ACCEPT from branch A: %q", b.GetQuery())
	}
	if strings.Contains(base.GetQuery(), "ACCEPT") || strings.Contains(base.GetQuery(), "REJECT") {
		t.Errorf("base mutated by branches: %q", base.GetQuery())
	}
}
