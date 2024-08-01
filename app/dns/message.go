package dns

type Message struct {
	Header    Header
	Questions []Question
	Answers   []Answer
}

// Marshal encodes the DNS Message into a byte slice.
// It marshals the Header, Question, and Answer fields of the Message
// and concatenates their byte representations into a single byte slice.
//
// Returns:
// - A byte slice containing the encoded DNS Message.
func (m *Message) Marshal() []byte {
	var encoded []byte
	encodedHeader := m.Header.Marshal()
	encoded = append(encoded, encodedHeader...)

	for _, quest := range m.Questions {
		encodedQuestion, _ := quest.Marshal()
		encoded = append(encoded, encodedQuestion...)
	}

	for _, ans := range m.Answers {
		encodedAnswer := ans.Marshal()
		encoded = append(encoded, encodedAnswer...)
	}

	return encoded
}

// UnMarshallMessage decodes a byte slice into a DNS Message struct.
// It unmarshal the Header and Question fields from the encoded byte slice
// and sets default values for the Answer field and some Header fields.
//
// Parameters:
// - encoded: A byte slice containing the encoded DNS message.
//
// Returns:
// - A pointer to a Message struct populated with the decoded values.
// - An error if any issue occurs during decoding.
func UnMarshallMessage(encoded []byte) (*Message, error) {
	header := UnmarshalHeader(encoded)
	questions, err := UnmarshalQuestions(encoded, header.QDCount)

	if err != nil {
		return nil, err
	}

	answers := make([]Answer, 0, len(questions))
	for _, quest := range questions {
		answer := FromQuestion(quest)
		answers = append(answers, answer)
	}

	header.QR = true
	header.ANCount = uint16(len(answers))

	if header.OpCode == 0 {
		header.RCode = 0
	} else {
		header.RCode = 4
	}

	return &Message{
		Header:    *header,
		Questions: questions,
		Answers:   answers,
	}, nil
}
