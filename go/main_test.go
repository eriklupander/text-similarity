package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/simdjson-go"
	"github.com/scylladb/go-set/strset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

var line1 = `Direction receive leader memory others none her. Car or long forget try modern water culture. Question maintain who section tonight unit. Affect bring operation. Dinner ok inside perhaps practice against. Particularly eight on despite service. Certain week rule school town thus tax herself. Record young protect social. Available forward pick situation. Chance produce let anything for. It lot wall. Report cover green customer onto. Way wrong skin learn once minute. Central free east one claim. Food degree message. Direction box recently take trouble discuss later. Purpose several whatever market. Recent education boy until. Explain very science improve player right. True threat commercial sign positive tree customer. College address idea you establish wall become seem. Item build small type collection hit. Network political wonder foreign contain son. Admit foreign he green learn soon. Church teacher it. Oil cut plan. Lawyer reflect identify happy particularly sometimes year. Public or hotel behavior wear moment whatever with. Hand yes feeling policy goal although term month. Statement decision respond card together house draw food. Particularly section cold seven. Box vote computer outside. Soon ball claim politics. Pattern sport world center pay travel. Attorney simply purpose word school worker sign. Against fast too service husband hour. Newspaper interview speech miss story other according. People pattern interview group maybe start. Protect cultural enjoy tree character firm. Price mean if leave we laugh great. Road professor shoulder chance anything interview. Song staff begin result street law. End heart sell particularly. He couple be head wind. Me individual scientist available probably. Work item their style. Reach property program reveal. Allow customer close your adult force. Wonder apply line amount. Later lawyer general movement watch manager them. Difficult respond nothing matter. True still dream focus cup. Under yard clear language order. Prove own stuff threat. Bill stuff how scene again administration. Several tell system mission yourself. Experience century source parent year similar. Figure card teach edge up tax peace. Blood set early reality growth rich prevent piece. Form moment check senior look research traditional. North ground Republican Democrat wind father. Particular boy land. Happy serve you camera hot remember air. Gun great population describe term. Seem their institution later medical despite. Defense population individual safe while office. Example worry movie key. Main fund tell strategy. Off risk green pass option writer newspaper. Agreement everything a few condition leader head. Four federal concern sea short read sister yard. Point prove simple factor decision. Star forward line first worker value. Training prove fund conference statement chair degree. Create likely pressure here arrive guy require such. Under upon possible seem former. Tax measure customer direction church trip war. Between hour little focus attack wide. Bed name between realize cause indeed ok. Both keep power go go. Positive material responsibility clear rather effort close. Seem trade another low. Everything along win receive. Government machine his certainly long bill start senior. Family management region age machine fish recognize. Read hot attorney general. Provide station consumer church focus tell. Before likely actually budget return. Own off bit ahead at. Teacher food necessary west win look owner. Operation statement get. Site movement enjoy specific consumer ahead need whether. Three race seem decade. Including structure arm consumer important. Nearly some responsibility all. Station say six conference. Wall member morning third break continue sing. Future woman say science event. Maybe actually cover position specific year six draw. College military purpose no but. Others political in local adult show vote. Resource matter sign late explain area movement. Result deep if sit pull argue push. Him collection compare just. Wish red specific industry. Process man moment enough total law seat. Fine sport water produce them political itself. Process today significant choose be region decision. Get room town third now. Big adult can center series. Stand those college oil. Well news drive Democrat leave. Wait return enter word while put. Fast once approach even term one card. Alone example building blood nothing risk. Administration here book role decade project. Environmental daughter himself hold foot off stand. Shoulder charge store fast. Feeling teach case drop pretty sometimes. Great have see left without. Pay her difficult. Watch actually always life follow market. Through dog easy who use. Rich hand future person show show appear. Then knowledge security. Turn fast concern some deal. Reality indeed board do himself radio music. Little itself author door. Protect company so call. Mouth ever nation charge usually. Yard agree use feel issue especially. Main maintain important matter church. South present rather not air avoid. Respond economic task star else blue. Dog recently big right hope. Outside song reach down run build. Else range man right. Message put ability body each wide. Spring either work follow citizen accept guy. Clearly join keep. Sit also small. Into then also. Entire a large yard. Make member sister Mr reflect. Certain best first agreement spend officer health. Which possible stuff simple civil. Common everyone report individual true do live how. No area morning improve. However rise actually goal. Scientist drop attorney agree every rather sense. System prove benefit much ever. Pm true behavior detail because exist business. Question note and community difference design base. Determine fall sign will test however. Degree development dog alone president enter. Opportunity middle along community product red recent. Decide fire hair. Magazine teach claim able wonder candidate. Among job feel per federal president go. Well blue reach science meeting dream. Why work analysis include argue. Trip town line same nature try. During order question TV. Across ability ever front citizen kid brother. Analysis trip small spring room. His be understand move movie where. Whole society treat something movement. Student laugh might week manager want. Impact activity company free will important. Receive late condition popular. Author painting if Congress return floor away. Themselves listen book agree company. Value nearly station capital. Career others alone book imagine side term. A less write under make including field. Game item success and president. Last grow audience wife once dark child end. Eye support few decision capital. Common guy market police night. Clearly such top Congress family. Each hope center pull ever. Four allow commercial nor organization ahead four. Education hold early newspaper deal certainly lose. Face soldier measure high. Bed director serve trouble. Agreement want seat class wall in nor. Expert various minute success economic. Financial sit clearly. Attention close management level country financial. White last find bring summer record expect. Can sometimes study very. Race they authority scientist police. Senior decision learn price audience. Wind summer sea too grow record. Himself letter only drive plant although. Already pay require should magazine. Explain vote center green against year church. Food low understand everything raise. International individual everyone with your company move. Simply TV receive hand away happen. Degree either several design around model. Service mouth everybody in. Give recent why amount. Myself guess forward game start four tax. This realize not though head. Last use reduce. This within leader light commercial worry someone memory. Happen customer but whom. Recognize government for crime follow once cause fight. Opportunity read identify resource. Seat within attack success. Big tough remember. Send performance old none. Mother describe carry it senior forget. Race approach science city ball. Everyone yet through six. Second no yes computer year. Perhaps believe serious professional. Around necessary blue crime political or truth. Age daughter four pattern measure ago. Country enough yeah natural possible else hospital. Indicate view at call according issue. Sport these however way city. Put office down key pressure have. Live game ever black. Plan understand face heavy. Authority quickly few consumer from soon. Nothing allow dog face everything. Mention then really quickly add Democrat. Bar fly institution than like. Teach speech reality manager this rather town. Still phone question find. Daughter join job. See education service read mind us. Center choose big race. Trial sense foot various difficult senior product. Billion military science stuff yard. Hold central environmental edge control opportunity commercial. One crime reflect group wish do effect year. She small positive training task discussion. As exist executive person PM must want since. Produce economy card after. Practice entire husband or politics wish court personal. Bit find forget. Voice increase represent figure manage information. Born hard hard trade. Western natural reason economy forward describe. Discussion effort six budget bill. Hard amount meet traditional reality. Individual language also east us somebody author. My heart charge walk. Traditional participant employee nothing sense name benefit. Respond realize success. Reveal trip every able authority claim. American manage imagine turn care. Fine yes wear class majority computer dog. Alone chance dark receive choice recent prepare. Suggest successful trade large window financial. Tell resource value top low three Congress. Game project oil place simple move page. Seem page national treatment tell. Lay wide region industry rule or. Not really real international. Many perform answer mission color paper admit. Remain professional water mind. White bit everyone challenge those. Laugh military instituti`
var line2 = `Television culture this air machine. Follow entire test key represent new out. Sea whether from. Run political drop. Process build lose. Cause affect war current. Within none size few strategy. Middle same gas central. Prove choose decision black financial fill. Side doctor well music economic security enough. Question last see stock upon effort player. Follow deep affect. Dream measure responsibility design view. Wrong large car himself adult add capital. Option country nature like election. International administration husband unit tax behavior. Front authority represent training make. Nearly final red get necessary home. Our might clear financial idea meeting. Be range support court cut some. Clearly win place who cup ten authority. East send sea room tough. Me goal blue know keep team economic husband. Likely second skill near president. As research avoid note use military political within. Sound bar new some six reduce. Hour ability production fall. Area prepare might. No those experience experience strategy sing always. Fish answer pattern eye. Billion too speech cultural series action charge bar. Arrive account early culture spring truth. Break pretty crime from here become choice. Cost statement relationship report plan. Rise son program friend always present. Own six network green skill. Knowledge term care not. Music rest if person fund standard fall stay. Foreign enough simply amount money decide attorney. Nothing close big travel culture section run sister. Window quickly officer affect she look while. Majority idea attention organization onto. Process especially decision avoid summer commercial beautiful. Choose war off hotel field pay try back. Effort out small soon country work. Visit list character anything month what learn. Board mission then upon myself. Past officer through entire case rest. Military star realize care blue more. Sister bring to speech sort life real. Say matter strategy rule world fall current. According region fine high benefit issue perhaps. Mr about bag room around magazine. Number hand establish step add. Leader always challenge story still past perhaps. Girl design accept agreement require Mr. Affect by drug matter blue ten. Value civil thought national network statement. They into likely stage least ago. Add only pull physical will state marriage I. Occur own woman arrive within property style. Word charge PM thing join popular. Time everybody night how yard hear daughter. Agent bar lead light today boy situation. Media goal side activity case very. Citizen billion picture find pick sound. Large industry probably relate dog movement. Something table develop. Thought clear they less note. Eye to somebody serve maybe. Enjoy international child whether half case some. Somebody gun professor ok social. Memory close ground agency system star. Race table call risk say ten nothing. Performance carry least rule employee. Free sometimes cultural buy myself support. Participant college wonder live language he too drop. Kid four out it. Impact best unit rich several. Building over skin bed southern. Benefit race enter speech number hard participant. Actually adult put source arm about. Decision color long. Collection around administration seven degree. Executive develop really Congress Mrs. Again their physical only machine away. On lay finish time yeah. Seem coach head especially question chance PM. Low water PM world sea this modern. Show opportunity wonder idea suggest couple. Change war there learn pretty natural. Check woman country people send group have. Significant design sell detail single its glass. Finally reality institution ever. Whose news reveal soon quality south six. Across lawyer attention policy no adult stock. Itself still reduce. Forward man professor few this maybe reason prevent. Environmental south table election. Ask enter security however kind old newspaper. Under day watch respond trip. Skill his suggest blood within. Sing big note class. Worker whose near really official once. They note official type. Prove very send care. Growth against less response cultural. Figure best man outside. Fall cut particularly only boy these. Minute those wear seat executive. Shoulder measure partner traditional growth key expert wide. Soon western development party gas. Color next kitchen kitchen. Word somebody smile. Goal action PM second. Instead husband hard feel agent move his. Heavy perform land several federal maintain. Cultural price painting. Office anyone must their leg reduce. Hair a heart hot success season. Accept whole explain physical local fill. Assume pick bed down. Whether admit all stop short. Seat much well task room night meet. Continue beautiful health customer interview. Short so matter back work allow. Foreign almost drug daughter it pick. Goal common material. Red professor become there. Thank occur task. Method education develop discover arrive animal. Note so follow skill door and government according. Few direction account read gun. Indeed other still site western admit wear. We sense record price science easy coach. Understand at air relationship. Consumer probably light Democrat win day claim. Enough follow physical during book test. Discussion at know nearly visit detail white. Kid cold I yes. Order see surface behind two edge. Around along local too skill. Guess phone because. Game three tend TV rate citizen claim. Sell measure mission difficult. Thought reveal expect television because. After course whether would top network. Young mother chair industry important. Small home born eye standard early court. Job finish over Mrs easy. Development stuff very. Generation city modern wonder return. Debate camera too important approach loss produce. Image those sell process kitchen discover knowledge. We west actually. Surface understand rate right. Artist skin view. Medical gas might out get add. Government development level myself for candidate design address. Water kind mouth action. My production local health pull should. Claim meet management occur mouth over. Source world hour past expert last paper. Dog where whether wonder spring similar good strong. Evening inside arrive might step campaign have this. Hold great reason school provide. First total reduce voice. Very themselves the once special involve difficult. Themselves son magazine but learn red page number. Simply step system blue whatever management. Town song structure sing significant discussion huge. Central state even forget appear address live. Meeting we bring. Court history others take own technology daughter. Tough the personal investment. Business sister laugh design inside may cell. Keep music hard history standard. Camera production movement everyone else century section next. Put player most treat thank suffer region event. Young site along pick. Decide film peace beyond along development. Congress television car behind month. Just own clearly trial point. Out everyone decision idea rather. Close hit company everyone not. Across week scene sound. Hand money owner simple fish. Show avoid tonight own hundred time because. Beyond serious fear fund either quite. By trouble friend mean system fine condition. Strong young beyond arm. Film structure eye as research brother. Mr commercial radio sit natural education. Toward material audience forget return make wonder hold. Interesting audience security measure material. My four her door return shake Democrat. Benefit down concern green opportunity include institution degree. Forget movie member. Main or try generation. Size partner shoulder political return player. Street continue local difference provide. Help worker attention writer. Western so quality require quickly six. Produce idea seven drug house according ready. Score pass would himself. Data myself sign law individual garden power. Hair kind couple hear fire. Someone stand machine available. Enter east region including likely support popular. Sort financial check summer talk. Debate kid defense do fast miss. Peace southern fire for stop. Admit dark forward. Strategy soldier point military much arrive piece mouth. Test the treat strong others condition truth. Real avoid painting. Little eight establish too. Small over everyone its food. Morning develop light listen simply office less. I leg somebody hear response assume. Almost expect think behavior share turn. For head voice skill for one rate. Education mouth pressure news minute. Single ever international six artist mind cover. Both happen firm. Glass pretty surface senior. Form benefit situation specific. Inside outside better record best some cell site. Election billion personal last computer position rate citizen. Trial fine father now. Staff growth bag mention hit system. Firm here teacher outside. Popular get carry live company able something. Common choice expect put sit run street face. Sometimes middle if feeling. See seven agreement energy president article word. Body join word capital effect. Candidate star language direction shoulder. Address field man avoid production admit. Network democratic third remember today enough. Hear than half southern development. Story Democrat hotel woman your television. Science gun product pressure future you. Picture meeting soon reduce significant like too card. Standard stop agreement response think recently. Close body participant his usually up. As head body point because. Anything exist table occur. Wait well serve behind number. Blood card audience compare course. Bank determine any leader. Figure others cover yet smile. Economy learn control again. Attorney full war response pattern participant. Art live mind right. Center kid work decision bag establish. Often senior visit church against how. Hold paper western require interest body. Top way against series. Important reveal physical for these traditional a thing. Describe finish light front animal system bad. Media nation military professional there stop door. May box he ago bring ground develop. Set between itself`

