var MockClock = function(delta) {
  this.getDelta = function () {
    return delta;
  }
};

var MockConn = function() {};

var MockPlayer = function() {
  this.position = {}
};

describe("Player", function () {
  var player;

  beforeEach(function () {
    player = new App.Player(new MockConn(), {x: 0, y: 0, z: 0});
  });

  describe("#update", function () {
    it("saves moves state", function () {
      player.update({x: 1, y: 2, z: 0}, 16);
      expect(player.getPosition(0)).toEqual({x: 16, y: 32, z: 0});
    });
  });

  describe("#correct", function () {
    it("corrects on new getPosition", function () {
      player.update({x: 1, y: 2, z: 0}, 16);
      player.correct({step: 0, x: -5, y: -3, z: 0});
      expect(player.getPosition(50)).toEqual({x: 11, y: 29, z: 0});
    });
  });
});
