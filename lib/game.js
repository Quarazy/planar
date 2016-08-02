var width = window.innerWidth;
var height = window.innerHeight;

var behind = new THREE.Vector3(0, -5, 5);

var mouse = {x: 0, y: 0};
var points = 0;

// Character variables
var sphereSize = 1.0;
var radius = Math.cbrt(sphereSize * 3/4 / Math.PI);
var speed = 0.1;

var position = new THREE.Vector3(0, 0, 0);
var coins = [];

var scene, renderer, character, camera;
var coinGeometry = new THREE.SphereGeometry( .5, 32, 32 );
var coinMaterial = new THREE.MeshBasicMaterial( { color: 0x00ff00 } );

function init() {
  scene = new THREE.Scene();
  camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 1000);

  // Adding the character
  var characterGeometry = new THREE.SphereGeometry(1, 32, 32);
  var characterMaterial = new THREE.MeshBasicMaterial({color: 0x825900});

  character = new THREE.Mesh(characterGeometry, characterMaterial);
  scene.add(character);

  renderer = new THREE.WebGLRenderer({ alpha: true });
  renderer.setSize( width, height );
  renderer.setClearColor( 0xffffff, 0);
  document.body.appendChild(renderer.domElement);

  // The ground
  var loader = new THREE.TextureLoader();
  loader.load("images/grass.jpg",
              function (texture) {
                texture.wrapS = THREE.RepeatWrapping;
                texture.wrapT = THREE.RepeatWrapping;
                texture.repeat.x = 256;
                texture.repeat.y = 256;
                var groundMat = new THREE.MeshBasicMaterial({map:texture});

                var groundGeo = new THREE.PlaneGeometry(400,400);
                var ground = new THREE.Mesh(groundGeo, groundMat);
                scene.add(ground);
              },
              function (xhr) {},
              function (xhr) {});

  // New map initialization of coins
  var x, y, z;
  for(var i = 0; i < 15; i++) {
    x = (Math.random() * 10) - 5
    y = (Math.random() * 10) - 5

    var coin = new THREE.Mesh(coinGeometry, coinMaterial );
    scene.add(coin);
    coin.position.set(x, y, 0);

    coins.push(coin);
  }

  camera.position.add(behind);
  camera.lookAt(new THREE.Vector3(0, 0, 0));

  document.addEventListener('mousemove', function (e) {
	  e.preventDefault();
	  mouse.x = (e.clientX / window.innerWidth) * 2 - 1;
	  mouse.y = - (e.clientY / window.innerHeight) * 2 + 1;

    // Make the sphere follow the mouse
    var vector = new THREE.Vector3(mouse.x, mouse.y, 0.5);
	  vector.unproject( camera );
	  var dir = vector.sub(camera.position).normalize();
	  var distance = - camera.position.z / dir.z;

	  position = camera.position.clone().add( dir.multiplyScalar( distance ) );
  }, false);

  render();
}

var conn = connect()
conn.onopen = function (e) {
  init();
}

// New map coin generation
function newCoins() {
  if (Math.random() > 0.98) {
    x = (Math.random() * 40) - 5
    y = (Math.random() * 40) - 5

    var coin = new THREE.Mesh(coinGeometry, coinMaterial);
    scene.add(coin);
    coin.position.set(x, y, 0);

    coins.push(coin);
  }
}

// Checks if there are any collisions between main Sphere and all others.
// Removes collided objects and increases point total for collisions
function checkCollisions() {
  var partitioned = _.partition(coins, function (coin) {
    return character.position.distanceTo(coin.position) < (radius + 0.5);
  });

  _.forEach(partitioned[0], function (coin) {
    scene.remove(coin);

    sphereSize += 0.1;
    radius = Math.cbrt(sphereSize * 3/4 / Math.PI);

    character.scale.set(radius, radius, radius);
  });

  coins = partitioned[1];
  points += partitioned[0].length;
  document.getElementById("score").innerHTML = points;
}

function animate() {
  var direction = position.clone().sub(character.position).normalize();

  var multiplier = -Math.log10(sphereSize/10+3)+.7;
  var newDir = direction.multiplyScalar(multiplier)
  character.position.add(newDir);

  camera.position.add(newDir);
  checkCollisions();
  newCoins();
}

// Game render Loop
function render() {
	requestAnimationFrame(render);

  animate();
	renderer.render(scene, camera);
}
