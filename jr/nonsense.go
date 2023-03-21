// modified from https://golang.org/doc%2Fcodewalk%2Fmarkov.go

package jr

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[Random.Intn(len(choices))]

		if i == n-1 {
			if strings.HasSuffix(next, ",") {
				next = strings.ReplaceAll(next, ",", ".")
			}
		}
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

func Lorem(size int) string {
	lorem := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. In ullamcorper non eros eget porta. Aliquam erat " +
		"volutpat. Mauris molestie lobortis dolor et cursus. Cras vulputate vitae urna et tristique. Nullam iaculis fringilla est, " +
		"vitae vulputate felis viverra suscipit. Nullam laoreet ornare tristique. Mauris porta, nisi sed laoreet scelerisque, nisi " +
		"velit eleifend nulla, sit amet pharetra mauris tortor cursus ipsum. Nulla dignissim nunc vel felis convallis ullamcorper. " +
		"Aliquam convallis nunc mi, ut tempor nibh efficitur eu.Nam aliquet elit ac eros facilisis, ac commodo magna porttitor. " +
		"Nunc ut lorem sit amet justo pulvinar aliquet accumsan sit amet turpis. Pellentesque habitant morbi tristique senectus et netus " +
		"et malesuada fames ac turpis egestas. In scelerisque maximus neque. Nunc in sapien vitae nunc eleifend pulvinar. Pellentesque " +
		"faucibus massa vel mauris molestie, non aliquam est maximus. Nam hendrerit purus a justo iaculis elementum. Aliquam interdum " +
		"scelerisque convallis. Vestibulum fringilla nunc ac sem ullamcorper, ac iaculis nulla posuere. Vivamus quis consequat ipsum. " +
		"In in neque dui. Mauris rutrum dapibus orci, id congue lectus pulvinar a. Proin in maximus dolor, id rutrum nulla. Etiam et " +
		"quam lacinia, porttitor urna ac, efficitur arcu. Suspendisse potenti. Lorem ipsum dolor sit amet, consectetur adipiscing elit." +
		"Fusce elit magna, lobortis nec semper non, aliquam at nisl. Vestibulum elementum suscipit justo et commodo. Etiam ultrices sem " +
		"non tellus molestie, ac elementum felis sollicitudin. Nam aliquet magna non nisi malesuada, ac varius mi mattis. Fusce ultricies " +
		"lorem id dolor malesuada, vitae tincidunt turpis finibus. Curabitur efficitur varius aliquam. Cras facilisis ultrices pellentesque. " +
		"Vivamus molestie nibh tincidunt, aliquet velit sed, ultrices magna. Suspendisse potenti.In dapibus, mauris ac lacinia ultricies, " +
		"nibh orci rhoncus felis, molestie pharetra urna justo non mauris. Praesent dignissim ex id lacinia ullamcorper. Donec varius eros " +
		"ex, at mollis magna imperdiet id. Cras non tristique tortor, eget placerat purus. Pellentesque quis enim interdum, blandit tortor " +
		"sit amet, laoreet ligula. Curabitur id nisl ut lorem commodo fringilla. Sed dapibus a libero a viverra. In odio ligula, tristique " +
		"a leo quis, imperdiet finibus mauris. Integer imperdiet justo vel mollis efficitur. Curabitur non felis accumsan ipsum vulputate " +
		"pharetra eu nec urna. Donec in mi sed libero dapibus pulvinar et sed tortor. Cras et accumsan risus, vitae porttitor felis." +
		"Aliquam bibendum, mi a tristique commodo, dolor orci blandit orci, id faucibus nisi neque eu erat. Nulla facilisi. Phasellus " +
		"maximus, augue at euismod euismod, nunc ante porta nibh, ac mattis nisl purus nec ex. In ut quam nisi. Nulla non placerat orci, " +
		"in suscipit lacus. Nulla porttitor sollicitudin bibendum. Vestibulum ipsum metus, sagittis ut sodales in, volutpat eu neque. " +
		"Nunc et felis diam. Integer hendrerit nisl metus, eu posuere tortor consequat a. Nunc eu sapien eu eros mollis suscipit. " +
		"Nam vitae rhoncus odio, vitae scelerisque augue. Maecenas elementum lacus vel sem pharetra, sed consectetur ipsum congue. " +
		"Proin nec diam purus. In sollicitudin feugiat sodales. Donec elementum volutpat nunc, sed ultricies diam mattis et. " +
		"Vivamus accumsan neque neque, et porta turpis finibus id."
	return Nonsense(1, size, lorem)
}

func Nonsense(prefixLen, numWords int, baseText string) string {
	c := NewChain(prefixLen)
	c.Build(strings.NewReader(baseText))
	return c.Generate(numWords)
}

func RandomString(min, max int) string {
	return RandomStringFromSource(min, max, alphabet)
}

func RandomStringFromSource(min, max int, source string) string {
	textb := make([]byte, min+Random.Intn(max-min))
	for i := range textb {
		textb[i] = source[Random.Intn(len(source))]
	}
	return string(textb)
}
