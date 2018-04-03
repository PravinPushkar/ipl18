'use strict';

var app = angular.module('ipl');

/**
 * Run
 * 
 * The run block for the application.
 */
app.run(['$rootScope', '$window', '$location', function ($rootScope, $window, $location) {
    // Restrict pages, except public pages, without a token
    $rootScope.$on('$locationChangeStart', function () {
        var publicPages = ['/login', '/register'];
        var restrictedPage = publicPages.indexOf($location.path()) === -1;
        var iNumber = $window.localStorage.getItem('iNumber');
        var token = $window.localStorage.getItem('token');
        if (restrictedPage && ((token === undefined || token === null) && (iNumber === undefined || iNumber === null))) {
            $location.path('/login');
        }
    });
}]);