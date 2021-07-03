package withdb

import (
	"github.com/ppdraga/go-shortener/internal/shortener/link/datatype"
	"gorm.io/gorm"
	"strings"
)

type WithDB struct {
	db *gorm.DB
}

func New(db *gorm.DB) *WithDB {
	return &WithDB{
		db: db,
	}
}

func (wdb *WithDB) ReadLink(id int64) (*datatype.Link, error) {
	var linkItem LinkModel
	err := wdb.db.First(&linkItem, id).Error
	if err != nil {
		return &datatype.Link{}, err
	}

	return &datatype.Link{
		ID:         &linkItem.ID,
		Resource:   &linkItem.Resource,
		ShortLink:  &linkItem.ShortLink,
		CustomName: &linkItem.CustomName,
	}, nil
}

func (wdb *WithDB) WriteLink(external *datatype.Link) error {
	var customName string
	if external.CustomName != nil {
		customName = *external.CustomName
	}
	shortLinkNum := getFreeNumber(customName, wdb.db)

	shortLink := toShort(shortLinkNum - 1)

	if customName != "" {
		if shortLinkNum == 1 {
			shortLink = customName + "_"
		} else {
			shortLink = customName + "_" + shortLink
		}
	}

	link := LinkModel{
		Resource:     *external.Resource,
		ShortLink:    shortLink,
		ShortLinkNum: shortLinkNum,
		CustomName:   customName,
	}
	err := wdb.db.Create(&link).Error

	return err
}

func getFreeNumber(customName string, db *gorm.DB) int64 {
	var result int64
	if customName == "" {
		db.Raw("SELECT COALESCE(MAX(short_link_num), 0) FROM links WHERE custom_name IS NULL OR custom_name = ''").Scan(&result)
	} else {
		db.Raw("SELECT COALESCE(MAX(short_link_num), 0) FROM links WHERE custom_name = ?", customName).Scan(&result)
	}
	return result + 1
}

func toShort(num int64) string {
	alphabetLen := len(alphabet)
	res := []string{}
	n := int(num)
	for n >= alphabetLen {
		remainder := n % alphabetLen
		n = n / alphabetLen
		res = append(res, alphabet[remainder])
	}
	res = append(res, alphabet[n])
	reverse(res)
	return strings.Join(res, "")
}

var alphabet = map[int]string{
	0: "a", 1: "b", 2: "c", 3: "d", 4: "e",
	5: "f", 6: "g", 7: "h", 8: "i", 9: "j",
	10: "k", 11: "l", 12: "m", 13: "n", 14: "o",
	15: "p", 16: "q", 17: "r", 18: "s", 19: "t",
	20: "u", 21: "v", 22: "w", 23: "x", 24: "y",
	25: "z",
	26: "A", 27: "B", 28: "C", 29: "D", 30: "E",
	31: "F", 32: "G", 33: "H", 34: "I", 35: "J",
	36: "K", 37: "L", 38: "M", 39: "N", 40: "O",
	41: "P", 42: "Q", 43: "R", 44: "S", 45: "T",
	46: "U", 47: "V", 48: "W", 49: "X", 50: "Y",
	51: "Z",
	52: "0", 53: "1", 54: "2", 55: "3", 56: "4",
	57: "5", 58: "6", 59: "7", 60: "8", 61: "9",
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}
