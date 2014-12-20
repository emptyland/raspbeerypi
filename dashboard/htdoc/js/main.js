//'use strict';

(function () {

var app = angular.module('dashboard', ['ngCookies']);

app.controller('NavController', function($scope, $cookies) {
    var navBars = {
        'dashboard': {}, 'console': {}, 'devices': {},

        current: 'dashboard'
    };

    $scope.switchNavBar = function (index) {

        for (var k in navBars) {
            navBars[k].active = '';
        }
        navBars[index].active = 'active';
        navBars.current = index;
        $cookies['nav-tab'] = index;
        $scope.navBars = navBars;
    };

    var initTab = $cookies['nav-tab'];
    if (typeof initTab != 'string') {
        initTab = 'dashboard';
    }

    $scope.switchNavBar(initTab);
});

app.controller('StateController', function($scope, $http) {
    $scope.state = {
        loadAvg: ["N/A", "N/A", "N/A"],
        cpuTemperature: "N/A",
        cpuPercent: "N/A"
    };

    var update = function () {
        $http.get('api/state').success(function(data) {
            $scope.state = data;
        });
    };

    var timer = setInterval(function () {
        $scope.$apply(update);
    }, 5000);

    update();
});

app.controller('MemoryController', function ($scope, $http) {
    $scope.memory = {
        total: "N/A",
        used: "N/A",
        swapUsed: "N/A"
    };

    var update = function () {
        $http.get('api/memory').success(function(data) {
            $scope.memory = data;
        });
    }

    var timer = setInterval(function () {
        $scope.$apply(update);
    }, 5000);

    update();
});

app.controller('DiskUsageController', function ($scope, $http) {
    $scope.usageEntries = [];

    $http.get('api/disk').success(function(data) {
        $scope.usageEntries = data.entries;

        for (var i in $scope.usageEntries) {
            var usage = $scope.usageEntries[i];

            if (usage.total == 0) {
                usage.usedPercent = 0;
            } else {
                usage.usedPercent = (usage.used / usage.total * 100).toFixed(2);
            }

            if (usage.usedPercent >= 0 && usage.usedPercent < 50) {
                usage.color = 'progress-bar-success';
            } else if (usage.usedPercent >= 50 && usage.usedPercent < 70) {
                usage.color = 'progress-bar-warning';
            } else {
                usage.color = 'progress-bar-danger';
            }
        }
    });
});

app.controller('JobController', function ($scope, $http) {

    var updateList = function () {
        $http.get('api/job/list').success(function(data) {
            jobs = {}

            for (var i in data.entries) {
                var entry = data.entries[i];
                jobs[entry.id] = entry;
            }

            $scope.jobs = jobs;
        });
    };
    updateList();


    $scope.isEditing = false;

    $scope.onEdit = function (job) {
        $scope.metadata = job;

        $scope.onLang(job.lang);
        $scope.editor.setValue(job.code);
        $scope.isEditing = true;
    };

    $scope.onEnable = function (index) {
        console.log('onEnable');
        console.log(index);
    }

    $scope.onDelete = function (index) {
        console.log('onDelete');
        console.log(index);
    }

    var langs = {
        'bash': {
            name: 'Bash',
            color: 'ace/mode/sh'
        },
        'python': {
            name: 'Python 2',
            color: 'ace/mode/python'
        },
        'node.js': {
            name: 'Node.js',
            color: 'ace/mode/javascript'
        }
    };

    $scope.langs = langs;

    var editor = ace.edit("nv-editor");
    editor.setTheme("ace/theme/monokai");
    editor.getSession().setMode("ace/mode/sh");

    $scope.editor = editor;

    $scope.lang = {
        current: 'bash'
    };

    $scope.metadata = {
        title: '',
        cron: ''
    };

    $scope.onLang = function (lang) {
        $scope.metadata.lang = lang;
        $scope.editor.getSession().setMode(langs[lang].color);
    };

    $scope.onSave = function (id) {
        $scope.metadata.code = $scope.editor.getValue();

        $http.post('/api/job/content',
            { entries: [ $scope.metadata ] }
        ).success(function (data) {
            $scope.jobs[id] = $scope.metadata;
            $scope.isEditing = false;

            console.log(data);
        });
    };

    $scope.onCancel = function (index) {
        $scope.isEditing = false;
    };

    $scope.onRun = function (index) {
        $http.post('/api/job/run',
            { entries: [ $scope.metadata ] }
        ).success(function (data) {
            console.log(data);
        }).error(function (data) {
            console.log(data);
        });
    }
});

}()); // end of module
