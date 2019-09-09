package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
)

var (
	totalParticipants int
	totalWinners      int
	randomSeedHex     string
	winningNumberSet  map[int]bool
	winningResult     []int
)

func init() {
	totalParticipants = 0
	totalWinners = 0
	randomSeedHex = ""
}

func SetParams(_totalParticipants, _totalWinners int, _randomSeedHex string) {
	totalParticipants = _totalParticipants
	totalWinners = _totalWinners
	randomSeedHex = _randomSeedHex
}

/**
 * @desc Sample a number.
 * @param nonce Nonce in this time.
 * @return
 */
func pickOneWinner(nonce int) int {
	var digitalNeeded int
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(randomSeedHex + strconv.Itoa(nonce)))
	var randomSeedHexLinkNonceHash = hex.EncodeToString(sha256Hash.Sum(nil))
	magnitude := math.Log10(float64(totalParticipants))
	if math.Ceil(magnitude) == magnitude {
		digitalNeeded = int(magnitude) + 1
	} else {
		digitalNeeded = int(math.Ceil(magnitude))
	}
	randomSeedParams, err := hex.DecodeString(randomSeedHexLinkNonceHash)
	if err != nil {
		panic(err)
	}
	params := new(big.Int)
	params.SetBytes(randomSeedParams)
	return getluckyNumber(params, digitalNeeded)
}

/**
 * @desc Sample some lottery numbers which is equal to the pre-fixed winners number.
 * @return
 */
func drawAllWinningNumbers() int {
	winningNumberSet = make(map[int]bool)
	winningResult = winningResult[0:0]
	var nonce = 0
	var winningNumber int
	for len(winningNumberSet) < totalWinners {
		winningNumber = pickOneWinner(nonce)
		if v, ok := winningNumberSet[winningNumber]; ok && v {
			nonce += 1
			continue
		}

		if winningNumber > totalParticipants || winningNumber == 0 {
			nonce += 1
			continue
		}
		winningNumberSet[winningNumber] = true
		winningResult = append(winningResult, winningNumber)
	}
	nonce += 1
	return nonce
}

/**
 * @desc Get least some decimal characters in a bigInteger in HEX. Get `digital` characters.
 * @param bigInteger
 * @param magnitude The number of characters to be get.
 * @return
 */
func getluckyNumber(bigInteger *big.Int, digitals int) int {
	reversedBigInteger := reverse(hex.EncodeToString(bigInteger.Bytes()))
	var trimedString = ""
	for i := 0; i < len(reversedBigInteger); i++ {
		charI := reversedBigInteger[i]
		if charI < 97 {
			trimedString += string(charI)
		}
	}

	if len(trimedString) < digitals {
		return -1
	}
	reversedTrimedString := reverse(trimedString[0:digitals])
	res, err := strconv.Atoi(reversedTrimedString)
	if err != nil {
		panic(err)
	}
	return res
}

func isCheckParams(_totalParticipants, _totalWinners int, _randomSeedHex string) error {
	if _totalParticipants <= 0 || _totalWinners <= 0 || _totalParticipants < _totalWinners || len(_randomSeedHex) == 0 {
		return errors.New("Arguments of the sampler is illegal. Please check! ")
	}
	return nil
}

func reverse(str string) string {
	reverseStr := []rune(str)
	strLen := len(reverseStr)
	for i := 0; i < strLen/2; i++ {
		reverseStr[i], reverseStr[strLen-i-1] = reverseStr[strLen-i-1], reverseStr[i]
	}
	return string(reverseStr)
}

func mapString(mapStr map[int]bool) string {
	result := ""
	for _, value := range winningResult {
		result += fmt.Sprintf("%d\n", value)
	}
	return result
}

func run() {
	drawAllWinningNumbers()
	fmt.Printf("Lottery sampler set block hash 0x%s.\n", randomSeedHex)
	fmt.Println("The winning numbers are : ")
	fmt.Printf(mapString(winningNumberSet))
}

func main() {
	var err error
	readerBuffer := bufio.NewReader(os.Stdin)
	fmt.Printf("TotalParticipants: ")
	participantsNumber, err := readerBuffer.ReadString('\n')
	if err != nil {
		panic(err)
	}
	totalParticipants, err = strconv.Atoi(participantsNumber[:len(participantsNumber)-1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("TotalWinners: ")
	winnersNumber, err := readerBuffer.ReadString('\n')
	if err != nil {
		panic(err)
	}
	totalWinners, err = strconv.Atoi(winnersNumber[:len(winnersNumber)-1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("RandomSeedHex: ")
	randomSeedHex, err = readerBuffer.ReadString('\n')
	if err != nil {
		panic(err)
	}
	if randomSeedHex[:2] == "0x" || randomSeedHex[:2] == "0X" {
		randomSeedHex = randomSeedHex[2 : len(randomSeedHex)-1]
	} else {
		randomSeedHex = randomSeedHex[:len(randomSeedHex)-1]
	}
	err = isCheckParams(totalParticipants, totalWinners, randomSeedHex)
	if err != nil {
		panic(err)
	}

	run()
}

//func main() {
//	resFile, err := os.OpenFile("winnerNumbersFrom6000001_6002000.txt", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
//	if err != nil {
//		panic(err)
//	}
//	defer resFile.Close()
//
//	openFile, err :=os.Open("./lbklottery/ethBlockHash6000001_6002000.txt")
//	if err != nil {
//		panic(err)
//	}
//	defer openFile.Close()
//
//	scanner := bufio.NewScanner(openFile)
//	for scanner.Scan() {
//		totalParticipants = 2000
//		totalWinners	  = 300
//		randomSeedHex =  scanner.Text()
//		fmt.Println(randomSeedHex)
//		if randomSeedHex == "" {
//			continue
//		}
//		if randomSeedHex[:2] == "0x" || randomSeedHex[:2] == "0X" {
//			randomSeedHex = randomSeedHex[2:]
//		} else {
//			randomSeedHex = randomSeedHex[:]
//		}
//		run()
//		writer := bufio.NewWriter(resFile)
//		writer.WriteString(mapString(winningNumberSet))
//		writer.Flush()
//	}
//}
