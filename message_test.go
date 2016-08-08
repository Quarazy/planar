package main

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const Epsilon = 1e-6

var _ = Describe("Message", func() {
	var message *Message
	var byt []byte

	BeforeEach(func() {
		message = &Message{}
		byt = []byte(`{"type":1,"id":13,"msec":16.66,"x":1.3,"y":-2.2,"z":0,"step":56,"buttons":0}`)
	})

	It("Serializes", func() {
		if err := json.Unmarshal(byt, message); err != nil {
			Fail(fmt.Sprintf("Error when unmarshalling: %s", err))
		}
		Expect(message.Type).To(Equal(User))
		Expect(message.X).To(BeNumerically("~", 1.3, Epsilon))
		Expect(message.Y).To(BeNumerically("~", -2.2, Epsilon))
		Expect(message.Z).To(BeNumerically("==", 0))
		Expect(message.Step).To(BeNumerically("==", 56))
	})

	It("Deserializes", func() {
		message = &Message{
			Type: 1,
			Id:   13,
			Msec: 16.66,
			X:    1.3,
			Y:    -2.2,
			Z:    0,
			Step: 56,
		}
		des, _ := json.Marshal(message)
		Expect(des).To(Equal(byt))
	})
})
