//'use strict';

(function () {

var app = angular.module('dashboard', ['ngCookies']);

app.controller('NavController', function($scope, $cookies) {
    var navBars = {
        'dashboard': {}, 'console': {}, 'devices': {},

        current: 'dashboard'
    };

    $scope.switchNavBar = function (index) {
        console.log(navBars);

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
    $scope.jobs = [
        {
            title: 'Crontab 1',
            index: 0,
            desc: 'Demo Job',
            lang: 'bash',
            cron: '1 * * * *',
            enable: true
        }, {
            title: "Crontab 2",
            index: 1,
            desc: "Demo Job",
            lang: 'python',
            cron: '1 1 * * *',
            enable: false
        }
    ];

    $scope.onEdit = function (index) {
        console.log('onEdit');
        console.log(index);
    };

    $scope.onEnable = function (index) {
        console.log('onEnable');
        console.log(index);
    }

    $scope.onDelete = function (index) {
        console.log('onDelete');
        console.log(index);
    }
});

}()); // end of module
