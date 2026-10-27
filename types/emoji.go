package types

type EmojiOne struct {
	Name      string
	Shortname string
	Keywords  []string
	ASCII     []string
}

type EmojiEngine struct {
	Keys []rune
}