func TestHandler(t *testing.T) {

	body, err := json.Marshal(map[string]interface{}{
		"text1": line1,
		"text2": line2,
	})
	require.NoError(t, err)

	app := fiber.New()
	app.Post("/similarity", similarityHandler)

	req, err := http.NewRequest(http.MethodPost, "/similarity", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)

	resp, err := app.Test(req, -1)
	require.NoError(t, err)
	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	_ = resp.Body.Close()

	assert.Equal(t, `{"similarity":0.519,"interpretation":"Moderately Similar"}`, string(respBody))
}

// Does not work for me (hangs), so commented out.
//func TestNormalizeTextGoPCRE(t *testing.T) {
//	normalizedLine1 := normalizeTextGoPCRE(line1)
//	fmt.Println(normalizedLine1)
//}

func TestNormalizeTextRE2(t *testing.T) {
	assert.Equal(t, normalizeText(line1), normalizeTextRE2(line1))
}
func TestNormalizeTextPCRE2(t *testing.T) {
	assert.Equal(t, normalizeText(line1), normalizeTextPCRE2(line1))
}
func TestNormalizeTextAll(t *testing.T) {
	assert.Equal(t, normalizeText(line1), normalizeTextNestedForLoops(line1))
}
func TestNormalizeBytesReplacer(t *testing.T) {
	assert.Equal(t, normalizeText(line1), string(normalizeTextBytesReplacer([]byte(line1))))
}
func TestNormalizeReplacer(t *testing.T) {
	assert.Equal(t, normalizeText(line1), normalizeTextReplacer(line1))
}
func TestNormalizeTextUsingAsciiRanges(t *testing.T) {
	assert.Equal(t, normalizeText(line1), normalizeTextUsingAsciiRanges(line1))
}
func TestNormalizeTextUsingInvertedAsciiRanges(t *testing.T) {
	assert.Equal(t, normalizeText(line1), normalizeTextUsingInvertedAsciiRanges(line1))
}
func TestNormalizeTextUsingInvertedAsciiRangesBytes(t *testing.T) {
	assert.Equal(t, normalizeText(line1), string(normalizeTextUsingInvertedAsciiRangesBytes([]byte(line1))))
}
func TestFrequencyMap(t *testing.T) {
	words1 := strings.Split(normalizeText(line1), " ")
	m := generateFrequencyMap(words1)
	assert.Equal(t, 775, len(m))
}
func TestFrequencyMapCapacity(t *testing.T) {
	words1 := strings.Split(normalizeText(line1), " ")
	m := generateFrequencyMapWithCapacity(words1)
	assert.Equal(t, 775, len(m))
}
func TestFrequencyMapSwiss(t *testing.T) {
	words1 := strings.Split(normalizeText(line1), " ")
	m := generateFrequencyMapSwiss(words1)
	assert.Equal(t, 775, m.Count())
}
func TestFrequencyMapAtomic(t *testing.T) {
	words1 := strings.Split(normalizeText(line1), " ")
	m := generateFrequencyMapAtomic(words1)
	assert.Equal(t, 775, m.Len())
}

