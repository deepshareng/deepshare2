module.exports = function (grunt) {
    grunt.initConfig({

    // define source files and their destinations
    uglify: {
        files: {
            src: 'deepshare-redirect.js',  // source files mask
            dest: 'fileserver/',    // destination folder
            expand: true,    // allow dynamic building
            flatten: true,   // remove all unnecessary nesting
            ext: '.min.js'   // replace .js to .min.js
        }
    },
    watch: {
        js:  { files: 'deepshare-redirect.js', tasks: [ 'uglify' ] },
    },
});

// load plugins
grunt.loadNpmTasks('grunt-contrib-watch');
grunt.loadNpmTasks('grunt-contrib-uglify');

// The task to add time stamp to html
grunt.registerTask('timestamp', 'Add timestamp to html', function() {
    var tpl = grunt.file.read(
        './tpl/sharelink_response_mobile.html',
        {encoding: 'utf8'}
    );
    var rendered = grunt.template.process(
        tpl, {
            data: {
                timeStamp: grunt.template.today('yyyymmddHHMMssl')
            }
    });
    grunt.file.write(
        './sharelink_response_mobile.html',
        rendered,
        {encoding: 'utf8'}
    );
});

// register at least this one task
grunt.registerTask('default', ['uglify', 'timestamp']);

};


