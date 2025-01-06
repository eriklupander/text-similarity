package main

import (
	"bytes"
	gopcre "github.com/GRbit/go-pcre"
	"github.com/buger/jsonparser"
	"github.com/dolthub/swiss"
	"github.com/scylladb/go-set/strset"
	"github.com/xaionaro-go/atomicmap"
	pcre2 "go.arsenm.dev/pcre"
	"log"
	"math"
	"net/http"
	"regexp"
	"strings"

	_ "net/http/pprof"

	"github.com/eriklupander/replacer"
	"github.com/gofiber/fiber/v2"
	"github.com/wasilibs/go-re2"
)

type SimilarityRequest struct {
	Text1 string `json:"text1"`
	Text2 string `json:"text2"`
}

type InterpretationResult string

const (
	InterpretationResultDissimilar InterpretationResult = "Dissimilar"
	InterpretationResultSlightly   InterpretationResult = "Slightly Similar"
	InterpretationResultModerately InterpretationResult = "Moderately Similar"
	InterpretationResultQuite      InterpretationResult = "Quite Similar"
	InterpretationResultHighly     InterpretationResult = "Highly Similar"
	InterpretationResultUnknown    InterpretationResult = "Unknown"
)

type SimilarityResponse struct {
	Similarity     float64 `json:"similarity"`
	Interpretation string  `json:"interpretation"`
}

// Pre-compiled regexps for the various regexp engines tested.
var (
	punctuationRegex = regexp.MustCompile(`[^\w\s]`)
	whitespaceRegex  = regexp.MustCompile(`\s+`)

	punctuationRegexRE2 = re2.MustCompile(`[^\w\s]`)
	whitespaceRegexRE2  = re2.MustCompile(`\s+`)

	punctuationPcre2 = pcre2.MustCompile("[^\\w\\s]")
	whitespacePcre2  = pcre2.MustCompile("\\s+")

	punctuationRegexPcre = gopcre.MustCompileJIT("[^\\w\\s]", 0, gopcre.STUDY_JIT_COMPILE)
	whitespaceRegexPcre  = gopcre.MustCompileJIT("\\s+", 0, gopcre.STUDY_JIT_COMPILE)
)

// Keys used with jsonparser
var (
	key1 = []byte(`text1`)
	key2 = []byte(`text2`)
)

var punctuationsAscii = []byte{'[', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '\\', '-', '.', '/', ':', ';', '<', '=', '>', '?', '@', '[', ']', '^', '_', '`', '{', '|', '}', '~', ']'}

var byteReplacer *replacer.ByteReplacer
var stringsReplacer = strings.NewReplacer(append(replacer.RemovePunctuationPairs, append(replacer.WhitespacesAsSpacesPairs, replacer.ToLowerReplacements...)...)...)

func init() {
	var err error
	byteReplacer, err = replacer.NewByteReplacerFromStringPairs(append(replacer.RemovePunctuationPairs, append(replacer.WhitespacesAsSpacesPairs, replacer.ToLowerReplacements...)...)...)
	if err != nil {
		panic(err.Error())
	}
}

func main() {

	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	app := fiber.New()
	app.Post("/similarity", similarityHandler)
	app.Listen(":8081")
}

func similarityHandler(c *fiber.Ctx) error {

	// Use jsonparser to efficiently extract the two texts from the input JSON.
	var text1, text2 []byte
	_ = jsonparser.ObjectEach(c.BodyRaw(), func(key []byte, value []byte, _ jsonparser.ValueType, _ int) error {
		if bytes.Equal(key, key1) {
			text1 = bytes.TrimSpace(byteReplacer.Replace(value))
		}
		if bytes.Equal(key, key2) {
			text2 = bytes.TrimSpace(byteReplacer.Replace(value))
		}
		return nil
	})

	words1 := strings.Split(string(text1), " ")
	words2 := strings.Split(string(text2), " ")

	uw := strset.NewWithSize(len(words1) + len(words2))
	uw.Add(append(words1, words2...)...)
	uniqueWords := uw.List()

	fm1 := generateFrequencyMapWithCapacity(words1)
	fm2 := generateFrequencyMapWithCapacity(words2)

	total1 := len(words1)
	total2 := len(words2)

	tf1 := calculateTFMultiply(uniqueWords, fm1, total1)
	tf2 := calculateTFMultiply(uniqueWords, fm2, total2)

	idf := calculateIDFOptimized(uniqueWords, fm1, fm2)

	tfidf1 := calculateTFIDF(tf1, idf)
	tfidf2 := calculateTFIDF(tf2, idf)

	similarity := calculateSimilarity(tfidf1, tfidf2)

	similarity = math.Round(similarity*1000) / 1000
	interpretation := interpretSimilarity(similarity)

	return c.JSON(SimilarityResponse{
		Similarity:     similarity,
		Interpretation: string(interpretation),
	})
}