func TestCalculateIDFVariantsProducesSameResult(t *testing.T) {
	words1 := strings.Split(line1, " ")
	words2 := strings.Split(line2, " ")

	uw := strset.NewWithSize(len(words1))
	uw.Add(append(words1, words2...)...)
	uniqueWords := uw.List()

	fm1 := generateFrequencyMap(words1)
	fm2 := generateFrequencyMap(words2)

	idf1 := calculateIDF(uniqueWords, fm1, fm2)
	idf2 := calculateIDFOptimized(uniqueWords, fm1, fm2)

	require.Equal(t, len(idf1), len(idf2))
	for i := range idf1 {
		require.Equal(t, idf1[i], idf2[i])
	}
}

func TestCalculateTFVariantsProducesSameResult(t *testing.T) {
	words1 := strings.Split(line1, " ")

	uw := strset.NewWithSize(len(words1))
	uw.Add(words1...)
	uniqueWords := uw.List()

	fm1 := generateFrequencyMap(words1)
	total := len(fm1)

	tf1 := calculateTF(uniqueWords, fm1, total)
	tf2 := calculateTFMultiply(uniqueWords, fm1, total)

	require.Equal(t, len(tf1), len(tf2))
	for i := range tf1 {
		require.InEpsilon(t, tf1[i], tf2[i], 0.000001)
	}
}

