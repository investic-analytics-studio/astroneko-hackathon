package history

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractJSONFromMessage_WithValidJSON(t *testing.T) {
	// Arrange
	message := `โอ้ ท่านผู้ใฝ่หาคำทำนาย ข้าจะหยิบไพ่จากห้วงจักรวาลมาให้ท่าน...

ไพ่ที่ปรากฏต่อหน้าท่านคือ
**THE_WHEEL_OF_FORTUNE**
โชคชะตาหมุนเวียน เปลี่ยนแปลงไม่สิ้นสุด วัฏจักรแห่งชีวิตนำพาท่านไปสู่โอกาสใหม่ ๆ บางครั้งขึ้น บางครั้งลง จงเตรียมใจรับมือกับความเปลี่ยนแปลง และเปิดรับสิ่งใหม่ที่จักรวาลจะมอบให้

Let's your Astro Shine !

` + "```json\n" + `{
"card": "THE_WHEEL_OF_FORTUNE",
"meaning": "Fate, cycles, and unexpected changes"
}
` + "```"

	expectedMessage := `โอ้ ท่านผู้ใฝ่หาคำทำนาย ข้าจะหยิบไพ่จากห้วงจักรวาลมาให้ท่าน...

ไพ่ที่ปรากฏต่อหน้าท่านคือ
**THE_WHEEL_OF_FORTUNE**
โชคชะตาหมุนเวียน เปลี่ยนแปลงไม่สิ้นสุด วัฏจักรแห่งชีวิตนำพาท่านไปสู่โอกาสใหม่ ๆ บางครั้งขึ้น บางครั้งลง จงเตรียมใจรับมือกับความเปลี่ยนแปลง และเปิดรับสิ่งใหม่ที่จักรวาลจะมอบให้

Let's your Astro Shine !`

	// Act
	cleanedMessage, card, meaning := ExtractJSONFromMessage(message)

	// Assert
	assert.Equal(t, expectedMessage, cleanedMessage)
	assert.Equal(t, "THE_WHEEL_OF_FORTUNE", card)
	assert.Equal(t, "Fate, cycles, and unexpected changes", meaning)
}

func TestExtractJSONFromMessage_WithoutJSON(t *testing.T) {
	// Arrange
	message := "This is a simple message without JSON"

	// Act
	cleanedMessage, card, meaning := ExtractJSONFromMessage(message)

	// Assert
	assert.Equal(t, message, cleanedMessage)
	assert.Equal(t, "", card)
	assert.Equal(t, "", meaning)
}

func TestExtractJSONFromMessage_WithInvalidJSON(t *testing.T) {
	// Arrange
	message := "Some text\n```json\n{invalid json}\n```\nMore text"

	// Act
	cleanedMessage, card, meaning := ExtractJSONFromMessage(message)

	// Assert
	assert.Equal(t, message, cleanedMessage) // Should return original if parsing fails
	assert.Equal(t, "", card)
	assert.Equal(t, "", meaning)
}

func TestExtractJSONFromMessage_WithEmptyJSON(t *testing.T) {
	// Arrange
	message := "Some text\n```json\n{}\n```\nMore text"

	// Act
	cleanedMessage, card, meaning := ExtractJSONFromMessage(message)

	// Assert
	assert.NotEqual(t, message, cleanedMessage) // JSON block should be removed
	assert.Equal(t, "", card)                   // Empty values
	assert.Equal(t, "", meaning)
}

func TestExtractJSONFromMessage_WithMultilineMessage(t *testing.T) {
	// Arrange
	message := `Line 1
Line 2
Line 3

` + "```json\n" + `{
"card": "TEST_CARD",
"meaning": "Test meaning"
}
` + "```" + `

Line 4
Line 5`

	// Act
	cleanedMessage, card, meaning := ExtractJSONFromMessage(message)

	// Assert
	assert.NotContains(t, cleanedMessage, "```json")
	assert.NotContains(t, cleanedMessage, "TEST_CARD")
	assert.Equal(t, "TEST_CARD", card)
	assert.Equal(t, "Test meaning", meaning)
	assert.Contains(t, cleanedMessage, "Line 1")
	assert.Contains(t, cleanedMessage, "Line 5")
}

func TestExtractJSONFromMessage_WithMalformedJSONMissingBraces(t *testing.T) {
	// Arrange
	message := "ท่านผู้ใฝ่หาคำทำนายจากจักรวาล ข้าจะหยิบไพ่จากห้วงดาราให้ท่าน...\n\nไพ่ที่ปรากฏคือ  \n**THE_EMPRESS**  \nนี่คือสัญลักษณ์แห่งความอุดมสมบูรณ์ การดูแลเอาใจใส่ และความคิดสร้างสรรค์ หากท่านกำลังเริ่มต้นสิ่งใหม่หรือหวังผลลัพธ์ที่ดี จงเชื่อมั่นว่าพลังแห่งการให้และความรักจะนำพาความสำเร็จมาสู่ท่าน\n\n, Let's your Astro Shine !\n\n```json\n\"card\": \"THE_EMPRESS\",\n\"meaning\": \"Abundance, nurturing, and creativity\"\n```"

	expectedMessage := "ท่านผู้ใฝ่หาคำทำนายจากจักรวาล ข้าจะหยิบไพ่จากห้วงดาราให้ท่าน...\n\nไพ่ที่ปรากฏคือ  \n**THE_EMPRESS**  \nนี่คือสัญลักษณ์แห่งความอุดมสมบูรณ์ การดูแลเอาใจใส่ และความคิดสร้างสรรค์ หากท่านกำลังเริ่มต้นสิ่งใหม่หรือหวังผลลัพธ์ที่ดี จงเชื่อมั่นว่าพลังแห่งการให้และความรักจะนำพาความสำเร็จมาสู่ท่าน\n\n, Let's your Astro Shine !"

	// Act
	cleanedMessage, card, meaning := ExtractJSONFromMessage(message)

	// Assert
	assert.Equal(t, expectedMessage, cleanedMessage)
	assert.Equal(t, "THE_EMPRESS", card)
	assert.Equal(t, "Abundance, nurturing, and creativity", meaning)
	assert.NotContains(t, cleanedMessage, "```json")
}
