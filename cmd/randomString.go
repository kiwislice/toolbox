// 該檔案定義了 randomString command
package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kiwislice/toolbox/tools"
	"github.com/spf13/cobra"
)

// randomStringSrc 是隨機字串的字元來源
// randomStringLength 是隨機字串的長度
// randomStringCount 是要產生的隨機字串數量
var (
	randomStringSrc    string
	randomStringLength int
	randomStringCount  int
)

// randomStringCmd 是 `toolbox randomString` command 的定義
var randomStringCmd = &cobra.Command{
	Use:   "randomString",
	Short: "產生隨機字串",
	Long:  `根據指定的字元、長度和數量來產生隨機字串`,
	Run: func(cmd *cobra.Command, args []string) {
		tools.Debug("開始產生隨機字串")
		defer tools.Debug("結束產生隨機字串")

		// 根據指定的數量，迴圈產生隨機字串
		for i := 0; i < randomStringCount; i++ {
			s := randomString(randomStringSrc, randomStringLength)
			fmt.Println(s)
		}
	},
}

// randomString 函式會根據指定的來源字串和長度，產生一個隨機字串
func randomString(src string, length int) string {
	// 使用時間和隨機數來當作亂數種子，避免每次產生的字串都一樣
	seed := time.Now().UnixNano() + rand.Int63()
	r := rand.New(rand.NewSource(seed))
	// 建立一個 byte slice，用來存放隨機字串的字元
	ar := make([]byte, 0, length)
	// 根據指定的長度，迴圈從來源字串中隨機挑選一個字元，加入到 slice 中
	for i := 0; i < length; i++ {
		index := r.Intn(len(src))
		ar = append(ar, src[index])
	}
	// 將 byte slice 轉換為字串
	return string(ar)
}

func init() {
	// 將 randomStringCmd 加入 rootCmd 中
	rootCmd.AddCommand(randomStringCmd)
	// 設定 randomStringCmd 的 flags
	randomStringCmd.Flags().StringVar(&randomStringSrc, "src", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", "character source")
	randomStringCmd.Flags().IntVar(&randomStringLength, "length", 8, "random string length")
	randomStringCmd.Flags().IntVar(&randomStringCount, "count", 1, "number of random strings")
}
