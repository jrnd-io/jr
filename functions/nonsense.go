// modified from https://golang.org/doc%2Fcodewalk%2Fmarkov.go

package functions

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

// Lorem generates a 'lorem ipsum' text of size words
func Lorem(size int) string {
	lorem := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. In ullamcorper non eros eget porta. Aliquam erat " +
		"volutpat. Mauris molestie lobortis dolor et cursus. Cras vulputate vitae urna et tristique. Nullam iaculis fringilla est, " +
		"vitae vulputate felis viverra suscipit. Nullam laoreet ornare tristique. Mauris porta, nisi sed laoreet scelerisque, nisi " +
		"velit eleifend nulla, sit amet pharetra mauris tortor cursus ipsum. Nulla dignissim nunc vel felis convallis ullamcorper. " +
		"Aliquam convallis nunc mi, ut tempor nibh efficitur eu.Nam aliquet elit ac eros facilisis, ac commodo magna porttitor. " +
		"Nunc ut Lorem sit amet justo pulvinar aliquet accumsan sit amet turpis. Pellentesque habitant morbi tristique senectus et netus " +
		"et malesuada fames ac turpis egestas. In scelerisque maximus neque. Nunc in sapien vitae nunc eleifend pulvinar. Pellentesque " +
		"faucibus massa vel mauris molestie, non aliquam est maximus. Nam hendrerit purus a justo iaculis elementum. Aliquam interdum " +
		"scelerisque convallis. Vestibulum fringilla nunc ac sem ullamcorper, ac iaculis nulla posuere. Vivamus quis consequat ipsum. " +
		"In in neque dui. Mauris rutrum dapibus orci, id congue lectus pulvinar a. Proin in maximus dolor, id rutrum nulla. Etiam et " +
		"quam lacinia, porttitor urna ac, efficitur arcu. Suspendisse potenti. Lorem ipsum dolor sit amet, consectetur adipiscing elit." +
		"Fusce elit magna, lobortis nec semper non, aliquam at nisl. Vestibulum elementum suscipit justo et commodo. Etiam ultrices sem " +
		"non tellus molestie, ac elementum felis sollicitudin. Nam aliquet magna non nisi malesuada, ac varius mi mattis. Fusce ultricies " +
		"Lorem id dolor malesuada, vitae tincidunt turpis finibus. Curabitur efficitur varius aliquam. Cras facilisis ultrices pellentesque. " +
		"Vivamus molestie nibh tincidunt, aliquet velit sed, ultrices magna. Suspendisse potenti.In dapibus, mauris ac lacinia ultricies, " +
		"nibh orci rhoncus felis, molestie pharetra urna justo non mauris. Praesent dignissim ex id lacinia ullamcorper. Donec varius eros " +
		"ex, at mollis magna imperdiet id. Cras non tristique tortor, eget placerat purus. Pellentesque quis enim interdum, blandit tortor " +
		"sit amet, laoreet ligula. Curabitur id nisl ut Lorem commodo fringilla. Sed dapibus a libero a viverra. In odio ligula, tristique " +
		"a leo quis, imperdiet finibus mauris. Integer imperdiet justo vel mollis efficitur. Curabitur non felis accumsan ipsum vulputate " +
		"pharetra eu nec urna. Donec in mi sed libero dapibus pulvinar et sed tortor. Cras et accumsan risus, vitae porttitor felis." +
		"Aliquam bibendum, mi a tristique commodo, dolor orci blandit orci, id faucibus nisi neque eu erat. Nulla facilisi. Phasellus " +
		"maximus, augue at euismod euismod, nunc ante porta nibh, ac mattis nisl purus nec ex. In ut quam nisi. Nulla non placerat orci, " +
		"in suscipit lacus. Nulla porttitor sollicitudin bibendum. Vestibulum ipsum metus, sagittis ut sodales in, volutpat eu neque. " +
		"Nunc et felis diam. Integer hendrerit nisl metus, eu posuere tortor consequat a. Nunc eu sapien eu eros mollis suscipit. " +
		"Nam vitae rhoncus odio, vitae scelerisque augue. Maecenas elementum lacus vel sem pharetra, sed consectetur ipsum congue. " +
		"Proin nec diam purus. In sollicitudin feugiat sodales. Donec elementum volutpat nunc, sed ultricies diam mattis et. " +
		"Vivamus accumsan neque neque, et porta turpis finibus id."
	return Nonsense(2, size, lorem)
}

