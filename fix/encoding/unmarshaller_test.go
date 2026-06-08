package encoding

import (
	"bytes"
	"strconv"
	"testing"

	"gitlab.b2broker.tech/highload/b2connect/libs/go/simplefix-go/fix"
)

const visibleDelimiter = "|"

func TestUnmarshalItems(t *testing.T) {
	var exampleRawMas = []byte("8=FIX.4.49=6135=A49=Client56=Server34=152=20210706-19:06:12.838108=510=206")

	testData := []struct {
		raw    []byte
		expect []byte
		items  fix.Items
	}{
		{
			raw: exampleRawMas,
			items: fix.Items{
				&fix.KeyValue{Key: "8", Value: &fix.String{}},
				&fix.KeyValue{Key: "9", Value: &fix.String{}},
				&fix.KeyValue{Key: "35", Value: &fix.String{}},
				&fix.KeyValue{Key: "49", Value: &fix.String{}},
				&fix.KeyValue{Key: "56", Value: &fix.String{}},
				&fix.KeyValue{Key: "34", Value: &fix.String{}},
				&fix.KeyValue{Key: "52", Value: &fix.String{}},
				&fix.KeyValue{Key: "108", Value: &fix.String{}},
				&fix.KeyValue{Key: "10", Value: &fix.String{}},
			},
			expect: exampleRawMas,
		},
		{
			raw: exampleRawMas,
			items: fix.Items{
				&fix.KeyValue{Key: "8", Value: &fix.String{}},
				&fix.KeyValue{Key: "35", Value: &fix.String{}},
				&fix.KeyValue{Key: "9", Value: &fix.String{}},
				&fix.KeyValue{Key: "56", Value: &fix.String{}},
				&fix.KeyValue{Key: "49", Value: &fix.String{}},
				&fix.KeyValue{Key: "34", Value: &fix.String{}},
				&fix.KeyValue{Key: "52", Value: &fix.String{}},
				&fix.KeyValue{Key: "100500", Value: &fix.String{}},
				&fix.KeyValue{Key: "108", Value: &fix.String{}},
				&fix.KeyValue{Key: "10", Value: &fix.String{}},
			},
			expect: []byte("8=FIX.4.435=A9=6156=Server49=Client34=152=20210706-19:06:12.838108=510=206"),
		},
		{
			raw: exampleRawMas,
			items: fix.Items{
				&fix.KeyValue{Key: "8", Value: &fix.String{}},
				&fix.KeyValue{Key: "9", Value: &fix.String{}},
				&fix.KeyValue{Key: "35", Value: &fix.String{}},
				&fix.KeyValue{Key: "56", Value: &fix.String{}},
				&fix.KeyValue{Key: "49", Value: &fix.String{}},
				&fix.KeyValue{Key: "34", Value: &fix.String{}},
				&fix.KeyValue{Key: "108", Value: &fix.String{}},
				&fix.KeyValue{Key: "52", Value: &fix.String{}},
				&fix.KeyValue{Key: "10", Value: &fix.String{}},
			},
			expect: []byte("8=FIX.4.49=6135=A56=Server49=Client34=1108=552=20210706-19:06:12.83810=206"),
		},
		{
			raw: []byte("8=FIX.4.49=6135=A56=Server161=4162=A163=1162=B163=2162=C162=DE163=310=206"),
			items: fix.Items{
				&fix.KeyValue{Key: "8", Value: &fix.String{}},
				&fix.KeyValue{Key: "9", Value: &fix.String{}},
				&fix.KeyValue{Key: "35", Value: &fix.String{}},
				&fix.KeyValue{Key: "56", Value: &fix.String{}},
				fix.NewGroup("161", &fix.KeyValue{Key: "162", Value: &fix.String{}}, &fix.KeyValue{Key: "163", Value: &fix.String{}}),
				&fix.KeyValue{Key: "10", Value: &fix.String{}},
			},
			expect: []byte("8=FIX.4.49=6135=A56=Server161=4162=A163=1162=B163=2162=C162=DE163=310=206"),
		},
		{
			raw: []byte("8=FIX.4.49=6135=A56=Server161=4162=A163=1162=B163=2162=C162=DE163=310=206"),
			items: fix.Items{
				&fix.KeyValue{Key: "8", Value: &fix.String{}},
				&fix.KeyValue{Key: "9", Value: &fix.String{}},
				&fix.KeyValue{Key: "35", Value: &fix.String{}},
				&fix.KeyValue{Key: "56", Value: &fix.String{}},
				fix.NewGroup("161", fix.NewComponent(
					&fix.KeyValue{Key: "162", Value: &fix.String{}},
					&fix.KeyValue{Key: "163", Value: &fix.String{}},
				)),
				&fix.KeyValue{Key: "10", Value: &fix.String{}},
			},
			expect: []byte("8=FIX.4.49=6135=A56=Server161=4162=A163=1162=B163=2162=C162=DE163=310=206"),
		},
		{
			raw: []byte("8=FIX.4.49=18335=i34=8649=XC22952=20220314-11:39:13.38556=Q048296=2302=62295=1299=0134=1.0135=1.0188=4588680190=4591680302=64295=1299=0134=100000135=100000188=1.39851190=1.3986510=085"),
			items: fix.Items{
				&fix.KeyValue{Key: "8", Value: &fix.String{}},
				&fix.KeyValue{Key: "9", Value: &fix.String{}},
				&fix.KeyValue{Key: "35", Value: &fix.String{}},
				&fix.KeyValue{Key: "34", Value: &fix.String{}},
				&fix.KeyValue{Key: "49", Value: &fix.String{}},
				&fix.KeyValue{Key: "52", Value: &fix.String{}},
				&fix.KeyValue{Key: "56", Value: &fix.String{}},
				fix.NewGroup("296", fix.NewComponent(
					&fix.KeyValue{Key: "302", Value: &fix.String{}},
					fix.NewComponent(
						&fix.KeyValue{Key: "295", Value: &fix.String{}},
						&fix.KeyValue{Key: "299", Value: &fix.String{}},
						&fix.KeyValue{Key: "134", Value: &fix.Float{}},
						&fix.KeyValue{Key: "135", Value: &fix.Float{}},
						&fix.KeyValue{Key: "188", Value: &fix.Float{}},
						&fix.KeyValue{Key: "190", Value: &fix.Float{}},
					),
				)),
				&fix.KeyValue{Key: "10", Value: &fix.String{}},
			},
			expect: []byte("8=FIX.4.49=18335=i34=8649=XC22952=20220314-11:39:13.38556=Q048296=2302=62295=1299=0134=1.0135=1.0188=4588680190=4591680302=64295=1299=0134=100000135=100000188=1.39851190=1.3986510=085"),
		},
	}

	for i, item := range testData {
		err := unmarshalItems(item.items, item.raw, true)
		if err != nil {
			t.Fatal(err)
		}

		res := item.items.ToBytes()
		if !bytes.Equal(item.expect, res) {
			t.Logf("%d. expect %s (%d)", i, showDelimiter(item.expect), len(item.expect))
			t.Logf("%d. result %s (%d)", i, showDelimiter(res), len(res))
			t.Fatal("the result is not equal to expected message")
		}
	}

}

