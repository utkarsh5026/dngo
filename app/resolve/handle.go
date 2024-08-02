package resolve

import (
	"context"
	"fmt"
	"github.com/codecrafters-io/dns-server-starter-go/app/dns"
	"net"
)

func HandleDnsResolution(dnsQuery []byte, resolver *net.Resolver) (*dns.Message, error) {
	header := dns.UnmarshalHeader(dnsQuery)
	questions, err := dns.UnmarshalQuestions(dnsQuery, header.QDCount)

	if err != nil {
		return nil, err
	}

	answers := make([]dns.Answer, 0, len(questions))
	for _, quest := range questions {
		fmt.Println("Resolving", quest.Name)
		answer, err := resolveQuestion(quest, resolver)

		if err != nil {
			continue
		}
		answers = append(answers, answer...)
	}

	header.QR = true
	header.ANCount = uint16(len(answers))
	header.QDCount = uint16(len(questions))
	header.NSCount = 0
	header.ARCount = 0

	if header.OpCode != 0 {
		header.RCode = 4
	}
	return &dns.Message{
		Header:    *header,
		Questions: questions,
		Answers:   answers,
	}, nil
}

func resolveQuestion(question dns.Question, resolver *net.Resolver) ([]dns.Answer, error) {
	answer := dns.Answer{
		Name:     question.Name,
		Type:     1,
		Class:    1,
		TTL:      60,
		RDLength: 4,
		RData:    []byte("\x08\x08\x08\x08"),
	}

	if resolver != nil {
		ips, err := resolver.LookupIPAddr(context.Background(), question.Name)
		if err != nil {
			return []dns.Answer{answer}, err
		}

		answers := make([]dns.Answer, 0, len(ips))
		for _, ip := range ips {
			fmt.Println("Resolved", question.Name, "to", ip.IP)
			newAns := answer
			newAns.RData = ip.IP
			answers = append(answers, newAns)
		}
		return answers, nil
	}

	return []dns.Answer{answer}, nil
}
