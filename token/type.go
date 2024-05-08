package token

type Type uint

//go:generate stringer -type=Type -linecomment
const (
	unknown Type = iota

	// punctuation
	LPar      // (
	RPar      // )
	LBrace    // {
	RBrace    // }
	LBrack    // [
	RBrack    // ]
	Comma     // ,
	Dot       // .
	Pipe      // |
	SemiColon // ;

	// aritmetic
	Plus  // +
	Minus // -
	Mult  // *
	Div   // /
	Mod   // %

	// logical
	Bang   // !
	NEQ    // !=
	Assign // =
	EQ     // ==
	GT     // >
	GTE    // >=
	LT     // <
	LTEQ   // <=

	// literals
	Ident  // ident
	String // string
	Int    // int
	Hex    // hex
	Oct    // oct
	Bin    // bin
	Float  // float

	Comment // comment

	// keywords
	keywordStart
	Class  // class
	Super  // super
	Self   // self
	Func   // fn
	Return // return
	Let    // let
	If     // if
	Else   // else
	And    // and
	Or     // or
	For    // for
	While  // while
	True   // true
	False  // false
	Nil    // nil
	keywordEnd

	EOF

	end
)

var keywords map[string]Type

func init() {
	keywords = make(map[string]Type, keywordEnd-keywordStart)
	for i := keywordStart + 1; i < keywordEnd; i++ {
		keywords[i.String()] = i
	}
}

func LookupIdent(ident string) Type {
	keyword, ok := keywords[ident]
	if ok {
		return keyword
	}
	return Ident
}
