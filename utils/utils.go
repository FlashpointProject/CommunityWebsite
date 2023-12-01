package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	letterBytes   = "abcdefghijklmnopqrstuvwxyz0123456789"
)

type ContentRef struct {
	ContentType string
	ContentID   string
}

type RealRandomString struct {
	src rand.Source
}

func NewRealRandomStringProvider() *RealRandomString {
	return &RealRandomString{
		src: rand.NewSource(time.Now().UnixNano()),
	}
}

func (r *RealRandomString) RandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, r.src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = r.src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func StringInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func RemoveSliceDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func ParseContentRef(raw string) ([]ContentRef, error) {
	parts := strings.Split(raw, ":")
	refs := make([]ContentRef, 0)
	for _, part := range parts {
		partSplit := strings.Split(part, "_")
		if len(partSplit) != 2 {
			return nil, fmt.Errorf("invalid content ref, invalid structure")
		}
		refs = append(refs, ContentRef{
			ContentType: partSplit[0],
			ContentID:   partSplit[1],
		})
	}
	if len(refs) == 0 {
		return nil, fmt.Errorf("invalid content ref, no refs found")
	}
	return refs, nil
}