func normalizeText(text string) string {
	lower := strings.ToLower(text)
	noPunctuation := punctuationRegex.ReplaceAllString(lower, "")
	cleanText := whitespaceRegex.ReplaceAllString(noPunctuation, " ")

	return strings.Trim(cleanText, " ")
}

func normalizeTextCheat(text string) string {
	noPunctuation := strings.ReplaceAll(strings.ToLower(text), ".", "")
	cleanText := strings.ReplaceAll(noPunctuation, "\t", " ")
	return strings.Trim(cleanText, " ")
}

func normalizeTextReplacer(text string) string {
	return strings.TrimSpace(stringsReplacer.Replace(text))
}

func normalizeTextBytesReplacer(text []byte) []byte {
	return bytes.TrimSpace(byteReplacer.Replace(text))
}

func normalizeTextRE2(text string) string {
	lower := strings.ToLower(text)
	noPunctuation := punctuationRegexRE2.ReplaceAllString(lower, "")
	cleanText := whitespaceRegexRE2.ReplaceAllString(noPunctuation, " ")
	return strings.Trim(cleanText, " ")
}

func normalizeTextPCRE2(text string) string {
	lower := strings.ToLower(text)
	noPunctuation := punctuationPcre2.ReplaceAllString(lower, "")
	cleanText := whitespacePcre2.ReplaceAllString(noPunctuation, " ")
	return strings.Trim(cleanText, " ")
}

// normalizeTextGoPCRE does not work for me, it just runs forever - probably due to a badly constructed regexp on my part.
func normalizeTextGoPCRE(text string) string {
	text = strings.ToLower(text)
	text = whitespaceRegexPcre.ReplaceAllString(text, " ", 0)

	textB := punctuationRegexPcre.
		ReplaceAllString(text, "", 0)
	return strings.Trim(textB, " ")
}

func normalizeTextNestedForLoops(text string) string {
	out := make([]byte, len(text))
	j := 0
OUTER:
	for i := 0; i < len(text); i++ {
		for _, c := range punctuationsAscii {
			if c == text[i] {
				continue OUTER // Skip
			}
		}

		out[j] = text[i]
		// Remove tabs etc.
		if text[i] > 7 && text[i] < 14 {
			out[j] = ' '
		}
		j++
	}

	return strings.TrimSpace(strings.ToLower(string(out[0:j])))
}

func normalizeTextUsingAsciiRanges(txt string) string {
	text := []byte(txt)
	out := make([]byte, len(text))
	j := 0

	for i := 0; i < len(text); i++ {
		if text[i] > 32 && text[i] < 48 || text[i] > 57 && text[i] < 65 || text[i] > 90 && text[i] < 97 || text[i] > 122 && text[i] < 128 {
			continue // Skip.
		}
		// Remove tabs etc.
		if text[i] < 14 && text[i] > 7 {
			out[j] = ' '
			j++
			break
		}

		// Check for upper-case
		if text[i] > 64 && text[i] < 91 {
			out[j] = text[i] + 32
		} else {
			out[j] = text[i]
		}
		j++
	}

	return strings.TrimSpace(string(out[0:j]))
}

func normalizeTextUsingInvertedAsciiRanges(txt string) string {
	text := []byte(txt)
	out := make([]byte, len(text))
	j := 0

	for i := 0; i < len(text); i++ {
		// If a-zA-Z or space
		if text[i] == 32 || (text[i] > 64 && text[i] < 91) || (text[i] > 96 && text[i] < 123) {
			// Check for upper-case
			if text[i] > 64 && text[i] < 91 {
				out[j] = text[i] + 32
			} else {
				out[j] = text[i]
			}
			j++
			continue
		} else if text[i] < 14 && text[i] > 7 {
			// Remove tabs etc.
			out[j] = ' '
			j++
			continue
		}
	}

	return strings.TrimSpace(string(out[0:j]))
}

func normalizeTextUsingInvertedAsciiRangesBytes(text []byte) []byte {

	out := make([]byte, len(text))
	j := 0

	for i := 0; i < len(text); i++ {
		// If a-zA-Z or space
		if text[i] == 32 || (text[i] > 64 && text[i] < 91) || (text[i] > 96 && text[i] < 123) {
			// Check for upper-case
			if text[i] > 64 && text[i] < 91 {
				out[j] = text[i] + 32
			} else {
				out[j] = text[i]
			}
			j++
			continue
		} else if text[i] < 14 && text[i] > 7 {
			// Remove tabs etc.
			out[j] = ' '
			j++
			continue
		}
	}
	return bytes.TrimSpace(out[0:j])
}