func showDelimiter(in []byte) []byte {
	return bytes.ReplaceAll(in, fix.Delimiter, []byte(visibleDelimiter))
}

// TestUnmarshalItems_EmptyRepeatingGroup: a zero-count group (146=0) must parse
// as empty, not fail "wrong items count: 0 != 1".
func TestUnmarshalItems_EmptyRepeatingGroup(t *testing.T) {
	// BodyLength/CheckSum aren't validated by unmarshalItems.
	raw := []byte("8=FIX.4.4\x019=22\x0135=y\x0156=Server\x01146=0\x0110=000\x01")

	group := fix.NewGroup("146", &fix.KeyValue{Key: "55", Value: &fix.String{}})
	items := fix.Items{
		&fix.KeyValue{Key: "8", Value: &fix.String{}},
		&fix.KeyValue{Key: "9", Value: &fix.String{}},
		&fix.KeyValue{Key: "35", Value: &fix.String{}},
		&fix.KeyValue{Key: "56", Value: &fix.String{}},
		group,
		&fix.KeyValue{Key: "10", Value: &fix.String{}},
	}

	if err := unmarshalItems(items, raw, true); err != nil {
		t.Fatalf("unmarshalItems returned error for empty group: %v", err)
	}
	if !group.IsEmpty() {
		t.Fatalf("expected empty group, got %d entries", len(group.Entries()))
	}
}

