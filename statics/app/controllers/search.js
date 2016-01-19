'use strict';

angular.module('LA').controller('SearchCtrl', function($scope, $location, $cookies, $routeParams, SearchService) {

    $scope.itemsPerPage = 100;
    $scope.currentPage = 1;
    $scope.pageCount = 0;
    $scope.total = 0;
    $scope.query = {};
    $scope.loadingLogs = false;

    $scope.findLogs = function() {
        $scope.query.s = ($scope.currentPage - 1) * $scope.itemsPerPage;
        $scope.query.l = $scope.itemsPerPage;
        $scope.query.from = $('#date-from-val').val();
        $scope.query.to = $('#date-to-val').val();

        $scope.loadingLogs = true;
        SearchService.findLogs($scope.query)
            .success(function(response) {
                $scope.logs = response.data;
                $scope.total = response.total;
                $scope.loadingLogs = false;
            })
            .error(function() {
                $scope.loadingLogs = false;
            });
    };

    $scope.resetSearch = function() {
        $scope.logs = undefined;
    };

    $scope.firstRun = true;
    $scope.$watch('currentPage', function() {
        if ($scope.firstRun) {
            $scope.firstRun = false;
            return;
        }

        $scope.findLogs();
    });

    $('#date-from').datetimepicker({locale: 'pl', format: 'YYYY-MM-DDTHH:mm:ss'});
    $('#date-to').datetimepicker({locale: 'pl', format: 'YYYY-MM-DDTHH:mm:ss'});
});
