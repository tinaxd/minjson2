package min

import (
	"io"
	"unicode"
	"unicode/utf8"
)

func runeToBytes(r rune, size int) []byte {
	bytes := make([]byte, size)
	utf8.EncodeRune(bytes, r)
	return bytes
}

func MinifyJSON(reader io.RuneReader, writer io.Writer) {
	forceChar := false
	inString := false

	for {
		ch, chSize, err := reader.ReadRune()
		if err != nil {
			break
		}

		if forceChar {
			writer.Write(runeToBytes(ch, chSize))
			forceChar = false
		} else if inString {
			if ch == '\\' {
				forceChar = true
				writer.Write(runeToBytes(ch, chSize))
			} else if ch == '"' {
				inString = false
				writer.Write(runeToBytes(ch, chSize))
			} else {
				writer.Write(runeToBytes(ch, chSize))
			}
		} else if unicode.IsSpace(ch) {
			// Skip
		} else {
			if ch == '"' {
				inString = true
				writer.Write(runeToBytes(ch, chSize))
			} else {
				writer.Write(runeToBytes(ch, chSize))
			}
		}
	}
}

type runeInfo struct {
	Rune rune
	Size int
}

func minifyJSON2(reader <-chan runeInfo, writer chan<- runeInfo) {
	forceChar := false
	inString := false

	for ch := range reader {
		if forceChar {
			writer <- ch
			forceChar = false
		} else if inString {
			if ch.Rune == '\\' {
				forceChar = true
				writer <- ch
			} else if ch.Rune == '"' {
				inString = false
				writer <- ch
			} else {
				writer <- ch
			}
		} else if unicode.IsSpace(ch.Rune) {
			// Skip
		} else {
			if ch.Rune == '"' {
				inString = true
				writer <- ch
			} else {
				writer <- ch
			}
		}
	}

	close(writer)
}

type PrettySetting struct {
	IndentWidth int
}

func PrettyJSON(reader io.RuneReader, writer io.Writer, setting PrettySetting) {
	const (
		objectCtx = iota
		arrayCtx
		stringCtx
		numberCtx
	)

	minified := make(chan runeInfo)
	go func() {
		sender := make(chan runeInfo)
		go minifyJSON2(sender, minified)
		for {
			ch, chSize, err := reader.ReadRune()
			if err != nil {
				break
			}
			sender <- runeInfo{ch, chSize}
		}
		close(sender)
	}()

	indentLevel := 0
	context := []int{objectCtx}
	forceChar := false
	lastComma := false

	for chInfo := range minified {
		ch := chInfo.Rune
		chSize := chInfo.Size

		writeChar := func() {
			//fmt.Printf("%v\n", runeToBytes(ch, chSize))
			writer.Write(runeToBytes(ch, chSize))
		}

		if forceChar {
			forceChar = false
			writeChar()
			lastComma = false
			continue
		}

		if context[len(context)-1] == stringCtx {
			if ch == '\\' {
				forceChar = true
			} else if ch == '"' {
				context = context[0 : len(context)-1]
			}
			writeChar()
			lastComma = false
			continue
		}

		if ch == '"' {
			writeChar()
			context = append(context, stringCtx)
			lastComma = false
			continue
		}

		if ch == ',' {
			writeChar()
			writer.Write([]byte{'\n'})
			for i := 0; i < indentLevel*setting.IndentWidth; i++ {
				writer.Write([]byte{' '})
			}
			lastComma = true
			continue
		}

		if ch == '[' || ch == '{' {
			writeChar()
			writer.Write([]byte{'\n'})
			indentLevel++
			for i := 0; i < indentLevel*setting.IndentWidth; i++ {
				writer.Write([]byte{' '})
			}
			switch ch {
			case '[':
				context = append(context, arrayCtx)
			case '{':
				context = append(context, objectCtx)
			}
			lastComma = false
			continue
		}

		if ch == ']' || ch == '}' {
			indentLevel--
			if !lastComma {
				writer.Write([]byte{'\n'})
				for i := 0; i < indentLevel*setting.IndentWidth; i++ {
					writer.Write([]byte{' '})
				}
			}
			writeChar()
			context = context[0 : len(context)-1]
			continue
		}

		if ch == ':' {
			writeChar()
			writer.Write([]byte{' '})
			continue
		}

		writeChar()
	}
	writer.Write([]byte{'\n'})
}
