'use strict';

var app = angular.module('ipl');

/**
 * Run
 * 
 * The run block for the application.
 */
app.run(function ($window, $http) {
    // add JWT token as default auth header
    $http.defaults.headers.common.Authorization = 'Bearer ' + $window.jwtToken;
    $http.defaults.headers.common.Accept = 'application/json';
    $http.defaults.headers.post['Content-Type'] = 'application/json';
});