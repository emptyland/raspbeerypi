(function () {

var app = angular.module('dashboard', []);

var navBars = {
    'dashboard': {
        visibility: 'visible',
        active: 'active',
        dom: null
    }, 'console': {
        visibility: 'hidden',
        active: '',
        dom: null
    }, 'devices': {
        visibility: 'hidden',
        active: '',
        dom: null
    }
};

app.controller('NavController', function($scope) {
    $scope.switchNavBar = function (index) {

        for (k in navBars) {
            if (navBars[k].active === 'active') {
                navBars[k].dom = $('#nav-' + k).remove()
            }

            navBars[k].visibility = 'hidden';
            navBars[k].active = '';
        }
        navBars[index].visibility = 'visible';
        navBars[index].active = 'active';

        $('#nav-parent').append(navBars[index].dom);

        $scope.navBars = navBars;
    };

    for (k in navBars) {
        if (navBars[k].active != 'active') {
            navBars[k].dom = $('#nav-' + k).remove()
        }
    }
    $scope.switchNavBar('dashboard');
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

        for (i in $scope.usageEntries) {
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
            title: "Crontab 1",
            desc: "Demo Job"
        }, {
            title: "Crontab 2",
            desc: "Demo Job"
        }
    ];
});

}()); // end of module