// SentencePrefix generates an 'alice in wonderland' text of size words with given prefixLen
func SentencePrefix(prefixLen, numWords int) string {
	alice := "Alice was beginning to get very tired of sitting by her sister on the bank, and of having nothing to do: " +
		"once or twice she had peeped into the book her sister was reading, but it had no pictures or conversations in it, " +
		"“and what is the use of a book,” thought Alice “without pictures or conversations?”" +
		"So she was considering in her own mind (as well as she could, for the hot day made her feel very sleepy and stupid), " +
		"whether the pleasure of making a daisy-chain would be worth the trouble of getting up and picking the daisies, when " +
		"suddenly a White Rabbit with pink eyes ran close by her." +
		"There was nothing so very remarkable in that; nor did Alice think it so very much out of the way to hear the Rabbit " +
		"say to itself, “Oh dear! Oh dear! I shall be late!” (when she thought it over afterwards, it occurred to her that she ought to have wondered at this, but at the time it all seemed quite natural); but when the Rabbit actually took a watch out of its waistcoat-pocket, and looked at it, and then hurried on, Alice started to her feet, for it flashed across her mind that she had never before seen a rabbit with either a waistcoat-pocket, or a watch to take out of it, and burning with curiosity, she ran across the field after it, and fortunately was just in time to see it pop down a large rabbit-hole under the hedge." +
		"In another moment down went Alice after it, never once considering how in the world she was to get out again." +
		"The rabbit-hole went straight on like a tunnel for some way, and then dipped suddenly down, so suddenly that " +
		"Alice had not a moment to think about stopping herself before she found herself falling down a very deep well." +
		"Either the well was very deep, or she fell very slowly, for she had plenty of time as she went down to look about " +
		"her and to wonder what was going to happen next. First, she tried to look down and make out what she was coming to, " +
		"but it was too dark to see anything; then she looked at the sides of the well, and noticed that they were filled " +
		"with cupboards and book-shelves; here and there she saw maps and pictures hung upon pegs. She took down a jar from " +
		"one of the shelves as she passed; it was labelled “ORANGE MARMALADE”, but to her great disappointment it was empty: " +
		"she did not like to drop the jar for fear of killing somebody underneath, so managed to put it into one of the cupboards " +
		"as she fell past it." +
		"“Well!” thought Alice to herself, “after such a fall as this, I shall think nothing of tumbling down stairs! How brave they’ll " +
		"all think me at home! Why, I wouldn’t say anything about it, even if I fell off the top of the house!” (Which was very likely " +
		"true.)" +
		"Down, down, down. Would the fall never come to an end? “I wonder how many miles I’ve fallen by this time?” she said aloud. " +
		"“I must be getting somewhere near the centre of the earth. Let me see: that would be four thousand miles down, I think—” (" +
		"for, you see, Alice had learnt several things of this sort in her lessons in the schoolroom, and though this was not a very " +
		"good opportunity for showing off her knowledge, as there was no one to listen to her, still it was good practice to say it over) " +
		"“—yes, that’s about the right distance—but then I wonder what Latitude or Longitude I’ve got to?” (Alice had no idea what Latitude" +
		"was, or Longitude either, but thought they were nice grand words to say.)" +
		"Presently she began again. “I wonder if I shall fall right through the earth! How funny it’ll seem to come out among the people " +
		"that walk with their heads downward! The Antipathies, I think—” (she was rather glad there was no one listening, this time, as it " +
		"didn’t sound at all the right word) “—but I shall have to ask them what the Name of the country is, you know. Please, Ma’am, is " +
		"this New Zealand or Australia?” (and she tried to curtsey as she spoke—fancy curtseying as you’re falling through the air! " +
		"Do you think you could manage it?) “And what an ignorant little girl she’ll think me for asking! No, it’ll never do to ask: " +
		"perhaps I shall see it written up somewhere.”" +
		"Down, down, down. There was nothing else to do, so Alice soon began talking again. “Dinah’ll miss me very much to-night, " +
		"I should think!” (Dinah was the cat.) “I hope they’ll remember her saucer of milk at tea-time. Dinah my dear! I wish you were " +
		"down here with me! There are no mice in the air, I’m afraid, but you might catch a bat, and that’s very like a mouse, you know." +
		"But do cats eat bats, I wonder?” And here Alice began to get rather sleepy, and went on saying to herself, in a dreamy sort of " +
		"way, “Do cats eat bats? Do cats eat bats?” and sometimes, “Do bats eat cats?” for, you see, as she couldn’t answer either question," +
		"it didn’t much matter which way she put it. She felt that she was dozing off, and had just begun to dream that she was walking " +
		"hand in hand with Dinah, and saying to her very earnestly, “Now, Dinah, tell me the truth: did you ever eat a bat?” when suddenly" +
		"thump! thump! down she came upon a heap of sticks and dry leaves, and the fall was over." +
		"Alice was not a bit hurt, and she jumped up on to her feet in a moment: she looked up, but it was all dark overhead; before her " +
		"was another long passage, and the White Rabbit was still in sight, hurrying down it. There was not a moment to be lost: away went " +
		"Alice like the wind, and was just in time to hear it say, as it turned a corner, “Oh my ears and whiskers, how late it’s getting!” " +
		"She was close behind it when she turned the corner, but the Rabbit was no longer to be seen: she found herself in a long, low hall" +
		" which was lit up by a row of lamps hanging from the roof." +
		"There were doors all round the hall, but they were all locked; and when Alice had been all the way down one side and up the other" +
		"trying every door, she walked sadly down the middle, wondering how she was ever to get out again." +
		"Suddenly she came upon a little three-legged table, all made of solid glass; there was nothing on it except a tiny golden key" +
		"and Alice’s first thought was that it might belong to one of the doors of the hall; but, alas! either the locks were too large " +
		"or the key was too small, but at any rate it would not open any of them. However, on the second time round, she came upon a low " +
		"curtain she had not noticed before, and behind it was a little door about fifteen inches high: she tried the little golden key " +
		"in the lock, and to her great delight it fitted!" +
		"Alice opened the door and found that it led into a small passage, not much larger than a rat-hole: she knelt down and looked along" +
		"the passage into the loveliest garden you ever saw. How she longed to get out of that dark hall, and wander about among those " +
		"beds of bright flowers and those cool fountains, but she could not even get her head through the doorway; “and even if my " +
		"head would go through,” thought poor Alice, “it would be of very little use without my shoulders. Oh, how I wish I could " +
		"shut up like a telescope! I think I could, if I only knew how to begin.” For, you see, so many out-of-the-way things had" +
		" happened lately, that Alice had begun to think that very few things indeed were really impossible." +
		"There seemed to be no use in waiting by the little door, so she went back to the table, half hoping she might find another " +
		"key on it, or at any rate a book of rules for shutting people up like telescopes: this time she found a little bottle on it" +
		"(“which certainly was not here before,” said Alice,) and round the neck of the bottle was a paper label, with the words " +
		"“DRINK ME,” beautifully printed on it in large letters." +
		"It was all very well to say “Drink me,” but the wise little Alice was not going to do that in a hurry. “No, I’ll look " +
		"first,” she said, “and see whether it’s marked ‘poison’ or not”; for she had read several nice little histories about " +
		"children who had got burnt, and eaten up by wild beasts and other unpleasant things, all because they would not remember" +
		"the simple rules their friends had taught them: such as, that a red-hot poker will burn you if you hold it too long; and" +
		"that if you cut your finger very deeply with a knife, it usually bleeds; and she had never forgotten that, if you drink" +
		"much from a bottle marked “poison,” it is almost certain to disagree with you, sooner or later." +
		"However, this bottle was not marked “poison,” so Alice ventured to taste it, and finding it very nice, (it had, in fact" +
		"a sort of mixed flavour of cherry-tart, custard, pine-apple, roast turkey, toffee, and hot buttered toast,) she very soon" +
		"finished it off."
	return Nonsense(prefixLen, numWords, alice)
}

// Sentence generates an 'alice in wonderland' text of size words
func Sentence(numWords int) string {
	return SentencePrefix(2, numWords)
}

// Nonsense generates a random Sentence of numWords wordsm using a prefixLen and a baseText to start from
func Nonsense(prefixLen, numWords int, baseText string) string {
	c := NewChain(prefixLen)
	c.Build(strings.NewReader(baseText))
	return c.Generate(numWords)
}

// RandomString returns a random string long between min and max characters
func RandomString(min, max int) string {
	return RandomStringVocabulary(min, max, alphabet)
}

// RandomStringVocabulary returns a random string long between min and max characters using a vocabulary
func RandomStringVocabulary(min, max int, source string) string {
	textb := make([]byte, min+Random.Intn(max-min))
	for i := range textb {
		textb[i] = source[Random.Intn(len(source))]
	}
	return string(textb)
}