func generateFrequencyMap(words []string) map[string]int {
	frequencyMap := make(map[string]int)
	for _, word := range words {
		frequencyMap[word]++
	}

	return frequencyMap
}
func generateFrequencyMapWithCapacity(words []string) map[string]int {
	frequencyMap := make(map[string]int, len(words))
	for _, word := range words {
		frequencyMap[word]++
	}

	return frequencyMap
}

func generateFrequencyMapSwiss(words []string) *swiss.Map[string, int] {
	frequencyMap := swiss.NewMap[string, int](uint32(len(words)))
	for _, word := range words {
		num, ok := frequencyMap.Get(word)
		if !ok {
			frequencyMap.Put(word, 1)
		} else {
			frequencyMap.Put(word, num+1)
		}
	}
	return frequencyMap
}

func generateFrequencyMapAtomic(words []string) atomicmap.Map {
	frequencyMap := atomicmap.New()

	for _, word := range words {
		num, err := frequencyMap.Get(word)
		if err != nil {
			_ = frequencyMap.Set(word, 1)
		} else {
			_ = frequencyMap.Set(word, num.(int)+1)
		}

	}
	return frequencyMap
}

func calculateTF(uniqueWords []string, frequencyMap map[string]int, total int) []float64 {
	tf := make([]float64, len(uniqueWords))
	for i, word := range uniqueWords {
		tf[i] = float64(frequencyMap[word]) / float64(total)
	}

	return tf
}

func calculateTFMultiply(uniqueWords []string, frequencyMap map[string]int, total int) []float64 {
	tf := make([]float64, len(uniqueWords))
	fraction := 1 / float64(total)
	for i, word := range uniqueWords {
		tf[i] = float64(frequencyMap[word]) * fraction
	}
	return tf
}

var oneOccurrenceIDF = math.Log(1.0 + 2.0/(float64(1)+1.0))
var twoOccurrencesIDF = math.Log(1.0 + 2.0/(float64(2)+1.0))

func calculateIDFOptimized(uniqueWords []string, fm1, fm2 map[string]int) []float64 {
	idf := make([]float64, len(uniqueWords))

	var occurrences = 0
	var ok = false
	for i, word := range uniqueWords {
		occurrences = 0
		if _, ok = fm1[word]; ok {
			occurrences++
		}
		if _, ok = fm2[word]; ok {
			occurrences++
		}
		if occurrences == 1 {
			idf[i] = oneOccurrenceIDF
		} else if occurrences == 2 {
			idf[i] = twoOccurrencesIDF
		}
	}
	return idf
}

func calculateIDF(uniqueWords []string, fm1, fm2 map[string]int) []float64 {
	docFreq := make(map[string]int)
	for _, word := range uniqueWords {
		count1, count2 := 0, 0
		if _, ok := fm1[word]; ok {
			count1 = 1
		}

		if _, ok := fm2[word]; ok {
			count2 = 1
		}

		docFreq[word] = count1 + count2
	}

	idf := make([]float64, len(uniqueWords))
	for i, word := range uniqueWords {
		idf[i] = math.Log(1.0 + 2.0/(float64(docFreq[word])+1.0))
	}

	return idf
}

func calculateTFIDF(tf, idf []float64) []float64 {
	tfidf := make([]float64, len(tf))
	for i := range len(tf) {
		tfidf[i] = tf[i] * idf[i]
	}

	return tfidf
}

func calculateSimilarity(tfidf1, tfidf2 []float64) float64 {
	dotProduct := 0.0
	for i := range len(tfidf1) {
		dotProduct += tfidf1[i] * tfidf2[i]
	}

	magnitude1 := 0.0
	for _, val := range tfidf1 {
		magnitude1 += val * val
	}
	magnitude1 = math.Sqrt(magnitude1)

	magnitude2 := 0.0
	for _, val := range tfidf2 {
		magnitude2 += val * val
	}
	magnitude2 = math.Sqrt(magnitude2)

	if magnitude1 <= 1e-9 || magnitude2 <= 1e-9 {
		return 0.0
	}

	return dotProduct / (magnitude1 * magnitude2)
}

func interpretSimilarity(similarity float64) InterpretationResult {
	if similarity <= 0.2 {
		return InterpretationResultDissimilar
	} else if similarity <= 0.4 {
		return InterpretationResultSlightly
	} else if similarity <= 0.6 {
		return InterpretationResultModerately
	} else if similarity <= 0.8 {
		return InterpretationResultQuite
	} else if similarity <= 1 {
		return InterpretationResultHighly
	}

	return InterpretationResultUnknown
}
