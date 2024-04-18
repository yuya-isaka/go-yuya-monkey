package evaluator

import (
	"testing"

	"github.com/yuya-isaka/go-yuya-monkey/object"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			`quote(5)`,
			`5`,
		},
		{
			`quote(5 + 8)`,
			`(5 + 8)`,
		},
		{
			`quote(foobar)`,
			`foobar`,
		},
		{
			`quote(foobar + barfoo)`,
			`(foobar + barfoo)`,
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		quote, ok := obj.(*object.QuoteObj)
		if !ok {
			t.Fatalf("expect *object.QuoteObj. got=%T (%+v)", obj, obj)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}

		if quote.Node.String() != tt.expect {
			t.Errorf("not equal. got=%q, want=%q", quote.Node.String(), tt.expect)
		}
	}
}

func TestQuoteUnQuote(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{
			`quote(unquote(4))`,
			`4`,
		},
		{
			`quote(unquote(4+4))`,
			`8`,
		},
		{
			`quote(8 + unquoe(4 + 4))`,
			`(8 + 8)`,
		},
		{
			`quote(unquote(4 + 4) * 9)`,
			`(8 * 9)`,
		},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		quote, ok := obj.(*object.QuoteObj)
		if !ok {
			t.Fatalf("expected *object.QuoteObj. got=%T (%+v)", obj, obj)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}

		if quote.Node.String() != tt.expect {
			t.Errorf("not equal. got=%q, want=%q", quote.Node.String(), tt.expect)
		}
	}
}
