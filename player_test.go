package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Player", func() {
	var p *player

	BeforeEach(func() {
		p = &player{
			Velocity: 10,
		}
	})

	var _ = Describe(".Update", func() {
		It("Modifies the player's position", func() {
			c := Command{
				Code:  Move,
				X:     1.3,
				Y:     -2.2,
				Z:     0,
				Delta: 0.1,
			}

			p.Update(c, 0.1)

			CompareVectors(Vector3D(p.Location), NewVector3D(1.3, -2.2, 0))
		})
	})
})

func CompareVectors(v1, v2 Vector3D) {
	Expect(v1.X).To(BeNumerically("~", v2.X, Epsilon))
	Expect(v1.Y).To(BeNumerically("~", v2.Y, Epsilon))
	Expect(v1.Z).To(BeNumerically("~", v2.Z, Epsilon))
}
