var AppUtil = AppUtil || {};

// Very simple Circular buffer.
// Supports constant pushes. When sizes is over the limit
// starts writing from the beginning.
AppUtil.CircularBuffer = function (size) {
  this.size = size;
  this.data = Array.apply(null, Array(size)).map(() => {});
  this.idx  = 0;
};

AppUtil.CircularBuffer.prototype = {
  get: function (i) {
    return this.data[i % this.size];
  },
  // This should be safe due to JS's concurrency model
  // https://developer.mozilla.org/en-US/docs/Web/JavaScript/EventLoop
  push: function (v) {
    this.data[this.idx % this.size] = v;
    this.idx = this.idx + 1 % this.size;
  }
};
