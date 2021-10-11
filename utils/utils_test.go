package utils_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/ahwhy/myGolang/utils"
)

func TestBinaryFormat(t *testing.T) {
	fmt.Println(utils.BinaryFormat(4095))
}

func TestHumanBytesLoaded(t *testing.T) {
	fmt.Println(utils.HumanBytesLoaded(1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024 * 1024 * 1024))
	fmt.Println(utils.HumanBytesLoaded(1024 * 1024 * 1024 * 1024 * 1024 * 1024))
}

func TestMd5(t *testing.T) {
	fmt.Println(utils.Md5("Atlantis"))
}

func TestMd5Salt(t *testing.T) {
	fmt.Println(utils.Md5Salt("Atlantis", ""))
}

func TestSplitMd5Salt(t *testing.T) {
	salt, md5 := utils.SplitMd5Salt(utils.Md5Salt("Atlantis", ""))
	fmt.Printf("Salt: %s \nMd5: %s\n", salt, md5)
}

func TestRandString(t *testing.T) {
	fmt.Println(utils.RandString(10))
}

// func daysAgo(t time.Time) int {
// 	return int(time.Since(t).Hours() / 24)
// }

func TestSearchIssues(t *testing.T) {
	// result, err := github.SearchIssues(os.Args[1:])
	result, err := utils.SearchIssues([]string{"repo:golang/go is:open json decoder"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)

	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}

	// 	templ := `
	// ----------------------------------------
	// Number:    {{daysAgo}}
	// User:      {{User}}
	// Title:     {{Title}}
	// Age:       {{Age}}
	// ----------------------------------------
	// `

	// 	var report = template.Must(template.New("issuelist").
	// 		Funcs(template.FuncMap{"daysAgo": daysAgo}).
	// 		Parse(templ))

	// 	if err := report.Execute(os.Stdout, result); err != nil {
	// 		log.Fatal(err)
	// 	}
}

func TestSnake(t *testing.T) {
	fmt.Println(utils.Snake("qwer1234ASDF"))
}