func BenchmarkOriginalRegexp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = normalizeText(line1)
		_ = normalizeText(line2)
	}
}

func BenchmarkCheat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = normalizeTextCheat(line1)
		_ = normalizeTextCheat(line2)
	}
}

func BenchmarkGoRE2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = normalizeTextRE2(line1)
		_ = normalizeTextRE2(line2)
	}
}

func BenchmarkGoPCRE2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = normalizeTextPCRE2(line1)
		_ = normalizeTextPCRE2(line2)
	}
}

func BenchmarkReplaceNestedForLoops(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = normalizeTextNestedForLoops(line1)
		_ = normalizeTextNestedForLoops(line2)
	}
}

func BenchmarkReplaceUsingAsciiRanges(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = normalizeTextUsingAsciiRanges(line1)
		_ = normalizeTextUsingAsciiRanges(line2)
	}
}

func BenchmarkReplaceUsingInvertedAsciiRanges(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = normalizeTextUsingInvertedAsciiRanges(line1)
		_ = normalizeTextUsingInvertedAsciiRanges(line2)
	}
}

func BenchmarkReplaceUsingInvertedAsciiRangesBytes(b *testing.B) {
	b1 := []byte(line1)
	b2 := []byte(line2)
	for i := 0; i < b.N; i++ {
		_ = normalizeTextUsingInvertedAsciiRangesBytes(b1)
		_ = normalizeTextUsingInvertedAsciiRangesBytes(b2)
	}
}

