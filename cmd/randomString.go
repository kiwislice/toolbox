package cmd

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kiwislice/toolbox/tools"
	"github.com/spf13/cobra"
)

var (
	randomStringSrc    string
	randomStringLength int
	randomStringCount  int
)

var randomStringCmd = &cobra.Command{
	Use:   "randomString",
	Short: "Generate random strings",
	Long:  `Generate random strings with specified characters, length, and count.`,
	Run: func(cmd *cobra.Command, args []string) {
		tools.Debug("開始產生隨機字串")
		defer tools.Debug("結束產生隨機字串")

		for i := 0; i < randomStringCount; i++ {
			s := randomString(randomStringSrc, randomStringLength)
			fmt.Println(s)
		}
	},
}

func randomString(src string, length int) string {
	seed := time.Now().UnixNano() + rand.Int63()
	r := rand.New(rand.NewSource(seed))
	ar := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		index := r.Intn(len(src))
		ar = append(ar, src[index])
	}
	return string(ar)
}

func init() {
	rootCmd.AddCommand(randomStringCmd)
	randomStringCmd.Flags().StringVar(&randomStringSrc, "src", "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz", "character source")
	randomStringCmd.Flags().IntVar(&randomStringLength, "length", 8, "random string length")
	randomStringCmd.Flags().IntVar(&randomStringCount, "count", 1, "number of random strings")
}
