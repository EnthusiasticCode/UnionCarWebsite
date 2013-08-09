'use strict';

module.exports = function (grunt) {
  // load all grunt tasks
  require('matchdep').filterDev('grunt-*').forEach(grunt.loadNpmTasks);

  // configurable paths
  var yeomanConfig = {
    app: 'app',
    dist: 'dist'
  };

  try {
    yeomanConfig.app = require('./bower.json').appPath || yeomanConfig.app;
  } catch (e) {}

  grunt.initConfig({
    yeoman: yeomanConfig,
    watch: {
      coffee: {
        files: ['<%= yeoman.app %>/scripts/{,*/}*.coffee'],
        tasks: ['coffee:dist']
      },
      compass: {
        files: [
          '<%= yeoman.app %>/styles/{,*/}*.{scss,sass}',
          '<%= yeoman.app %>/sprites/{,*/}*.png'
        ],
        tasks: ['compass']
      },
      views: {
        files: ['<%= yeoman.app %>/{views,controllers,models}/*.{php,html}'],
        tasks: ['copy:develop']
      },
      styles: {
        files: ['{.tmp,<%= yeoman.app %>}/styles/{,*/}*.css'],
        tasks: ['cssmin']
      },
      scripts: {
        files: ['{.tmp,<%= yeoman.app %>}/scripts/{,*/}*.js'],
        tasks: ['concat']
      },
      images: {
        files: ['{.tmp,<%= yeoman.app %>}/images/{,*/}*.{png,jpg,jpeg,gif,webp,svg}'],
        tasks: ['imagemin']
      }
    },
    clean: {
      dist: {
        files: [{
          dot: true,
          src: [
            '.tmp',
            '<%= yeoman.dist %>/*',
            '!<%= yeoman.dist %>/.git*'
          ]
        }]
      },
      develop: '.tmp'
    },
    karma: {
      unit: {
        configFile: 'karma.conf.js',
        singleRun: true
      }
    },
    coffee: {
      dist: {
        files: [{
          expand: true,
          cwd: '<%= yeoman.app %>/scripts',
          src: '{,*/}*.coffee',
          dest: '.tmp/scripts',
          ext: '.js'
        }]
      },
    },
    compass: {
      options: {
        sassDir: '<%= yeoman.app %>/styles',
        cssDir: '.tmp/styles',
        imagesDir: '<%= yeoman.app %>/sprites',
        javascriptsDir: '<%= yeoman.app %>/scripts',
        fontsDir: '<%= yeoman.app %>/styles/fonts',
        importPath: '<%= yeoman.app %>/components',
        raw: 'http_images_path = "./sprites/"\ngenerated_images_dir = ".tmp/images"\nhttp_generated_images_path = "../images/"',
        relativeAssets: false
      },
      dist: {},
      develop: {
        options: {
          debugInfo: true
        }
      }
    },
    concat: {
      dist: {
        files: {
          '<%= yeoman.dist %>/scripts/scripts.js': [
            '<%= yeoman.app %>/scripts/angular/*.js',
            '<%= yeoman.app %>/scripts/angular/modules/*.js',

            '.tmp/scripts/{,*/}*.js'
          ]
        }
      }
    },
    imagemin: {
      dist: {
        files: [{
          expand: true,
          cwd: '<%= yeoman.app %>/images',
          src: '*.{png,jpg,jpeg}',
          dest: '<%= yeoman.dist %>/images'
        }, {
          expand: true,
          cwd: '.tmp/images',
          src: '{,*/}*.png',
          dest: '<%= yeoman.dist %>/images'
        }]
      }
    },
    cssmin: {
      dist: {
        files: {
          '<%= yeoman.dist %>/styles/main.css': [
            '.tmp/styles/{,*/}*.css',
            '<%= yeoman.app %>/styles/{,*/}*.css'
          ]
        }
      }
    },
    ngmin: {
      dist: {
        files: [{
          expand: true,
          cwd: '<%= yeoman.dist %>/scripts',
          src: '*.js',
          dest: '<%= yeoman.dist %>/scripts'
        }]
      }
    },
    uglify: {
      dist: {
        files: {
          '<%= yeoman.dist %>/scripts/scripts.js': [
            '<%= yeoman.dist %>/scripts/scripts.js'
          ]
        }
      }
    },
    copy: {
      dist: {
        files: [{
          expand: true,
          dot: true,
          cwd: '<%= yeoman.app %>',
          dest: '<%= yeoman.dist %>',
          src: [
            'views/*.{php,html}',
            'images/{,*/}*.{gif,webp}',
            'config/*',
            'controllers/*',
            'errors/*',
            'core/*',
            'models/*'
          ]
        }]
      },
      develop: {
        files: [{
          expand: true,
          cwd: '<%= yeoman.app %>',
          dest: '<%= yeoman.dist %>',
          src: [
            'views/*.{php,html}',
            'controllers/*',
            'models/*'
          ]
        }]
      }
    }
  });

  grunt.renameTask('regarde', 'watch');

  grunt.registerTask('develop', [
    'clean:dist',
    'coffee:dist',
    'compass:develop',
    'imagemin',
    'cssmin',
    'concat',
    'copy',
    'watch'
  ]);

  grunt.registerTask('build', [
    'clean:dist',
    'coffee',
    'compass:dist',
    'imagemin',
    'cssmin',
    'concat',
    'copy',
    'ngmin',
    // 'uglify'
  ]);

  grunt.registerTask('default', ['build']);
};
