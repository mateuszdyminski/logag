'use strict';

angular.module('LA').controller('RealTimeCtrl', function($scope, $timeout, RealtimeService) {
    $scope.logs = [];
    $scope.active_filter = {};
    $scope.filter = {};

    // taken from http://stackoverflow.com/questions/105034/create-guid-uuid-in-javascript
    function guid() {
        function s4() {
            return Math.floor((1 + Math.random()) * 0x10000)
                .toString(16)
                .substring(1);
        }
        return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
            s4() + '-' + s4() + s4() + s4();
    }


    function getId() {
        var conn_id = window.sessionStorage['logag_connection_id'];
        if (!conn_id) {
            conn_id = guid();
            window.sessionStorage['logag_connection_id'] = conn_id;
        }

        $scope.id = conn_id;

        return conn_id;
    }

    var ws = new WebSocket('ws://127.0.0.1:8001/wsapi/ws/' + getId());

    ws.onmessage = function(event) {
        $timeout(function() {
            $scope.logs.push(JSON.parse(event.data));
        });
    };

    $scope.registerFilter = function() {
        if ($scope.filter.keywordsText) {
            $scope.filter.keywords  = $scope.filter.keywordsText.split(",");
        }
        $scope.filter.id = $scope.id;

        RealtimeService.registerFilter($scope.id, $scope.filter)
            .success(function() {
                angular.copy($scope.filter, $scope.active_filter);
            })
            .error(function(err) {
                console.log(err);
            });
    };

    $scope.unregisterFilter = function() {
        RealtimeService.unregisterFilter($scope.id)
            .success(function() {
                $scope.active_filter = {};
                $scope.filter.keywordsText = "";
                $scope.filter.level = "";
            })
            .error(function(err) {
                console.log(err);
            });
    };

    $scope.clearResults = function() {
        $scope.logs.length = 0;
    };
});
