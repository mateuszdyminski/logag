'use strict';

var LA = angular.module('LA', ['ngRoute', 'ngAnimate', 'ngCookies', 'ui.bootstrap', 'toastr'])
    .config(function($routeProvider) {
        $routeProvider
            .when('/search', {
                templateUrl: 'app/views/logSearch.html',
                controller: 'SearchCtrl'
            })
            .when('/real-time', {
                templateUrl: 'app/views/realTime.html',
                controller: 'RealTimeCtrl'
            })
            .otherwise({
                redirectTo: '/search'
            });
    }).filter('formatDate', function() {
        return function(date) {
            return moment(date).format("YYYY-MM-DD HH:mm:ss.SSS");
        }
    });