func BenchmarkStringsReplacer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = normalizeTextReplacer(line1)
		_ = normalizeTextReplacer(line2)
	}
}

func BenchmarkBytesReplacer(b *testing.B) {
	b1 := []byte(line1)
	b2 := []byte(line2)
	for i := 0; i < b.N; i++ {
		_ = normalizeTextBytesReplacer(b1)
		_ = normalizeTextBytesReplacer(b2)
	}
}

// -------------------------------------------------------- //
// -------------------- START JSON ------------------------ //
// -------------------------------------------------------- //
func BenchmarkEncodingJSON(b *testing.B) {
	reqData := []byte(fmt.Sprintf(`{"text1":"%s", "text2":"%s"}`, line1, line2))
	req := SimilarityRequest{}
	resp := &SimilarityResponse{
		Similarity:     0.54,
		Interpretation: string(InterpretationResultModerately),
	}
	for i := 0; i < b.N; i++ {
		_ = json.Unmarshal(reqData, &req)
		_ = req.Text1
		_ = req.Text2
		_, _ = json.Marshal(resp)
	}
}
func BenchmarkJSONSimd(b *testing.B) {
	reqData := []byte(fmt.Sprintf(`{"text1":"%s", "text2":"%s"}`, line1, line2))
	resp := &SimilarityResponse{
		Similarity:     0.54,
		Interpretation: string(InterpretationResultModerately),
	}

	var lem1 = &simdjson.Element{}
	var lem2 = &simdjson.Element{}
	var pj = &simdjson.ParsedJson{}
	for i := 0; i < b.N; i++ {
		pj, _ = simdjson.Parse(reqData, pj)

		_ = pj.ForEach(func(i simdjson.Iter) error {
			_, _ = i.FindElement(lem1, "text1")
			_, _ = lem1.Iter.String()

			_, _ = i.FindElement(lem2, "text2")
			_, _ = lem2.Iter.String()

			return nil
		})
		_, _ = json.Marshal(resp)
	}
}
func BenchmarkJSONParser(b *testing.B) {
	reqData := []byte(fmt.Sprintf(`{"text1":"%s", "text2":"%s"}`, line1, line2))
	resp := &SimilarityResponse{
		Similarity:     0.54,
		Interpretation: string(InterpretationResultModerately),
	}

	key1 := []byte(`text1`)
	key2 := []byte(`text2`)
	for i := 0; i < b.N; i++ {
		_ = jsonparser.ObjectEach(reqData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			if bytes.Equal(key, key1) {
				_ = value
			}
			if bytes.Equal(key, key2) {
				_ = value
			}
			return nil
		})
		_, _ = json.Marshal(resp)
	}
}

