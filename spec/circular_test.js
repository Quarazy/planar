describe("CircularBuffer", function () {
  beforeEach(function() {
    this.subject = function() {
      var cbuf = new AppUtil.CircularBuffer(3);
      cbuf.push(1);
      cbuf.push(2);
      cbuf.push(3);

      return cbuf;
    };
  });

  describe("#get", function () {
    it("gets values", function () {
      expect(this.subject().get(0)).toBe(1);
      expect(this.subject().get(1)).toBe(2);
      expect(this.subject().get(2)).toBe(3);
    });

    it("gets values after wrap", function () {
      var cbuf = this.subject();

      for(var i = 4; i <= 10; i++) {
        cbuf.push(i);
      }
      expect(cbuf.get(7)).toBe(8);
      expect(cbuf.get(8)).toBe(9);
      expect(cbuf.get(9)).toBe(10);
    });
  });
});
