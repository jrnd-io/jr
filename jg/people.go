package jg

import (
	"strings"
)

type Person struct {
	Name    string
	Surname string
	Gender  string
}

type City struct {
	Name  string
	Zip   string
	State string
}

var maleAll string = "James\nRobert\nJohn\nMichael\nDavid\nWilliam\nRichard\nJoseph\nThomas\nCharles\nChristopher\nDaniel" +
	"\nMatthew\nAnthony\nMark\nDonald\nSteven\nPaul\nAndrew\nJoshua\nKenneth\nKevin\nBrian\nGeorge\nTimothy\nRonald\nEdward" +
	"\nJason\nJeffrey\nRyan\nJacob\nGary\nNicholas\nEric\nJonathan\nStephen\nLarry\nJustin\nScott\nBrandon\nBenjamin\nSamuel" +
	"\nGregory\nAlexander\nFrank\nPatrick\nRaymond\nJack\nDennis\nJerry\nTyler\nAaron\nJose\nAdam\nNathan\nHenry\nDouglas" +
	"\nZachary\nPeter\nKyle\nEthan\nWalter\nNoah\nJeremy\nChristian\nKeith\nRoger\nTerry\nGerald\nHarold\nSean\nAustin\nCarl" +
	"\nArthur\nLawrence\nDylan\nJesse\nJordan\nBryan\nBilly\nJoe\nBruce\nGabriel\nLogan\nAlbert\nWillie\nAlan\nJuan\nWayne\nElijah" +
	"\nRandy\nRoy\nVincent\nRalph\nEugene\nRussell\nBobby\nMason\nPhilip\nLouis"
var males []string = strings.Split(maleAll, "\n")
var femaleAll string = "Mary\nPatricia\nJennifer\nLinda\nElizabeth\nBarbara\nSusan\nJessica\nSarah\nKaren\nLisa\nNancy\nBetty" +
	"\nMargaret\nSandra\nAshley\nKimberly\nEmily\nDonna\nMichelle\nCarol\nAmanda\nDorothy\nMelissa\nDeborah\nStephanie\nRebecca" +
	"\nSharon\nLaura\nCynthia\nKathleen\nAmy\nAngela\nShirley\nAnna\nBrenda\nPamela\nEmma\nNicole\nHelen\nSamantha\nKatherine" +
	"\nChristine\nDebra\nRachel\nCarolyn\nJanet\nCatherine\nMaria\nHeather\nDiane\nRuth\nJulie\nOlivia\nJoyce\nVirginia\nVictoria" +
	"\nKelly\nLauren\nChristina\nJoan\nEvelyn\nJudith\nMegan\nAndrea\nCheryl\nHannah\nJacqueline\nMartha\nGloria\nTeresa\nAnn\nSara" +
	"\nMadison\nFrances\nKathryn\nJanice\nJean\nAbigail\nAlice\nJulia\nJudy\nSophia\nGrace\nDenise\nAmber\nDoris\nMarilyn\nDaniell" +
	"e\nBeverly\nIsabella\nTheresa\nDiana\nNatalie\nBrittany\nCharlotte\nMarie\nKayla\nAlexis\nLori"
var females []string = strings.Split(femaleAll, "\n")
var surnamesAll string = "Smith\nJohnson\nWilliams\nBrown\nJones\nGarcia\nMiller\nDavis\nRodriguez\nMartinez\nHernandez" +
	"\nLopez\nGonzalez\nWilson\nAnderson\nThomas\nTaylor\nMoore\nJackson\nMartin\nLee\nPerez\nThompson\nWhite\nHarris" +
	"\nSanchez\nClark\nRamirez\nLewis\nRobinson\nWalker\nYoung\nAllen\nKing\nWright\nScott\nTorres\nNguyen\nHill\nFlores" +
	"\nGreen\nAdams\nNelson\nBaker\nHall\nRivera\nCampbell\nMitchell\nCarter\nRoberts\nGomez\nPhillips\nEvans\nTurner\nDia" +
	"z\nParker\nCruz\nEdwards\nCollins\nReyes\nStewart\nMorris\nMorales\nMurphy\nCook\nRogers\nGutierrez\nOrtiz\nMorgan" +
	"\nCooper\nPeterson\nBailey\nReed\nKelly\nHoward\nRamos\nKim\nCox\nWard\nRichardson\nWatson\nBrooks\nChavez\nWood" +
	"\nJames\nBennett\nGray\nMendosa\nRuiz\nHughes\nPrice\nAlvarez\nCastillo\nSanders\nPatel\nMyers\nLong\nRoss\nFoster\nJimenez"
