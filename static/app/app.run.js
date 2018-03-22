'use strict';

var app = angular.module('ipl');

/**
 * Run
 * 
 * The run block for the application.
 */
app.run(function ($http, $rootScope, $window, $location, displayName, token) {
    $http.defaults.headers.common.Accept = 'application/json';
    $http.defaults.headers.post['Content-Type'] = 'application/json';

    // Restrict pages, except public pages, without a token
    // $rootScope.$on('$locationChangeStart', function () {
    //     var user = {};
    //     var publicPages = ['/login', '/register'];
    //     var restrictedPage = publicPages.indexOf($location.path()) === -1;
    //     user.name = $window.localStorage.getItem(displayName);
    //     user.token = $window.localStorage.getItem(token);
    //     if (restrictedPage && ((user.token === undefined || user.token === null) && (user.name === undefined || user.name === null))) {
    //         $location.path('/login');
    //     }
    // });
});