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

func UnMarshallMessage(encoded []byte) (*Message, error) {
	header := UnmarshalHeader(encoded)
	question, _, err := UnMarshallQuestion(encoded[HeaderSize:])

	if err != nil {
		return nil, err
	}

	var answer Answer

	header.ANCount = 1
	header.QDCount = 1
	header.QR = true

	if header.OpCode == 0 {
		header.RCode = 0
	} else {
		header.RCode = 4
	}

	return &Message{
		Header:   *header,
		Question: *question,
		Answer:   answer,
	}, nil
}

func (m *Message) FormAnswer() *Answer {
	return &Answer{
		Name:  m.Question.Name,
		Type:  1,
		Class: 1,
		TTL:   60,
	}
}
