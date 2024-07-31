package dns

type Message struct {
	Header   Header
	Question Question
}

func (m *Message) Marshal() []byte {
	encodedHeader := m.Header.Marshal()
	encodedQuestion, _ := m.Question.Marshal()

	encoded := make([]byte, len(encodedHeader)+len(encodedQuestion))
	copy(encoded, encodedHeader)
	copy(encoded[len(encodedHeader):], encodedQuestion)

	return encoded
}
