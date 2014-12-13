(function () {

var app = angular.module('dashboard', []);

app.controller('StateController', function($scope, $http) {
    $scope.state = {
        load: "N/A",
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

// app.controller('NameController', function($scope) {
//     $scope.name = 'shit';
// });

// app.controller('AjaxController', function($scope, $http) {
//     $http.get('api/hello').success(function(data) {
//         $scope.res = data;
//     });
// });

// app.controller('CommitController', function($scope, $http) {
//     $scope.commit = function (key) {
//         $http.post('api/hello', {
//             'key': key
//         }).success(function (data, status) {
//             console.log(data)
//         });
//     }
// });

}()); // end of module
