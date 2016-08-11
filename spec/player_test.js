var MockClock = function(delta) {
  this.getDelta = function () {
    return delta;
  };
};

var MockConn = function() {};

var MockPlayer = function() {
  this.position = {};
};

describe("Player", function () {
  var player;

  beforeEach(function () {
    player = new App.Player({x: 0, y: 0, z: 0});
  });

  describe("#update", function () {
    it("saves moves state", function () {
      player.update({x: 1, y: 2, z: 0}, 16);
      expect(player.getPosition(0)).toEqual({x: 64, y: 128, z: 0});
    });
  });

  describe("#correct", function () {
    it("corrects on new getPosition", function () {
      player.update({x: 1, y: 2, z: 0}, 16);
      player.correct({step: 0, x: -5, y: -3, z: 0});
      expect(player.getPosition(50)).toEqual({x: 59, y: 125, z: 0});
    });
  });
});
