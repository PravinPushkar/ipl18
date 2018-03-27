'use strict';

var app = angular.module('ipl');

/**
 * Run
 * 
 * The run block for the application.
 */
app.run(function ($rootScope, $window, $location) {
    // Restrict pages, except public pages, without a token
    // $rootScope.$on('$locationChangeStart', function () {
    //     var user = {};
    //     var publicPages = ['/login', '/register'];
    //     var restrictedPage = publicPages.indexOf($location.path()) === -1;
    //     user.name = $window.localStorage.getItem('displayName');
    //     user.token = $window.localStorage.getItem('token');
    //     if (restrictedPage && ((user.token === undefined || user.token === null) && (user.name === undefined || user.name === null))) {
    //         $location.path('/login');
    //     } else if (!restrictedPage && ((user.token !== undefined || user.token !== null) && (user.name !== undefined || user.name !== null))) {
    //         $location.path('/profile');
    //     }
    // });
});