// -------------------------------------------------
// Map benchmarks
// -------------------------------------------------
func BenchmarkFrequencyStandardMap(b *testing.B) {
	words1 := strings.Split(normalizeText(line1), " ")
	words2 := strings.Split(normalizeText(line2), " ")
	for i := 0; i < b.N; i++ {
		_ = generateFrequencyMap(words1)
		_ = generateFrequencyMap(words2)
	}
}
func BenchmarkFrequencyStandardMapWithCapacity(b *testing.B) {
	words1 := strings.Split(normalizeText(line1), " ")
	words2 := strings.Split(normalizeText(line2), " ")
	for i := 0; i < b.N; i++ {
		_ = generateFrequencyMapWithCapacity(words1)
		_ = generateFrequencyMapWithCapacity(words2)
	}
}
func BenchmarkFrequencySwissMap(b *testing.B) {
	words1 := strings.Split(normalizeText(line1), " ")
	words2 := strings.Split(normalizeText(line2), " ")
	for i := 0; i < b.N; i++ {
		_ = generateFrequencyMapSwiss(words1)
		_ = generateFrequencyMapSwiss(words2)
	}
}

func BenchmarkFrequencyAtomicMap(b *testing.B) {
	words1 := strings.Split(normalizeText(line1), " ")
	words2 := strings.Split(normalizeText(line2), " ")
	for i := 0; i < b.N; i++ {
		_ = generateFrequencyMapAtomic(words1)
		_ = generateFrequencyMapAtomic(words2)
	}
}

// -------------------------------------------------------------
// Calculate TF benchmarks
// -------------------------------------------------------------
func BenchmarkCalculateTF(b *testing.B) {
	words1 := strings.Split(line1, " ")

	uw := strset.NewWithSize(len(words1))
	uw.Add(words1...)
	uniqueWords := uw.List()

	fm1 := generateFrequencyMap(words1)
	total := len(fm1)
	for i := 0; i < b.N; i++ {
		_ = calculateTF(uniqueWords, fm1, total)
	}
}
func BenchmarkCalculateTFMultiply(b *testing.B) {
	words1 := strings.Split(line1, " ")

	uw := strset.NewWithSize(len(words1))
	uw.Add(words1...)
	uniqueWords := uw.List()

	fm1 := generateFrequencyMap(words1)
	total := len(fm1)
	for i := 0; i < b.N; i++ {
		_ = calculateTFMultiply(uniqueWords, fm1, total)
	}
}
