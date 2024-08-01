package dns

type Message struct {
	Header   Header
	Question Question
	Answer   Answer
}

func (m *Message) Marshal() []byte {
	encodedHeader := m.Header.Marshal()
	encodedQuestion, _ := m.Question.Marshal()
	encodedAnswer := m.Answer.Marshal()

	byteCount := len(encodedHeader) + len(encodedQuestion) + len(encodedAnswer)
	encoded := make([]byte, byteCount)
	copy(encoded, encodedHeader)
	copy(encoded[len(encodedHeader):], encodedQuestion)
	copy(encoded[len(encodedHeader)+len(encodedQuestion):], encodedAnswer)

	return encoded
}
func CreateDnsMessage(header *Header, question *Question, answer *Answer) *Message {
	return &Message{
		Header:   *header,
		Question: *question,
		Answer:   *answer,
	}
}
