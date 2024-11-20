package main

import (
	"math"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/maps"
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

var (
	punctuationRegex = regexp.MustCompile(`[^\w\s]`)
	whitespaceRegex  = regexp.MustCompile(`\s+`)
)

func similarityHandler(c *fiber.Ctx) error {
	var req SimilarityRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad Request",
		})
	}

	text1 := normalizeText(req.Text1)
	text2 := normalizeText(req.Text2)

	words1 := strings.Split(text1, " ")
	words2 := strings.Split(text2, " ")

	fm1 := generateFrequencyMap(words1)
	fm2 := generateFrequencyMap(words2)

	uw := make(map[string]any, 0)
	for word := range fm1 {
		uw[word] = struct{}{}
	}
	for word := range fm2 {
		uw[word] = struct{}{}
	}

	uniqueWords := maps.Keys(uw)

	total1 := len(words1)
	total2 := len(words2)

	tf1 := calculateTF(uniqueWords, fm1, total1)
	tf2 := calculateTF(uniqueWords, fm2, total2)

	idf := calculateIDF(uniqueWords, fm1, fm2)

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

func main() {
	app := fiber.New()
	app.Post("/similarity", similarityHandler)
	app.Listen(":8082")
}

func normalizeText(text string) string {
	lower := strings.ToLower(text)
	noPunctuation := punctuationRegex.ReplaceAllString(lower, "")
	cleanText := whitespaceRegex.ReplaceAllString(noPunctuation, " ")

	return strings.Trim(cleanText, " ")
}

func generateFrequencyMap(words []string) map[string]int {
	frequencyMap := make(map[string]int)
	for _, word := range words {
		frequencyMap[word]++
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
