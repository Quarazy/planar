var concat       = require('gulp-concat'),
    gulp         = require('gulp'),
    livereload   = require('gulp-livereload'),
    runSequence  = require('run-sequence'),
    stylus       = require('gulp-stylus');

var PATHS = {
  stylesheets : 'assets/stylus/**/*.styl',
  sources : 'lib/**/*.js',
  outDir : 'public'
}

function stylesheets(paths, outfile) {
  var outfile = outfile + '.css'
  return gulp.src(paths)
    .pipe(stylus())
    .pipe(concat(outfile))
    .pipe(gulp.dest('./'))
    .pipe(livereload());
}

function sources(paths, outfile) {
  var outfile = outfile + '.js'
  return gulp.src(paths)
    .pipe(concat(outfile))
    .pipe(gulp.dest('./'))
    .pipe(livereload());
}

gulp.task('default', ['development'])

gulp.task('development', function (callback) {
  runSequence(
    'build'
    , 'watch'
  )
})

gulp.task('build', [
  'stylesheets',
  'sources'
])

gulp.task('watch', function () {
  livereload.listen();
  gulp.watch(PATHS.stylesheets, ['stylesheets']);
  gulp.watch(PATHS.sources, ['sources']);
})

gulp.task('stylesheets', function () {
  return stylesheets(PATHS.stylesheets, PATHS.outDir + '/style')
})

gulp.task('sources', function () {
  return sources(PATHS.sources, PATHS.outDir + '/app')
})
