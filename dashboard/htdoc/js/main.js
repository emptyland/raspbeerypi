
var app = angular.module('dashboard', []);

app.controller('NameController', function($scope) {
    $scope.name = 'shit';
});

app.controller('AjaxController', function($scope, $http) {
    $http.get('api/hello').success(function(data) {
        $scope.res = data;
    });
});

app.controller('CommitController', function($scope, $http) {
    $scope.commit = function (key) {
        $http.post('api/hello', {
            'key': key
        }).success(function (data, status) {
            console.log(data)
        });
    }
});
