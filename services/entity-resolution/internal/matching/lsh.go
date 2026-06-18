package matching

import (
	"encoding/binary"
	"hash/fnv"
	"math"
	"math/rand"
	"strings"
	"sync"
)

const lshPrime = 2147483647

var (
	coeffsA []int
	coeffsB []int
	coeffMu sync.Mutex
)

func getCoeffs(numHashes int) ([]int, []int) {
	coeffMu.Lock()
	defer coeffMu.Unlock()

	if len(coeffsA) != numHashes {
		coeffsA = make([]int, numHashes)
		coeffsB = make([]int, numHashes)
		for i := 0; i < numHashes; i++ {
			coeffsA[i] = rand.Intn(lshPrime-1) + 1
			coeffsB[i] = rand.Intn(lshPrime)
		}
	}

	return coeffsA, coeffsB
}

type IndexableIdentity struct {
	ID        string
	NNU       string
	FullName  string
	FirstName string
	LastName  string
	DOB       string
	TaxID     string
}

type LSHIndex struct {
	numHashes   int
	numBands    int
	rowsPerBand int
	buckets     map[int][]string
	signatures  map[string][]int
	threshold   float64
	mu          sync.RWMutex
}

func NewLSHIndex(numHashes, numBands, rowsPerBand int) *LSHIndex {
	return &LSHIndex{
		numHashes:   numHashes,
		numBands:    numBands,
		rowsPerBand: rowsPerBand,
		buckets:     make(map[int][]string),
		signatures:  make(map[string][]int),
		threshold:   1.0 / math.Pow(float64(numBands), 1.0/float64(rowsPerBand)),
	}
}

func (idx *LSHIndex) Insert(id string, features []float64) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	sig := minHashSignature(features, idx.numHashes)
	idx.signatures[id] = sig

	for band := 0; band < idx.numBands; band++ {
		start := band * idx.rowsPerBand
		end := start + idx.rowsPerBand
		if end > len(sig) {
			end = len(sig)
		}
		if start < end {
			bh := hashBand(sig[start:end], band)
			idx.buckets[bh] = append(idx.buckets[bh], id)
		}
	}
}

func (idx *LSHIndex) Query(features []float64) []string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	sig := minHashSignature(features, idx.numHashes)
	seen := make(map[string]bool)

	for band := 0; band < idx.numBands; band++ {
		start := band * idx.rowsPerBand
		end := start + idx.rowsPerBand
		if end > len(sig) {
			end = len(sig)
		}
		if start >= end {
			continue
		}
		bh := hashBand(sig[start:end], band)
		for _, id := range idx.buckets[bh] {
			seen[id] = true
		}
	}

	result := make([]string, 0, len(seen))
	for id := range seen {
		result = append(result, id)
	}
	return result
}

func (idx *LSHIndex) Remove(id string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	sig, ok := idx.signatures[id]
	if !ok {
		return
	}

	for band := 0; band < idx.numBands; band++ {
		start := band * idx.rowsPerBand
		end := start + idx.rowsPerBand
		if end > len(sig) {
			end = len(sig)
		}
		if start < end {
			bh := hashBand(sig[start:end], band)
			bucket := idx.buckets[bh]
			for i, bid := range bucket {
				if bid == id {
					idx.buckets[bh] = append(bucket[:i], bucket[i+1:]...)
					break
				}
			}
		}
	}

	delete(idx.signatures, id)
}

func (idx *LSHIndex) Build(identities []IndexableIdentity) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.buckets = make(map[int][]string)
	idx.signatures = make(map[string][]int)

	for _, ident := range identities {
		set := extractFeatureIndices(ident)
		sig := minHashFromSet(set, idx.numHashes)
		idx.signatures[ident.ID] = sig

		for band := 0; band < idx.numBands; band++ {
			start := band * idx.rowsPerBand
			end := start + idx.rowsPerBand
			if end > len(sig) {
				end = len(sig)
			}
			if start < end {
				bh := hashBand(sig[start:end], band)
				idx.buckets[bh] = append(idx.buckets[bh], ident.ID)
			}
		}
	}
}

func minHashSignature(features []float64, numHashes int) []int {
	set := make([]int, 0)
	for i, v := range features {
		if v > 0 {
			set = append(set, i)
		}
	}
	return minHashFromSet(set, numHashes)
}

func minHashFromSet(featureSet []int, numHashes int) []int {
	a, b := getCoeffs(numHashes)

	sig := make([]int, numHashes)
	for i := range sig {
		sig[i] = lshPrime
	}

	for h := 0; h < numHashes; h++ {
		ah := a[h]
		bh := b[h]
		for _, f := range featureSet {
			fm := f % lshPrime
			if fm < 0 {
				fm += lshPrime
			}
			hv := (ah*fm + bh) % lshPrime
			if hv < 0 {
				hv += lshPrime
			}
			if hv < sig[h] {
				sig[h] = hv
			}
		}
	}

	return sig
}

func hashBand(band []int, bandIndex int) int {
	h := fnv.New32a()
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(bandIndex))
	h.Write(buf)
	for _, v := range band {
		binary.LittleEndian.PutUint32(buf, uint32(v))
		h.Write(buf)
	}
	return int(h.Sum32())
}

func extractFeatureIndices(id IndexableIdentity) []int {
	seen := make(map[int]bool)

	addStr := func(prefix, val string) {
		fnvHash := fnv.New32a()
		fnvHash.Write([]byte(strings.ToUpper(prefix + val)))
		pos := int(fnvHash.Sum32()) & math.MaxInt32
		seen[pos] = true
	}

	full := strings.ToUpper(strings.TrimSpace(id.FullName))
	for i := 0; i < len(full)-1; i++ {
		addStr("BG_", full[i:i+2])
	}
	for _, w := range strings.Fields(full) {
		addStr("W_", w)
	}

	if id.FirstName != "" {
		addStr("FN_", id.FirstName)
	}
	if id.LastName != "" {
		addStr("LN_", id.LastName)
	}
	if id.DOB != "" {
		addStr("DOB_", id.DOB)
	}
	if id.TaxID != "" {
		addStr("TAX_", id.TaxID)
	}
	if id.NNU != "" {
		addStr("NNU_", id.NNU)
	}

	result := make([]int, 0, len(seen))
	for p := range seen {
		result = append(result, p)
	}
	return result
}

func IdentitiesToIndexable(ids []IdentityStoreRecord) []IndexableIdentity {
	result := make([]IndexableIdentity, len(ids))
	for i, id := range ids {
		result[i] = IndexableIdentity{
			ID:        id.GetID(),
			NNU:       id.GetNNU(),
			FullName:  id.GetFullName(),
			FirstName: id.GetFirstName(),
			LastName:  id.GetLastName(),
			DOB:       id.GetDOB(),
			TaxID:     id.GetTaxID(),
		}
	}
	return result
}

type IdentityStoreRecord interface {
	GetID() string
	GetNNU() string
	GetFullName() string
	GetFirstName() string
	GetLastName() string
	GetDOB() string
	GetTaxID() string
}