var surnames []string = strings.Split(surnamesAll, "\n")
var usStatesAll = "Alabama\nAlaska\nArizona\nArkansas\nCalifornia\nColorado\nConnecticut\nDelaware\nFlorida\nGeorgia" +
	"\nHawaii\nIdaho\nIllinois\nIndiana\nIowa\nKansas\nKentucky\nLouisiana\nMaine\nMaryland\nMassachusetts\nMichigan" +
	"\nMinnesota\nMississippi\nMissouri\nMontana\nNebraska\nNevada\nNew Hampshire\nNew Jersey\nNew Mexico\nNew York" +
	"\nNorth Carolina\nNorth Dakota\nOhio\nOklahoma\nOregon\nPennsylvania\nRhode Island\nSouth Carolina\nSouth Dakota" +
	"\nTennessee\nTexas\nUtah\nVermont\nVirginia\nWashington\nWest Virginia\nWisconsin\nWyoming"
var usStates = strings.Split(usStatesAll, "\n")
var usStatesShortAll = "AL\nAK\nAZ\nAR\nCA\nCO\nCT\nDE\nFL\nGA\nHI\nID\nIL\nIN\nIA\nKS\nKY\nLA\nME\nMD\nMA\nMI\nMN\nMS" +
	"\nMO\nMT\nNE\nNV\nNH\nNJ\nNM\nNY\nNC\nND\nOH\nOK\nOR\nPA\nRI\nSC\nSD\nTN\nTX\nUT\nVT\nVA\nWA\nWV\nWI\nWY"
var usStatesShort = strings.Split(usStatesShortAll, "\n")
var usCapitalsAll = "Montgomery\nJuneau\nPhoenix\nLittle Rock\nSacramento\nDenver\nHartford\nDover\nTallahassee" +
	"\nAtlanta\nHonolulu\nBoise\nSpringfield\nIndianapolis\nDes Moines\nTopeka\nFrankfort\nBaton Rouge\nAugusta" +
	"\nAnnapolis\nBoston\nLansing\nSt.Paul\nJackson\nJefferson City\nHelena\nLincoln\nCarson City\nConcord\nTrenton" +
	"\nSanta Fe\nAlbany\nRaleigh\nBismarck\nColumbus\nOklahoma City\nSalem\nHarrisburg\nProvidence\nColumbia\nPierre" +
	"\nNashville\nAustin\nSalt Lake City\nMontpelier\nRichmond\nOlympia\nCharleston\nMadison\nCheyenne"
var usCapitals = strings.Split(usCapitalsAll, "\n")
var zipCodesAll = "36104\n99801\n85001\n72201\n95814\n80202\n6103\n19901\n32301\n30303\n96813\n83702\n62701\n46225" +
	"\n50309\n66603\n40601\n70802\n4330\n21401\n02201\n48933\n55102\n39205\n65101\n59623\n68502\n89701\n3301\n8608\n87501" +
	"\n12207\n27601\n58501\n43215\n73102\n97301\n17101\n2903\n29217\n57501\n37219\n78701\n84111\n5602\n23219\n98507" +
	"\n25301\n53703\n82001" //\n96799\n20001\n96941\n96910\n96960\n96950\n96939\n901\n802"
var zipCodes = strings.Split(zipCodesAll, "\n")

func name() string {
	s := Random.Intn(2)
	if s == 0 {
		return males[Random.Intn(len(males))]

	} else {
		return females[Random.Intn(len(males))]
	}
}

func nameM() string {
	return males[Random.Intn(len(males))]
}

func nameF() string {
	return females[Random.Intn(len(males))]
}

func surname() string {
	return surnames[Random.Intn(len(surnames))]
}

func middlename() string {
	middles := []string{"M", "J", "K", "P", "T", "S"}
	return middles[Random.Intn(len(middles))]
}

func address() string {
	addresses := []string{"80", "81", "443", "22", "631"}
	return addresses[Random.Intn(len(addresses))]
}

func state() string {
	return stateAt(Random.Intn(len(usStates)))
}

func stateAt(index int) string {
	return usStates[index]
}

func stateShort() string {
	return stateAt(Random.Intn(len(usStatesShort)))
}

func stateShortAt(index int) string {
	return usStatesShort[index]
}

func capital() string {
	return capitalAt(Random.Intn(len(usCapitals)))
}

func capitalAt(index int) string {
	return usCapitals[index]
}

func zip() string {
	return zipAt(Random.Intn(len(zipCodes)))
}

func zipAt(index int) string {
	return zipCodes[index]
}

func company() string {
	companies := []string{"Acme Corporation", "Globex Corporation", "Soylent Corp", "Initech", "Umbrella Corporation",
		"Hooli", "Veement Capital Partners", "Massive Dynamics", "Evil Partners", "Angels Investors", "Boston Static"}
	return companies[Random.Intn(len(companies))]
}
