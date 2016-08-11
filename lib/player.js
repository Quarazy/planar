// sleep delays an action. This is just used to simulate latency right now
function sleep (time) {
  return new Promise((resolve) => setTimeout(resolve, time));
}

// send passes user action back to the server
function send(conn, data) {
  sleep(1).then(() => {
    conn.send(JSON.stringify(data));
  });
}

var App = App || {};
const velocity = 4;

App.Player = function (position) {
  this.cbuf = new AppUtil.CircularBuffer(1000)
  this.confirmedBuffer = new AppUtil.CircularBuffer(30)

  this.position = {
    x: position.x || 0,
    y: position.y || 0,
    z: position.z || 0
  };
  this.lastConfirmed = {
    idx: -1,
    x: position.x || 0,
    y: position.y || 0,
    z: position.z || 0
  };
  this.lastApplied = -1;

  this.clientIdx = 0
};

// correct appends to a buffer of a remote peer's position
App.Player.prototype.correct = function(msg) {
  this.lastConfirmed = {
    idx: msg.step,
    msec: msg.msec,
    x: msg.x,
    y: msg.y,
    z: msg.z
  };

  this.confirmedBuffer.push(this.lastConfirmed)
  this.clientIdx += 1
};

// Adds new move to buffer
App.Player.prototype.update = function(move, delta) {
  var move = {
    x: move.x*velocity*delta || 0,
    y: move.y*velocity*delta || 0,
    z: move.z*velocity*delta || 0
  };
  this.cbuf.push(move);

  this.position.x += move.x;
  this.position.y += move.y;
  this.position.z += move.z;
};

// Gets the position of the player
// This recalculates from the last confirmed position
App.Player.prototype.getPosition = function(step) {
  // FIXME(quarazy): This isn't smooth, mostly likely due to
  // an index offset
  if (this.lastApplied < this.lastConfirmed.idx) {
    var position = {x: this.lastConfirmed.x,
                    y: this.lastConfirmed.y,
                    z: this.lastConfirmed.z};

    var move;
    for(var i = this.lastConfirmed.idx; i <= step; i++) {
      move = this.cbuf.get([i]);

      if (move !== undefined) {
        position.x += move.x;
        position.y += move.y;
        position.z += move.z;
      }
    }

    this.lastApplied = this.lastConfirmed.idx;
    this.position = position;
  }

  return this.position;
};

// interp using linear interpolation to estimate the position of an object
// it rounds from nearby snapshots using the current client time.
App.Player.prototype.interp = function(time) {
  var step1 = this.confirmedBuffer.get(this.clientIdx - 2);
  var step2 = this.confirmedBuffer.get(this.clientIdx - 3);

  var slope = (time-step1.msec)/(step2.msec - step1.msec);

  return {
    x: slope*(step2.x - step1.x) + step1.x,
    y: slope*(step2.y - step1.y) + step1.y,
    z: 0
  };
};
