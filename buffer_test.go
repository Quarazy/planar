package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Buffer", func() {
	var buf *Buffer

	BeforeEach(func() {
		buf = NewBuffer()
	})

	Describe("Append", func() {
		It("Adds new node", func() {
			buf.Append(3)
			Expect(buf.tail.value.(int)).To(Equal(3))

			val, _ := buf.PopLeft()
			Expect(val.(int)).To(Equal(3))
		})
	})

	Describe("PopLeft", func() {
		BeforeEach(func() {
			buf.Append(3)
			buf.Append(4)
		})

		It("Returns the first added", func() {
			val, _ := buf.PopLeft()
			Expect(val.(int)).To(Equal(3))
		})
		It("Returns the second added on second pop", func() {
			buf.PopLeft()
			val, _ := buf.PopLeft()
			Expect(val.(int)).To(Equal(4))
		})
		It("Errors if popped on an empty buffer", func() {
			buf.PopLeft()
			buf.PopLeft()
			val, err := buf.PopLeft()
			Expect(val).To(BeNil())
			Expect(err.Error()).To(Equal("Pop on an empty buffer"))
		})
	})

	Describe("PopAll", func() {
		BeforeEach(func() {
			buf.Append(3)
			buf.Append(4)
		})

		It("Returns everything", func() {
			var elems, elemsOne, elemsTwo []interface{}
			done := make(chan struct{})

			go func() {
				elemsOne = buf.PopAll()
				done <- struct{}{}
			}()

			go func() {
				elemsTwo = buf.PopAll()
				done <- struct{}{}
			}()

			<-done
			<-done

			elems = append(elemsOne, elemsTwo...)
			Expect(len(elems)).To(Equal(2))
			Expect(elems).To(ConsistOf(3, 4))
		})

		Measure("is fast for a huge queue of tasks", func(b Benchmarker) {
			for i := 0; i < 50000; i++ {
				buf.Append(Command{
					X: 1.5,
					Y: -2.2,
					Z: 0,
				})
			}

			runtime := b.Time("runtime", func() {
				position := NewVector3D(0, 0, 0)

				for _, v := range buf.PopAll() {
					com, _ := v.(Command)

					position.X += com.X
					position.Y += com.Y
					position.Z += com.Z
				}

				Expect(position.X).To(BeNumerically("==", 75000))
				Expect(position.Y).To(BeNumerically("~", -110000, Epsilon))
			})

			Expect(runtime.Seconds()).To(BeNumerically("<", 0.66))
		}, 10)
	})
})
