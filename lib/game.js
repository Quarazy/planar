var App = App || {};

var points = 0;

// Character variables
var sphereSize = 1.0;
var radius = Math.cbrt(sphereSize * 3/4 / Math.PI);
var speed = 0.1;

var characterGeometry = new THREE.SphereGeometry(1, 32, 32);
var characterMaterial = new THREE.MeshBasicMaterial({color: 0x825900});
var clock = new THREE.Clock();

// Test Game with Spheres
App.Game = function (conn, width, height) {
  this.conn = conn
  this.step = 0
  this.playerId = -1
  this.players = {}

  this.init = function () {
    this.scene = new THREE.Scene();
    this.camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 1000);
    this.camera.position.z = 5;

    this.controls = new App.Controls(document);

    this.renderer = new THREE.WebGLRenderer({alpha: true});
    this.renderer.setSize(width, height);
    this.renderer.setClearColor(0xffffff, 0);

    // Add renderer to document body
    document.body.appendChild(this.renderer.domElement);
  };

  this.start = function () {
    this.render();
  };

  this.onmessage = function(msg) {
    if (msg.type == 1) {
      this.playerId = msg.id;

      // Add players
      msg.players.forEach(function (playerId) {
        var player = new App.Player({});
        var playerObj = new THREE.Mesh(characterGeometry, characterMaterial);
        console.log("Adding ", playerId);
        this.scene.add(playerObj);

        this.players[playerId] = {
          player: player,
          obj: playerObj
        };
      }, this);

    } else if (msg.type == 0) {
      // Got response to make a move
      if (this.players[msg.id] != undefined) {
        this.players[msg.id]['player'].correct(msg);
      }
    } else if (msg.type == 2) {
      var player = new App.Player({});
      var playerObj = new THREE.Mesh(characterGeometry, characterMaterial);
      console.log("Somebody just joined ", msg.id);
      this.scene.add(playerObj);

      this.players[msg.id] = {
        player: player,
        obj: playerObj
      };
    }
  };

  // Game render loop
  this.render = function() {
    this.update();
	  this.renderer.render(this.scene, this.camera);
  };

  this.update = function () {
    var accel = this.controls.getAccel()
    if (this.playerId >= 0) {
      var delta = clock.getDelta();

      send(this.conn,
           {step: this.step,
            id: this.playerId,
            msec: delta,
            x: accel.x,
            y: accel.y,
            z: 0,
            buttons: this.controls.getButtons()});

      var player;
      var pos;
      var now = Date.now() - (2000 / tickRate);
      for (var playerId in this.players) {
        player = this.players[playerId];

        if (this.playerId == playerId) {
          player['player'].update(accel, delta);
          pos = player['player'].getPosition(this.step);
        } else {
          pos = player['player'].interp(now);
        }

        player['obj'].position.set(pos.x, pos.y, pos.z);
      }

      this.step += 1;
    };
  };
};


// TODO(quarazy): Need another setTimeout to send back user commands. This is
// separate from the rendering loop
var tickRate = 25;
var intervalID = window.setInterval(emitInput, 1000 / tickRate);

function emitInput() {
  // TODO(quarazy): Send something back
}

// animate starts the render loop
function animate() {
  requestAnimationFrame(animate);

	game.render();
}

var conn = connect();
var game = new App.Game(conn, window.innerWidth, window.innerHeight);

conn.onopen = function (e) {
  game.init();
  animate()
};

// Receives direction from the server, so now just need to apply this direction
conn.onmessage = function (e) {
  var message = JSON.parse(e.data);

  game.onmessage(message)
};
