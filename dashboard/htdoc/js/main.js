//'use strict';

(function () {

var app = angular.module('dashboard', ['ngCookies']);

app.controller('NavController', function($scope, $cookies) {
    var navBars = {
        'dashboard': {}, 'console': {}, 'file': {}, 'devices': {},

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
    $scope.isRunning = false;

    $scope.onNew = function () {
        var job = {
            id: -1,
            title: '',
            desc: '',
            cron: '',
            user: 'root',
            lang: 'bash',
            enable: true
        };

        $scope.onEdit(job)
    };

    $scope.onEdit = function (job) {
        $scope.metadata = job;

        $scope.onLang(job.lang);
        $scope.editor.setValue(job.code);

        $scope.onEditing(true);
    };

    $scope.onEnable = function (job) {
        job.enable = !job.enable
        $http.post('/api/job/content',
            { entries: [ job ] }
        ).error(function (data) {
            job.enable = !job.enable; // rollback

            $scope.result = {
                category: 'danger',
                output: [data]
            };
        });
    }

    $scope.onDelete = function (job) {
        $http.post('/api/job/delete',
            { entries: [ job ] }
        ).success(function (data) {
            delete $scope.jobs[job.id];
        }).error(function (data) {

            $scope.result = {
                category: 'danger',
                output: [data]
            };
        });
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

        if (typeof $scope.metadata.title != 'string' ||
            $scope.metadata.title === '') {

            console.log($scope.metadata);
            $scope.result = {
                category: 'warning',
                output: ['title can not be empty.']
            };
            return;
        }

        if (typeof $scope.metadata.cron != 'string' ||
            $scope.metadata.cron === '') {

            $scope.result = {
                category: 'warning',
                output: ['cron can not be empty.']
            };
            return;
        }

        $http.post('/api/job/content',
            { entries: [ $scope.metadata ] }
        ).success(function (data) {
            if (id === -1) {
                updateList();
            } else {
                $scope.jobs[id] = $scope.metadata;
            }

            $scope.onEditing(false);
        });
    };

    $scope.onCancel = function (index) {
        $scope.onEditing(false);
    };

    $scope.onEditing = function (isEditing) {
        $scope.result = {
            category: 'none',
            msg: ''
        };
        $scope.isEditing = isEditing;
    };

    $scope.onRun = function (index) {
        $scope.metadata.code = $scope.editor.getValue();

        $scope.result = {
            category: "info",
            output: ["Running ..."]
        };
        $scope.isRunning = true;

        $http.post('/api/job/run',
            { entries: [ $scope.metadata ] }
        ).success(function (data) {
            $scope.result = {
                category: data.ok ? "success" : "danger",
                output: data.output
            };

            $scope.isRunning = false;
        }).error(function (data) {
            $scope.result = {
                category: 'danger',
                output: [data]
            };

            $scope.isRunning = false;
        });
    };

    $scope.onEditing(false);
});

app.controller('FileController', function($scope, $http) {
    $scope.pathDirs = [
        {
            index: 1,
            name: 'Home'
        }
    ];

    // $scope.fileEntries = [
    //     {
    //         name: 'bull',
    //         type: 'Directory',
    //         isDir: true,
    //         containNum: 14
    //     }, {
    //         name: 'balls',
    //         type: 'File',
    //         isDir: false,
    //         wtime: '2015-01-01 12:00'
    //     }
    // ];
    $scope.fileEntries = [];

    $scope.current = '';

    $scope.onPathDir = function (index) {
        var newDirs = []
        for (var i = 0; i < index; i++) {
            newDirs[i] = $scope.pathDirs[i];
        }

        $scope.pathDirs = newDirs;
        $scope.readDir(index)
    };

    $scope.onReadDir = function (entry) {
        if (!entry.isDir) {
            return;
        }

        var name = entry.name;
        var newDirs = []
        for (var i in $scope.pathDirs) {
            newDirs[i] = $scope.pathDirs[i];
        }
        var index = $scope.pathDirs.length + 1
        newDirs[$scope.pathDirs.length] = {
            'index': index,
            'name': name
        };

        $scope.pathDirs = newDirs;
        $scope.onPathDir(index);
    };

    $scope.onSwitch = function (entry) {
        $scope.current = entry.name;
    };

    $scope.readDir = function (index) {
        var path = ''
        for (var i = 1; i < index; i++) {
            path += ('/' + $scope.pathDirs[i].name);
        }

        if (index == 0) {
            path = '/'
        }

        $http.get('api/file/list' + path).success(function(data) {
            $scope.current = data.entries[0].name;
            $scope.fileEntries = data.entries;
        });
    };

    $scope.readDir($scope.pathDirs[0].index);
});

}()); // end of module
