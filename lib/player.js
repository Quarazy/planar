// send passes user action back to the server
function send(conn, data) {
  conn.send(JSON.stringify(data));
}

var App = App || {};

App.Player = function (conn, position) {
  this.cbuf = new AppUtil.CircularBuffer(100)
  this.position = position;
  this.lastConfirmed = {idx: -1, x: 0, y: 0, z: 0};
  this.lastApplied = -1;

  conn.onmessage = function (e) {
    var message = JSON.parse(e.data);
  }
};

// Update updates the player position based on
// the message from the server
// player only cares about position messages
App.Player.prototype.correct = function(msg) {
  this.lastConfirmed = {
    idx: msg.step,
    x: msg.x,
    y: msg.y,
    z: msg.z
  }
};

// Adds new move to buffer
App.Player.prototype.update = function(move, delta) {
  var move = {x: move.x*delta, y: move.y*delta, z: move.z*delta};
  this.cbuf.push(move);

  this.position.x += move.x;
  this.position.y += move.y;
  this.position.z += move.z;
};

// Gets the position of the player
// This recalculates from the last confirmed position
App.Player.prototype.getPosition = function(step) {
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
    this.position = position
  }

  return this.position;
};

// Creates a THREE.Vector3 from object's x, y, z fields
function vec3FromOBJ(obj) {
  return new THREE.Vector3(obj.x, obj.y, obj.z);
}
