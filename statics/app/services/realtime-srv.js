'use strict';

LA.service('RealtimeService', ['$http',
    function($http) {
        this.registerFilter = function(id, filter) {
            return $http({
                url: '/wsapi/filter/' + id,
                method: "POST",
                data: filter
            });
        };

        this.unregisterFilter = function(id) {
            return $http({
                url: '/wsapi/filter/' + id,
                method: "DELETE"
            });
        };
    }
]);