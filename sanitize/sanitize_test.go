package sanitize

import (
	"strings"
	"testing"
)

func TestTrimWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no whitespace", "hello", "hello"},
		{"leading spaces", "  hello", "hello"},
		{"trailing spaces", "hello  ", "hello"},
		{"both sides", "  hello  ", "hello"},
		{"tabs", "\thello\t", "hello"},
		{"newlines", "\nhello\n", "hello"},
		{"mixed whitespace", " \t\nhello \t\n", "hello"},
		{"empty string", "", ""},
		{"only whitespace", "   ", ""},
		{"internal spaces preserved", "  hello world  ", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TrimWhitespace(tt.input)
			if got != tt.want {
				t.Errorf("TrimWhitespace(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeSpaces(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"single spaces", "hello world", "hello world"},
		{"multiple spaces", "hello   world", "hello world"},
		{"tabs to space", "hello\tworld", "hello world"},
		{"newlines to space", "hello\nworld", "hello world"},
		{"mixed whitespace", "hello \t\n world", "hello world"},
		{"leading trailing", "  hello world  ", "hello world"},
		{"empty string", "", ""},
		{"only whitespace", "   \t\n   ", ""},
		{"many words", "one   two   three   four", "one two three four"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeSpaces(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeSpaces(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestStripHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no HTML", "hello world", "hello world"},
		{"simple tag", "<p>hello</p>", "hello"},
		{"nested tags", "<div><p>hello</p></div>", "hello"},
		{"tag with attributes", "<a href=\"test\">link</a>", "link"},
		{"self-closing tag", "hello<br/>world", "helloworld"},
		{"script tag", "<script>alert('xss')</script>", "alert('xss')"},
		{"multiple tags", "<b>bold</b> and <i>italic</i>", "bold and italic"},
		{"empty tags", "<div></div>", ""},
		{"malformed tag", "<div>unclosed", "unclosed"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StripHTML(tt.input)
			if got != tt.want {
				t.Errorf("StripHTML(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no special chars", "hello world", "hello world"},
		{"ampersand", "rock & roll", "rock &amp; roll"},
		{"less than", "a < b", "a &lt; b"},
		{"greater than", "a > b", "a &gt; b"},
		{"double quote", "say \"hello\"", "say &quot;hello&quot;"},
		{"single quote", "it's", "it&#39;s"},
		{"all special", "<script>alert('xss')</script>", "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"},
		{"empty string", "", ""},
		{"mixed content", "Price: $5 < $10 & discount > 0", "Price: $5 &lt; $10 &amp; discount &gt; 0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapeHTML(tt.input)
			if got != tt.want {
				t.Errorf("EscapeHTML(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase", "john doe", "John Doe"},
		{"uppercase", "JOHN DOE", "John Doe"},
		{"mixed case", "jOhN dOe", "John Doe"},
		{"single word", "john", "John"},
		{"with extra spaces", "  john   doe  ", "John Doe"},
		{"empty string", "", ""},
		{"single char", "j", "J"},
		{"portuguese name", "joão silva", "João Silva"},
		{"already correct", "John Doe", "John Doe"},
		{"with numbers", "john123", "John123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeName(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase", "test@example.com", "test@example.com"},
		{"uppercase", "TEST@EXAMPLE.COM", "test@example.com"},
		{"mixed case", "Test@Example.COM", "test@example.com"},
		{"with spaces", "  test@example.com  ", "test@example.com"},
		{"empty string", "", ""},
		{"complex email", "User.Name+Tag@Sub.Domain.COM", "user.name+tag@sub.domain.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeEmail(tt.input)
			if got != tt.want {
				t.Errorf("NormalizeEmail(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRemoveNonPrintable(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"printable only", "hello world", "hello world"},
		{"with null byte", "hello\x00world", "helloworld"},
		{"with control chars", "hello\x01\x02world", "helloworld"},
		{"keeps newline", "hello\nworld", "hello\nworld"},
		{"keeps tab", "hello\tworld", "hello\tworld"},
		{"keeps carriage return", "hello\rworld", "hello\rworld"},
		{"empty string", "", ""},
		{"unicode printable", "olá mundo", "olá mundo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveNonPrintable(tt.input)
			if got != tt.want {
				t.Errorf("RemoveNonPrintable(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRemoveControlChars(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no control chars", "hello world", "hello world"},
		{"with null byte", "hello\x00world", "helloworld"},
		{"with bell", "hello\x07world", "helloworld"},
		{"keeps newline", "hello\nworld", "hello\nworld"},
		{"keeps tab", "hello\tworld", "hello\tworld"},
		{"keeps carriage return", "hello\rworld", "hello\rworld"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveControlChars(tt.input)
			if got != tt.want {
				t.Errorf("RemoveControlChars(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToUppercase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"lowercase", "hello", "HELLO"},
		{"mixed", "Hello World", "HELLO WORLD"},
		{"already upper", "HELLO", "HELLO"},
		{"with numbers", "hello123", "HELLO123"},
		{"empty string", "", ""},
		{"unicode", "olá", "OLÁ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToUppercase(tt.input)
			if got != tt.want {
				t.Errorf("ToUppercase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestToLowercase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"uppercase", "HELLO", "hello"},
		{"mixed", "Hello World", "hello world"},
		{"already lower", "hello", "hello"},
		{"with numbers", "HELLO123", "hello123"},
		{"empty string", "", ""},
		{"unicode", "OLÁ", "olá"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToLowercase(tt.input)
			if got != tt.want {
				t.Errorf("ToLowercase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRemoveDigits(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no digits", "hello", "hello"},
		{"with digits", "hello123", "hello"},
		{"only digits", "123456", ""},
		{"mixed", "a1b2c3", "abc"},
		{"empty string", "", ""},
		{"spaces and digits", "hello 123 world", "hello  world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveDigits(tt.input)
			if got != tt.want {
				t.Errorf("RemoveDigits(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestKeepDigits(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"only digits", "123456", "123456"},
		{"with letters", "abc123def", "123"},
		{"no digits", "hello", ""},
		{"phone format", "+258 84 123 4567", "258841234567"},
		{"empty string", "", ""},
		{"spaces and special", "1-2-3 4/5/6", "123456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KeepDigits(tt.input)
			if got != tt.want {
				t.Errorf("KeepDigits(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestKeepAlphanumeric(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"alphanumeric", "abc123", "abc123"},
		{"with special chars", "hello@world.com", "helloworldcom"},
		{"with spaces", "hello world 123", "helloworld123"},
		{"only special", "!@#$%", ""},
		{"empty string", "", ""},
		{"unicode letters", "olá123", "olá123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := KeepAlphanumeric(tt.input)
			if got != tt.want {
				t.Errorf("KeepAlphanumeric(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestChain(t *testing.T) {
	tests := []struct {
		name  string
		input string
		fns   []Func
		want  string
	}{
		{
			name:  "no functions",
			input: "hello",
			fns:   []Func{},
			want:  "hello",
		},
		{
			name:  "single function",
			input: "  hello  ",
			fns:   []Func{TrimWhitespace},
			want:  "hello",
		},
		{
			name:  "multiple functions",
			input: "  <b>HELLO</b>  ",
			fns:   []Func{StripHTML, TrimWhitespace, ToLowercase},
			want:  "hello",
		},
		{
			name:  "order matters",
			input: "  hello  ",
			fns:   []Func{TrimWhitespace, ToUppercase},
			want:  "HELLO",
		},
		{
			name:  "complex chain",
			input: "  <p>João   Silva</p>  ",
			fns:   []Func{StripHTML, NormalizeSpaces, NormalizeName},
			want:  "João Silva",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Chain(tt.input, tt.fns...)
			if got != tt.want {
				t.Errorf("Chain(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestSanitizer(t *testing.T) {
	t.Run("empty sanitizer", func(t *testing.T) {
		s := NewSanitizer()
		result := s.Apply("hello")
		if result != "hello" {
			t.Errorf("expected 'hello', got %q", result)
		}
	})

	t.Run("chained operations", func(t *testing.T) {
		s := NewSanitizer().
			StripHTML().
			TrimWhitespace().
			NormalizeSpaces().
			ToLowercase()

		input := "  <b>HELLO   WORLD</b>  "
		want := "hello world"
		got := s.Apply(input)
		if got != want {
			t.Errorf("Apply(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("custom function", func(t *testing.T) {
		reverse := func(s string) string {
			runes := []rune(s)
			for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
				runes[i], runes[j] = runes[j], runes[i]
			}
			return string(runes)
		}

		s := NewSanitizer().
			TrimWhitespace().
			Custom(reverse)

		input := "  hello  "
		want := "olleh"
		got := s.Apply(input)
		if got != want {
			t.Errorf("Apply(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("all methods", func(t *testing.T) {
		// Test that all builder methods work
		s := NewSanitizer().
			TrimWhitespace().
			NormalizeSpaces().
			StripHTML().
			RemoveNonPrintable().
			RemoveControlChars().
			ToUppercase()

		result := s.Apply("  <b>hello</b>  ")
		if result != "HELLO" {
			t.Errorf("expected 'HELLO', got %q", result)
		}
	})

	t.Run("name normalization chain", func(t *testing.T) {
		s := NewSanitizer().
			StripHTML().
			NormalizeName()

		input := "<b>john</b> <i>doe</i>"
		want := "John Doe"
		got := s.Apply(input)
		if got != want {
			t.Errorf("Apply(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("email normalization chain", func(t *testing.T) {
		s := NewSanitizer().
			TrimWhitespace().
			NormalizeEmail()

		input := "  TEST@EXAMPLE.COM  "
		want := "test@example.com"
		got := s.Apply(input)
		if got != want {
			t.Errorf("Apply(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("keep digits chain", func(t *testing.T) {
		s := NewSanitizer().
			KeepDigits()

		input := "+258 84 123 4567"
		want := "258841234567"
		got := s.Apply(input)
		if got != want {
			t.Errorf("Apply(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("keep alphanumeric chain", func(t *testing.T) {
		s := NewSanitizer().
			KeepAlphanumeric()

		input := "hello@world.123"
		want := "helloworld123"
		got := s.Apply(input)
		if got != want {
			t.Errorf("Apply(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("escape HTML chain", func(t *testing.T) {
		s := NewSanitizer().
			TrimWhitespace().
			EscapeHTML()

		input := "  <script>alert('xss')</script>  "
		want := "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"
		got := s.Apply(input)
		if got != want {
			t.Errorf("Apply(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("to lowercase chain", func(t *testing.T) {
		s := NewSanitizer().
			ToLowercase()

		input := "HELLO WORLD"
		want := "hello world"
		got := s.Apply(input)
		if got != want {
			t.Errorf("Apply(%q) = %q, want %q", input, got, want)
		}
	})
}

func TestPrebuiltSanitizers(t *testing.T) {
	t.Run("TextSanitizer", func(t *testing.T) {
		s := TextSanitizer()
		input := "  <b>Hello\x00   World</b>  "
		want := "Hello World"
		got := s.Apply(input)
		if got != want {
			t.Errorf("TextSanitizer.Apply(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("NameSanitizer", func(t *testing.T) {
		s := NameSanitizer()
		input := "<b>JOHN\x00 doe</b>"
		want := "John Doe"
		got := s.Apply(input)
		if got != want {
			t.Errorf("NameSanitizer.Apply(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("EmailSanitizer", func(t *testing.T) {
		s := EmailSanitizer()
		input := "  TEST@EXAMPLE.COM  "
		want := "test@example.com"
		got := s.Apply(input)
		if got != want {
			t.Errorf("EmailSanitizer.Apply(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("PhoneSanitizer", func(t *testing.T) {
		s := PhoneSanitizer()
		input := "+258 84-123-4567"
		want := "258841234567"
		got := s.Apply(input)
		if got != want {
			t.Errorf("PhoneSanitizer.Apply(%q) = %q, want %q", input, got, want)
		}
	})
}

func TestInputNotModified(t *testing.T) {
	// Verify that all functions return new values and don't modify input.
	original := "  <b>HELLO</b>  "
	input := original

	_ = TrimWhitespace(input)
	if input != original {
		t.Error("TrimWhitespace modified input")
	}

	_ = NormalizeSpaces(input)
	if input != original {
		t.Error("NormalizeSpaces modified input")
	}

	_ = StripHTML(input)
	if input != original {
		t.Error("StripHTML modified input")
	}

	_ = EscapeHTML(input)
	if input != original {
		t.Error("EscapeHTML modified input")
	}

	_ = NormalizeName(input)
	if input != original {
		t.Error("NormalizeName modified input")
	}

	_ = NormalizeEmail(input)
	if input != original {
		t.Error("NormalizeEmail modified input")
	}

	_ = ToUppercase(input)
	if input != original {
		t.Error("ToUppercase modified input")
	}

	_ = ToLowercase(input)
	if input != original {
		t.Error("ToLowercase modified input")
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("very long string", func(t *testing.T) {
		input := strings.Repeat("a", 10000)
		result := TrimWhitespace(input)
		if len(result) != 10000 {
			t.Errorf("expected length 10000, got %d", len(result))
		}
	})

	t.Run("unicode handling", func(t *testing.T) {
		input := "日本語テスト"
		result := TrimWhitespace(input)
		if result != input {
			t.Errorf("expected %q, got %q", input, result)
		}
	})

	t.Run("emoji handling", func(t *testing.T) {
		input := "Hello 👋 World 🌍"
		result := NormalizeSpaces(input)
		if result != input {
			t.Errorf("expected %q, got %q", input, result)
		}
	})

	t.Run("mixed unicode and HTML", func(t *testing.T) {
		input := "<p>Olá Mundo 🌍</p>"
		want := "Olá Mundo 🌍"
		result := StripHTML(input)
		if result != want {
			t.Errorf("expected %q, got %q", want, result)
		}
	})
}
