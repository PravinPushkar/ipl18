'use strict';

var app = angular.module('ipl');

/**
 * Requests factory service
 * 
 * This service contains all the http requests made
 * to the back-end.
 */
app.factory('requestsService', function ($http) {
    var service = {};

    service.registerUser = registerUser;
    return service;

    function success(res) {
        console.log('success', res);
    }

    function error(err) {
        console.log('error', err);
    }

    function registerUser(data) {
        return $http.post('/pub/register', data)
            .then(success, error);
    }
});