package speller

import (
	"fmt"
	"testing"
)

type testCase struct {
	number   int64
	spelling string
}

func TestSpell(t *testing.T) {
	for _, tc := range []testCase{
		{number: 0, spelling: "zero"},
		{number: 1, spelling: "one"},
		{number: 2, spelling: "two"},
		{number: 3, spelling: "three"},
		{number: 4, spelling: "four"},
		{number: 5, spelling: "five"},
		{number: 6, spelling: "six"},
		{number: 7, spelling: "seven"},
		{number: 8, spelling: "eight"},
		{number: 9, spelling: "nine"},
		{number: 10, spelling: "ten"},
		{number: 11, spelling: "eleven"},
		{number: 12, spelling: "twelve"},
		{number: 13, spelling: "thirteen"},
		{number: 14, spelling: "fourteen"},
		{number: 20, spelling: "twenty"},
		{number: 22, spelling: "twenty-two"},
		{number: 50, spelling: "fifty"},
		{number: 100, spelling: "one hundred"},
		{number: 123, spelling: "one hundred twenty-three"},
		{number: 1000, spelling: "one thousand"},
		{number: 1234, spelling: "one thousand two hundred thirty-four"},
		{number: 1000000, spelling: "one million"},
		{number: 1001000, spelling: "one million one thousand"},
		{number: 1002345, spelling: "one million two thousand three hundred forty-five"},
		{number: 1000000000, spelling: "one billion"},
		{number: 999999999999, spelling: "nine hundred ninety-nine billion nine hundred ninety-nine million nine hundred ninety-nine thousand nine hundred ninety-nine"},
		{number: 987654321123, spelling: "nine hundred eighty-seven billion six hundred fifty-four million three hundred twenty-one thousand one hundred twenty-three"},
		{number: 791947779410, spelling: "seven hundred ninety-one billion nine hundred forty-seven million seven hundred seventy-nine thousand four hundred ten"},
		{number: 223082153551, spelling: "two hundred twenty-three billion eighty-two million one hundred fifty-three thousand five hundred fifty-one"},
		{number: 611666145821, spelling: "six hundred eleven billion six hundred sixty-six million one hundred forty-five thousand eight hundred twenty-one"},
		{number: 794235010051, spelling: "seven hundred ninety-four billion two hundred thirty-five million ten thousand fifty-one"},
		{number: 616287113937, spelling: "six hundred sixteen billion two hundred eighty-seven million one hundred thirteen thousand nine hundred thirty-seven"},
		{number: 724549167320, spelling: "seven hundred twenty-four billion five hundred forty-nine million one hundred sixty-seven thousand three hundred twenty"},
		{number: 647632969758, spelling: "six hundred forty-seven billion six hundred thirty-two million nine hundred sixty-nine thousand seven hundred fifty-eight"},
		{number: 317331776148, spelling: "three hundred seventeen billion three hundred thirty-one million seven hundred seventy-six thousand one hundred forty-eight"},
		{number: 949183117216, spelling: "nine hundred forty-nine billion one hundred eighty-three million one hundred seventeen thousand two hundred sixteen"},
		{number: 40480279449, spelling: "forty billion four hundred eighty million two hundred seventy-nine thousand four hundred forty-nine"},
		{number: 750760398084, spelling: "seven hundred fifty billion seven hundred sixty million three hundred ninety-eight thousand eighty-four"},
		{number: 64263669287, spelling: "sixty-four billion two hundred sixty-three million six hundred sixty-nine thousand two hundred eighty-seven"},
		{number: 410884491574, spelling: "four hundred ten billion eight hundred eighty-four million four hundred ninety-one thousand five hundred seventy-four"},
		{number: 875414458836, spelling: "eight hundred seventy-five billion four hundred fourteen million four hundred fifty-eight thousand eight hundred thirty-six"},
		{number: 871211445515, spelling: "eight hundred seventy-one billion two hundred eleven million four hundred forty-five thousand five hundred fifteen"},
		{number: 483838182873, spelling: "four hundred eighty-three billion eight hundred thirty-eight million one hundred eighty-two thousand eight hundred seventy-three"},
		{number: 275472644968, spelling: "two hundred seventy-five billion four hundred seventy-two million six hundred forty-four thousand nine hundred sixty-eight"},
		{number: 474910584091, spelling: "four hundred seventy-four billion nine hundred ten million five hundred eighty-four thousand ninety-one"},
		{number: 610539110790, spelling: "six hundred ten billion five hundred thirty-nine million one hundred ten thousand seven hundred ninety"},
		{number: 113853353331, spelling: "one hundred thirteen billion eight hundred fifty-three million three hundred fifty-three thousand three hundred thirty-one"},
		{number: 113853353331, spelling: "one hundred thirteen billion eight hundred fifty-three million three hundred fifty-three thousand three hundred thirty-one"},
		{number: -1, spelling: "minus one"},
		{number: -987654321123, spelling: "minus nine hundred eighty-seven billion six hundred fifty-four million three hundred twenty-one thousand one hundred twenty-three"},
		{number: -225682062215, spelling: "minus two hundred twenty-five billion six hundred eighty-two million sixty-two thousand two hundred fifteen"},
		{number: -894297259382, spelling: "minus eight hundred ninety-four billion two hundred ninety-seven million two hundred fifty-nine thousand three hundred eighty-two"},
		{number: -993892754487, spelling: "minus nine hundred ninety-three billion eight hundred ninety-two million seven hundred fifty-four thousand four hundred eighty-seven"},
		{number: -659484130377, spelling: "minus six hundred fifty-nine billion four hundred eighty-four million one hundred thirty thousand three hundred seventy-seven"},
		{number: -380785607485, spelling: "minus three hundred eighty billion seven hundred eighty-five million six hundred seven thousand four hundred eighty-five"},
		{number: -938823202882, spelling: "minus nine hundred thirty-eight billion eight hundred twenty-three million two hundred two thousand eight hundred eighty-two"},
		{number: -852162447597, spelling: "minus eight hundred fifty-two billion one hundred sixty-two million four hundred forty-seven thousand five hundred ninety-seven"},
		{number: -659617715691, spelling: "minus six hundred fifty-nine billion six hundred seventeen million seven hundred fifteen thousand six hundred ninety-one"},
		{number: -260967139958, spelling: "minus two hundred sixty billion nine hundred sixty-seven million one hundred thirty-nine thousand nine hundred fifty-eight"},
		{number: -369028252198, spelling: "minus three hundred sixty-nine billion twenty-eight million two hundred fifty-two thousand one hundred ninety-eight"},
		{number: -423328317660, spelling: "minus four hundred twenty-three billion three hundred twenty-eight million three hundred seventeen thousand six hundred sixty"},
		{number: -115689502993, spelling: "minus one hundred fifteen billion six hundred eighty-nine million five hundred two thousand nine hundred ninety-three"},
		{number: -671972563340, spelling: "minus six hundred seventy-one billion nine hundred seventy-two million five hundred sixty-three thousand three hundred forty"},
		{number: -784539984240, spelling: "minus seven hundred eighty-four billion five hundred thirty-nine million nine hundred eighty-four thousand two hundred forty"},
		{number: -529813952648, spelling: "minus five hundred twenty-nine billion eight hundred thirteen million nine hundred fifty-two thousand six hundred forty-eight"},
		{number: -792121679800, spelling: "minus seven hundred ninety-two billion one hundred twenty-one million six hundred seventy-nine thousand eight hundred"},
		{number: -608075137246, spelling: "minus six hundred eight billion seventy-five million one hundred thirty-seven thousand two hundred forty-six"},
		{number: -607721170055, spelling: "minus six hundred seven billion seven hundred twenty-one million one hundred seventy thousand fifty-five"},
		{number: -103492360166, spelling: "minus one hundred three billion four hundred ninety-two million three hundred sixty thousand one hundred sixty-six"},
		{number: -988549293674, spelling: "minus nine hundred eighty-eight billion five hundred forty-nine million two hundred ninety-three thousand six hundred seventy-four"},
		{number: -999999999999, spelling: "minus nine hundred ninety-nine billion nine hundred ninety-nine million nine hundred ninety-nine thousand nine hundred ninety-nine"},
	} {
		t.Run(fmt.Sprintf("%d", tc.number), func(t *testing.T) {
			if spelling := Spell(tc.number); spelling != tc.spelling {
				t.Errorf("%d should be spelled as %q; got %q", tc.number, tc.spelling, spelling)
			}
		})
	}
}
