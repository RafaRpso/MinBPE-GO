package main

import (
	"fmt"
	"os"
)

type Config struct {
	vocabSize int
	numMerges int
}

type PairUint16 struct {
	p1 uint16
	p2 uint16
}

func GetVocab() ([]uint16, error) {

	file, err := os.ReadFile("./files/vocab.txt")
	if err != nil {
		return nil, fmt.Errorf("Erro ao ler o arquivo: %v", err)
	}
	fileContent := string(file)

	vocab := make([]uint16, len(fileContent))
	for i, char := range fileContent {
		vocab[i] = uint16(char)
	}

	return vocab, nil
}

func FromVocabGetPairs(vocab []uint16) (map[PairUint16]int, error) {

	var pairs []PairUint16
	pairsDc := make(map[PairUint16]int)

	for i := 0; i < len(vocab); i++ {
		if i+1 < len(vocab) {
			pairs = append(pairs, PairUint16{vocab[i], vocab[i+1]})
		}
	}

	for i := 0; i < len(pairs); i++ {
		if val, ok := pairsDc[pairs[i]]; ok {
			pairsDc[pairs[i]] = val + 1
		} else {
			pairsDc[pairs[i]] = 1
		}

	}

	return pairsDc, nil
}

func MergeBPE(tokens []uint16, topPair PairUint16, newToken uint16) ([]uint16, error) {
	var newTokens []uint16
	for i := 0; i < len(tokens); {
		if i+2 < len(tokens) && uint16(topPair.p1) == tokens[i] && uint16(topPair.p2) == tokens[i+1] {
			newTokens = append(newTokens, newToken)
			i = i + 2
		} else {
			newTokens = append(newTokens, uint16(tokens[i]))
			i++
		}
	}

	return newTokens, nil
}

func GetTopPair(pairs map[PairUint16]int) (PairUint16, int, error) {
	idx := -1
	var topPair PairUint16
	for pair, i := range pairs {
		if i > idx {
			idx = i
			topPair = pair
		}
	}
	return topPair, idx, nil
}

func VocabToUint16(vocab []byte) []uint16 {
	if len(vocab)%2 != 0 {
		panic("O slice de bytes deve ter um número par de elementos")
	}

	result := make([]uint16, len(vocab)/2)

	for i := 0; i < len(result); i++ {
		result[i] = uint16(vocab[2*i]) | uint16(vocab[2*i+1])<<8
	}

	return result
}
func ByteToUint16(b []byte) []uint16 {
	var buinit16 []uint16

	for i := range b {
		buinit16 = append(buinit16, uint16(b[i]))
	}

	return buinit16
}
func MergePairs(config Config) map[PairUint16]int {

	var merges = make(map[PairUint16]int)
	ids, _ := GetVocab()
	for i := 0; i < config.numMerges; i++ {
		idx := uint16(256 + i)
		nids, _ := FromVocabGetPairs(ids)
		topPair, _, _ := GetTopPair(nids)
		ids, _ = MergeBPE(ids, topPair, idx)
		pairUint16 := PairUint16{uint16(topPair.p1), uint16(topPair.p2)}

		merges[pairUint16] = int(idx)

	}

	return merges

}

func GetVocabMerges(merges map[PairUint16]int) map[int][]byte {

	vocab := make(map[int][]byte, 256)

	for idx := 0; idx < 256; idx++ {
		vocab[idx] = []byte{byte(idx)}
	}

	for pair, idx := range merges {
		p0, p1 := pair.p1, pair.p2
		vocab[idx] = append(vocab[int(p0)], vocab[int(p1)]...)
	}

	return vocab
}

func Decode(vocab map[int][]byte) {
	tokens := ""
	for i := 0; i < len(vocab); i++ {
		tokens := string(vocab[i])

	}
	// socorro.

}

func Encode() {

}

func main() {
	//TODO: GET VOCAB
	//TODO: DECODE
	//TODO: ENCODE
	// TODO: REGEX (GPT4)
	// TODO: IMPLEMENTS TWO CHILDS : TOKENIZER CLASSIC, REGEX TOKENIZER

	vsize := 276
	config := Config{
		vocabSize: vsize,
		numMerges: vsize - 256,
	}
	merge := MergePairs(config)
	fmt.Println(GetVocabMerges(merge))
}
