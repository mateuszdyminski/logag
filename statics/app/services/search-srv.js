'use strict';

LA.service('SearchService', ['$http',
    function($http) {
        this.findLogs = function(query) {
            return $http({
                url: '/api/logs',
                method: "GET",
                params: query
            });
        };
    }
]);