// TestUnmarshal_EmptyRepeatingGroupFullMessage exercises the full strict
// Unmarshal path (BodyLength + CheckSum validation) for a message whose
// repeating group is empty (146=0), as a counterparty message would arrive.
func TestUnmarshal_EmptyRepeatingGroupFullMessage(t *testing.T) {
	body := "35=y\x0149=Sender\x0156=Target\x0134=1\x01146=0\x01"
	head := "8=FIX.4.4\x019=" + strconv.Itoa(len(body)) + "\x01"
	pre := head + body
	sum := fix.CalcCheckSumOptimized([]byte(pre[:len(pre)-1]))
	raw := []byte(pre + "10=" + string(sum) + "\x01")

	group := fix.NewGroup("146", &fix.KeyValue{Key: "55", Value: &fix.String{}})
	msg := fix.NewMessage("8", "9", "10", "35", "FIX.4.4", "y")
	msg.SetHeader(fix.NewComponent())
	msg.SetTrailer(fix.NewComponent())
	msg.SetBody(
		&fix.KeyValue{Key: "49", Value: &fix.String{}},
		&fix.KeyValue{Key: "56", Value: &fix.String{}},
		&fix.KeyValue{Key: "34", Value: &fix.String{}},
		group,
	)

	if err := Unmarshal(msg, raw); err != nil {
		t.Fatalf("Unmarshal failed for empty group: %v", err)
	}
	if !group.IsEmpty() {
		t.Fatalf("expected empty group, got %d entries", len(group.Entries()))
	}
}

// TestUnmarshal_ValidationErrors covers the strict validateRaw rejection paths.
func TestUnmarshal_ValidationErrors(t *testing.T) {
	body := "35=y\x0149=Sender\x0156=Target\x0134=1\x01146=0\x01"
	head := "8=FIX.4.4\x019=" + strconv.Itoa(len(body)) + "\x01"
	pre := head + body
	validSum := string(fix.CalcCheckSumOptimized([]byte(pre[:len(pre)-1])))

	newMsg := func() *fix.Message {
		m := fix.NewMessage("8", "9", "10", "35", "FIX.4.4", "y")
		m.SetHeader(fix.NewComponent())
		m.SetTrailer(fix.NewComponent())
		m.SetBody(
			&fix.KeyValue{Key: "49", Value: &fix.String{}},
			&fix.KeyValue{Key: "56", Value: &fix.String{}},
			fix.NewGroup("146", &fix.KeyValue{Key: "55", Value: &fix.String{}}),
		)
		return m
	}

	tests := []struct {
		name string
		raw  []byte
	}{
		{"checksum mismatch", []byte(pre + "10=000\x01")},
		{"body length mismatch", []byte("8=FIX.4.4\x019=" + strconv.Itoa(len(body)+7) + "\x01" + body + "10=" + validSum + "\x01")},
		{"non-numeric body length", []byte("8=FIX.4.4\x019=abc\x0135=y\x0156=Server\x0110=000\x01")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := Unmarshal(newMsg(), tc.raw); err == nil {
				t.Fatalf("expected error for %q, got nil", tc.name)
			}
		})
	}
}
